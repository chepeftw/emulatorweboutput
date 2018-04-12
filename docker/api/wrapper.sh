#!/usr/bin/env sh

mkdir -p /var/log/golang
echo "START" > /var/log/golang/wrapper.log

echo "starting ... "
echo "---------------------------------------------"
echo "Starting API ... " >> /var/log/golang/wrapper.log

/root/myapp &
status=$?
if [ $status -ne 0 ]; then
  echo "Failed to start router process: $status"
  exit $status
fi

# Naive check runs checks once a minute to see if either of the processes exited.
# This illustrates part of the heavy lifting you need to do if you want to run
# more than one service in a container. The container will exit with an error
# if it detects that either of the processes has exited.
# Otherwise it will loop forever, waking up every 60 seconds

sleep 2

while /bin/true; do

  ps aux | grep myapp | grep -v grep
  P1_STATUS=$?

  echo "PROCESS1 STATUS = $P1_STATUS |"

  echo "PROCESS1 STATUS = $P1_STATUS " >> /var/log/golang/wrapper.log

  # If the greps above find anything, they will exit with 0 status
  # If they are not both 0, then something is wrong
  if [ $P1_STATUS -ne 0 ]; then
    echo "API has already exited."
    exit -1
  fi

  sleep 60
done