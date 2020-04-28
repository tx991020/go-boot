#!/usr/bin/env bash
set -e
set -x
go run main.go new -t mysql -H 118.24.156.209 -P 3306 -u root -p 123456 -n ladder  -d gm1 --http_port 12000



#!/usr/bin/env bash



#

#CUR_DIR="$( cd "$( dirname "$0"  )" && pwd  )"
#
#CMD="WebGenerator.exe"
#if [ $(uname -s) = 'Linux' ]; then
#	CMD="WebGenerator"
#elif [ $(uname -s) = 'Darwin' ]; then
#    CMD="WebGenerator_mac"
#fi
#
#$CUR_DIR/WebGenerator_mac new -t mysql -H 118.24.156.209 -P 3306 -u root -p 123456 -n gm-go  -d gm1 --http_port 12000