# Configuration de Flux

Flux lit ses paramètres depuis les variables d’environnement préfixées `FLUX_`. Au démarrage, le point d’entrée tente de charger [`flux/.env`](.env) avec `godotenv` ; son absence n’empêche pas le démarrage si les variables requises sont déjà définies dans l’environnement. Depuis la racine du dépôt, copie [`flux/.env.example`](.env.example) vers `flux/.env` pour le développement local. Remplace `FLUX_DB_PASSWORD=change-me-locally` par une valeur locale. Le fichier `.env` est ignoré par Git ; ne commite jamais de vraie valeur secrète.

Les valeurs déjà définies dans l’environnement du processus prennent le pas sur les valeurs du fichier `.env`.

Avec Docker Compose, `flux/.env` est nécessaire pour l’interpolation des variables PostgreSQL et pour `env_file` du service applicatif. `compose.yaml` impose `FLUX_DB_HOST=flux-db` à l’application ; la base n’expose pas de port vers l’hôte. Le volume nommé `postgres_data` conserve les données après `docker compose down`. Les commandes de démarrage et d’arrêt sont dans le [README](../README.md#démarrage-local).

## Variables

| Variable | Valeur d’exemple | Règle |
| --- | --- | --- |
| `FLUX_ENV` | `development` | Environnement d’exécution |
| `FLUX_LOG_LEVEL` | `info` | Niveau de journalisation |
| `FLUX_DB_HOST` | `flux-db` avec Compose ; `localhost` hors Compose si la base est locale | Hôte PostgreSQL requis |
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

La configuration est validée au démarrage. Les valeurs par défaut de développement ne contiennent aucun identifiant secret ; `FLUX_DB_PASSWORD` doit être fourni localement. Le processus actuel reste actif jusqu’à `SIGINT` ou `SIGTERM`, mais ne se connecte pas encore à PostgreSQL ni à Kraken.
