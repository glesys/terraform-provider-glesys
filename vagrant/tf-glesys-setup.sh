#!/bin/bash

set -e

code_dir="$HOME/code/"
project="terraform-provider-glesys"
tf_version="0.12.29"
tf_url="https://releases.hashicorp.com/terraform/${tf_version}/terraform_${tf_version}_linux_amd64.zip"
repo_url="https://github.com/norrland/terraform-provider-glesys.git"


echo "Install dependencies"

sudo apt-get update && sudo apt-get install -y \
    curl \
    git \
    golang \
    unzip \
    make

echo "Setup directories for terraform-provider-glesys"

mkdir -p ~/code

echo "Fetch code"

if [ ! -d "${code_dir}/${project}/.git" ]; then
    git clone ${repo_url} "${code_dir}/${project}"
else
    echo "Repo already exists. Trying to fetch new code"
    cd "${code_dir}/${project}"
    git pull
fi

echo "Build provider"

cd ~/code/terraform-provider-glesys
make

if [ -f ~/go/bin/terraform-provider-glesys ]; then
    echo "Setting up provider"
    mkdir -p ~/.terraform.d/plugins
    if [ ! -f ~/.terraform.d/plugins/${project} ]; then
        ln -s ~/go/bin/terraform-provider-glesys ~/.terraform.d/plugins/
    fi
fi

echo "Installing Terraform"
if [ -f "$HOME/bin/terraform" ]; then
    echo "Terraform already installed"
else
    mkdir -p "$HOME/bin"
    curl -s ${tf_url} -o tf_binary.zip
    unzip tf_binary.zip -d "$HOME/bin"
fi
