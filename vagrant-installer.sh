#!/bin/bash

NOTIFY_PORT="${1}"

if [ -z $NOTIFY_PORT ]; then
  echo 'No server port was provided to the `notify-send` HTTP client!' 1>&2
  exit 1
fi

if [ -x /usr/bin/notify-send ]; then
  echo 'notify-send HTTP client already installed, skipping'
  exit 0
fi

curl -sL https://github.com/fgrehm/notify-send-http/releases/download/v0.2.0/client-linux_amd64 | sudo tee /usr/bin/notify-send &>/dev/null
sudo chmod +x /usr/bin/notify-send

cat <<-STR >> /home/vagrant/.bashrc
SERVER_IP=\$(ip route|awk '/default/ { print \$3 }')
export NOTIFY_SEND_URL="http://\${SERVER_IP}:${NOTIFY_PORT}"
STR
