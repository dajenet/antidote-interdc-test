#! /bin/bash

go build -o main .
chmod +x main

cd $(dirname $0)
cd dc2n1
docker-compose up -d
cd ../monitoring
docker-compose up -d
cd ..
sleep 30


docker run --rm -d --name pumba  -v /var/run/docker.sock:/var/run/docker.sock gaiaadm/pumba netem --duration 1h --target 172.28.0.3 loss --percent 50 dc1n1
docker run --rm -d --name pumba2 -v /var/run/docker.sock:/var/run/docker.sock gaiaadm/pumba netem --duration 1h --target 172.28.0.2 loss --percent 50 dc2n1

main
