.PHONY: mysql redis crd api gc ga

domain := harbor.domain.com
project := lunara-common
mysql_image := "$(domain)/$(project)/mysql-slave:v1.0.0"
redis_image := "$(domain)/$(project)/redis-slave:v1.0.0"

mysql:
	cd dockerfile/mysql && docker build -t $(mysql_image) .
	docker push $(mysql_image)

redis:
	cd dockerfile/redis && docker build -t $(redis_image) .
	docker push $(redis_image)

crd:
	go mod vendor
	bash ./make.sh

api:
	go mod vendor
	cd api/cmd && bash ./make.sh

# gen crd
gc:
	go mod vendor
	bash ./gen.sh crd

# gen api
ga:
	go mod vendor
	bash ./gen.sh api

