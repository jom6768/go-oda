# Go, PostgreSQL, Docker, Minikube for TMForum ODA project

[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://paypal.me/jom6768)
[![Code Go](https://img.shields.io/badge/Code-Go-007F9F)](https://go.dev)
[![DB PostgreSQL](https://img.shields.io/badge/DB-PostgreSQL-336791)](https://www.postgresql.org)
[![Container Docker](https://img.shields.io/badge/Container-Docker-0DB7ED)](https://www.docker.com)
[![K8s Minikube](https://img.shields.io/badge/K8s-Minikube-306EE5)](https://minikube.sigs.k8s.io)

This repository contains a Go programs with Gin library that
implement TMForum Open Digital Architecture (ODA) standard
using PostgreSQL database and deploy on Docker and Minikube.

## Clone the project

```
$ git clone https://github.com/jom6768/go-oda
$ cd go-oda
```
 

## [Database configuration](oda/config)

A database configuration is in go-oda/oda/config/database_table.sql
 

## [TMF629](oda/tmf629) and/or [TMF632](oda/tmf632)

### Run locally:

* main.go in [TMF629](oda/tmf629) and/or [TMF632](oda/tmf632)
* Change database connection string in main.go at func initDB() to your connection

```
connStr := "postgresql://myuser:mypass@localhost:5432/go_oda?sslmode=disable"
```

* Run command at folder go-oda

```
$ go run ./oda/tmf629/main.go
```
 

### Run on Docker:

* main.go in [TMF629](oda/tmf629) and [TMF632](oda/tmf632)
* Change database connection string in main.go at func initDB() to your connection

```
connStr := "postgresql://myuser:mypass@host.docker.internal:5432/go_oda?sslmode=disable"
```

* docker-compose.yml
* Change image to your Docker user (make sure that <your_docker_user> is replaced by your Docker user)

```
image: <your_docker_user>/go-oda-tmf629:latest
```

* Run command at folder go-oda

```
$ docker-compose up --build -d
```

*** This will run both TMF629 on port 8629 and TMF632 on port 8632 ***
 

### Run on Minikube:

* Change database connection string in main.go at func initDB() to your connection

```
connStr := "postgresql://myuser:mypass@host.minikube.internal:5432/go_oda?sslmode=disable"
```

* deployment-tmf629.yaml and deployment-tmf632.yaml in [k8s](k8s)
* Change image to your Docker user (make sure that <your_docker_user> is replaced by your Docker user)

```
image: <your_docker_user>/go-oda-tmf629:latest
```

* Run command at folder go-oda (make sure that <your_docker_user> is replaced by your Docker user and <your_docker_key> is replaced by your Docker key)

```
$ docker build . --no-cache -t <your_docker_user>/go-oda-tmf629 -f ./oda/tmf629/Dockerfile
$ docker build . --no-cache -t <your_docker_user>/go-oda-tmf632 -f ./oda/tmf632/Dockerfile
$ echo "<your_docker_key>" | docker login -u <your_docker_user> --password-stdin
$ docker push <your_docker_user>/go-oda-tmf629
$ docker push <your_docker_user>/go-oda-tmf632

$ minikube start
$ kubectl apply -f k8s/deployment-tmf629.yaml
$ kubectl apply -f k8s/deployment-tmf632.yaml
$ kubectl wait --for=condition=ready pod -l app=tmf629 --timeout=30s
$ kubectl wait --for=condition=ready pod -l app=tmf632 --timeout=30s

$ nohup kubectl port-forward svc/tmf629 8629:8629 > tmf629.log 2>&1 &
$ echo $! > tmf629.pid
$ nohup kubectl port-forward svc/tmf632 8632:8632 > tmf632.log 2>&1 &
$ echo $! > tmf632.pid
```

*** This will run both TMF629 on port 8629 and TMF632 on port 8632 ***
 

## Testing

#### [Try to call webservice](oda/config)

A sample request is in go-oda/oda/config/sample_request.txt
 

#### [Run TMForum CTK Test Kit](https://www.tmforum.org/oda/open-apis/directory)

You can download CTK Test Kit from TMForum website and test.
* [TMF629](https://www.tmforum.org/oda/open-apis/directory/customer-management-api-TMF629/v5.0)
* [TMF632](https://www.tmforum.org/oda/open-apis/directory/party-management-api-TMF632/v5.0)
 

## Reference

* [Go](https://go.dev)
* [PostgreSQL](https://www.postgresql.org)
* [Docker](https://www.docker.com)
* [Minikube](https://minikube.sigs.k8s.io)

You can check Go learning in my youtube channel: [PJ Sweet Home](https://www.youtube.com/@PJSweetHome)

[![Go - A Tour of Go (Completed version with answer of all examples) [English & Thai Subtitle]](https://img.youtube.com/vi/Drrm470ta_Q/0.jpg)](https://www.youtube.com/watch?v=Drrm470ta_Q)
 

## Donation

:sparkling_heart::sparkling_heart::sparkling_heart: If this project help you to develop, you can give me a cup of coffee :sparkling_heart::sparkling_heart::sparkling_heart:

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://paypal.me/jom6768)
 

## License

Copyright (c) 2024 by Methee Suriyakrai. Some rights reserved.
