GO_VERSION := 1.27.3
GO_IMAGE := golang:$(GO_VERSION)
.DEFAULT_GOAL := help

DOCKER_GO := docker run --rm \
	-v "$(CURDIR):/app" \
	-w /app \
	$(GO_IMAGE)

.PHONY: run build test test-race fmt vet clean help

run:
	$(DOCKER_GO) go run ./cmd/flux

build:
	$(DOCKER_GO) go build -o flux ./cmd/flux

test:
	$(DOCKER_GO) go test ./...

test-race:
	$(DOCKER_GO) go test -race ./...

fmt:
	$(DOCKER_GO) go fmt ./...

vet:
	$(DOCKER_GO) go vet ./...

clean:
	rm -f flux

help:
	@echo "Commandes disponibles :"
	@echo ""
	@echo "  run          Démarrer Flux"
	@echo "  build        Compiler Flux"
	@echo "  test         Exécuter les tests"
	@echo "  test-race    Exécuter les tests avec le Race Detector"
	@echo "  fmt          Formater le code Go"
	@echo "  vet          Analyser statiquement le code"
	@echo "  clean        Supprimer l'exécutable compilé"
	@echo "  help         Afficher cette aide"
