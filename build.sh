export DOCKER_NAME=lpoms
export DOCKER_TAG=${1}

echo "Building docker image " ${DOCKER_NAME} ${DOCKER_TAG}
make 

