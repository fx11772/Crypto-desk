# Décision 0001 : langage et version du runtime

- **Statut :** acceptée
- **Date :** 2026-10-20
- **Périmètre :** service Flux

## Contexte

Flux est un service réseau de longue durée. Son architecture comprend un Scheduler, une Queue et un Worker Pool, ainsi que des intégrations HTTP avec Kraken et PostgreSQL.

## Décision

Le langage de développement de Flux est **Go**, avec **Go 1.27** comme version de runtime choisie.

## Raisons

- Le modèle de concurrence de Go convient au Scheduler, à la Queue et au Worker Pool.
- Go est adapté aux services réseau qui fonctionnent en continu.
- Le runtime est léger et le service peut être déployé sous la forme d’un binaire compilé.
- La bibliothèque standard fournit les outils nécessaires pour HTTP, le temps et la concurrence.
- L’écosystème Go offre une intégration PostgreSQL mature.

## Conséquences

- Le code applicatif et les outils de développement de Flux utiliseront Go 1.27.
- Les choix de framework, de driver PostgreSQL, de client HTTP et de bibliothèque de tests restent à documenter dans les tâches P0.2 à P0.4.
