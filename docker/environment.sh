#!/bin/sh

IMAGE='tfenv:0.3.2'

docker pull ${IMAGE}

if ! docker inspect --type=image ${IMAGE} 1> /dev/null; then
  echo "No docker env, please build https://github.com/glesys/terraform-provider-glesys"
  (cd .. && docker build . -t tfenv:0.3.2 -f docker/Dockerfile)
fi


echo ""
read -p "Username:" username

echo ""
read -s -p "Secret: " password

TOKEN=$(curl -sb -X POST -H "Accept: application/json" --data-urlencode "username=${username}" --data-urlencode "password=${password}" https://api.glesys.com/user/login/ |jq -r '.response.login.apikey')

if [ "${TOKEN}" = 'null' ]
then
      echo "Unable to fetch Glesys token"
      exit 1
else
      echo "export GLESYS_TOKEN"
fi

export GLESYS_TOKEN="$TOKEN"

docker run -it --rm -e GLESYS_TOKEN \
  -v ${PWD}:/home \
  ${IMAGE}
