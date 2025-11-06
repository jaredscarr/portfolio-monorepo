# AWS Deployment Guide

This guide covers the complete deployment process for the Portfolio Monorepo on AWS EC2.

## Prerequisites

- AWS EC2 instance (t3.micro recommended)
- Domain name with DNS management access
- SSH access to EC2 instance
- Ports 80 (HTTP) and 443 (HTTPS) open in security group

### Supported Operating Systems

The deployment script automatically detects and supports:
- **Amazon Linux 2023** (recommended) - Uses `dnf` package manager
- **Amazon Linux 2** - Uses `yum` package manager  
- **Ubuntu 22.04 LTS or later** - Uses `apt` package manager

**Recommended AMI:** When launching your EC2 instance, choose:
- Amazon Linux 2023 AMI (most recent, best support)
- Or Ubuntu Server 22.04 LTS AMI

The script will automatically detect the OS and use the appropriate package manager.

### Docker and Docker Compose

The deployment script automatically installs:
- **Docker** - Container runtime (installed via official Docker install script)
- **Docker Compose** - Legacy binary (`docker-compose` command)

**Note:** The script uses the legacy `docker-compose` binary (with hyphen) rather than the newer Docker Compose plugin (`docker compose` without hyphen). This ensures compatibility with the existing `docker-compose.yml` file and commands used throughout the documentation.

## DNS Configuration

Before deploying, configure the following DNS records to point to your EC2 instance's public IP:

### Required DNS Records

| Type | Name | Value | Purpose |
|------|------|-------|---------|
| A | `@` or root domain | EC2 Public IP | Main portfolio UI |
| A | `api` | EC2 Public IP | API endpoints |
| A | `outbox` | EC2 Public IP | Direct outbox API access (optional) |
| A | `flags` | EC2 Public IP | Direct feature flags API access (optional) |
| A | `metrics` | EC2 Public IP | Direct observability API access (optional) |

**Example for domain `example.com`:**
- `example.com` → EC2 IP
- `api.example.com` → EC2 IP
- `outbox.example.com` → EC2 IP
- `flags.example.com` → EC2 IP
- `metrics.example.com` → EC2 IP

### DNS Propagation

DNS changes can take up to 48 hours to propagate globally, though typically complete within a few hours. Verify DNS resolution before proceeding:

```bash
# Check DNS resolution
dig example.com +short
dig api.example.com +short
```

## Security Group Configuration

Configure your EC2 security group to allow inbound traffic:

| Type | Protocol | Port Range | Source | Purpose |
|------|----------|------------|--------|---------|
| HTTP | TCP | 80 | 0.0.0.0/0 | HTTP traffic (redirected to HTTPS) |
| HTTPS | TCP | 443 | 0.0.0.0/0 | HTTPS traffic |
| SSH | TCP | 22 | Your IP | SSH access (restrict to your IP) |

**Note:** Ports 3000, 4000, 8080, 8081, and 5432 should NOT be exposed publicly. They are only accessible within the Docker network.

## TLS/SSL Certificate Management

Caddy automatically manages SSL certificates using Let's Encrypt:

