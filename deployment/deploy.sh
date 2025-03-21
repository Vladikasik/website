#!/bin/bash
# Deployment script for aynshteyn.dev (both frontend and backend)

set -e  # Exit on any error

# Define variables
SERVER_IP="95.163.223.236"
DOMAIN="aynshteyn.dev"
BACKEND_DIR="/var/www/$DOMAIN/backend"
FRONTEND_DIR="/var/www/$DOMAIN/www"
SERVICE_NAME="aynshteyn-backend"

# Archive the Next.js output for transfer
echo "Preparing frontend files..."
cd ../my-app && tar -czf ../deployment/frontend.tar.gz .next public

# Back to deployment directory
cd ../deployment

echo "Connecting to server and setting up directories..."
ssh root@$SERVER_IP "mkdir -p $BACKEND_DIR $FRONTEND_DIR"

echo "Deploying backend..."
# Copy backend binary, service file, and nginx config
scp aynshteyn-backend root@$SERVER_IP:$BACKEND_DIR/
scp aynshteyn-backend.service root@$SERVER_IP:/etc/systemd/system/$SERVICE_NAME.service
scp nginx.conf root@$SERVER_IP:/etc/nginx/sites-available/$DOMAIN

# Deploy the frontend
echo "Deploying frontend..."
scp frontend.tar.gz root@$SERVER_IP:$FRONTEND_DIR/

# Set up everything on the server
echo "Setting up services on server..."
ssh root@$SERVER_IP "
  # Set up backend service
  chmod +x $BACKEND_DIR/aynshteyn-backend
  chown -R www-data:www-data $BACKEND_DIR

  # Extract frontend files
  cd $FRONTEND_DIR
  tar -xzf frontend.tar.gz
  rm frontend.tar.gz
  chown -R www-data:www-data $FRONTEND_DIR

  # Set up Nginx
  if [ ! -f /etc/nginx/sites-enabled/$DOMAIN ]; then
    ln -s /etc/nginx/sites-available/$DOMAIN /etc/nginx/sites-enabled/
  fi

  # Test Nginx configuration
  nginx -t

  # Reload services
  systemctl daemon-reload
  systemctl enable $SERVICE_NAME
  systemctl restart $SERVICE_NAME
  systemctl reload nginx

  # Show service status
  systemctl status $SERVICE_NAME --no-pager
"

echo "Deployment completed successfully!"
echo "Frontend: https://$DOMAIN"
echo "API: https://api.$DOMAIN"
echo ""
echo "To check subscribers, use: curl -u admin:aynshteyn-secure-password https://api.$DOMAIN/api/v1/admin/subscribers" 