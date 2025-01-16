#!/bin/bash

DOCKER_COMPOSE_FILE="docker-compose.yml"

PREV_IP=""

get_current_ip() {
  ip=$(ifconfig | grep "inet " | grep -v 127.0.0.1 | awk '{print $2}')
  echo $ip
}

redeploy() {
    echo "IP changed! Redeploying docker containers..."
    docker-compose -f $DOCKER_COMPOSE_FILE down 
    env HOST_IP=$CURRENT_IP docker-compose -f $DOCKER_COMPOSE_FILE up -d
    echo "Docker containers redeployed."
}

# Main loop to monitor IP changes
while true; do
    CURRENT_IP=$(get_current_ip)

    # Check if the IP has changed
    if [[ "$CURRENT_IP" != "$PREV_IP" && ! -z "$CURRENT_IP" ]]; then
        echo "Detected IP change: $PREV_IP -> $CURRENT_IP"
        PREV_IP=$CURRENT_IP
        redeploy
    fi

    # Wait for 5 seconds before checking again
    sleep 5
done
