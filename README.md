# API GoLangEstudos

Esta é uma API desenvolvida em Go (Golang) para estudos, com funcionalidades relacionadas a usuários, publicações, autenticação e interações sociais, como seguir usuários e curtir publicações.

## 📋 Requisitos

- Go 1.23.5 ou superior.
- MySQL (ou outro banco de dados compatível com `go-sql-driver/mysql`).
- Git (para clonar o repositório).

## 🚀 Como executar o projeto

### 1. Clonar o repositório

```bash
git clone https://github.com/VictorBriske/Dev-book
cd Dev-book
```

### 2. Configurar o ambiente

Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis de ambiente:

```env
DB_USER=golang
DB_PASSWORD=golang
DB_NAME=devbook
API_PORT=5000
SECRET_KEY=Ma/YbSgmMe8vw/9Lh2dbgfIiGikW6QnkG5mXewBLhfsatVAAZf+X5T3R30TewQqPsKkyYAAxRlLCDFTJlnealw==
```

### 3. Configurar o banco de dados

Execute o script SQL abaixo para criar o banco de dados e as tabelas necessárias:

```sql
CREATE DATABASE IF NOT EXISTS devbook;
USE devbook;

DROP TABLE IF EXISTS publications;
DROP TABLE IF EXISTS followers;
DROP TABLE IF EXISTS users;

CREATE TABLE users(
    id int auto_increment primary key,
    name varchar(50) not null,
    nickname varchar(50) not null unique,
    email varchar(50) not null unique,
    password varchar(120) not null,
    created_at timestamp default current_timestamp()
) ENGINE=INNODB;

CREATE TABLE followers(
    user_id int not null,
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,
    follower_id int not null,
    FOREIGN KEY (follower_id)
    REFERENCES users(id)
    ON DELETE CASCADE,
    primary key(user_id, follower_id)
) ENGINE=INNODB;

CREATE TABLE publications(
    id int auto_increment primary key,
    title varchar(50) not null,
    content varchar(300) not null,
    author_id int not null,
    FOREIGN KEY (author_id)
    REFERENCES users(id)
    ON DELETE CASCADE,
    likes int default 0,
    created_at timestamp default current_timestamp
) ENGINE=INNODB;
```

### 4. Instalar dependências

Certifique-se de que todas as dependências do projeto estão instaladas. Execute:

```bash
go mod tidy
```

### 5. Executar a API

Para iniciar o servidor da API, execute:

```bash
go run main.go
```

A API estará disponível em `http://localhost:5000`.

---

## 📚 Endpoints da API

Aqui estão os principais endpoints disponíveis na API:

### Users

- **Create user**: `POST /users`
  - Request body:
    ```json
    {
      "name": "User 1",
      "email": "user1@gmail.com",
      "nickname": "user1",
      "password": "123"
    }
    ```

- **Login**: `POST /login`
  - Request body:
    ```json
    {
      "email": "user1@gmail.com",
      "password": "123"
    }
    ```

- **Get users**: `GET /users?user=test`
- **Get user by ID**: `GET /users/{id}`
- **Update user**: `PUT /users/{id}`
- **Update user password**: `PUT /users/{id}/password`
- **Delete user**: `DELETE /users/{id}`

### Followers

- **Follow user**: `POST /users/{id}/follow`
- **Unfollow user**: `DELETE /users/{id}/follow`
- **Get followers of a user**: `GET /users/{id}/followers`
- **Get users a user follows**: `GET /users/{id}/following`

### Publications

- **Create publication**: `POST /publications`
  - Request body:
    ```json
    {
      "title": "My publication",
      "content": "Publication content"
    }
    ```

- **Get publication by ID**: `GET /publications/{id}`
- **Get all publications**: `GET /publications`
- **Update publication**: `PUT /publications/{id}`
- **Delete publication**: `DELETE /publications/{id}`
- **Get publications of a user**: `GET /users/{id}/publications`

### Likes

- **Like publication**: `POST /publications/{id}/like`
- **Unlike publication**: `DELETE /publications/{id}/like`

---

## 🛠️ Bibliotecas utilizadas

- **[Gorilla Mux](https://github.com/gorilla/mux)**: Para roteamento HTTP.
- **[JWT Go](https://github.com/dgrijalva/jwt-go)**: Para autenticação via tokens JWT.
- **[Go MySQL Driver](https://github.com/go-sql-driver/mysql)**: Para conexão com o banco de dados MySQL.
- **[Godotenv](https://github.com/joho/godotenv)**: Para gerenciamento de variáveis de ambiente.
- **[Badoux Checkmail](https://github.com/badoux/checkmail)**: Para validação de e-mails.
- **[Golang Crypto](https://pkg.go.dev/golang.org/x/crypto)**: Para criptografia de senhas.
