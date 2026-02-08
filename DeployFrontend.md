
# Steps to update React project in S3

```bash
# Build frontend with production API URL
cd frontend

# Create production environment file
echo "VITE_API_URL=http://<YOUR_EC2_PUBLIC_IP>:8080" > .env.production

# Build the app
npm run build

# Install AWS CLI if needed
winget install Amazon.AWSCLI

# Configure AWS credentials
aws configure

# Create S3 bucket and upload
aws s3 mb s3://trash-taste-search-frontend --region us-east-1
aws s3 sync dist/ s3://trash-taste-search-frontend --delete

# Enable static website hosting
aws s3 website s3://trash-taste-search-frontend --index-document index.html --error-document index.html
```