1. **Automatic Certificate Provisioning**: On first start, Caddy detects the domain configuration and automatically requests certificates from Let's Encrypt
2. **Certificate Renewal**: Caddy automatically renews certificates before expiration (Let's Encrypt certificates are valid for 90 days)
3. **Email Requirement**: An email address is required for Let's Encrypt notifications (configured in Caddyfile via `{{EMAIL}}` placeholder)

### Certificate Storage

Certificates are stored in the `caddy_data` Docker volume and persist across container restarts.

### Troubleshooting Certificate Issues

If certificate provisioning fails:

1. **Check DNS**: Ensure DNS records are properly configured and propagated
2. **Check Ports**: Verify ports 80 and 443 are open in the security group
3. **Check Logs**: Review Caddy logs: `docker-compose logs caddy`
4. **Rate Limits**: Let's Encrypt has rate limits (50 certificates per registered domain per week). If exceeded, wait or use staging environment

## Secrets Management

This section describes strategies for managing sensitive configuration values like database passwords and API keys.

### Environment Variables (Simple Approach)

For a portfolio project, the simplest approach is to use environment variables set directly in `docker-compose.yml` or via a `.env.production` file:

1. **Create `.env.production` from template**:
   ```bash
   cp env.production.template .env.production
   # Edit .env.production with your actual values
   ```

2. **Update `docker-compose.yml`** to use environment variables:
   ```yaml
   environment:
     DB_PASSWORD: ${DB_PASSWORD}
   ```

3. **Load from file** (modify `deploy.sh` to source the file):
   ```bash
   source .env.production
   docker-compose up -d
   ```

**Security Considerations:**
- Never commit `.env.production` to version control
- Restrict file permissions: `chmod 600 .env.production`
- Rotate passwords regularly
- Use strong, unique passwords

### AWS Secrets Manager (Recommended for Production)

For production deployments, consider using AWS Secrets Manager:

1. **Store secrets in AWS Secrets Manager**:
   ```bash
   aws secretsmanager create-secret \
     --name portfolio/db-password \
     --secret-string "your-strong-password"
   ```

2. **Retrieve secrets at runtime**:
   ```bash
   DB_PASSWORD=$(aws secretsmanager get-secret-value \
     --secret-id portfolio/db-password \
     --query SecretString --output text)
   ```

3. **Update `deploy.sh`** to fetch secrets before starting containers:
   ```bash
   # Fetch secrets from AWS Secrets Manager
   export DB_PASSWORD=$(aws secretsmanager get-secret-value \
     --secret-id portfolio/db-password \
     --query SecretString --output text)
   
   docker-compose up -d
   ```

**Benefits:**
- Centralized secret management
- Automatic rotation support
- Audit logging via CloudTrail
- IAM-based access control

### AWS Systems Manager Parameter Store (Alternative)

Parameter Store is a simpler alternative to Secrets Manager:

1. **Store secrets**:
   ```bash
   aws ssm put-parameter \
     --name /portfolio/db-password \
     --value "your-strong-password" \
     --type SecureString
   ```

2. **Retrieve at runtime**:
   ```bash
   DB_PASSWORD=$(aws ssm get-parameter \
     --name /portfolio/db-password \
     --with-decryption \
     --query Parameter.Value --output text)
   ```

**When to Use:**
- Parameter Store: Free, simpler, good for non-critical secrets
- Secrets Manager: Better for production, supports rotation, costs $0.40/secret/month

### Required Secrets

The following values should be treated as secrets and managed securely:

| Variable | Service | Purpose | Recommended Storage |
|----------|---------|---------|---------------------|
| `DB_PASSWORD` | PostgreSQL | Database authentication | Secrets Manager |
| `DB_USER` | PostgreSQL | Database user (if changed from default) | Parameter Store or env var |
| Any API keys | Various | External service authentication | Secrets Manager |

### Secret Rotation

**Manual Rotation:**
1. Update secret in Secrets Manager/Parameter Store
2. Restart affected services: `docker-compose restart outbox-api postgres`
3. Verify services are healthy: `docker-compose ps`

**Automated Rotation** (Secrets Manager):
- Configure automatic rotation using Lambda functions
- See [AWS Secrets Manager Rotation Documentation](https://docs.aws.amazon.com/secretsmanager/latest/userguide/rotating-secrets.html)

### IAM Permissions

If using AWS Secrets Manager or Parameter Store, ensure the EC2 instance has appropriate IAM permissions:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "secretsmanager:GetSecretValue",
        "ssm:GetParameter",
        "ssm:GetParameters"
      ],
      "Resource": [
        "arn:aws:secretsmanager:region:account:secret:portfolio/*",
        "arn:aws:ssm:region:account:parameter/portfolio/*"
      ]
    }
  ]
}
```

## Deployment Steps

1. **SSH into EC2 instance**
   ```bash
   ssh -i your-key.pem ec2-user@your-ec2-ip
   ```

2. **Run deployment script**
   ```bash
   bash deploy.sh
   ```

3. **Follow prompts**:
   - Enter your domain name (e.g., `example.com`)
   - Enter your email for SSL certificates

4. **Wait for services to start** (approximately 30-60 seconds)

5. **Verify deployment**:
   ```bash
   docker-compose ps
   curl https://your-domain.com/health
   ```

## Post-Deployment Verification

### Health Checks

Verify all services are healthy:

```bash
# Check container status
docker-compose ps

