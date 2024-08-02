#!/bin/bash

cat <<EOF > /etc/systemd/system/sp.service
[Unit]
Description=48Club Soul Point Service
After=network.target

[Service]
ExecStart=$(go env GOPATH)/bin/sp
Restart=always
LimitNOFILE=infinity
WorkingDirectory=${PWD}

[Install]
WantedBy=multi-user.target

EOF

systemctl daemon-reload

systemctl enable sp --now