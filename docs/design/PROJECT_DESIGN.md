# Backend System Design Documentation

## Overview

This document outlines the design of the backend system for a **Problem Tracking & Challenge System**. The system is designed as a **monolithic architecture** using **Go**. While the system will start as a monolith, it is structured to allow a transition to microservices in the future.

### System Overview:
- **Monolithic Backend**: Initially, all services are bundled into a single application but are modular to allow separation into microservices later.
- **Go-Based Backend**: The system is built using **Go** with REST APIs (initially), with the possibility of adopting gRPC later.
- **Scalable Design**: The system will be structured to allow future horizontal scaling.
- **Key Modules**: The monolithic system consists of multiple modules, starting with the **Auth Service**.

---

## 1. Functional Requirements

- Users should be able to **view problems** from multiple sources (Codeforces, Atcoder, etc.) and perform multiple types of queries.
- Problem information should be updated periodically via a **batch job**.
- Users register their **public Codeforces/Atcoder accounts**, and the system keeps track of their submissions.
- Users can **set a challenge** (solve a problem within a limited time). The result is logged in their history, and their **virtual rating** is updated.
- Users should be able to **visualize their rating history**.

---

## 2. Architecture Overview

### High-Level Architecture

The system will have the following core modules:

1. **Auth Module**: Handles user authentication, including email/password login, OAuth login, and token management.
2. **Problem Module**: Manages problem data from multiple sources (e.g., Codeforces, Atcoder).
3. **Challenge Module**: Manages challenges that users can engage in, tracks challenge status, and updates virtual ratings.
4. **Rating Module**: Manages user virtual ratings based on their performance in challenges.
5. **Notification Module**: Sends notifications to users for various updates (e.g., challenge results).
6. **History & Stats Module**: Tracks users' challenge history and provides visualization for user performance over time.

---

### System Design Approach

- **Monolithic Design**: The backend is implemented as a **monolith**, where all services will be in a single codebase. Each module is structured for a future migration to microservices.
- **REST API First**: The initial implementation will use REST APIs. Future migration to **gRPC** can be considered when transitioning to microservices.
- **Service Modularity**: Each feature is separated into its own module (e.g., `auth`, `problem`, `challenge`, etc.) to allow easy refactoring into microservices later.

---

## 3. First Module: **Auth Module**

### Overview

The **Auth Module** handles **user authentication** and **authorization** for the system. It includes features such as **user registration**, **login via email/password**, and **OAuth login** (such as Google login). It also handles **JWT token generation** for secure authentication.

### Features of Auth Module

1. **User Registration**: Allow users to sign up using an email and password.
2. **User Login**: Authenticate users using their email/password and issue JWT tokens.
3. **OAuth Login**: Allow users to authenticate via third-party OAuth providers (e.g., Google).
4. **Token Validation**: Provide endpoints to validate JWT tokens.
5. **User Profile Management**: Allow retrieval of user information after authentication.

---

## 6. Future Expansion

As the system evolves, the **Auth Module** can be expanded to include:
- **Password Reset**: Allow users to reset their passwords.
- **Multi-factor Authentication (MFA)**: Add support for enhanced security (e.g., SMS or email-based verification).
- **Social Media Login**: Integrate with other OAuth providers such as Facebook, Twitter, etc.

---

## 7. Future Plans for Microservices

While the system starts as a monolith, we will structure it in a way that makes future microservice adoption easier:
- **Encapsulated Modules**: Each module will have a well-defined API boundary.
- **Future gRPC Support**: REST APIs can later be adapted into gRPC-based microservices.
- **Database Per Service Approach**: As the system scales, we may move towards separate databases for each microservice.

---

## Conclusion

This design begins with a **monolithic backend** using **Go** and **REST APIs**, where the **Auth Module** is the first module. The structure is modular and easily extensible to support additional services as the system grows. 

The **Auth Module** provides basic authentication features, including registration, login (password-based and OAuth), token validation, and user info retrieval, forming the core foundation for secure access to the system. The system is structured to transition smoothly to microservices when required.
