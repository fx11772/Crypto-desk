GO_VERSION := 1.27.3
GO_IMAGE := golang:$(GO_VERSION)
.DEFAULT_GOAL := help

CONTAINER_GO := podman run --rm \
	-v "$(CURDIR):/app" \
	-w /app \
	$(GO_IMAGE)

.PHONY: run build test test-race test-integration fmt vet clean db-start migrate-up help

run:
	$(CONTAINER_GO) go run ./cmd/flux

build:
	$(CONTAINER_GO) go build -o flux ./cmd/flux

test:
	$(CONTAINER_GO) go test ./...

test-race:
	$(CONTAINER_GO) go test -race ./...

test-integration:
	./scripts/run-integration-tests.sh

fmt:
	$(CONTAINER_GO) go fmt ./...

vet:
	$(CONTAINER_GO) go vet ./...

clean:
	rm -f flux

db-start:
	podman compose up -d flux-db

migrate-up:
	podman compose run --rm flux-migrate up

help:
	@echo "Commandes disponibles :"
	@echo ""
	@echo "  run          Démarrer Flux"
	@echo "  build        Compiler Flux"
	@echo "  test         Exécuter les tests"
	@echo "  test-race    Exécuter les tests avec le Race Detector"
	@echo "  test-integration  Exécuter les tests d'intégration"
	@echo "  fmt          Formater le code Go"
	@echo "  vet          Analyser statiquement le code"
	@echo "  clean        Supprimer l'exécutable compilé"
	@echo "  db-start     Démarrer PostgreSQL localement"
	@echo "  migrate-up   Appliquer les migrations en attente"
	@echo "  help         Afficher cette aide"
