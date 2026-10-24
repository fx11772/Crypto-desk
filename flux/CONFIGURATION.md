# Configuration de Flux

Flux lit ses paramètres depuis les variables d’environnement préfixées `FLUX_`. `make -C flux run` charge aussi `flux/internal/.env` avec `godotenv`. Copie [`internal/.env.example`](internal/.env.example) vers `internal/.env` pour le développement local. Le fichier `.env` est ignoré par Git ; ne commite jamais de vraie valeur secrète.

Les valeurs déjà définies dans l’environnement du processus prennent le pas sur les valeurs du fichier `.env`.

## Variables

| Variable | Valeur d’exemple | Règle |
| --- | --- | --- |
| `FLUX_ENV` | `development` | Environnement d’exécution |
| `FLUX_LOG_LEVEL` | `info` | Niveau de journalisation |
| `FLUX_DB_HOST` | `localhost` | Hôte PostgreSQL requis |
| `FLUX_DB_PORT` | `5432` | Entier entre 1 et 65535 |
| `FLUX_DB_NAME` | `flux` | Base PostgreSQL requise |
| `FLUX_DB_USER` | `flux` | Utilisateur PostgreSQL requis |
| `FLUX_DB_PASSWORD` | `change-me-locally` | Obligatoire ; remplacer par le mot de passe local, ne pas utiliser cette valeur comme secret |
| `FLUX_DB_SSLMODE` | `disable` | `disable`, `allow`, `prefer`, `require`, `verify-ca` ou `verify-full` |
| `FLUX_KRAKEN_BASE_URL` | `https://api.kraken.com` | URL HTTPS requise |
| `FLUX_KRAKEN_TIMEOUT_SECONDS` | `15` | Entier positif |
| `FLUX_WORKER_COUNT` | `3` | Entier positif |
| `FLUX_CANDLE_CLOSE_DELAY_SECONDS` | `10` | Entier non négatif |
| `FLUX_HISTORY_START` | *(vide)* | Optionnel, date RFC3339 ; convertie en UTC |

La configuration est validée au démarrage. Les valeurs par défaut de développement ne contiennent aucun identifiant secret ; `FLUX_DB_PASSWORD` doit être fourni localement.
