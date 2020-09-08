.PHONY: mysql redis crd api gc ga

HARBOR_DOMAIN := $(shell echo ${HARBOR})
PROJECT := lunara-common
MYSQL_IMAGE := "$(HARBOR_DOMAIN)/$(PROJECT)/mysql-slave:v1.0.0"
REDIS_IMAGE := "$(HARBOR_DOMAIN)/$(PROJECT)/redis-slave:v1.0.0"
API_SERVER_IMAGE := "$(HARBOR_DOMAIN)/$(PROJECT)/api-server:latest"
MULTIPLEX_CRD_IMAGE := "$(HARBOR_DOMAIN)/$(PROJECT)/multiplex-crd:latest"

mysql:
	-i docker image rm $(MYSQL_IMAGE)
	cd dockerfile/mysql && docker build -t $(MYSQL_IMAGE) .
	docker push $(MYSQL_IMAGE)

redis:
	-i docker image rm $(REDIS_IMAGE)
	cd dockerfile/redis && docker build -t $(REDIS_IMAGE) .
	docker push $(REDIS_IMAGE)

crd:
	-i docker image rm $(MULTIPLEX_CRD_IMAGE)
	cd scripts && bash ./make.sh multiplexcrd-binary
	cd cmd/multiplex && docker build -t $(MULTIPLEX_CRD_IMAGE) .
	rm -rf cmd/multiplex/multiplex-crd
	docker push $(MULTIPLEX_CRD_IMAGE)

api:
	-i docker image rm $(API_SERVER_IMAGE)
	cd scripts && bash ./make.sh api-binary
	cd api && docker build -t $(API_SERVER_IMAGE) .
	rm -rf api/api-server
	docker push $(API_SERVER_IMAGE)

# gen crd
gc:
	go mod vendor
	CRD=mysqloperator bash ./gen.sh crd
	CRD=redisoperator bash ./gen.sh crd

# gen api
ga:
	go mod vendor
	bash ./gen.sh api

