# Auth Service Design Document

## Overview
This document outlines the design of the **Auth Service** for a Problem Tracking & Challenge System. The service handles **user authentication** via OAuth and manages JWT-based session handling.

## Database Schema
### Tables

#### `users`
Stores user information.
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),  -- Unique internal user ID
    email TEXT UNIQUE NOT NULL,                     -- User email from OAuth provider
    codeforces_account TEXT,                        -- User's Codeforces handle (non-unique)
    atcoder_account TEXT,                           -- User's Atcoder handle (non-unique)
    oauth_provider TEXT NOT NULL,                   -- e.g., 'google', 'github'
    oauth_provider_id TEXT UNIQUE NOT NULL,         -- Unique ID from OAuth provider
    created_at TIMESTAMP DEFAULT NOW(),             -- Timestamp of account creation
    updated_at TIMESTAMP DEFAULT NOW()              -- Timestamp of last update
);
```

#### `tokens`
Stores refresh tokens.
```sql
CREATE TABLE tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),   -- Unique token ID
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,  -- Links to user
    refresh_token TEXT NOT NULL,                    -- Stores the refresh token
    expires_at TIMESTAMP NOT NULL,                  -- Expiration time of refresh token
    created_at TIMESTAMP DEFAULT NOW()              -- Timestamp of token creation
);
```

## API Endpoints

### **1️⃣ OAuth Login Initiation**
**Endpoint:** `GET /auth/oauth/login`
- **Description:** Redirects the user to the OAuth provider's login page.
- **Query Parameters:**
  - `provider` (string) → e.g., `google`
- **Response:** Redirect to OAuth provider

### **2️⃣ OAuth Callback**
**Endpoint:** `POST /auth/oauth/callback`
- **Description:** Handles the callback from OAuth provider, exchanges authorization code for tokens.
- **Request Body:**
```json
{
    "code": "AUTHORIZATION_CODE"
}
```
- **Response:**
```json
{
    "access_token": "JWT_ACCESS_TOKEN",
    "refresh_token": "JWT_REFRESH_TOKEN"
}
```

### **3️⃣ Refresh Token**
**Endpoint:** `POST /auth/refresh`
- **Description:** Generates a new access token using the refresh token.
- **Request Body:**
```json
{
    "refresh_token": "JWT_REFRESH_TOKEN"
}
```
- **Response:**
```json
{
    "access_token": "NEW_JWT_ACCESS_TOKEN"
}
```

### **4️⃣ Logout**
**Endpoint:** `POST /auth/logout`
- **Description:** Invalidates the refresh token.
- **Request Body:**
```json
{
    "refresh_token": "JWT_REFRESH_TOKEN"
}
```
- **Response:**
```json
{
    "message": "Logged out successfully"
}
```

## Security Considerations
- **CORS:** Backend should allow requests from the frontend domain only.
- **HTTPS:** All requests should be encrypted.
- **State Parameter:** Used in OAuth login initiation to prevent CSRF.
- **HttpOnly Cookies:** Refresh tokens should be stored in HttpOnly cookies for security.

## Future Extensions
- Add support for more OAuth providers (e.g., GitHub, Facebook)
- Implement role-based access control (RBAC)
- Support multi-session token revocation
