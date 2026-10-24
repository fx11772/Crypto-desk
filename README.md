# Flux

Flux collecte, normalise et conserve des données de marché. Les décisions d’architecture du projet sont consignées dans [`docs/decisions/`](docs/decisions/).

## Langage et runtime

Flux est développé en **Go 1.27**. Go a été retenu pour son modèle de concurrence adapté au Scheduler, à la Queue et au Worker Pool, son adéquation aux services réseau longue durée, son runtime léger et son déploiement sous forme de binaire compilé. Sa bibliothèque standard couvre les besoins HTTP, temps et concurrence du projet, et son écosystème offre une intégration PostgreSQL mature.

## Framework et modèle d’exécution

La proposition P0.2 retient la bibliothèque standard Go pour l’infrastructure applicative et un service autonome à processus unique. Le Scheduler, la Collection Queue et le Worker Pool y fonctionneront ensemble ; le pool démarrera avec trois workers concurrents, selon une limite configurable. Voir la [décision P0.2](docs/decisions/0002-framework-and-execution-model.md).

## PostgreSQL et migrations

Flux utilisera **pgx v5.11.0** avec **pgxpool** et l’API native pgx. Les migrations seront des fichiers SQL versionnés gérés par **Goose v3.29.0**, exécutés par la CLI avant le démarrage de Flux. Aucun ORM ne sera utilisé. Voir la [décision P0.3](docs/decisions/0003-postgresql-driver-and-migrations.md).

## HTTP et tests

Flux utilisera `net/http` et `context.Context` avec des timeouts configurables. Les tests s’appuieront sur `testing`, `testify` v1 (`assert`, `require`) et `net/http/httptest`. Les tests de concurrence seront exécutés avec `go test -race ./...`. Aucun framework HTTP externe ne sera ajouté. Voir la [décision P0.4](docs/decisions/0004-http-client-and-testing.md).

## Commandes de développement

Les commandes Makefile s’exécutent depuis la racine du dépôt avec `make -C flux <cible>`. La liste est disponible via `make -C flux help` : `run`, `build`, `test`, `test-race`, `fmt`, `vet` et `clean`. `run` exécute actuellement le point d’entrée au premier plan dans un conteneur temporaire ; le programme se termine après son message de démarrage, donc aucun arrêt séparé n’est nécessaire.

## Configuration locale

L’exemple sans secret se trouve dans [`flux/internal/.env.example`](flux/internal/.env.example). Pour le démarrage local, copie-le vers `flux/internal/.env`, puis remplace `FLUX_DB_PASSWORD=change-me-locally` par le mot de passe de ta base PostgreSQL locale. Le fichier `.env` local est ignoré par Git. La liste des variables et leurs règles est détaillée dans [`flux/CONFIGURATION.md`](flux/CONFIGURATION.md).
