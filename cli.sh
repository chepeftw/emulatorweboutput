#!/bin/bash

S='docker-compose.yml'

function dock_build {
    docker-compose -f $S build --no-cache
}

function dock_start {
    docker-compose -f $S up -d
}

function dock_stop {
    docker-compose -f $S down
}

function gitupdate {
    git pull origin master
}



while [[ $# -gt 1 ]]
do
key="$1"

case $key in
    -a|--action)
        ACTION="$2"
        shift # past argument
    ;;
    *)
        # unknown option
    ;;
esac
shift # past argument or value
done

case "$ACTION" in
    build)
        dock_build
        ;;

    start)
        dock_start
        ;;

    stop)
        dock_stop
        ;;

    update)
        gitupdate
        dock_build
        dock_start
        ;;

    *)
        echo $"Usage: $0 -a {start|stop|build}"
        exit 1
esac
