#!/bin/sh

IMAGE='terraform-glesys:latest'

#docker pull ${IMAGE}

if ! docker inspect --type=image ${IMAGE} 1> /dev/null; then
  echo "No docker env, please build https://github.com/glesys/terraform-provider-glesys"
  (cd .. && docker build . -t tfenv:0.3.2 -f docker/Dockerfile)
fi


echo ""
read -p "Username:" username

echo ""
read -s -p "Secret: " password

export GLESYS_USERID="${username}"
export GLESYS_TOKEN="${password}"

docker run -it --rm  -e GLESYS_USERID -e GLESYS_TOKEN \
  -v ${PWD}:/home \
  ${IMAGE}
