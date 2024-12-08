#!/bin/bash

function install_docker {
    echo "Docker installation..."	
    apt-get update
    apt-get install -y docker.io docker-compose wget git
}

function install_ss {
    echo "Outline installation..."	
    wget  https://raw.githubusercontent.com/Jigsaw-Code/outline-apps/master/server_manager/install_scripts/install_server.sh
    chmod u+x ./install_server.sh
}

function install_ck_server {
    # It needs for key generation in system
    wget https://github.com/cbeuw/Cloak/releases/download/v2.7.0/ck-server-linux-amd64-v2.7.0 -O ck-server
    chmod +x ck-server
    mv ck-server /bin/ck-server # sudo permissions!
}

function clone_repo {
    echo "Repo cloning..."	
    git clone https://github.com/DobbyVPN/dobbyvpn-server.git
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
    sed -i "s/<domain-name>/$1/g" ./Caddyfile
    sed -i "s/<special-url>/$2/g" ./Caddyfile
    sed -i "s/<cloak-server-port>/$3/g" ./Caddyfile
}

function replace_cloak_holders {
    # Accepts 6 args:
    # $1 - keys-port for ss
    # $2 - cloak-server port
    # $3 - bypassUID
    # $4 - adminUID
    # $5 - domain-name (for RedirAddr)
    # $6 - cloak private key
    
    cp ./cloak-server-template.conf ./cloak-server.conf
    sed -i "s/<keys-port>/$1/g" ./cloak-server.conf
    sed -i "s/<cloak-server-port>/$2/g" ./cloak-server.conf
    sed -i "s/<user-UID>/$3/g" ./cloak-server.conf
    sed -i "s/<admin-UID>/$4/g" ./cloak-server.conf
    sed -i "s/<domain-name>/$5/g" ./cloak-server.conf
    sed -i "s/<cloak-private-key>/$6/g" ./cloak-server.conf
}

function save_credentials {
   # Function saves sensitive data to file
   # $1 - filename
   # $2 - associative array
    local file=$1
    local array=$2

    echo "Saving credentials"
    # if [ -e "$file" ]; then
    #     echo "$file already exists."
    #     read -e -p "Do you want to override it?(Y/n): " choice
    #     case "$choice" in
	#         y|Y)
    #             rm $file
    #             for key in "${!array[@]}"; do
    #                 echo "$key => ${array[$key]}" >> "$file"
    #             done
    #             return
    #    	        ;;
	#         n|N)
    #             return
    # 	        ;;
    #     esac	 	     
	     
    # fi


   for key in "${!array[@]}"; do
    echo "$key => ${array[$key]}"
    echo "$key => ${array[$key]}" >> "$file"
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

    KEYPAIRS=$(bin/ck-server -key)
    CLOAK_PRIVATE_KEY=$(echo $KEYPAIRS | cut -d" " -f13)
    CLOAK_PUBLIC_KEY=$(echo $KEYPAIRS | cut -d" " -f5)
    USER_UID=$(bin/ck-server -uid | cut -d" " -f4)
    ADMIN_UID=$(bin/ck-server -uid | cut -d" " -f4)
}

function main {

    install_ck_server
    readArgs	

    declare -A array_creds
    install_docker
    install_ss
    clone_repo
    run_ss $OUTLINE_API_PORT $OUTLINE_KEYS_PORT

    URL=$(generate_url)
    replace_caddy_holders $DOMAIN_NAME $URL $CLOAK_PORT
    replace_cloak_holders $OUTLINE_KEYS_PORT $CLOAK_PORT $USER_UID $ADMIN_UID $DOMAIN_NAME $CLOAK_PRIVATE_KEY

    docker-compose -f ./docker-compose.yaml up -d

    filename="creds.txt"
    array_creds["Special-url"]=$URL
    array_creds["Cloak-public-key"]=$CLOAK_PUBLIC_KEY
    array_creds["Cloak-private-key"]=$CLOAK_PRIVATE_KEY
    array_creds["User-uid"]=$USER_UID
    array_creds["Admin-uid"]=$ADMIN_UID
    save_credentials $filename $array_creds

    echo "All credentials are saved in $filename"
    echo "Done!"
}

main

