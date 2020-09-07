default:
	go build

docker:
	docker-compose build
	docker-compose push