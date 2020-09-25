#!/bin/bash
: '
This is the first setup script of the server. You have to do some things before that:
- you have to add generate an SSH key and add it to the github account that has access to the projects
'

export DEBIAN_FRONTEND=noninteractive

# Setup docker reporsitory

apt-get update;
apt-get -y install \
	git \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common;
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -;
apt-key fingerprint 0EBFCD88;
add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable";

# Install docker engine

apt-get update;
apt-get -y install docker-ce docker-ce-cli containerd.io;

# Install nginx

deb http://nginx.org/packages/ubuntu/ xenial nginx
deb-src http://nginx.org/packages/ubuntu/ xenial nginx
apt-get update
apt-get -y install nginx

# Install certbot

add-apt-repository universe
add-apt-repository ppa:certbot/certbot -y
apt-get update
apt-get install certbot python-certbot-nginx -y

# Finish

echo "Start hacking and deploy projects with Boiler!"
