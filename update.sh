#!/usr/bin/env bash

go build > /dev/null 2>&1
echo "Commit name"
# shellcheck disable=SC2162
read message
echo "Sending update"
git add .
git commit -am "$message"
git push
ssh -i ssh.key root@129.213.112.128 "systemctl stop classroom.service; runuser -l api -c 'cd /home/api/Student-Checkout; git reset --hard; git pull; chmod +x checkout'; systemctl start classroom.service"
echo "Sent update"