# Aynshteyn.dev Deployment

This directory contains everything needed to deploy the Aynshteyn.dev website and API.

## What's Included

- `aynshteyn-backend`: The compiled Go binary for the API server
- `aynshteyn-backend.service`: Systemd service file for running the API server
- `nginx.conf`: Nginx configuration for both the website and API
- `deploy.sh`: Deployment script to automate the process

## Manual Deployment Instructions

If you don't want to use the automated script, follow these steps:

### Backend Deployment

1. SSH to the server:
   ```
   ssh root@95.163.223.236
   ```

2. Create necessary directories:
   ```
   mkdir -p /var/www/aynshteyn.dev/backend
   mkdir -p /var/www/aynshteyn.dev/www
   ```

3. Copy the backend binary and make it executable:
   ```
   scp aynshteyn-backend root@95.163.223.236:/var/www/aynshteyn.dev/backend/
   ssh root@95.163.223.236 "chmod +x /var/www/aynshteyn.dev/backend/aynshteyn-backend"
   ```

4. Copy the systemd service file and enable it:
   ```
   scp aynshteyn-backend.service root@95.163.223.236:/etc/systemd/system/aynshteyn-backend.service
   ssh root@95.163.223.236 "systemctl daemon-reload && systemctl enable aynshteyn-backend && systemctl restart aynshteyn-backend"
   ```

### Frontend Deployment

1. Build the Next.js app:
   ```
   cd my-app && npm run build
   ```

2. Copy the built files to the server:
   ```
   tar -czf frontend.tar.gz .next public
   scp frontend.tar.gz root@95.163.223.236:/var/www/aynshteyn.dev/www/
   ssh root@95.163.223.236 "cd /var/www/aynshteyn.dev/www/ && tar -xzf frontend.tar.gz && rm frontend.tar.gz"
   ```

### Nginx Configuration

1. Copy the Nginx configuration:
   ```
   scp nginx.conf root@95.163.223.236:/etc/nginx/sites-available/aynshteyn.dev
   ```

2. Enable the configuration:
   ```
   ssh root@95.163.223.236 "ln -sf /etc/nginx/sites-available/aynshteyn.dev /etc/nginx/sites-enabled/ && nginx -t && systemctl reload nginx"
   ```

## Automated Deployment

Simply run the deployment script:

```
./deploy.sh
```

## Checking Subscribers

To view all subscribers, use the admin API:

```
curl -u admin:aynshteyn-secure-password https://api.aynshteyn.dev/api/v1/admin/subscribers
```

## Security Information

The backend includes several security features:

1. Rate limiting to prevent abuse (60 requests per minute per IP)
2. Input validation for all API endpoints
3. CORS protection allowing only the main website domain
4. HMAC verification for request integrity
5. Secure cookie handling
6. TLS encryption (HTTPS)
7. Admin API protected with HTTP Basic Authentication
8. IP address logging for security monitoring
9. Email hashing for privacy
10. Protection against common attacks (XSS, CSRF, etc.)

## Troubleshooting

- Check backend logs: `systemctl status aynshteyn-backend`
- Check Nginx logs: `tail -f /var/log/nginx/error.log`
- Database location: `/var/www/aynshteyn.dev/backend/subscribers.db` 