.PHONY: mysql redis crd api

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
	bash ./make.sh

api:
	cd api/cmd && bash ./make.sh
