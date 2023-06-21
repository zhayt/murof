# Forum

This is a web forum created using Go, HTML, and SQLite. The forum allows users to create posts and comments, associate categories with posts, and like/dislike posts and comments. The project also includes a filter mechanism that allows users to filter posts by categories, created posts, and liked posts.

## Installation

1. Clone the repository:
```bash
git clone https://01.alem.school/git/Madi856/forum.git
```

2. Build the Docker image:
```bash
docker build -t forum .
```

3. Run the Docker container:
```bash
docker run --name new-forum --rm -p8080:8080 forum
```

4. Open the forum in your web browser:
```
http://localhost:8080
```

## Usage

### Registration

To register as a new user on the forum, enter your email, username, and password. If the email or name is already taken, an error response will be returned. 
There is also a password requirement (read about password requirement in the top of the topic).

### Login

To access the forum, login using your credentials. Only registered users can create posts and comments, like/dislike posts and comments, and filter posts.

### Posts and Comments

Registered users can create posts and comments, and associate categories with posts. The posts and comments are visible to all users, including non-registered users.

### Likes and Dislikes

Only registered users can like or dislike posts and comments. The number of likes and dislikes is visible to all users.

### Filter

Users can filter posts by categories, created posts, and liked posts. Non-registered users can only see posts and comments.

## Database

The database includes the following tables:

### user

The user table stores user data, including ID, password, login, username, token, and token duration.

### post

The post table stores post data, including ID, user ID, title, text, and date.

### category

The category table stores category data, including ID and name.

### post_category

The post_category table stores the relationship between posts and categories.

### comment

The comment table stores comment data, including ID, user ID, post ID, text, and date.

### like

The like table stores like data, including ID, user ID, post/comment ID, and value.

## Packages

The project uses the following packages:

- `database/sql` - database/sql provides a generic interface around SQL (or SQL-like) databases.
- `github.com/mattn/go-sqlite3` - sqlite3 driver for Go using database/sql.
- `golang.org/x/crypto/bcrypt` - package bcrypt implements Provos and Mazières's bcrypt adaptive hashing algorithm.
- `github.com/google/uuid` - the uuid package generates and inspects UUIDs based on RFC 4122 and DCE 1.1: Authentication and Security Services.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://i.imgur.com/pUverGl.jpeg)
