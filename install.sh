#!/bin/bash

function social {
    echo "Installing social ..."

    apt-get install jq -y

    cd /etc/social

    LATEST_VERSION=$(curl --silent "https://api.github.com/repos/voocel/social/releases/latest" | jq '.tag_name' | sed -E 's/.*"([^"]+)".*/\1/' | tr -d v)

    curl -sL https://github.com/voocel/social/releases/download/v{$LATEST_VERSION}/social_{$LATEST_VERSION}_Linux_x86_64.tar.gz | tar xz

    echo "[Unit]
Description=Social
Documentation=https://github.com/voocel/social
[Service]
ExecStart=/etc/social/social server -c /etc/social/config.yml
Restart=on-failure
RestartSec=2
[Install]
WantedBy=multi-user.target" > /etc/systemd/system/social.service

    systemctl daemon-reload
    systemctl enable social.service
    systemctl start social.service

    echo "Social installation done!"
}

social