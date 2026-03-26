# API Documentation

## Base URL
Default local development: `http://localhost:8081`

## Authentication
Authentication is handled via JWT tokens stored in an [Auth](file:///mnt/d/code/go_lang/go_36_gin_jwt/middlewares/auth.go#14-71) HTTP-only, SameSite=Lax cookie.
The server automatically sets this cookie on successful Login or Signup.
Protected routes will return a `401 Unauthorized` if the cookie is missing, missing claims, or expired.

---

## Authentication Endpoints

### 1. User Signup
Create a new user account. Upon success, the server automatically sets the [Auth](file:///mnt/d/code/go_lang/go_36_gin_jwt/middlewares/auth.go#14-71) cookie.

**Endpoint:** `POST /signup`

**Request Body (JSON):**
```json
{
  "email": "user@example.com",
  "password": "yourpassword123"
}
```

**Responses:**
- `200 OK`:
  ```json
  {
    "message": "Signup successful"
  }
  ```
- `400 Bad Request`: (e.g., if the user already exists or the payload is invalid)
  ```json
  {
    "error": "Email missing or already taken"
  }
  ```

---

### 2. User Login
Authenticate an existing user. Upon success, the server automatically sets the [Auth](file:///mnt/d/code/go_lang/go_36_gin_jwt/middlewares/auth.go#14-71) cookie.

**Endpoint:** `POST /login`

**Request Body (JSON):**
```json
{
  "email": "user@example.com",
  "password": "yourpassword123"
}
```

**Responses:**
- `200 OK`:
  ```json
  {
    "message": "Login successful"
  }
  ```
- `401 Unauthorized`:
  ```json
  {
    "error": "Unauthorized!"
  }
  ```

---

### 3. User Logout
Invalidate the current session by clearing the authentication cookie.

**Endpoint:** `DELETE /logout`

**Responses:**
- `200 OK`:
  ```json
  {
    "message": "Logged out successfully"
  }
  ```

---

### 4. Get Current User (Me)
Retrieve details of the currently authenticated user.

**Endpoint:** `GET /me`

> [!IMPORTANT]
> This endpoint requires the [Auth](file:///mnt/d/code/go_lang/go_36_gin_jwt/middlewares/auth.go#14-71) cookie to be present in the request.

**Responses:**
- `200 OK`:
  ```json
  {
    "user": {
      "ID": 1,
      "CreatedAt": "2026-03-15T16:00:00Z",
      "UpdatedAt": "2026-03-15T16:00:00Z",
      "DeletedAt": null,
      "Email": "user@example.com"
    }
  }
  ```
- `401 Unauthorized`: If the auth cookie is missing or invalid.
