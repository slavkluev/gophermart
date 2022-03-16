up:
	docker-compose up -d database

build:
	docker-compose up -d --build database

down:
	docker-compose down --remove-orphans

ps:
	docker-compose ps

stop:
	docker-compose stop

down-clear:
	docker-compose down -v --remove-orphans

pull:
	docker-compose pull