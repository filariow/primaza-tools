#!/bin/env bash

set -e

TENANT=primaza-mytenant
MAIN_CLUSTER=primaza-tools-main

## Create cluster
kind delete cluster --name "$MAIN_CLUSTER"
kind create cluster --name "$MAIN_CLUSTER"

## Install Primaza
kubectl apply \
    -f "https://github.com/cert-manager/cert-manager/releases/download/v1.12.0/cert-manager.yaml"
kubectl rollout status -n cert-manager deploy/cert-manager-webhook -w --timeout=120s

n=0
until [ "$n" -ge 5 ]
do
    primazactl create tenant "$TENANT" --version latest && break
    n=$((n+1))
    sleep 5
done



INTERNAL_URL="https://$(docker container inspect $MAIN_CLUSTER-control-plane --format {{.NetworkSettings.Networks.kind.IPAddress}}):6443"
echo "$INTERNAL_URL"

primazactl join cluster \
        --version latest \
        --tenant "$TENANT" \
        --cluster-environment self-demo \
        --environment demo \
        --internal-url "$INTERNAL_URL"

primazactl create application-namespace applications \
        --version latest \
        --tenant "$TENANT" \
        --cluster-environment self-demo \
        --tenant-internal-url "$INTERNAL_URL"

## Seed Primaza tenant
( cd config/base && kustomize edit set namespace "$TENANT" )
kustomize build config/base | kubectl apply -f -

( cd config/myapp && kustomize edit set namespace "applications" )
n=0
until [ "$n" -ge 5 ]
do
    kustomize build config/myapp | kubectl apply -f - && break
    n=$((n+1))
    sleep 5
done
