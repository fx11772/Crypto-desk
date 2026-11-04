---

type: implementation-plan
status: ready
tags:
  - flux
  - implementation
  - github

---

# Plan d’implémentation de Flux

Cette checklist est ordonnée pour être convertie en issues GitHub. Chaque case correspond à une tâche indépendante avec un résultat vérifiable. Créer un milestone par phase et conserver l’ordre indiqué à l’intérieur de chaque phase.

## Ordre et dépendances

```text
P0 Fondations
 ├── P1 Persistance PostgreSQL ─┐
 └── P2 Intégration Kraken ─────┴── P3 Collecteur
                                      ↓
                                   P4 Orchestration
                                      ↓
                                   P5 Résilience
                                      ↓
                                   P6 Exploitation et livraison
```

P1 et P2 peuvent être développées en parallèle après P0. P3 dépend des deux. Les autres phases sont séquentielles.

---

## P0 — Fondations

- [x] **P0.1 — Choisir le langage et la version de runtime.** Consigner la décision et les raisons dans la documentation du projet.
- [x] **P0.2 — Choisir le framework d’application et le modèle d’exécution.** Définir comment Flux démarre et s’arrête proprement.
- [x] **P0.3 — Choisir le driver PostgreSQL et l’outil de migrations.** Vérifier leur compatibilité avec le runtime choisi.
- [x] **P0.4 — Choisir le client HTTP et la bibliothèque de tests.** Consigner les choix dans le manifeste de dépendances.
- [x] **P0.5 — Créer le squelette applicatif minimal.** Ajouter le point d’entrée, la configuration du projet et les dossiers de code et tests selon la stack retenue.
- [x] **P0.6 — Ajouter les commandes standard du projet.** Documenter les cibles Makefile disponibles : `run`, `build`, `test`, `test-race`, `fmt`, `vet`, `clean` et `help`. `run` démarre le programme au premier plan ; il s’arrête sur `Ctrl+C`.
- [x] **P0.7 — Définir la configuration par environnement.** Prévoir les valeurs de PostgreSQL, Kraken et Flux sans secrets codés en dur.
- [x] **P0.8 — Ajouter un exemple de configuration sans secret.** Documenter les variables requises et leurs valeurs locales de démonstration.
- [x] **P0.9 — Ajouter un README de démarrage.** Expliquer les prérequis, l’installation, la configuration et le démarrage local.
- [x] **P0.10 — Ajouter une CI minimale.** Installer les dépendances et exécuter les contrôles automatisés disponibles sur chaque pull request.

**Sortie de phase :** le projet démarre avec sa configuration locale et les vérifications de base tournent dans la CI.

## P1 — Persistance PostgreSQL

