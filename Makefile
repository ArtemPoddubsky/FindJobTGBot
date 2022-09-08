all:
	@mkdir -p /var/tmp/docker/postgresql
	@docker compose --env-file ./.env build
	@docker compose --env-file ./.env up -d

app:
	@docker compose --env-file ./.env up -d bot

db:
	@docker compose --env-file ./.env up -d postgresql

start:
	@docker compose --env-file ./.env up -d

stop:
	@docker compose stop

clean: stop
	docker compose down
	@-docker volume rm $$(docker volume ls -q)
	@-docker rmi $$(docker images -q)

fclean: clean
	@-rm -rf /var/tmp/docker

re: clean all