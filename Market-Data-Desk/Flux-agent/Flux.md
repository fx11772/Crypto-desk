---

type: system
status: decided
tags:
  - flux
  - market-data
  - data-collection

---

# Flux

## Rôle

**Flux** collecte, normalise et conserve continuellement des données de marché fiables. Il ne produit ni indicateur, ni signal, ni décision de trading : ces usages appartiennent au futur Indicator Engine.

## Architecture

```text
Kraken API → KrakenClient → MarketData → Collector → CandleRepository → PostgreSQL
```

Les données internes sont des [[Candle]]. Le fournisseur est isolé par [[KrakenClient]] et [[MarketData]] ; le stockage est isolé par [[CandleRepository]].

## Orchestration

```text
PostgreSQL → CollectionTargetRepository → Scheduler → CollectionQueue → WorkerPool → Collector
```

Un [[CollectionTarget]] identifie une série à collecter. [[Scheduler]] détermine quand le vérifier, [[CollectionQueue]] garde le travail en attente, [[WorkerPool]] limite la concurrence et [[Collector]] détermine quelles Candles manquent réellement.

## Composants et sources de vérité

| Sujet | Note canonique |
| --- | --- |
| Donnée OHLCV, normalisation et clôture | [[Candle]] |
| API Kraken et erreurs propres au fournisseur | [[Kraken API]], [[KrakenClient]] |
| Abstraction des fournisseurs et normalisation | [[MarketData]] |
| Rattrapage, Candles terminées et retry de récupération | [[Collector]] |
| Persistance des Candles | [[CandleRepository]], [[candles]] |
| Configuration des séries à collecter | [[CollectionTarget]], [[CollectionTargetRepository]], [[collection_targets]] |
| Configuration par environnement | [[Configuration]] |
| Échéances et exécution concurrente | [[Scheduler]], [[CollectionQueue]], [[WorkerPool]] |
| États globaux et indisponibilité PostgreSQL | [[FluxState]] |
| Journalisation opérationnelle | [[Observability]] |

## Décisions d’architecture

- Langage et runtime : Go 1.27. Son modèle de concurrence convient au Scheduler, à la Queue et au Worker Pool ; il est adapté aux services réseau longue durée, produit un binaire léger à déployer, fournit une bibliothèque standard solide pour HTTP, le temps et la concurrence, et dispose d’une intégration PostgreSQL mature.
- **P0.2 — Décision acceptée sur le framework et le modèle d’exécution :** Flux utilise la bibliothèque standard de Go pour son infrastructure applicative, sans framework externe. Le service cible fonctionne en continu dans un processus unique, où le Scheduler, la Collection Queue et le Worker Pool cohabitent. La concurrence utilise les goroutines pour les workers, les channels pour la queue, `context.Context` pour les annulations et timeouts, et les mécanismes de synchronisation natifs lorsque nécessaire. Le pool démarre avec trois workers concurrents, selon une limite configurable. La queue n’est pas persistante ; une seule instance active est prévue initialement, sans système distribué dans la première version. Des bibliothèques spécialisées, notamment pour PostgreSQL, restent possibles. Le point d’entrée actuel charge et valide sa configuration, puis attend `SIGINT` ou `SIGTERM` pour s’arrêter proprement ; le Scheduler, la Queue et le Worker Pool seront réalisés dans les phases d’orchestration.
- **P0.4 — Client HTTP et tests :** utiliser `net/http` (Go 1.27) avec `context.Context` et des timeouts configurables. Les tests utiliseront `testing`, `testify` v1 (`assert`, `require`) et `net/http/httptest` pour simuler les échanges sans appeler Kraken. Les tests de concurrence seront lancés avec `go test -race ./...`. Aucun framework HTTP externe ne sera utilisé. Ces choix visent à limiter les dépendances tout en couvrant les erreurs HTTP, les timeouts, les réponses invalides et les comportements concurrents du Scheduler et du Worker Pool.
- **P0.7 — Configuration par environnement :** les valeurs de PostgreSQL, Kraken et Flux sont lues depuis les variables `FLUX_*`; la configuration est validée au démarrage et aucun secret n’est codé en dur. Les détails et valeurs par défaut sont dans [[Configuration]].
- Première source : Kraken ; premier target : `kraken + BTC/USD + 15m`.
- Une Candle est identifiée par `exchange + symbol + timeframe + timestamp`; son timestamp est UTC et représente le début de sa période.
- Seules les Candles terminées sont persistées. PostgreSQL est l’unique source de vérité de progression.
- La collecte est incrémentale : la reprise part de la dernière Candle réellement persistée, ou du timestamp global de démarrage lors d’un premier lancement.
- Les échéances sont absolues, alignées sur les frontières UTC des timeframes, avec un délai de fermeture configurable de 10 secondes par défaut.
- Les doublons de Candles sont ignorés ; les doublons de targets sont des erreurs explicites de configuration.

## Exécution locale actuelle

`flux/compose.yaml` lance `flux-app` et PostgreSQL (`flux-db`) dans Docker Compose. La base dispose d’un healthcheck et du volume persistant `postgres_data` ; l’application reçoit `FLUX_DB_HOST=flux-db` et attend que la base soit saine avant de démarrer. `flux/.env.example` sert de modèle à `flux/.env`, ignoré par Git. Le point d’entrée reste actif jusqu’à un signal d’arrêt, mais ne se connecte pas encore à PostgreSQL ni à Kraken. Voir [[Configuration]] et le README du dépôt pour les commandes.

## Frontières

Flux ne contient pas de logique d’analyse, d’indicateurs ou de trading. Les détails de contrat et de comportement sont maintenus exclusivement dans les notes liées ci-dessus.

La checklist ordonnée de réalisation se trouve dans [[Implementation-Tasks]].
