.PHONY: mysql redis crd api gc ga

HARBOR_DOMAIN := $(shell echo ${HARBOR})
PROJECT := lunara-common
MYSQL_IMAGE := "$(HARBOR_DOMAIN)/$(PROJECT)/mysql-slave:v1.0.0"
REDIS_IMAGE := "$(HARBOR_DOMAIN)/$(PROJECT)/redis-slave:v1.0.0"
API_SERVER_IMAGE := "$(HARBOR_DOMAIN)/$(PROJECT)/api-server:latest"
MULTIPLE_CRD_IMAGE := "$(HARBOR_DOMAIN)/$(PROJECT)/mupliple-crd:latest"

mysql:
	cd dockerfile/mysql && docker build -t $(MYSQL_IMAGE) .
	docker push $(MYSQL_IMAGE)

redis:
	cd dockerfile/redis && docker build -t $(REDIS_IMAGE) .
	docker push $(REDIS_IMAGE)

crd:
	go mod vendor
	bash ./make.sh

api:
	cd scripts && bash ./make.sh
	cd api && docker build -t $(API_SERVER_IMAGE) .
	rm -rf api/api-server
	docker push $(API_SERVER_IMAGE)

# gen crd
gc:
	go mod vendor
	bash ./gen.sh crd

# gen api
ga:
	go mod vendor
	bash ./gen.sh api

