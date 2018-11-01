all: dev

install:
	cd api && make install
	cd app && make install
	cd manager && make install

dev:
	make start

start:
	make start-gateway
	make start-services -j

start-gateway:
	cd gateway && make

start-services: start-api start-app start-asset start-manager start-docs

start-api:
	cd api && make

start-app:
	cd app && make

start-asset:
	cd asset && make

start-manager:
	cd manager && make

start-docs:
	cd docs && make

stop:
	make stop-services -j
	make stop-gateway

stop-services: stop-api stop-app stop-asset stop-manager stop-docs

stop-api:
	cd api && make stop

stop-app:
	cd app && make stop

stop-asset:
	cd asset && make stop

stop-manager:
	cd manager && make stop

stop-docs:
	cd docs && make stop

stop-gateway:
	cd gateway && make stop

restart:
	make stop
	make start
