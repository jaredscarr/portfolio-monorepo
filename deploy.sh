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
    echo "⚠️  .env.production not found"
    echo "   Creating from template (you should review and update values)..."
    if [ -f "env.production.template" ]; then
        cp env.production.template .env.production
        echo "   Created .env.production from template"
        echo "   ⚠️  IMPORTANT: Review .env.production and update all placeholder values"
        echo "   ⚠️  Change DB_PASSWORD and other secrets before deploying!"
        read -p "   Press Enter to continue after reviewing .env.production..."
    else
        echo "❌ env.production.template not found"
        echo "   Please create .env.production manually with required environment variables"
        exit 1
    fi
fi

# Update Caddyfile with your domain
echo "🌐 Updating Caddyfile..."
read -p "Enter your domain name (e.g., example.com): " DOMAIN
read -p "Enter your email for SSL certificates: " EMAIL

# Replace placeholders in Caddyfile
sed -i "s/{{DOMAIN}}/$DOMAIN/g" Caddyfile
sed -i "s/{{EMAIL}}/$EMAIL/g" Caddyfile

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
