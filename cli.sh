#!/bin/bash
# Use -gt 1 to consume two arguments per pass in the loop (e.g. each
# argument has a corresponding value to go with it).
# Use -gt 0 to consume one or more arguments per pass in the loop (e.g.
# some arguments don't have a corresponding value to go with it such
# as in the --default example).

S='docker-compose.yml'
FORCE_FLAG=false

# Docker Build, Start and Stop
# This 3 functions will heavily rely on docker-compose and the yml compose conf file.
# The intention is to build, start or stop the instances.
# The start function, not only start the docker infrastructure, but also assigns some
# permissions to some nginx folder that caused me problems in the past.
# Lastly it updates symfony schema.
function dock_build {
    docker-compose -f $S build --no-cache
}

function dock_start {
    docker-compose -f $S up -d
}

function dock_stop {
    docker-compose -f $S down
}



while [[ $# -gt 1 ]]
do
key="$1"

case $key in
    -e|--environment)
        ENV="$2"
        shift # past argument
    ;;
    -a|--action)
        ACTION="$2"
        shift # past argument
    ;;
    -f|--force)
        FORCE_FLAG=$2
        echo "Force ..."
        shift # past argument
    ;;
    *)
        # unknown option
    ;;
esac
shift # past argument or value
done

#ENV="dev"

if [ "$ENV" == "dev" ];
then
    echo "Dev ..."
    S='docker-compose.dev.yml'
fi


case "$ACTION" in
    start)
        dock_start
        ;;

    stop)
        dock_stop
        ;;

    build)
        dock_build
        ;;

    *)
        echo $"Usage: $0 -a {start|stop|build}"
        exit 1
esac