# Test main UI
curl https://your-domain.com/health

# Test API endpoints
curl https://api.your-domain.com/outbox/health
curl https://api.your-domain.com/flags/health
curl https://api.your-domain.com/observability/health
```

### SSL Certificate Verification

Check that SSL certificates are properly configured:

```bash
# View Caddy logs
docker-compose logs caddy | grep -i certificate

# Test SSL connection
openssl s_client -connect your-domain.com:443 -servername your-domain.com
```

## Updating the Deployment

To update the application:

1. **Pull latest changes**
   ```bash
   cd /opt/portfolio/portfolio-monorepo
   git pull
   ```

2. **Rebuild and restart**
   ```bash
   docker-compose up --build -d
   ```

3. **Verify services**
   ```bash
   docker-compose ps
   ```

## Monitoring and Logs

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f caddy
docker-compose logs -f portfolio
docker-compose logs -f outbox-api
```

### Caddy Access Logs

Caddy logs are stored in Docker volumes:
- UI access: `/var/log/caddy/access.log` (inside container)
- API access: `/var/log/caddy/api.log` (inside container)

To access logs:
```bash
docker-compose exec caddy cat /var/log/caddy/access.log
docker-compose exec caddy cat /var/log/caddy/api.log
```

## Troubleshooting

### Services Not Starting

1. Check container logs: `docker-compose logs [service-name]`
2. Verify Docker is running: `docker ps`
3. Check disk space: `df -h`
4. Verify ports aren't in use: `netstat -tulpn | grep [port]`

### DNS Not Resolving

1. Verify DNS records are correct: `dig your-domain.com`
2. Check DNS propagation: Use online tools like `whatsmydns.net`
3. Verify EC2 public IP hasn't changed (use Elastic IP for static IP)

### SSL Certificate Issues

1. Check Caddy logs: `docker-compose logs caddy`
2. Verify DNS is pointing to correct IP
3. Ensure ports 80 and 443 are open
4. Check Let's Encrypt rate limits (if exceeded, wait or use staging)

### Database Connection Issues

1. Verify PostgreSQL is running: `docker-compose ps postgres`
2. Check database logs: `docker-compose logs postgres`
3. Verify connection string in environment variables

## Security Considerations

1. **Change Default Passwords**: Update `DB_PASSWORD` in `.env.production` or use AWS Secrets Manager (see [Secrets Management](#secrets-management) section)
2. **Restrict SSH Access**: Limit SSH (port 22) to your IP address in security group
3. **Use Secrets Management**: See the [Secrets Management](#secrets-management) section for best practices on managing sensitive values
4. **Never Commit Secrets**: Ensure `.env.production` is in `.gitignore` and never commit it to version control
5. **Regular Updates**: Keep Docker images and system packages updated
6. **Backup Database**: Implement regular backups of PostgreSQL data volume
7. **Rotate Secrets**: Regularly rotate passwords and API keys (see [Secret Rotation](#secret-rotation))

## Cost Optimization

- **t3.micro**: Free tier eligible, sufficient for portfolio site
- **Elastic IP**: Use Elastic IP to avoid IP changes (free when attached to running instance)
- **Volume Storage**: Default EBS volume size is typically sufficient
- **Data Transfer**: First 1GB/month free, then $0.09/GB

## Additional Resources

- [Caddy Documentation](https://caddyserver.com/docs/)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/)
- [AWS EC2 Documentation](https://docs.aws.amazon.com/ec2/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)

