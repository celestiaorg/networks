# Attempt to create network on 3 DO Droplets

Following Instructions from [here](https://github.com/celestiaorg/networks/blob/a378c7cddb91a71db533631d7bbc2b67cb956d5c/README.md)

## Run on all Nodes
```
git clone https://github.com/celestiaorg/celestia-app.git
cd celestia-app 
apt install make
cd ..
wget https://go.dev/dl/go1.17.4.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.17.4.linux-amd64.tar.gz && rm go1.17.4.linux-amd64.tar.gz
echo "PATH=$PATH:/usr/local/go/bin" >> ~/.profile
echo "PATH=$PATH:~/go/bin" >> ~/.profile
source ~/.profile
cd celestia-app
make install
git clone https://github.com/celestiaorg/networks.git ~/
cd ../networks/scripts/
./firewall_ubuntu.sh
node_name="MightyValidator-${node-num}"
./1_create_key.sh $node_name
# Copy Account Address to file
echo "celes1san3ge6n37jspwsmjl0llh9dveyvylwnrypnff" > ~/account-address-$node_num.txt
```

## Run on a terminal not ssh'd to a node
```
mkdir /tmp/orchestrate && cd /tmp/orchestrate
node_0_ip=188.166.21.209
node_1_ip=188.166.91.191
node_2_ip=174.138.11.26
scp root@$node_1_ip:/root/account-address-1.txt ./account-address-1.txt
scp root@$node_2_ip:/root/account-address-2.txt ./account-address-2.txt
scp ./account-address-1.txt root@$node_0_ip:/root/account-address-1.txt
scp ./account-address-2.txt root@$node_0_ip:/root/account-address-2.txt
```

## Run on orchestrator node
```
cd ~
apt install moreutils -y
apt install jq -y
# Init genesis
celestia-appd init ReiAyanami --chain-id devnet-2
mv account-address.txt account-address-0.txt
acc_addr_0=$(cat account-address-0.txt)
acc_addr_1=$(cat account-address-1.txt)
acc_addr_2=$(cat account-address-2.txt)
celestia-appd add-genesis-account $acc_addr_0 800000000000celes
celestia-appd add-genesis-account $acc_addr_1 800000000000celes
celestia-appd add-genesis-account $acc_addr_2 800000000000celes
networks/scripts/fix_genesis.sh ~/.celestia-app/config/genesis.json celes
```

## Run on a terminal not ssh'd to a node
```
# retrieve genesis.json from node0
scp root@$node_0_ip:/root/.celestia-app/config/genesis.json ./genesis.json 
# Send genesis.json to node1
scp ./genesis.json root@$node_1_ip:/root/.celestia-app/config/genesis.json
# Send genesis.json to node2
scp ./genesis.json root@$node_2_ip:/root/.celestia-app/config/genesis.json
```

## Run on all nodes
```
# Validate genesis.json is same
cat ~/.celestia-app/config/genesis.json | jq .app_state.bank.balances
celestia-appd gentx $node_name 5000000000celes --keyring-backend=test --chain-id devnet-2
```

## Run on a terminal not ssh'd to a node
```
mkdir gentx
scp root@$node_0_ip:/root/.celestia-app/config/gentx/* ./gentx/.
scp root@$node_1_ip:/root/.celestia-app/config/gentx/* ./gentx/.
scp root@$node_2_ip:/root/.celestia-app/config/gentx/* ./gentx/.
scp -r  gentx/ root@$node_0_ip:/root/.celestia-app/config/.
```

## Run on orchestrator node
```
celestia-appd collect-gentxs
```

## Run on a terminal not ssh'd to a node
```
rm genesis.json
# Copy final genesis.json from node0
scp root@$node_0_ip:/root/.celestia-app/config/genesis.json ./golden-genesis.json
scp ./golden-genesis.json root@$node_1_ip:/root/.celestia-app/config/genesis.json
scp ./golden-genesis.json root@$node_2_ip:/root/.celestia-app/config/genesis.json
```

## Run on all nodes
```
celestia-appd start
```