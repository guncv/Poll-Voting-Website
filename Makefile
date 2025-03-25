.PHONY: run docker_build down clean build-proto clean-proto

info:
	docker-compose ps

run:
	docker-compose up

down:
	docker-compose down

build:
	docker-compose build

clean:
	docker-compose down --rmi all --volumes --remove-orphans

logs:
	docker-compose logs -f

restart:
	docker-compose restart

ps:
	docker-compose ps

rebuild: clean docker_build run
