#!/bin/env bash

set -e

[ $(command -v primazactl) > /dev/null ] || ( cat << EOF && exit 1
please install primazactl:

mkdir -p ./bin
curl -sL0 https://github.com/primaza/primazactl/releases/download/latest/primazactl > ./bin/primazactl
export PATH=\$PATH:$(realpath ./bin)
EOF
)

TENANT=primaza-mytenant
MAIN_CLUSTER=primaza-tools-main
MAX_RETRY=10
PRIMAZA_BRANCH=main
TMP_DIR=./out/tmp
PRIMAZA_DIR="$TMP_DIR/primaza"

## Create cluster
kind delete cluster --name "$MAIN_CLUSTER"
kind create cluster --name "$MAIN_CLUSTER"

# Build Primaza
mkdir -p "$TMP_DIR"
# rm -rf "$PRIMAZA_DIR"
( cd "$TMP_DIR" && git clone git@github.com:primaza/primaza.git && git checkout "$PRIMAZA_BRANCH") || \
    ( cd "$PRIMAZA_DIR" && git fetch origin && ( git reset --hard "origin/$PRIMAZA_BRANCH" || git reset --hard "$PRIMAZA_BRANCH" ))

IMG=ghcr.io/primaza/primaza:latest && ( cd "$PRIMAZA_DIR" && IMG=$IMG make primaza docker-build && kind load docker-image "$IMG" --name "$MAIN_CLUSTER" )
IMG=ghcr.io/primaza/primaza-agentapp:latest && ( cd "$PRIMAZA_DIR" && IMG=$IMG make agentapp docker-build && kind load docker-image "$IMG" --name "$MAIN_CLUSTER" )
IMG=ghcr.io/primaza/primaza-agentsvc:latest && ( cd "$PRIMAZA_DIR" && IMG=$IMG make agentsvc docker-build && kind load docker-image "$IMG" --name "$MAIN_CLUSTER" )


## Install Primaza
kubectl apply \
    -f "https://github.com/cert-manager/cert-manager/releases/download/v1.12.0/cert-manager.yaml"
kubectl rollout status -n cert-manager deploy/cert-manager-webhook -w --timeout=240s

n=0
until [ "$n" -ge "$MAX_RETRY" ]
do
    primazactl create tenant "$TENANT" --version latest && break
    n=$((n+1))
    sleep 10
done

INTERNAL_URL="https://$(docker container inspect $MAIN_CLUSTER-control-plane --format {{.NetworkSettings.Networks.kind.IPAddress}}):6443"
echo "$INTERNAL_URL"

n=0
until [ "$n" -ge "$MAX_RETRY" ]
do
    primazactl join cluster \
        --version latest \
        --tenant "$TENANT" \
        --cluster-environment self-demo \
        --environment demo \
        --internal-url "$INTERNAL_URL" && break
    n=$((n+1))
    sleep 10
done

primazactl create application-namespace applications \
        --version latest \
        --tenant "$TENANT" \
        --cluster-environment self-demo \
        --tenant-internal-url "$INTERNAL_URL"

primazactl create service-namespace services \
        --version latest \
        --tenant "$TENANT" \
        --cluster-environment self-demo \
        --tenant-internal-url "$INTERNAL_URL"

( cd "$PRIMAZA_DIR" && make primaza install )

## Seed Primaza tenant
( cd config/base && kustomize edit set namespace "$TENANT" )
kustomize build config/base | kubectl apply -f -
kustomize build config/services | kubectl create -f - -n "services"

queue_uid=$(kubectl get queues.sqs.services.k8s.aws my-queue -n services --output jsonpath='{.metadata.uid}')
kubectl annotate -n primaza-mytenant registeredservices.primaza.io sqs-queue-reader "primaza.io/service-uid=$queue_uid" --overwrite
kubectl annotate -n primaza-mytenant registeredservices.primaza.io sqs-queue-writer "primaza.io/service-uid=$queue_uid" --overwrite

n=0
until [ "$n" -ge "$MAX_RETRY" ]
do
    kustomize build config/myapp | kubectl apply -f - -n applications && break
    n=$((n+1))
    sleep 10
done

kubectl config set-context --current --namespace primaza-mytenant
