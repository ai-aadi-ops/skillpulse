# SkillPulse — Skill Tracker & Learning Dashboard

SkillPulse is a modern, lightweight skill-tracking application designed to help developers log study sessions, track learning goals, and visualize progress dynamically. 

It is designed to be highly portable and supports two deployment modes:
1. **Local Development (Containerized):** Run a full-stack Docker Compose setup locally with Nginx, Go (Gin), and MySQL.
2. **Cloud-Native (Serverless):** Deploy globally for free with Vercel Serverless Functions (Go runtime) and TiDB Cloud Serverless (MySQL).

---

## Tech Stack & Architecture

### Local Architecture (Docker Compose)
* **Frontend:** Vanilla HTML5, CSS3, and JavaScript served by Nginx.
* **Backend:** REST API using Go (Gin Web Framework) served on port `8080`.
* **Database:** MySQL 8.0 container initialized with seed schemas.
* **Reverse Proxy:** Nginx routing requests from port `80` (routing static files to the browser and API requests to the Go backend).

### Serverless Architecture (Vercel + TiDB Cloud)
* **Frontend:** Static assets served globally via Vercel's Edge Network.
* **Backend:** Go REST API compiled on-demand inside Vercel Serverless Functions.
* **Database:** Cloud MySQL-compatible database hosted on TiDB Cloud Serverless with TLS/SSL encryption.
* **Routing:** Declared declaratively inside `vercel.json`.

---

## 🚀 Getting Started (Run Locally)

### Prerequisites
Make sure you have [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/) installed on your machine.

### Quick Start
1. Navigate to the project directory:
   ```bash
   cd project/skillpulse
   ```
2. Copy the example environment variables:
   ```bash
   cp .env.example .env
   ```
3. Boot the environment:
   ```bash
   docker-compose up --build
   ```
4. Access the web interface at **[http://localhost](http://localhost)**.

To stop the containers:
```bash
docker-compose down      # Stops and removes containers
docker-compose down -v   # Stops, removes containers, and clears DB volume data
```

---

## ☁️ Deploying to Vercel & TiDB Cloud

For zero-cost, serverless production hosting, you can deploy the stack over the internet in a few steps.

### Step 1: Spin Up TiDB Cloud Serverless
1. Log in to the [TiDB Cloud Console](https://tidbcloud.com) and click **Create Resource**.
2. Select the **Starter** plan ($0/month) and choose your nearest region.
3. In the **Connect** modal, copy the connection details: `Host`, `Port`, `Username` (starts with your account prefix), and `Password`.
4. Update the TiDB Firewall IP Access List to allow `0.0.0.0` to `255.255.255.255`. *(Required since Vercel utilizes dynamic IP addresses)*.
5. In the **SQL Editor** tab, select the default `test` database, paste the initialization script (`mysql/init.sql`), and run it to set up the schema.

### Step 2: Configure Vercel Deployment
1. Log in to Vercel and import your repository.
2. Under **Project Settings > Environment Variables**, add the variables specified below in the environment reference (matching your TiDB credentials).
3. Set `DB_SSL` to `skip-verify` to ensure secure handshake compatibility with TiDB Cloud.
4. Click **Deploy**. Vercel will build the static frontend and compile the Go handlers automatically.
5. Open your deployment URL to view your live cloud app!

---

## ⚙️ Environment Variables Reference

| Variable Name | Description | Local (Default) | Vercel Serverless (TiDB Cloud) |
| :--- | :--- | :--- | :--- |
| `DB_HOST` | Hostname of the MySQL server | `db` | `gateway01.ap-southeast-1.prod.aws.tidbcloud.com` |
| `DB_PORT` | Port of the database connection | `3306` | `4000` |
| `DB_USER` | Database connection username | `skillpulse` | `<your-tidb-generated-username>` |
| `DB_PASSWORD` | Database connection password | `skillpulse123` | `<your-tidb-generated-password>` |
| `DB_NAME` | Name of the database schema | `skillpulse` | `test` |
| `DB_SSL` | Enable database SSL certificate verification | `false` | `skip-verify` |

*Note: For local development, copy `.env.example` to `.env` and fill the variables. For Vercel, inject them via the Vercel Dashboard.*

---

## 🔌 API Endpoints Reference

All API routes are prefixed with `/api`.

| HTTP Method | Route | Description | Response Content |
| :--- | :--- | :--- | :--- |
| **GET** | `/api/skills` | List all tracked skills & current logged progress | JSON Array of Skills |
| **POST** | `/api/skills` | Register a new skill to track | JSON Object of Created Skill |
| **GET** | `/api/skills/:id` | Fetch details of a specific skill and its log entries | Detailed JSON Skill Object |
| **DELETE** | `/api/skills/:id` | Remove a skill and its corresponding study logs | Status Confirmation JSON |
| **POST** | `/api/skills/:id/log` | Register a study session (hours + notes) | JSON Object of Created Log |
| **GET** | `/api/dashboard` | Retrieve overall progress metrics for dashboard | Metrics JSON (Total Hours, count, etc.) |
| **GET** | `/health` | Application health status check | `{"status": "healthy"}` |

---

## 📂 Repository Structure

* `api/index.go`: Serverless handler entry point designed for Vercel Go Runtime.
* `database/db.go`: MySQL connection pool library with automated retry logic and SSL support.
* `handlers/`: Go API handlers for skills, logs, and dashboard metrics.
* `models/`: Data structures defining the SkillPulse domain.
* `frontend/`: Single-page app HTML/CSS/JavaScript static assets.
* `nginx/`: Nginx proxy configuration for routing traffic locally.
* `mysql/`: DB initialization schemas and seed datasets.
* `vercel.json`: Declarative router and builder settings for Vercel.