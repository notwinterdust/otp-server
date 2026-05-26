# OTP Authenticator — Self-Host Server

This is the sync server for Offline Authenticator.  
It lets you keep your accounts in sync across multiple devices, all data stays on your own server.

> **Account registration is server-side only.** There is no sign-up screen in the app.  
> You create users via the steps below, then log in from the app.

---

## Requirements

- Docker + Docker Compose
- A server reachable from your devices (home server, VPS, etc.)
- **Strongly recommended but optionnal:** a domain with HTTPS (via a reverse proxy like Caddy or nginx)

---

## Quick start

### 1. Clone this repository

```bash
git clone https://github.com/notwinterdust/otp-server.git
cd otp-server
```

### 2. Edit `docker-compose.yml`

Open `docker-compose.yml` and set the following environment variables:

| Variable           | Description                                      |
|--------------------|--------------------------------------------------|
| `JWT_SECRET`       | A strong random string — run `openssl rand -hex 32` |
| `INITIAL_EMAIL`    | Email for the first user account                 |
| `INITIAL_PASSWORD` | Password for the first user account              |

> After the first start the `INITIAL_EMAIL` / `INITIAL_PASSWORD` values are no longer used  
> (the user is only created once). You can remove them from the file for cleanliness.

### 3. Start the server

```bash
docker compose up -d
```

The server listens on port **8080** by default.

### 4. Verify it's running

```bash
curl http://localhost:8080/api/v1/health
# {"status":"ok","version":"2.0.0"}
```

### 5. Connect the app

In the OTP Authenticator app → **Settings → Self-host**:

- **Address** — your server URL, e.g. `https://otp.example.com` or `http://192.168.1.10:8080`
- **Email** — the email you set in step 2
- **Password** — the password you set in step 2

Tap **Test connection** to verify, then **Save**.

---

## Adding more users

Registration can only be done from the server side:

```bash
docker exec otp-server /app/otp-server add-user \
  --email=another@example.com \
  --password=strongpassword
```

---

## HTTPS (strongly recommended)

Your OTP secrets are sensitive. Use a reverse proxy with a valid TLS certificate.

### Example — Caddy

```
otp.example.com {
    reverse_proxy localhost:8080
}
```

Run `caddy run` and Caddy will obtain a Let's Encrypt certificate automatically.

### Example — nginx + Certbot

```nginx
server {
    listen 443 ssl;
    server_name otp.example.com;

    ssl_certificate     /etc/letsencrypt/live/otp.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/otp.example.com/privkey.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

## Data & persistence

Account data is stored in a SQLite database at `/data/otp.db` inside the container,  
mounted as a Docker volume (`otp_data`). Backup this volume to keep your data safe.

```bash
# Back up the database
docker run --rm -v otp_data:/data -v $(pwd):/backup debian:bookworm-slim \
  cp /data/otp.db /backup/otp-backup.db
```

---

## API reference

| Method | Path                     | Auth | Description                     |
|--------|--------------------------|------|---------------------------------|
| POST   | `/api/v1/auth/login`     | —    | Get JWT token                   |
| GET    | `/api/v1/health`         | —    | Health check                    |
| GET    | `/api/v1/accounts`       | JWT  | Pull all accounts from server   |
| POST   | `/api/v1/accounts/sync`  | JWT  | Push all accounts to server     |

---

## Environment variables

| Variable           | Default        | Description                              |
|--------------------|----------------|------------------------------------------|
| `JWT_SECRET`       | **required**   | HMAC secret for signing tokens           |
| `INITIAL_EMAIL`    | *(optional)*   | Seed first user on first start           |
| `INITIAL_PASSWORD` | *(optional)*   | Seed first user on first start           |
| `DB_PATH`          | `/data/otp.db` | Path to the SQLite database file         |
| `PORT`             | `8080`         | HTTP port the server listens on          |
