#!/bin/bash

# Ensure root user
sudo su

# Update system
yum update -y

# Install Python 3, pip, cloudwatch logs agent, collectd for system logs
yum install -y python3 python3-pip amazon-cloudwatch-agent jq

# Install dependencies
pip3 install flask gunicorn

# Install aws-cli
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install --bin-dir /usr/local/bin --install-dir /usr/local/aws-cli --update

# Fetch application code
/usr/local/bin/aws s3 cp s3://${s3_bucket} /home/ec2-user/ --recursive
chmod +x /home/ec2-user/app.py

# Setup Web Server to run as daemon process
echo '
[Unit]
Description=Gunicorn instance for a simple web server
After=network.target

[Service]
User=root
Group=root
WorkingDirectory=/home/ec2-user
ExecStart=/usr/local/bin/gunicorn -b 0.0.0.0:80 app:app

Restart=always

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
