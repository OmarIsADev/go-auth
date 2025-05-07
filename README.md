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

   OR
   
   ```bash
   make
   ```

4. **API Endpoints:**

   - <a href="./handlers/auth_handler.go">Authrization:</a>
     - `POST /register`: Register a new user and receive a JWT token
     - `POST /login`: Login and receive a JWT token
     - `POST /reset-password`: Reset user's password
   - <a href="./handlers/auth_handler.go">Token Management:</a>
     - `POST /refresh-token`: Generate new JWT token with refresh token
     - `POST /logout`: Logout user by deleting refresh token
   - <a href="./handlers/userdata.go">User's CRUD:</a>
     - `GET /user`: Retrieve user information (protected route, requires JWT)

## Description of the App

The `auth` is a backend application that manages user authentication and session handling using JWT. It provides a simple interface for registering new users, logging them in, and retrieving user data securely through protected API routes. The application uses SQLite as the database to store user credentials securely with hashed passwords.
