all: build docker

build:
	go build -tags debug -o ./server/main -gcflags "all=-N -l" ./server/main.go

docker:
	docker build -t 884647653201.dkr.ecr.ap-southeast-1.amazonaws.com/lp_oms:${DOCKER_TAG} .
	aws ecr get-login-password --region ap-southeast-1 | docker login --username AWS --password-stdin 884647653201.dkr.ecr.ap-southeast-1.amazonaws.com
	docker push 884647653201.dkr.ecr.ap-southeast-1.amazonaws.com/lp_oms:${DOCKER_TAG}	
	docker run --name order_management -p 5003:5003 -d -e versafleet_pushOrders="http://a21-dev.integrum.global:8000/push_results" -e versafleet_pullOrders="http://a21-dev.integrum.global:8000/pull_orders" -e versafleet_unassignTasks="http://a21-dev.integrum.global:8000/unassign_tasks" 884647653201.dkr.ecr.ap-southeast-1.amazonaws.com/lp_oms:${DOCKER_TAG}

clean:
	rm -rf ./server/main
	docker rmi -f lp_oms:${DOCKER_TAG}
	
push:
	docker tag lp_oms:${DOCKER_TAG} 13.229.233.143:5003/lp_oms:${DOCKER_TAG} 
#	docker push 13.229.233.143:5001/lp_ums:${DOCKER_TAG}
