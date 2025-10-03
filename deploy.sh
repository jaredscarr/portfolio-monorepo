#!/bin/bash

# Portfolio Monorepo AWS EC2 Deployment Script
# Run this script on your EC2 instance

set -e

echo "🚀 Starting Portfolio Monorepo Deployment..."

# Update system
echo "📦 Updating system packages..."
sudo apt update && sudo apt upgrade -y

# Install Docker
echo "🐳 Installing Docker..."
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    rm get-docker.sh
fi

# Install Docker Compose
echo "🔧 Installing Docker Compose..."
if ! command -v docker-compose &> /dev/null; then
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
fi

# Create application directory
echo "📁 Setting up application directory..."
sudo mkdir -p /opt/portfolio
sudo chown $USER:$USER /opt/portfolio
cd /opt/portfolio

# Clone repository (replace with your actual repo URL)
echo "📥 Cloning repository..."
if [ ! -d "portfolio-monorepo" ]; then
    git clone https://github.com/jared-scarr/portfolio-monorepo.git
fi
cd portfolio-monorepo

# Set up environment
echo "⚙️ Setting up environment..."
if [ ! -f ".env.production" ]; then
    echo "❌ Please create .env.production file from env.production.template"
    echo "   Update domain and email values before continuing"
    exit 1
fi

# Update Caddyfile with your domain
echo "🌐 Updating Caddyfile..."
read -p "Enter your domain name: " DOMAIN
read -p "Enter your email for SSL: " EMAIL

sed -i "s/your-domain.com/$DOMAIN/g" Caddyfile
sed -i "s/your-email@example.com/$EMAIL/g" Caddyfile

# Start services
echo "🚀 Starting services..."
docker-compose up --build -d

# Wait for services to be healthy
echo "⏳ Waiting for services to be healthy..."
sleep 30

# Check service status
echo "✅ Checking service status..."
docker-compose ps

echo "🎉 Deployment complete!"
echo "🌐 Your portfolio should be available at: https://$DOMAIN"
echo "📊 API endpoints available at: https://api.$DOMAIN"
echo ""
echo "📋 Useful commands:"
echo "  View logs: docker-compose logs -f"
echo "  Stop services: docker-compose down"
echo "  Restart services: docker-compose restart"
echo "  Update: git pull && docker-compose up --build -d"
