#!/bin/bash
## --------------------------------------------------------------------------------------------------------------------

set -e

echo "Checking ask.sh is installed ..."
if [ ! /home/chilts/bin/ask.sh ]; then
    echo "Please put ask.sh into ~/bin (should already be in your path from ~/.profile):"
    echo ""
    echo "    mkdir ~/bin"
    echo "    wget -O ~/bin/ask.sh https://gist.githubusercontent.com/chilts/6b547307a6717d53e14f7403d58849dd/raw/ecead4db87ad4e7674efac5ab0e7a04845be642c/ask.sh"
    echo "    chmod +x ~/bin/ask.sh"
    echo ""
    exit 2
fi
echo

# General
WHO=`whoami`
CSSMINIFIER_PORT=`ask.sh cssminifier CSSMINIFIER_PORT 'Which local port should the server listen on (e.g. 8420):'`
CSSMINIFIER_APEX=`ask.sh cssminifier CSSMINIFIER_APEX 'What is the apex (e.g. localhost:8420 or cssminifier.com) :'`
CSSMINIFIER_BASE_URL=`ask.sh cssminifier CSSMINIFIER_BASE_URL 'What is the base URL (e.g. http://localhost:1234 or https://cssminifier.com) :'`
CSSMINIFIER_DIR=`ask.sh cssminifier CSSMINIFIER_DIR 'What is the storage dir (e.g. /var/lib/cssminifier/raw) :'`
CSSMINIFIER_GOOGLE_ANALYTICS=`ask.sh cssminifier CSSMINIFIER_GOOGLE_ANALYTICS 'What is the Google Analytics code (e.g. UA-123-4) :'`

echo "Building code ..."
gb build
echo

echo "Minifying assets ..."
make minify
echo

echo "Creating storage dir ..."
sudo mkdir -p $CSSMINIFIER_DIR
sudo chown ${WHO}.${WHO} $CSSMINIFIER_DIR
echo

# copy the supervisor script into place
echo "Copying supervisor config ..."
m4 \
    -D __CSSMINIFIER_PORT__=$CSSMINIFIER_PORT \
    -D __CSSMINIFIER_APEX__=$CSSMINIFIER_APEX \
    -D __CSSMINIFIER_BASE_URL__=$CSSMINIFIER_BASE_URL \
    -D __CSSMINIFIER_DIR__=$CSSMINIFIER_DIR \
    -D __CSSMINIFIER_GOOGLE_ANALYTICS__=$CSSMINIFIER_GOOGLE_ANALYTICS \
    etc/supervisor/conf.d/cssminifier.conf.m4 | sudo tee /etc/supervisor/conf.d/cssminifier.conf
echo

# restart supervisor
echo "Restarting supervisor ..."
sudo systemctl restart supervisor.service
echo

# copy the caddy conf
echo "Copying Caddy config config ..."
m4 \
    -D __CSSMINIFIER_PORT__=$CSSMINIFIER_PORT \
    -D __CSSMINIFIER_APEX__=$CSSMINIFIER_APEX \
    etc/caddy/vhosts/cssminifier.conf.m4 | sudo tee /etc/caddy/vhosts/cssminifier.conf
echo

# restarting Caddy
echo "Restarting caddy ..."
sudo systemctl restart caddy.service
echo

## --------------------------------------------------------------------------------------------------------------------
