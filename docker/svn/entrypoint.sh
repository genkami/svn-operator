#!/bin/bash

set -e

: "${APACHE_CONFDIR:=/etc/apache2}"
: "${APACHE_PID_FILE:=${APACHE_RUN_DIR:=/var/run/apache2}/apache2.pid}"

source /etc/apache2/envvars

export SERVER_NAME="test.example.com"

mkdir -p /svn
chown -R www-data:www-data /svn

sudo -u www-data -g www-data mkdir -p /svn/repos
sudo -u www-data -g www-data /work/server-updater &

exec apache2 -DFOREGROUND "$@"
