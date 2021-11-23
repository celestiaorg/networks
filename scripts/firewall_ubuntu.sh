#!/bin/bash

sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow OpenSSH
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 60000:61000/udp
sudo ufw allow 26657/tcp
sudo ufw allow 26656/tcp
sudo ufw allow 2121/tcp
sudo ufw allow 1317/tcp
sudo ufw allow 8080/tcp
sudo ufw --force enable