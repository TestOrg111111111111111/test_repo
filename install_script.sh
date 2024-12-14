#!/bin/bash

function install_docker {
    echo "Docker installation..."	
    #apt-get update
    #apt-get install -y docker.io docker-compose wget git
    # Add Docker's official GPG key:
    sudo apt-get update
    sudo apt-get install ca-certificates curl
    sudo install -m 0755 -d /etc/apt/keyrings
    sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    sudo chmod a+r /etc/apt/keyrings/docker.asc

    # Add the repository to Apt sources:
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
      $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
      sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    sudo apt-get update

    sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
}

function install_ss {
    echo "Outline installation..."	
    if ! [ -f "install_server.sh" ]
    then
        wget  https://raw.githubusercontent.com/Jigsaw-Code/outline-apps/master/server_manager/install_scripts/install_server.sh
        chmod u+x ./install_server.sh
    fi
}

function generate_url {
    echo `head /dev/urandom | tr -dc A-Za-z0-9 | head -c40`
}


function run_ss {
    # Accepts 2 args: $1 - api-port, $2 - keys-port
    echo "Outline server runninng..."   
    ./install_server.sh --api-port $1 --keys-port $2 
}

function replace_caddy_holders {
    # Accepts 3 args: $1 - domain name, $2 - secret-url, $3 - cloak-server port
    rm -rf Caddyfile
    cp "Caddyfile-template" "Caddyfile"
    sed -i "s|<domain-name>|${1}|" "Caddyfile"
    sed -i "s|<special-url>|${2}|" "Caddyfile"
    sed -i "s|<cloak-server-port>|${3}|" "Caddyfile"
}


function save_credentials {
   # Function saves sensitive data to file
   # $1 - filename

    echo "Saving credentials"
    if [ -f "$1" ]
    then
        echo "$1 already exists."
        read -e -p "Do you want to override it?(Y/n): " choice
        case "$choice" in
	        y|Y)
                rm $1
                for key in "${!array_creds[@]}"; do
                    echo "$key => ${array_creds[$key]}" >> "$1"
                done
                return
       	        ;;
	        n|N)
                return
    	        ;;
        esac 
    fi

    for key in "${!array_creds[@]}"; do
        echo "$key => ${array_creds[$key]}" >> "$1"
    done
}

function readArgs {
    read -e -p "Enter Cloak Port: " -i 8443 CLOAK_PORT
    read -e -p "Enter Api Port(outline): " -i 11111 OUTLINE_API_PORT
    read -e -p "Enter Keys Port(outline): " -i 22222 OUTLINE_KEYS_PORT
    read -e -p "Enter Domain Name: " DOMAIN_NAME

    if [ -z "$DOMAIN_NAME" ]; then
        echo "Error: you didn't enter domain name!" >&2
        exit 1
    fi
}

function stop_and_remove_caddy_cloak {
    CONTAINERS=("caddy" "ck-server")

    for CONTAINER in "${CONTAINERS[@]}"; do
    CONTAINER_ID=$(docker ps -aq --filter "name=^${CONTAINER}$")
    
    if [ -n "$CONTAINER_ID" ]; then
        echo "Found container: $CONTAINER (ID: $CONTAINER_ID)"
        docker stop "$CONTAINER_ID" && echo "Stopped container: $CONTAINER"
        docker rm "$CONTAINER_ID" && echo "Removed container: $CONTAINER"
    else
        echo "No container found with name: $CONTAINER"
    fi
    done
}

function replace_holders_cloak_start {
    # replaces all holders in start script for cloak
    # $1 - keys-port for outline server
    # $2 - bind port for cloak server
    # $3 - domain name
    sed -i "s|<1>|$1|" "cloak_start.sh"
    sed -i "s|<2>|$2|" "cloak_start.sh"
    sed -i "s|<3>|$3|" "cloak_start.sh"
}

function main {

    readArgs	

    install_docker
    install_ss
    run_ss $OUTLINE_API_PORT $OUTLINE_KEYS_PORT

    URL=$(generate_url)
    replace_caddy_holders $DOMAIN_NAME $URL $CLOAK_PORT

    stop_and_remove_caddy_cloak

    rm -rf cloak_start.sh
    cp "cloak_start_template.sh" "cloak_start.sh"
	
    replace_holders_cloak_start $OUTLINE_KEYS_PORT $CLOAK_PORT $DOMAIN_NAME
    chmod a+x ./cloak_start.sh

    rm -rf cloak-server.conf
    cp "cloak-server-template.conf" "cloak-server.conf"

    docker-compose -f docker-compose.yaml up -d


    filename="creds.txt"
    declare -A array_creds
    array_creds["Special-url"]=$URL

    save_credentials $filename

    echo "All credentials are saved in $filename"
    echo "Done!"
}

main