- [x] **P1.0 — Mettre en place l’environnement Docker Compose de Flux** ([issue #97](https://github.com/fx11772/Crypto-desk/issues/97)). Ajouter `compose.yaml` et le Dockerfile de `flux-app`, avec le runtime Go retenu, un processus long-running et un arrêt propre. Orchestrer `flux-app` et `flux-db` avec configuration par environnement (`FLUX_DB_HOST=flux-db`), healthcheck PostgreSQL, volume persistant, réseau interne et politique de redémarrage. Réutiliser l’exemple existant à la racine comme unique `.env.example`, puis documenter `up`, `down` et `logs`. Dépend de P0.1, P0.2 et P0.5 à P0.8 ; précède P1.1 à P1.21.
- [ ] **P1.1 — Définir la configuration de connexion PostgreSQL.** Construire la connexion depuis la configuration sans exposer les identifiants dans les logs. Dépend de P1.0.
- [ ] **P1.2 — Ajouter la commande de démarrage PostgreSQL local.** Documenter un moyen reproductible de lancer la base pour le développement. Dépend de P1.0.
- [ ] **P1.3 — Créer la migration initiale de `candles`.** Appliquer les types et colonnes définis dans [[candles]].
- [ ] **P1.4 — Ajouter la clé primaire technique des candles.** Utiliser un identifiant généré par PostgreSQL.
- [ ] **P1.5 — Ajouter la contrainte unique métier des candles.** Garantir l’unicité `(exchange, symbol, timeframe, timestamp)`.
- [ ] **P1.6 — Créer la migration de `collection_targets`.** Appliquer les colonnes et types définis dans [[collection_targets]].
- [ ] **P1.7 — Ajouter les contraintes des targets.** Garantir l’unicité `(exchange, symbol, timeframe)` et les champs obligatoires.
- [ ] **P1.8 — Ajouter la commande de migration.** Permettre d’appliquer les migrations sur une base vierge et une base existante. Dépend de P1.0.
- [ ] **P1.9 — Définir le modèle interne `Candle`.** Représenter les champs OHLCV, timestamp UTC, exchange, symbol et timeframe.
- [ ] **P1.10 — Définir le modèle interne `CollectionTarget`.** Représenter l’identifiant, les champs métier, l’activation et les dates.
- [ ] **P1.11 — Implémenter `CandleRepository.get_last()`.** Retourner la candle la plus récente du target, ou l’absence de résultat.
- [ ] **P1.12 — Implémenter l’insertion transactionnelle de candles.** Insérer un lot atomiquement et ignorer les conflits sur l’identité métier.
- [ ] **P1.13 — Retourner `SaveCandlesResult`.** Renseigner `inserted_count` et `duplicate_count` à partir du résultat réel de la base.
- [ ] **P1.14 — Traduire les erreurs du driver dans `CandleRepository`.** Exposer `PersistenceTemporaryError` ou `PersistencePermanentError`.
- [ ] **P1.15 — Implémenter `CollectionTargetRepository.get_enabled()`.** Retourner uniquement les targets activés.
- [ ] **P1.16 — Implémenter `get_all()` et `get_by_id()`.** Fournir la lecture de toute la configuration et d’un target précis.
- [ ] **P1.17 — Implémenter création et mise à jour de targets.** Faire retourner l’objet persistant et actualiser `updated_at` sur modification.
- [ ] **P1.18 — Implémenter `set_enabled()`.** Modifier l’activation sans supprimer le target ni son historique.
- [ ] **P1.19 — Signaler explicitement les targets dupliqués.** Convertir la violation d’unicité en erreur de configuration interne.
- [ ] **P1.20 — Traduire les erreurs du driver dans `CollectionTargetRepository`.** Utiliser les catégories internes de persistance.
- [ ] **P1.21 — Vérifier les migrations et repositories sur PostgreSQL.** Couvrir le démarrage vierge, les lectures, les insertions, les doublons et le rollback transactionnel. Dépend de P1.0.

**Sortie de phase :** `flux-app` et `flux-db` fonctionnent dans Docker Compose avec configuration documentée, PostgreSQL sain et données persistantes ; les deux tables sont créées par migration et leurs repositories fournissent les contrats documentés.

## P2 — Intégration Kraken et données internes

- [ ] **P2.1 — Définir les unités internes de timeframe.** Convertir les formes supportées (`1m`, `15m`, `1h`, etc.) vers une durée exploitable.
- [ ] **P2.2 — Définir les types d’erreurs de marché internes.** Distinguer erreur temporaire, erreur permanente et réponse invalide.
- [ ] **P2.3 — Créer le client HTTP Kraken.** Centraliser URL, timeout et transport dans un client injectable.
- [ ] **P2.4 — Implémenter la requête Kraken OHLC.** Construire les paramètres à partir d’un target et d’un point de départ.
- [ ] **P2.5 — Décoder la réponse Kraken dans un type brut privé au client.** Rejeter les réponses mal formées.
- [ ] **P2.6 — Traduire les erreurs Kraken.** Mapper réseau/serveur/rate limit en temporaire, paramètres invalides en permanent, payload incohérent en invalide.
- [ ] **P2.7 — Normaliser le symbole vers `BASE/QUOTE`.** Ne pas laisser le format fournisseur fuiter dans le modèle interne.
- [ ] **P2.8 — Normaliser les timestamps en UTC.** Interpréter chaque timestamp comme le début de la période.
- [ ] **P2.9 — Convertir les valeurs OHLCV sans perte de précision.** Employer un type décimal exact adapté à la stack.
- [ ] **P2.10 — Implémenter `MarketData`.** Exposer au Collector une interface fournisseur indépendante de Kraken.
- [ ] **P2.11 — Rejeter un lot si une candle ne peut pas être normalisée.** Ne retourner aucun résultat partiel ni candle fabriquée.
- [ ] **P2.12 — Vérifier la normalisation avec des réponses Kraken représentatives.** Couvrir réponses valides, anomalies de schéma, timestamps invalides et valeurs invalides.
- [ ] **P2.13 — Vérifier la traduction des erreurs avec un transport HTTP simulé.** Couvrir erreur temporaire, erreur permanente, rate limit et timeout.

**Sortie de phase :** `MarketData` retourne des candles internes valides ou des erreurs internes sans exposer le format Kraken.

## P3 — Collecteur incrémental

- [ ] **P3.1 — Définir le contrat d’appel de collecte d’un target.** Le traitement reçoit un `CollectionTarget` et une horloge injectable.
- [ ] **P3.2 — Implémenter le calcul de la prochaine période attendue.** Partir de la dernière candle persistée lorsqu’elle existe.
- [ ] **P3.3 — Implémenter le départ initial global.** Utiliser une configuration globale lorsque le repository ne retourne aucune candle.
- [ ] **P3.4 — Déterminer la borne de récupération.** Utiliser une heure UTC cohérente pour l’ensemble d’un traitement.
- [ ] **P3.5 — Demander les données manquantes à `MarketData`.** Le Collector ne connaît pas les paramètres ni payloads Kraken.
- [ ] **P3.6 — Filtrer les candles encore ouvertes.** Persister uniquement celles dont `timestamp + timeframe <= now UTC`.
- [ ] **P3.7 — Envoyer le lot valide à `CandleRepository`.** Ne transmettre aucun lot déclaré invalide par MarketData.
- [ ] **P3.8 — Ne pas faire avancer le point de reprise après échec.** La progression future doit toujours relire PostgreSQL.
- [ ] **P3.9 — Gérer une réponse vide.** Terminer le traitement sans créer de candle artificielle.
- [ ] **P3.10 — Vérifier premier démarrage et reprise.** Simuler absence de données, retard de plusieurs périodes, chevauchement et redémarrage.
- [ ] **P3.11 — Vérifier la clôture temporelle.** Couvrir exactement avant, à et après la frontière de timeframe.

**Sortie de phase :** le Collector rattrape l’historique manquant, ne sauvegarde que des candles clôturées et utilise l’état persistant comme référence.

## P4 — Configuration des targets et orchestration

- [ ] **P4.1 — Charger les targets activés au démarrage.** Utiliser `CollectionTargetRepository.get_enabled()`.
- [ ] **P4.2 — Implémenter le calcul de la prochaine échéance UTC.** Aligner chaque target sur les frontières de son timeframe.
- [ ] **P4.3 — Ajouter le délai de fermeture configurable.** Utiliser 10 secondes comme valeur initiale.
- [ ] **P4.4 — Implémenter la boucle du Scheduler.** Calculer les prochaines échéances absolues sans dérive cumulative.
- [ ] **P4.5 — Implémenter la queue interne de targets dus.** Ajouter et consommer des targets sans y stocker les candles.
- [ ] **P4.6 — Implémenter un worker.** Prendre un target et appeler le Collector une fois.
- [ ] **P4.7 — Créer le pool configurable.** Utiliser trois workers par défaut.
- [ ] **P4.8 — Empêcher deux traitements simultanés du même target.** Définir et vérifier le comportement de déduplication en queue/en cours.
- [ ] **P4.9 — Gérer l’arrêt propre de l’orchestrateur.** Arrêter le Scheduler et laisser les traitements en cours terminer proprement.
- [ ] **P4.10 — Vérifier échéances simultanées et saturation.** Les travaux excédant la concurrence restent en attente.
- [ ] **P4.11 — Vérifier les targets désactivés.** Ils ne sont ni planifiés ni supprimés avec leur historique.

**Sortie de phase :** les targets configurés sont vérifiés à échéance UTC et traités avec une concurrence plafonnée.

## P5 — Résilience et états globaux

- [ ] **P5.1 — Implémenter le calcul de backoff configurable.** Valeurs initiales 2, 4 et 8 secondes, jitter ±20 %.
- [ ] **P5.2 — Respecter `Retry-After`.** Il remplace le délai calculé sans offrir de tentative supplémentaire.
- [ ] **P5.3 — Ajouter le compteur maximal de retries fournisseur.** Trois retries après la tentative initiale ; ne pas réessayer les erreurs permanentes ou invalides.
- [ ] **P5.4 — Implémenter les transitions d’état global.** Définir les changements autorisés entre `RUNNING`, `DEGRADED` et `ERROR`.
- [ ] **P5.5 — Déclencher `DEGRADED` après épuisement des retries PostgreSQL temporaires.** La transition est globale à Flux.
- [ ] **P5.6 — Suspendre l’orchestration en `DEGRADED`.** Le Scheduler n’ajoute plus, les workers ne démarrent pas et la queue en attente est abandonnée.
- [ ] **P5.7 — Vérifier l’état avant toute persistance après un appel externe.** Abandonner le lot si Flux est devenu `DEGRADED`.
- [ ] **P5.8 — Ajouter le health check PostgreSQL global.** Requête légère, intervalle initial configurable de 10 secondes.
- [ ] **P5.9 — Recharger les targets actifs après un health check réussi.** Ne pas reprendre avec une configuration mémorisée avant la panne.
- [ ] **P5.10 — Reprendre l’orchestration depuis PostgreSQL.** Recalculer les manques depuis la dernière candle persistée.
- [ ] **P5.11 — Déclencher `ERROR` pour les erreurs permanentes/intégrité.** Arrêter les nouvelles collectes et ne pas récupérer automatiquement.
- [ ] **P5.12 — Vérifier indisponibilité au démarrage.** Si les retries PostgreSQL échouent, démarrer en `DEGRADED` sans charger ni planifier de targets.
- [ ] **P5.13 — Vérifier reprise après panne PostgreSQL.** Couvrir retour réussi, erreur pendant rechargement et transition vers `ERROR`.
- [ ] **P5.14 — Vérifier l’isolation des erreurs de marché.** Un target en échec ne bloque pas les autres workers.

**Sortie de phase :** les pannes temporaires PostgreSQL sont récupérables sans progression non persistée ; les erreurs permanentes arrêtent Flux de façon explicite.

## P6 — Exploitation et livraison

- [ ] **P6.1 — Ajouter les logs structurés minimum.** Inclure horodatage UTC, niveau, composant et événement.
- [ ] **P6.2 — Ajouter les champs target aux événements concernés.** Inclure exchange, symbol et timeframe.
- [ ] **P6.3 — Journaliser les transitions et événements de retry.** Distinguer retry de marché et retry de persistance.
- [ ] **P6.4 — Ajouter un contrôle de santé applicatif.** Exposer l’état global et l’état de disponibilité PostgreSQL selon l’interface retenue.
- [ ] **P6.5 — Ajouter un scénario de bout en bout local.** Démarrer PostgreSQL, appliquer les migrations, traiter le target initial et vérifier les lignes persistées.
- [ ] **P6.6 — Documenter le démarrage et l’arrêt opérationnels.** Décrire configuration, migrations, commandes, logs et récupération après panne.
- [ ] **P6.7 — Vérifier le parcours complet en CI.** Exécuter les contrôles unitaires, intégration et bout en bout disponibles.
- [ ] **P6.8 — Préparer le premier déploiement.** Définir secrets, migrations au déploiement, arrêt propre et stratégie de retour arrière.
- [ ] **P6.9 — Valider les critères de sortie v1.** Vérifier collecte continue, reprise, contraintes d’unicité, états et logs avant de clôturer le milestone.

**Sortie de phase :** Flux est observable, documenté et validé de la base jusqu’à la collecte réelle du target initial.

## Utilisation dans GitHub

- Créer les milestones `P0 — Fondations` à `P6 — Exploitation et livraison`.
- Créer une issue par case à cocher, en conservant l’identifiant (`P3.6`) dans le titre.
- Ajouter `area:*` selon le sous-système et `type:feature`, `type:test`, `type:docs` ou `type:chore` selon le travail.
- Garder les dépendances de phase indiquées en tête ; au sein d’une phase, exécuter les tâches dans l’ordre et ouvrir le travail parallèle uniquement lorsque les dépendances le permettent.
