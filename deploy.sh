#!/bin/bash

# Portfolio Monorepo AWS EC2 Deployment Script
# Run this script on your EC2 instance
# Supports: Amazon Linux 2023, Amazon Linux 2, Ubuntu

set -e

echo "🚀 Starting Portfolio Monorepo Deployment..."

# Detect OS and package manager
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$ID
else
    echo "❌ Cannot detect OS. This script requires Amazon Linux or Ubuntu."
    exit 1
fi

echo "📋 Detected OS: $OS"

# Update system packages based on OS
echo "📦 Updating system packages..."
if [[ "$OS" == "amzn" ]] || [[ "$OS" == "amazon" ]]; then
    # Amazon Linux 2023 uses dnf, Amazon Linux 2 uses yum
    if command -v dnf &> /dev/null; then
        sudo dnf update -y
    else
        sudo yum update -y
    fi
elif [[ "$OS" == "ubuntu" ]] || [[ "$OS" == "debian" ]]; then
    sudo apt update && sudo apt upgrade -y
else
    echo "⚠️  Unsupported OS: $OS. Attempting generic package update..."
    echo "   This script is tested on Amazon Linux and Ubuntu."
fi

# Install Docker
echo "🐳 Installing Docker..."
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    rm get-docker.sh
fi

# Install git if not present (required for cloning repo)
if ! command -v git &> /dev/null; then
    echo "📥 Installing git..."
    if [[ "$OS" == "amzn" ]] || [[ "$OS" == "amazon" ]]; then
        if command -v dnf &> /dev/null; then
            sudo dnf install -y git
        else
            sudo yum install -y git
        fi
    elif [[ "$OS" == "ubuntu" ]] || [[ "$OS" == "debian" ]]; then
        sudo apt install -y git
    fi
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
docker-compose --env-file .env.production -f docker-compose.yml -f docker-compose.prod.yml up --build -d

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
