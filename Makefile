.SILENT:
include .env # provides PORT variable

name = makoto-go-websocket

docker-clear:
	sudo docker rm -f ${name}

docker-run: docker-clear
	@echo 'Start container "${name}" on port ${PORT}'
	sudo docker run --name=${name} -d -p ${PORT}:${PORT} ${name}

docker-stop:
	@echo 'Stop container "${name}"';
	sudo docker stop ${name}

docker-build:
	@echo 'Build container "${name}" with port ${PORT}';
	sudo docker build -t ${name} .
