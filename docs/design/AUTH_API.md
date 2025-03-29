## 4. API Design for Auth Module (REST API)

### **1. Register (User Registration)**
- **Endpoint**: `POST /auth/register`
- **Request**:
    ```json
    {
      "username": "user123",
      "email": "user@example.com",
      "password": "securepassword"
    }
    ```
- **Response**:
    ```json
    {
      "user_id": "uuid",
      "token": "jwt_token"
    }
    ```

### **2. Login (User Login with Email & Password)**
- **Endpoint**: `POST /auth/login`
- **Request**:
    ```json
    {
      "email": "user@example.com",
      "password": "securepassword"
    }
    ```
- **Response**:
    ```json
    {
      "user_id": "uuid",
      "token": "jwt_token"
    }
    ```

### **3. OAuth Login (Google OAuth)**
- **Endpoint**: `POST /auth/oauth`
- **Request**:
    ```json
    {
      "provider": "google",
      "oauth_token": "google_oauth_token"
    }
    ```
- **Response**:
    ```json
    {
      "user_id": "uuid",
      "token": "jwt_token"
    }
    ```

### **4. Validate Token**
- **Endpoint**: `POST /auth/validate`
- **Request**:
    ```json
    {
      "token": "jwt_token"
    }
    ```
- **Response**:
    ```json
    {
      "valid": true,
      "user_id": "uuid"
    }
    ```

### **5. Get User Info**
- **Endpoint**: `GET /auth/user/{user_id}`
- **Response**:
    ```json
    {
      "user_id": "uuid",
      "username": "user123",
      "email": "user@example.com",
      "oauth_provider": "google"
    }
    ```