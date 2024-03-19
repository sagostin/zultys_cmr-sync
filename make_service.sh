#!/bin/bash

# Check if the script is run as root
if [ "$(id -u)" -ne 0 ]; then
  echo "This script must be run as root"
  exit 1
fi

# Check if an argument is provided
if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <config_name>"
  exit 1
fi

CONFIG_NAME=$1

# Copy the service template to systemd directory
cp /opt/crm_sync/zultys_crm-sync@.service /etc/systemd/system/zultys_crm-sync@$CONFIG_NAME.service

# Reload systemd daemon to recognize the new service file
systemctl daemon-reload

# Enable the new service to start on boot
systemctl enable zultys_crm-sync@$CONFIG_NAME
systemctl start zultys_crm-sync@$CONFIG_NAME

echo "Service zultys_crm-sync@$CONFIG_NAME has been created and enabled."