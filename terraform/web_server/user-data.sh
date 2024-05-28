#!/bin/bash

# Ensure root user
sudo su

# Update system
yum update -y

# Install cloudwatch logs agent, collectd for system logs
yum install -y amazon-cloudwatch-agent jq git

# Install Golang
wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz
tar -xvf go1.22.3.linux-amd64.tar.gz
mv go /usr/local
rm go1.22.3.linux-amd64.tar.gz

# Configure Go environment
export GOROOT=/usr/local/go
export GOPATH=/home/ec2-user
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >>/home/root/.bashrc

# Install aws-cli
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install --bin-dir /usr/local/bin --install-dir /usr/local/aws-cli --update

# Fetch application code
mkdir -p /home/ec2-user/src
/usr/local/bin/aws s3 cp s3://${s3_bucket} /home/ec2-user/src --recursive
chmod +x /home/ec2-user/src/server

# Setup Web Server to run as daemon process
echo '
[Unit]
Description=Go web server
After=network.target

[Service]
Type=simple
ExecStart=/home/ec2-user/src/server
WorkingDirectory=/home/ec2-user/src
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
' >/etc/systemd/system/web-server.service

# Reload daemon processes
systemctl daemon-reload

# Start web server
systemctl start web-server

# Enable daemon in systemd to run on startup
systemctl enable web-server

# Print status to logfile: /var/log/cloud-init-output.log
systemctl status web-server

# Fetch cloudwatch log configuration
/usr/local/bin/aws ssm get-parameter --name "${ssm_parameter}" --region us-west-1 | jq -r '.Parameter.Value | fromjson' >/opt/aws/amazon-cloudwatch-agent/bin/config.json

# Start cloudwatct worked that workedh logs agent
/opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl -a fetch-config -m ec2 -c file:/opt/aws/amazon-cloudwatch-agent/bin/config.json -s

# Output cloudwatch agent status
/opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl -a status
