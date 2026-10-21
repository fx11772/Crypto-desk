# Décision 0004 : client HTTP et bibliothèque de tests

- **Statut :** acceptée
- **Date :** 2026-10-20
- **Périmètre :** appels HTTP et tests de Flux

## Décision

| Composant | Choix |
| --- | --- |
| Client HTTP | `net/http` de Go 1.27 |
| Timeouts et annulations | `context.Context` et timeouts configurables |
| Framework de tests | `testing` de la bibliothèque standard Go |
| Assertions | `testify` v1 (`assert`, `require`) |
| Simulation HTTP | `net/http/httptest` |
| Vérification des conditions de concurrence | `go test -race ./...` |
| Framework HTTP externe | Aucun |

## Justifications

- Limiter les dépendances externes et utiliser les fonctionnalités natives de Go.
- Tester `KrakenClient` sans communiquer avec Kraken.
- Couvrir les erreurs HTTP, les timeouts et les réponses invalides.
- Vérifier les comportements concurrents du Scheduler et du Worker Pool.

## Conséquences

- Les appels HTTP seront construits avec `net/http`, avec des délais configurables et des contextes d’annulation.
- Les tests HTTP utiliseront des serveurs simulés par `httptest`.
- La suite de tests pourra être lancée avec le détecteur de courses de données Go via `go test -race ./...`.
