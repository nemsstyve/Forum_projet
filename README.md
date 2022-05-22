# Forum_projet
Ce projet consiste à créer un forum web qui permet :

    communication entre les utilisateurs.
    associer des catégories aux messages.
    aimer et ne pas aimer les messages et les commentaires.
    filtrage des messages.


## SQLite

Afin de stocker les données de votre forum (comme les utilisateurs, les messages, les commentaires, etc.), vous utiliserez la bibliothèque de base de données SQLite.

SQLite est un choix populaire en tant que logiciel de base de données intégré pour le stockage local/client dans les logiciels d'application tels que les navigateurs Web. Il vous permet de créer une base de données ainsi que de la contrôler à l'aide de requêtes.

Pour structurer votre base de données et obtenir de meilleures performances, nous vous conseillons vivement de consulter le diagramme de relation d'entité et d'en créer un basé sur votre propre base de données.

    Vous devez utiliser au moins une requête SELECT, une CREATE et une INSERT.

Pour en savoir plus sur SQLite, vous pouvez consulter la page SQLite.


## Authentification

Dans ce segment, le client doit pouvoir s'inscrire en tant que nouvel utilisateur sur le forum, en saisissant ses informations d'identification. Vous devez également créer une session de connexion pour accéder au forum et pouvoir ajouter des messages et des commentaires.

Vous devez utiliser des cookies pour permettre à chaque utilisateur d'avoir une seule session ouverte. Chacune de ces sessions doit contenir une date d'expiration. C'est à vous de décider combien de temps le cookie reste "vivant". L'utilisation de l'UUID est une tâche bonus.

Instructions pour l'enregistrement de l'utilisateur :

    Doit demander un e-mail
        Lorsque l'e-mail est déjà pris, renvoyez une réponse d'erreur.
    Doit demander un nom d'utilisateur
    Doit demander un mot de passe
        Le mot de passe doit être crypté lorsqu'il est stocké (il s'agit d'une tâche bonus)

Le forum doit être en mesure de vérifier si l'e-mail fourni est présent dans la base de données et si tous les identifiants sont corrects. Il vérifiera si le mot de passe est le même que celui fourni et, si le mot de passe n'est pas le même, il renverra une réponse d'erreur.


## Communication

Pour que les utilisateurs puissent communiquer entre eux, ils devront pouvoir créer des posts et des commentaires.

    Seuls les utilisateurs enregistrés pourront créer des publications et des commentaires.
    Lorsque les utilisateurs enregistrés créent une publication, ils peuvent y associer une ou plusieurs catégories.
        La mise en œuvre et le choix des catégories vous appartient.
    Les publications et les commentaires doivent être visibles par tous les utilisateurs (inscrits ou non).
    Les utilisateurs non enregistrés ne pourront voir que les publications et les commentaires.


## Aime et n'aime pas

Seuls les utilisateurs enregistrés pourront aimer ou ne pas aimer les publications et les commentaires.

Le nombre de likes et dislikes doit être visible par tous les utilisateurs (inscrits ou non).


## Filtre

Vous devez implémenter un mécanisme de filtrage, qui permettra aux utilisateurs de filtrer les publications affichées par :

    catégories
    messages créés
    messages aimés

Vous pouvez regarder le filtrage par catégories en tant que sous-forums. Un sous-forum est une section d'un forum en ligne dédiée à un sujet spécifique.

Notez que les deux derniers ne sont disponibles que pour les utilisateurs enregistrés et doivent faire référence à l'utilisateur connecté.


## Docker

Pour le projet de forum, vous devez utiliser Docker. Vous pouvez en savoir plus sur les bases de Docker dans le sujet ascii-art-web-dockerize.


## Des instructions

    Vous devez utiliser SQLite.
    Vous devez gérer les erreurs de site Web, le statut HTTP.
    Vous devez gérer toutes sortes d'erreurs techniques.
    Le code doit respecter les bonnes pratiques.
    Il est recommandé d'avoir des fichiers de test pour les tests unitaires.


## Forfaits autorisés

    Tous les forfaits Go standard sont autorisés.
    sqlite3
    bcrypt
    UUID

Ce projet vous aidera à découvrir :

    Les bases du web :
        HTML
        HTTP
        Sessions et cookies
    Utiliser et configurer Docker
        Conteneuriser une application
        Compatibilité/Dépendance
        Créer des images
    Langage SQL
        Manipulation de bases de données
    Les bases du cryptage