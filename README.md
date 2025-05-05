# auth

This is a simple auth built with Go and Fiber framework. It provides endpoints for user registration, login, and retrieving user information with JWT-based authentication.

## Features

- User Registration
- User Login
- JWT Authentication
- Access protected routes with valid JWT

## Prerequisites

- Go 1.16 or later
- CGO
- SQLite

## Instructions for Running

1. **Clone the repository:**

   ```bash
   git clone https://github.com/omarisadev/go-auth.git
   cd auth
   ```

2. **Set up environment variables:**

   Ensure you have a `.env` file with the following content:

   ```env
   JWT_SECRET=your_secret_key
   ```

   *Incase you didn't the app will run normally, because there is a test key*

3. **Run the application:**

   ```bash
   go run cmd/main.go
   ```

4. **API Endpoints:**

   - <a href="./handlers/auth_handler.go">`POST /register`</a>: Register a new user
   - <a href="./handlers/auth_handler.go">`POST /login`</a>: Login and receive a JWT token
   -  <a href="./handlers/userdata.go">`GET /user`</a>: Retrieve user information (protected route, requires JWT)
     - Check the <a href="./middleware/auth_middleware.go">middleware</a> function.

## Description of the App

The `auth` is a backend application that manages user authentication and session handling using JWT. It provides a simple interface for registering new users, logging them in, and retrieving user data securely through protected API routes. The application uses SQLite as the database to store user credentials securely with hashed passwords.
