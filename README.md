# Chirpy

Chirpy is a lightweight social media platform designed for quick and simple interactions. Users can create, read, and manage short posts (chirps) while enjoying features like user authentication, token-based authorization, and profanity filtering.

## Why Chirpy?

Chirpy is perfect for developers looking to explore:
- Building RESTful APIs with Go.
- Implementing JWT-based authentication.
- Managing relational databases with PostgreSQL.
- Using SQL migrations and code generation tools like `sqlc`.

## Installation

1. **Clone the repository**:
   ```shell
   git clone https://github.com/ahnaftahmid39/chirpy.git
   cd chirpy
   ```

2. **Set up the environment**:

  - Copy .env.example to .env:
    ```shell
    cp .env.example .env
    ```
  - Update the .env file with your database credentials and secrets.

3. **Install dependencies**:
    ```shell
    go mod tidy
    ```

4. **Run database migrations**:
Ensure PostgreSQL is running, then apply migrations using goose:
    ```shell
    goose -dir sql/schema postgres "$DB_URL" up
    ```

5. **Start the server**:
    ```shell
    go run main.go
    ```
## Usage
- Access the app at http://localhost:<PORT> (default: 8069).
- Use the provided RESTful API endpoints for user management, chirps, and admin operations.

## Features
- **User Authentication**: Secure login with hashed passwords and JWTs.
- **Chirps**: Create, read, and delete short posts with profanity filtering.
- **Admin Tools**: Reset data and monitor metrics.
- **Refresh Tokens**: Long-lived sessions with token revocation support.

Happy coding with Chirpy!
