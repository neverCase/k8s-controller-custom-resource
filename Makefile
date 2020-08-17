.PHONY: mysq redis

push := true
domain := harbor.domain.com
mysql_image := "$(domain)/helix_saga/mysql-slave:v1.0.0"
redis_image := "$(domain)/helix_saga/redis-slave:v1.0.0"

mysql:
	cd dockerfile/mysql && docker build -t $(mysql_image) .
	docker push $(mysql_image)

redis:
	cd dockerfile/redis && docker build -t $(redis_image) .
	docker push $(redis_image)