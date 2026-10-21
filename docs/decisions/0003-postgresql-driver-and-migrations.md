# Décision 0003 : driver PostgreSQL et migrations

- **Statut :** acceptée
- **Date :** 2026-10-20
- **Périmètre :** accès PostgreSQL et gestion du schéma de Flux

## Décision

| Élément | Choix |
| --- | --- |
| Driver PostgreSQL | pgx v5.11.0 |
| Gestion des connexions | `pgxpool` |
| API PostgreSQL | API native pgx |
| Outil de migrations | Goose v3.29.0 |
| Format des migrations | Fichiers SQL versionnés |
| Exécution des migrations | CLI, avant le démarrage de Flux |
| ORM | Aucun |

## Conséquences

- Les accès PostgreSQL de Flux utiliseront l’API native de pgx et le pool `pgxpool`.
- Le schéma sera modifié au moyen de migrations SQL versionnées exécutées avec Goose avant le lancement du service.
- Le projet ne reposera pas sur un ORM.
