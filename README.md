# Flux

Flux collecte, normalise et conserve des données de marché. Les décisions d’architecture du projet sont consignées dans [`docs/decisions/`](docs/decisions/).

## Démarrage local

### Prérequis

- Docker et Docker Compose disponibles, avec le daemon en cours d’exécution.
- GNU Make pour les commandes de développement.
- Une connexion réseau au premier lancement pour télécharger les images et les modules Go.

Go n’a pas besoin d’être installé sur la machine hôte. Le service actuel charge et valide sa configuration, puis attend un signal d’arrêt. Il ne se connecte pas encore à PostgreSQL ni à Kraken.

### Préparer la configuration

Depuis la racine du dépôt :

```sh
cp flux/.env.example flux/.env
```

Remplacer `FLUX_DB_PASSWORD=change-me-locally` par un mot de passe local. Cette valeur d’exemple n’est pas un secret utilisable. Le fichier `flux/.env` est ignoré par Git et exclu du contexte de construction Docker. Les variables sont détaillées dans [`flux/CONFIGURATION.md`](flux/CONFIGURATION.md).

### Démarrer avec Docker Compose

```sh
cd flux
docker compose up -d --build
docker compose ps
docker compose logs -f flux-app
```

Compose lance `flux-db` (PostgreSQL 17) et `flux-app`. L’application attend que le healthcheck PostgreSQL réussisse. Elle reçoit `FLUX_DB_HOST=flux-db` ; PostgreSQL reste accessible uniquement sur le réseau Compose. Les données sont conservées dans le volume nommé `postgres_data`.

Pour arrêter les services sans supprimer les données :

```sh
docker compose down
```

`docker compose down -v` supprimerait aussi le volume PostgreSQL et ses données. Les logs de la base sont disponibles avec `docker compose logs -f flux-db`.

### Commandes de développement

Depuis la racine du dépôt :

```sh
make -C flux run
make -C flux test
make -C flux test-race
make -C flux vet
make -C flux help
```

`make -C flux run` lance le point d’entrée au premier plan dans un conteneur Go temporaire ; utiliser `Ctrl+C` pour l’arrêter. Le processus tente de charger `flux/.env` avec `godotenv`, sans remplacer les variables déjà définies. Le fichier peut être absent si l’environnement fournit les valeurs requises. Cette commande ne lance pas PostgreSQL ; utiliser Compose pour démarrer les deux services ensemble.

## Langage et runtime

Flux est développé en **Go 1.27**. Go a été retenu pour son modèle de concurrence adapté au Scheduler, à la Queue et au Worker Pool, son adéquation aux services réseau longue durée, son runtime léger et son déploiement sous forme de binaire compilé. Sa bibliothèque standard couvre les besoins HTTP, temps et concurrence du projet, et son écosystème offre une intégration PostgreSQL mature.

## Framework et modèle d’exécution

La décision P0.2 retient la bibliothèque standard Go pour l’infrastructure applicative et un service autonome à processus unique. Le Scheduler, la Collection Queue et le Worker Pool y fonctionneront ensemble ; le pool démarrera avec trois workers concurrents, selon une limite configurable. Voir la [décision P0.2](docs/decisions/0002-framework-and-execution-model.md).

## PostgreSQL et migrations

Flux utilisera **pgx v5.11.0** avec **pgxpool** et l’API native pgx. Les migrations seront des fichiers SQL versionnés gérés par **Goose v3.29.0**, exécutés par la CLI avant le démarrage de Flux. Aucun ORM ne sera utilisé. Voir la [décision P0.3](docs/decisions/0003-postgresql-driver-and-migrations.md).

## HTTP et tests

Flux utilisera `net/http` et `context.Context` avec des timeouts configurables. Les tests s’appuieront sur `testing`, `testify` v1 (`assert`, `require`) et `net/http/httptest`. Les tests de concurrence seront exécutés avec `go test -race ./...`. Aucun framework HTTP externe ne sera ajouté. Voir la [décision P0.4](docs/decisions/0004-http-client-and-testing.md).
