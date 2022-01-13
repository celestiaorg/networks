# 1. Running a Node
Here we describe how to run abd sync a node using celestia-appd.
## 1.2 Installing Dependencies
First, make sure to update and upgrade de OS:
```
sudo apt update && sudo apt upgrade -y
```
Install essential packages:
```
sudo apt install curl tar wget clang pkg-config libssl-dev jq build-essential git make ncdu -y
```
Install GO:
```
ver="1.17.2"
cd $HOME
wget "https://golang.org/dl/go$ver.linux-amd64.tar.gz"
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf "go$ver.linux-amd64.tar.gz"
rm "go$ver.linux-amd64.tar.gz"
echo "export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin" >> $HOME/.bash_profile
source $HOME/.bash_profile
```
To check if Go was installed correctly run:
```
go version
```
Output should be the version installed:
```
go version go1.17.2 linux/amd64
```
## 1.3 Downloading and Compiling Celestia-App
```
cd $HOME
rm -rf celestia-app
git clone https://github.com/celestiaorg/celestia-app.git
cd celestia-app/
make install
```
The steps above will create a binary file named celestia-appd inside $HOME/go/bin folder which will be used later to run the node.
## 1.4 Setting up the Network
First clone the networks repository:
```
cd $HOME
rm -rf networks
git clone https://github.com/celestiaorg/networks.git
```
To initialize the network pick a "node-name" the describes your node. The --chain-id parameter we are using here is "devnet-2". Keep in mind that this might change if a new testnet is deployed.
```
celestia-appd init "node-name" --chain-id devnet-2
```
Copy the genesis file. For devnet-2 we are using:
```
cp $HOME/networks/devnet-2/genesis.json $HOME/.celestia-app/config
```
Set seeds and peers:
```
SEEDS="74c0c793db07edd9b9ec17b076cea1a02dca511f@46.101.28.34:26656"
PEERS="34d4bfec8998a8fac6393a14c5ae151cf6a5762f@194.163.191.41:26656"
sed -i.bak -e "s/^seeds *=.*/seeds = \"$SEEDS\"/; s/^persistent_peers *=.*/persistent_peers = \"$PEERS\"/" $HOME/.celestia-app/config/config.toml
```
Reset network:
```
celestia-appd unsafe-reset-all
```
## 1.5 Run Ceslestia-App using Systemd
Create Celestia-App systemd file:
```
sudo tee <<EOF >/dev/null /etc/systemd/system/celestia-appd.service
[Unit]
Description=celestia-appd Cosmos daemon
After=network-online.target

[Service]
User=$USER
ExecStart=$HOME/go/bin/celestia-appd start
Restart=on-failure
RestartSec=3
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target
EOF
```
If the file was created succesfully you will be able to see its content:
```
cat /etc/systemd/system/celestia-appd.service
```
Enable and start celestia-appd daemon:
```
sudo systemctl enable celestia-appd
sudo systemctl daemon-reload
sudo systemctl restart celestia-appd
```
Check if the daemon has been started correctly:
```
sudo systemctl status celestia-appd
```
Check the daemon logs in real time:
```
sudo journalctl -u celestia-appd.service -f
```
Wait for the node to sync before doing next steps.
