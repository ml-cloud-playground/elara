BASEDIR = $(shell pwd)
PORT_FE=8080
PORT_DB=3306
PORT_API=9000

dev: db redis api fe

fe: cleanfe
	cd frontend && docker build -t elara-html .
	docker run --name elara-html --expose $(PORT_FE) -p $(PORT_FE):80 \
	-v $(BASEDIR)/frontend/www:/usr/share/nginx/html elara-html

db: cleandb
	cd database && docker build -t elara-mysql .
	docker run --name elara-mysql -p $(PORT_DB):$(PORT_DB) \
	-e MYSQL_ROOT_PASSWORD=password -e MYSQL_ROOT_HOST=% -d elara-mysql

api: cleanapi
	cd middleware && docker build -t elara-goapi .
	docker run --name elara-goapi --expose $(PORT_API) \
	-p $(PORT_API):$(PORT_API)  -e PORT=$(PORT_API) -e elara_user=root \
	-e elara_pass=password -e elara_host=host.docker.internal -e elara_name=elara  \
	-e REDISPORT=6379 -e REDISHOST=host.docker.internal -d elara-goapi	

cleanfe:
	-docker stop elara-html
	-docker rm elara-html

cleandb:
	-docker stop elara-mysql
	-docker rm elara-mysql

cleanapi:
	-docker stop elara-goapi
	-docker rm elara-goapi		


redis: cleanredis
	docker run --name elara-redis -p 6379:6379 -d redis	

cleanredis:
	-docker stop elara-redis
	-docker rm elara-redis	