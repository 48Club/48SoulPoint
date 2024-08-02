#!/bin/bash

cat <<EOF > /etc/systemd/system/sp.service
[Unit]
Description=48Club Soul Point Service
After=network.target

[Service]
ExecStart=$(go env GOPATH)/bin/sp
Restart=always
LimitNOFILE=infinity
WorkingDirectory=${HOME}/.config/sp/

[Install]
WantedBy=multi-user.target

EOF

systemctl daemon-reload

mkdir ${HOME}/.config/sp/
cp -r config.json ${HOME}/.config/sp/

systemctl enable sp --now