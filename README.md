# Flux

Flux collecte, normalise et conserve des données de marché. Les décisions d’architecture du projet sont consignées dans [`docs/decisions/`](docs/decisions/).

## Langage et runtime

Flux est développé en **Go 1.27**. Go a été retenu pour son modèle de concurrence adapté au Scheduler, à la Queue et au Worker Pool, son adéquation aux services réseau longue durée, son runtime léger et son déploiement sous forme de binaire compilé. Sa bibliothèque standard couvre les besoins HTTP, temps et concurrence du projet, et son écosystème offre une intégration PostgreSQL mature.

## Framework et modèle d’exécution

La proposition P0.2 retient la bibliothèque standard Go pour l’infrastructure applicative et un service autonome à processus unique. Le Scheduler, la Collection Queue et le Worker Pool y fonctionneront ensemble ; le pool démarrera avec trois workers concurrents, selon une limite configurable. Voir la [décision P0.2](docs/decisions/0002-framework-and-execution-model.md).
