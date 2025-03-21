# Running Aynshteyn.dev (Private Instructions)

Simple instructions for setting up and running the aynshteyn.dev website and API, assuming the repo is already on the server.

## Quick Setup

1. **Build the backend:**
   ```bash
   cd backend
   go build -o aynshteyn-backend ./cmd/api
   ```

2. **Build the frontend:**
   ```bash
   cd my-app
   npm install
   npm run build
   ```

3. **Run the backend:**
   ```bash
   cd backend
   ./aynshteyn-backend
   ```
   
   Or run as a service:
   ```bash
   sudo cp aynshteyn-backend.service /etc/systemd/system/
   sudo systemctl daemon-reload
   sudo systemctl enable aynshteyn-backend
   sudo systemctl start aynshteyn-backend
   ```

4. **Setup Nginx:**
   ```bash
   sudo cp deployment/nginx.conf /etc/nginx/sites-available/aynshteyn.dev
   sudo ln -sf /etc/nginx/sites-available/aynshteyn.dev /etc/nginx/sites-enabled/
   sudo nginx -t
   sudo systemctl reload nginx
   ```

5. **Copy frontend files:**
   ```bash
   sudo mkdir -p /var/www/aynshteyn.dev/www
   sudo cp -r my-app/.next my-app/public /var/www/aynshteyn.dev/www/
   ```

## Checking Subscribers

View all subscribers:
```bash
curl -u admin:aynshteyn-secure-password http://localhost:4000/api/v1/admin/subscribers
```

Or directly from the database:
```bash
sqlite3 backend/subscribers.db "SELECT * FROM subscribers;"
```

## Troubleshooting

- **Check backend status:**
  ```bash
  systemctl status aynshteyn-backend
  ```

- **View backend logs:**
  ```bash
  journalctl -u aynshteyn-backend
  ```

- **Check Nginx status:**
  ```bash
  systemctl status nginx
  ```

- **View Nginx logs:**
  ```bash
  tail -f /var/log/nginx/error.log
  ```

## Security Notes

- Default admin credentials are: `admin:aynshteyn-secure-password` (change in production)
- The backend listens on port 4000 by default
- The SQLite database is stored at: `backend/subscribers.db` 