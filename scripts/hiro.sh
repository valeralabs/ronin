# Shell script to initialize a Ubuntu 20.04 LTS server with Hiro's API

export MAINNET_ARCHIVE=https://docker.06815d71-a2bc-4176-98bb-dccd6c237f84.uk-lon1.upcloudobjects.com/mainnet.tar.gz

echo "Updating APT..."
sudo apt-get update -qq

echo "Installing clang, JQ, PostgreSQL client & PV..."
sudo apt-get install clang postgresql-client-common jq pv -y -qq

if [ -x "$(command -v docker)" ]; then
    echo "Docker is installed, skipping install..."
else
    echo "Installing Docker..."
    curl https://get.docker.com -sSf | sh
fi

if [ -x "$(command -v cargo)" ]; then
    echo "Cargo is installed, skipping install..."
else
    echo "Installing Cargo..."
    curl https://sh.rustup.rs -sSf | sh -s -- -y
    source $HOME/.cargo/env
fi

echo "Installing b3sum, toml-cli & bottom..."
cargo install b3sum toml-cli bottom

VERSION=$(curl --silent https://api.github.com/repos/docker/compose/releases/latest | jq .name -r)
echo "Installing Docker Compose $VERSION"
DESTINATION=/usr/local/bin/docker-compose
sudo curl -L --silent https://github.com/docker/compose/releases/download/${VERSION}/docker-compose-$(uname -s)-$(uname -m) -o $DESTINATION
sudo chmod 755 $DESTINATION

echo "Checking for existing archive..."
if [ -f "mainnet.tar.gz" ]; then
    echo "Archive found, skipping download..."
else
    echo "Downloading node + snapshot..."
    curl $MAINNET_ARCHIVE -o mainnet.tar.gz
fi

echo "Verifiying integrity..."

export EXPECTED=a43298502be0f3ab5e8b2dfe76ed24c8826c00a7935f85db2392454e794d30fe

if pv mainnet.tar.gz | b3sum --no-names | grep -q $EXPECTED; then 
    echo "Integrity verified."
else
    echo "Integrity failure. Removing archive. Please try again."
    rm mainnet.tar.gz
    exit 1
fi

echo "Extracting node..."

pv mainnet.tar.gz | tar -xz

cd stacks-local-dev

echo "###############################
## Stacks-Node-API
##
NODE_ENV=production
GIT_TAG=master
PG_HOST=postgres
PG_PORT=5432
PG_USER=postgres
PG_PASSWORD=postgres
PG_DATABASE=postgres
STACKS_CHAIN_ID=2147483648
V2_POX_MIN_AMOUNT_USTX=90000000260
STACKS_CORE_EVENT_PORT=3700
STACKS_CORE_EVENT_HOST=0.0.0.0
STACKS_BLOCKCHAIN_API_PORT=3999
STACKS_BLOCKCHAIN_API_HOST=0.0.0.0
STACKS_BLOCKCHAIN_API_DB=pg
STACKS_CORE_RPC_HOST=stacks-blockchain
STACKS_CORE_RPC_PORT=20443
STACKS_EXPORT_EVENTS_FILE=/tmp/event-replay/stacks-node-events.tsv
#BNS_IMPORT_DIR=/bns-data

###############################
## Postgres
##
# Make sure the password is the same as PG_PASSWORD above.
# note to document: this is set in the sql for postgres. if the above is changed, that needs to change as well. 
POSTGRES_PASSWORD=postgres

###############################
## Docker image versions
## 
STACKS_BLOCKCHAIN_VERSION=2.05.0.1.0-A
STACKS_BLOCKCHAIN_API_VERSION=1.0.7

# version of the postgres image to use (if there is existing data, set to this to version 13)
POSTGRES_VERSION=13" > .env
echo "Configured environment"
