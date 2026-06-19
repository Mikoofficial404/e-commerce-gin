# 🛒 E-Commerce Gin: Microservice-Ready Backend

A highly robust, scalable, and production-ready E-Commerce Backend REST API built with **Golang (Gin Framework)**. This project implements **Clean Architecture** principles and integrates modern enterprise-grade technologies to ensure high performance, security, and reliability.

---

## 🚀 Tech Stack & Features

### Core Technologies
*   **Golang (Gin)**: Extremely fast HTTP web framework.
*   **PostgreSQL & GORM**: Robust relational database with an advanced ORM.
*   **Clean Architecture**: Separation of concerns (`Delivery/Handler` ➡️ `Service` ➡️ `Repository` ➡️ `Database`).

### Enterprise Features
1.  **JWT Authentication & RBAC (Role-Based Access Control)**
    *   Secure login/register system.
    *   `AuthMiddleware` to protect user routes.
    *   `AdminMiddleware` to restrict sensitive actions (e.g., adding products) to admins only.
2.  **RabbitMQ (Asynchronous Task Queue)**
    *   A separate background `worker` process consumes queues.
    *   Handles heavy/slow tasks like sending **Welcome Emails** and **Invoice Emails** so the main API responds in milliseconds.
3.  **Redis (Caching & Rate Limiting)**
    *   **Caching**: Caches the Product Catalog (`FindAllProduct`) to reduce database load. Invalidation triggers automatically when a new product is added.
    *   **Rate Limiting**: Protects `/login` and `/register` endpoints from Brute Force/DDoS attacks.
4.  **AWS S3 / MinIO (Cloud Object Storage)**
    *   Product images are uploaded directly to MinIO (Local S3 clone) via the AWS SDK.
    *   Ensures storage persists even if the backend server restarts.
5.  **Xendit Payment Gateway**
    *   Integrated with Xendit API to generate virtual accounts and handle real-time Webhook payments.
6.  **Structured Logging (Logrus)**
    *   Replaced standard logs with `logrus` to output structured JSON logs, making it ready for analytics tools (Datadog, Elasticsearch, etc.).
7.  **Swagger UI API Documentation**
    *   Beautiful, interactive API documentation accessible via `/swagger/index.html`.
8.  **Graceful Shutdown**
    *   Ensures the server completes ongoing requests before shutting down during a restart.
9.  **Docker Compose**
    *   One-command infrastructure setup containing PostgreSQL, Redis, RabbitMQ, and MinIO.

---



If you want to visualize the flow, here is the architectural diagram layout for your **RabbitMQ** setup:

### 1. The RabbitMQ Async Workflow 
![RabbitMQ Async Workflow](https://github.com/user-attachments/assets/4eb9ac61-0b7b-418e-b48f-c027d1938f7f)

### 2. The Redis Caching Workflow
*   **[ User ]** ➡️ Request `GET /products` ➡️ **[ API ]**
*   **[ API ]** ➡️ Checks **[ Redis Cache ]**
    *   *If Hit:* Return data immediately (Ultra Fast).
    *   *If Miss:* Query **[ PostgreSQL ]** ➡️ Save to **[ Redis Cache ]** for 5 mins ➡️ Return data to User.

---

## 🛠️ How to Run

1. **Start the Infrastructure**
   ```bash
   docker-compose up -d
   ```
2. **Install Dependencies**
   ```bash
   go mod tidy
   ```
3. **Run the Main API Server**
   ```bash
   go run cmd/api/main.go
   ```
4. **Run the RabbitMQ Background Worker** (Open a new terminal)
   ```bash
   go run cmd/worker/main.go
   ```

Enjoy building the future of E-Commerce! 🚀
