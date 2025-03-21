#!/bin/bash
# Deployment script for aynshteyn.dev backend

set -e  # Exit on any error

# Define variables
SERVER_IP="95.163.223.236"
DOMAIN="aynshteyn.dev"
BACKEND_DIR="/var/www/$DOMAIN/backend"
SERVICE_NAME="aynshteyn-backend"

# Build the binary
echo "Building for production..."
GOOS=linux GOARCH=amd64 go build -o aynshteyn-backend ./cmd/api

echo "Copying files to server..."
# Create directory if it doesn't exist
ssh root@$SERVER_IP "mkdir -p $BACKEND_DIR"

# Copy binary and service file
scp aynshteyn-backend root@$SERVER_IP:$BACKEND_DIR/
scp aynshteyn-backend.service root@$SERVER_IP:/etc/systemd/system/$SERVICE_NAME.service

# Set up systemd service
echo "Setting up systemd service..."
ssh root@$SERVER_IP "
  # Set permissions
  chmod +x $BACKEND_DIR/aynshteyn-backend
  chown -R www-data:www-data $BACKEND_DIR

  # Set up and start service
  systemctl daemon-reload
  systemctl enable $SERVICE_NAME
  systemctl restart $SERVICE_NAME
  systemctl status $SERVICE_NAME
"

# Set up Nginx (if needed)
echo "Setting up Nginx proxy..."
ssh root@$SERVER_IP "
  cat > /etc/nginx/sites-available/$DOMAIN.api << 'EOL'
server {
    listen 80;
    server_name api.$DOMAIN;
    
    location / {
        return 301 https://\$host\$request_uri;
    }
}

server {
    listen 443 ssl http2;
    server_name api.$DOMAIN;
    
    ssl_certificate /etc/letsencrypt/live/$DOMAIN/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/$DOMAIN/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    
    location / {
        proxy_pass http://localhost:4000;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOL

  # Create symlink if it doesn't exist
  if [ ! -f /etc/nginx/sites-enabled/$DOMAIN.api ]; then
    ln -s /etc/nginx/sites-available/$DOMAIN.api /etc/nginx/sites-enabled/
  fi

  # Test Nginx configuration
  nginx -t
  
  # Restart Nginx if test was successful
  systemctl reload nginx
"

echo "Deployment completed successfully."
echo "API is available at https://api.$DOMAIN" 