# Décision 0002 : framework et modèle d’exécution

- **Statut :** proposition
- **Date :** 2026-10-20
- **Périmètre :** service Flux

## Contexte

Flux est un service autonome de collecte continue. Il orchestre un Scheduler, une Collection Queue et un Worker Pool pour exécuter des tâches réseau.

## Proposition

### Framework

Flux utilisera exclusivement la bibliothèque standard de Go pour son infrastructure applicative, sans framework externe. Des bibliothèques spécialisées pourront être ajoutées lorsque nécessaire, notamment pour PostgreSQL.

### Modèle d’exécution

Flux fonctionnera en continu comme un service autonome dans un processus unique. Le Scheduler, la Collection Queue et le Worker Pool seront réunis dans ce processus.

La concurrence s’appuiera sur les mécanismes natifs de Go :

- les goroutines pour les workers ;
- les channels pour la Collection Queue ;
- `context.Context` pour les annulations et les timeouts ;
- les mécanismes de synchronisation natifs lorsque nécessaire.

Le Worker Pool démarrera avec trois workers concurrents. Cette limite restera configurable.

## Justifications

- Ce modèle convient aux tâches réseau et à la collecte continue.
- Un processus unique limite la complexité opérationnelle initiale.
- L’infrastructure applicative ne dépend pas d’un framework externe.
- La concurrence peut être contrôlée sans infrastructure externe.
- Le modèle est compatible avec les mécanismes de reprise définis pour Flux.
- L’architecture pourra évoluer ultérieurement si les besoins le justifient.

## Limites acceptées

- La Collection Queue n’est pas persistante.
- Une seule instance de Flux sera active initialement.
- La première version n’introduira pas de système distribué.
