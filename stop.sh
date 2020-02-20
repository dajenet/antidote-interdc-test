#! /bin/bash

cd $(dirname $0)
cd dc2n1
docker-compose down
cd ../monitoring
docker-compose down
cd ..

docker stop pumba

docker stop pumba2
