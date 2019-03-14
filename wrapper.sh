#!/bin/bash

#Kill previous versions of this proxy. Remove this to run multiple copies at the same time
pkill emp-proxy

#Credentials provided by EMP specific to an endpoint
export EMP_USERNAME=""
export EMP_PASSWORD=""

#EMP Endpoint to proxy data to
#Do not include the 'channel' at the end of the provided URL, but do include a trailing slash
export EMP_PATH=""

#The port to run the proxy server on localhost.
export PROXY_PORT="8080"

LOG=emp-proxy.log

#Start the proxy. Redirect stderr and stdout to the log
./emp-proxy > $LOG 2>&1 &

#Let the process run in the background
disown
