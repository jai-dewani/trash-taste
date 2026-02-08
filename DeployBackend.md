1. Go to EC2 Dashboard → Click Launch Instance

2. Configure the instance:  

| Setting       | Value                     |
| ------------- | -------------             |
| Name          | trash-taste-search        |
| AMI	Amazon  | Linux 2023                |
| Instance type | t3.micro (Free Tier)      |
| Key pair      | new → Download .pem file  |

3.  Network Settings → Click Edit and configure Security Group:

| Type	        | Port	| Source	            | Purpose       |
| -----         | ----- | ---------             | -------       |
| SSH           | 22    | My IP	                | SSH access    |
| HTTP          | 80    | Anywhere (0.0.0.0/0)	| Web traffic   |
| Custom TCP    | 8080  | Anywhere (0.0.0.0/0)	| Go API        |


4. Connect to Your Instance
```bash
# Move your key file to a secure location
Move-Item Downloads\your-key.pem ~/.ssh/

# Fix permissions (PowerShell)
icacls .\.ssh\Trash-Taste-Search.pem /inheritance:r /grant:r "$($env:USERNAME):(R)"

# Connect via SSH
ssh -i .\.ssh\trash-taste-search-ec2.pem ec2-user@ec2-18-232-174-140.compute-1.amazonaws.com
```

5.  Install Dependencies on EC2
```bash 
# Update system packages
sudo dnf update -y

# Install Go
sudo dnf install -y golang

# Install Python 
sudo dnf install -y python

# Install Tmux 
sudo dnf install -y tmux

# Install Git
sudo dnf install -y git

# Verify installations
go version    # Should show go1.21+
git --version
```

6. Build the database & Deploy Your Application
```bash
# Clone your repository
cd ~
git clone https://github.com/jai-dewani/trash-taste.git
cd trash-taste-search

# Build the sqlite db
python scripts/generate_database.py

# Create 2GB swap file
sudo dd if=/dev/zero of=/swapfile bs=128M count=16
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile

# Make it permanent
echo '/swapfile swap swap defaults 0 0' | sudo tee -a /etc/fstab

# Build the Go backend
cd backend
go mod download
go build -x -o server cmd/server/main.go

# Test run (Ctrl+C to stop)
./server
```

7. Set Up as a System Service

This ensures your app runs on boot and auto-restarts on crashes:
```bash 
# Create systemd service file
sudo nano /etc/systemd/system/trash-taste-api.service
```

```bash
[Unit]
Description=Trash Taste Search API
After=network.target

[Service]
Type=simple
User=ec2-user
WorkingDirectory=/home/ec2-user/trash-taste/backend
ExecStart=/home/ec2-user/trash-taste/backend/server
Restart=always
RestartSec=5
Environment=PORT=8080
Environment=DATABASE_PATH=/home/ec2-user/trash-taste/backend/data/trash_taste.db

[Install]
WantedBy=multi-user.target
```

```bash 
# Enable and start the service
sudo systemctl daemon-reload
sudo systemctl enable trash-taste-api
sudo systemctl start trash-taste-api

# Check status
sudo systemctl status trash-taste-api

# View logs if needed
sudo journalctl -u trash-taste-api -f
```

8. Verify Deployment 
```bash 
# Test locally on EC2
curl http://localhost:8080/api/health

# Test from your local machine (PowerShell)
curl http://<YOUR_EC2_PUBLIC_IP>:8080/api/health
```