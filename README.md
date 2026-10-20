# Flux

Flux collecte, normalise et conserve des données de marché. Les décisions d’architecture du projet sont consignées dans [`docs/decisions/0001-language-and-runtime.md`](docs/decisions/0001-language-and-runtime.md).

## Langage et runtime

Flux est développé en **Go 1.27**. Go a été retenu pour son modèle de concurrence adapté au Scheduler, à la Queue et au Worker Pool, son adéquation aux services réseau longue durée, son runtime léger et son déploiement sous forme de binaire compilé. Sa bibliothèque standard couvre les besoins HTTP, temps et concurrence du projet, et son écosystème offre une intégration PostgreSQL mature.
