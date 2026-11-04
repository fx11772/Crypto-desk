---

type: component
status: decided
tags:
  - flux
  - configuration
  - environment

---

# Configuration

Flux lit sa configuration depuis les variables d’environnement préfixées `FLUX_`. Les valeurs non secrètes ont des valeurs par défaut ; le mot de passe PostgreSQL doit être fourni explicitement. La configuration est validée au démarrage.

Au démarrage, le point d’entrée tente de charger `flux/.env` avec `godotenv`. Le fichier local est ignoré par Git et exclu du contexte Docker ; l’exemple versionné sans secret est `flux/.env.example`. Copie-le vers `flux/.env` et remplace `FLUX_DB_PASSWORD` par un mot de passe local. Les variables déjà présentes dans l’environnement du processus ne sont pas remplacées. L’absence du fichier `.env` est acceptée si les variables requises sont définies autrement ; Docker Compose exige ce fichier pour `env_file` et l’interpolation de PostgreSQL. Le processus reste actif jusqu’à `SIGINT` ou `SIGTERM`.

## Variables

| Variable | Défaut | Règle |
| --- | --- | --- |
| `FLUX_ENV` | `development` | Environnement d’exécution |
| `FLUX_LOG_LEVEL` | `info` | Niveau de journalisation |
| `FLUX_DB_HOST` | `flux-db` avec Compose ; `localhost` hors Compose si la base est locale | Hôte PostgreSQL requis |
| `FLUX_DB_PORT` | `5432` | Entier entre 1 et 65535 |
| `FLUX_DB_NAME` | `flux` | Base PostgreSQL requise |
| `FLUX_DB_USER` | `flux` | Utilisateur PostgreSQL requis |
| `FLUX_DB_PASSWORD` | aucun | Obligatoire ; ne pas committer sa valeur |
| `FLUX_DB_SSLMODE` | `disable` | `disable`, `allow`, `prefer`, `require`, `verify-ca` ou `verify-full` |
| `FLUX_KRAKEN_BASE_URL` | `https://api.kraken.com` | URL HTTPS requise |
| `FLUX_KRAKEN_TIMEOUT_SECONDS` | `15` | Entier positif |
| `FLUX_WORKER_COUNT` | `3` | Entier positif |
| `FLUX_CANDLE_CLOSE_DELAY_SECONDS` | `10` | Entier non négatif |
| `FLUX_HISTORY_START` | aucun | Optionnel, au format RFC3339 ; converti en UTC |

## Chargement et validation

`config.Load()` construit les configurations PostgreSQL, Kraken et Flux depuis l’environnement, applique les défauts ci-dessus et rejette les valeurs invalides. Le mot de passe n’est jamais inclus dans les messages de validation ni dans les logs de démarrage.

## Environnement Docker Compose

`flux/compose.yaml` démarre `flux-app` et `flux-db` sur le réseau Compose. L’application attend que PostgreSQL soit sain, reçoit `FLUX_DB_HOST=flux-db` et utilise les variables de `flux/.env`. La base n’expose pas de port hôte ; le volume nommé `postgres_data` conserve ses données après `docker compose down`. Le point d’entrée actuel ne se connecte pas encore à PostgreSQL ni à Kraken. Les commandes `up`, `logs` et `down` sont dans le README du dépôt.
