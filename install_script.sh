#!/bin/bash

function install_docker {
    apt-get update
    apt-get install -y docker.io docker-compose wget git
}

function install_ss {
    wget  https://raw.githubusercontent.com/Jigsaw-Code/outline-apps/master/server_manager/install_scripts/install_server.sh
    chmod u+x ./install_server.sh
}

function clone_repo {
    git clone https://github.com/DobbyVPN/dobbyvpn-server.git
}

function generate_url {
    URL=`head /dev/urandom | tr -dc A-Za-z0-9 | head -c40`
    echo $(( $URL ))
}


function run_ss {
    ./install_server.sh --api-port $1 --keys-port $2
}

install_docker
install_ss
clone_repo

