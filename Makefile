.PHONY: html
html:
	go run ./cmd/primaza-adm/main.go list dependencies -o html primaza-mytenant > out/graph.html

.PHONY: vet
vet:
	go vet ./...
