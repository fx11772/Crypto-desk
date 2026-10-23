GO_VERSION := 1.27.3
GO_IMAGE := golang:$(GO_VERSION)
.DEFAULT_GOAL := help

DOCKER_GO := docker run --rm \
	-v "$(CURDIR):/app" \
	-w /app \
	$(GO_IMAGE)
DOCKER_APP := docker run --rm --name flux-dev \
	-v "$(CURDIR):/app" \
	-w /app \
	$(GO_IMAGE)

.PHONY: install start run build verify test test-race fmt vet stop clean help

install:
	$(DOCKER_GO) go mod download

start:
	$(DOCKER_APP) go run ./cmd/flux

run: start

build:
	$(DOCKER_GO) go build -o flux ./cmd/flux

test:
	$(DOCKER_GO) go test ./...

verify: test vet

test-race:
	$(DOCKER_GO) go test -race ./...

fmt:
	$(DOCKER_GO) go fmt ./...

vet:
	$(DOCKER_GO) go vet ./...

clean:
	rm -f flux

stop:
	@docker stop flux-dev 2>/dev/null || echo "Flux n’est pas en cours d’exécution."

help:
	@echo "Commandes disponibles :"
	@echo ""
	@echo "  install      Télécharger les dépendances Go"
	@echo "  start        Démarrer Flux au premier plan"
	@echo "  stop         Arrêter le conteneur Flux"
	@echo "  build        Compiler Flux"
	@echo "  verify       Exécuter les tests et go vet"
	@echo "  test         Exécuter les tests"
	@echo "  test-race    Exécuter les tests avec le Race Detector"
	@echo "  fmt          Formater le code Go"
	@echo "  vet          Analyser statiquement le code"
	@echo "  clean        Supprimer l'exécutable compilé"
	@echo "  help         Afficher cette aide"
