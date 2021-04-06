#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
docker build -t registry.cn-shenzhen.aliyuncs.com/satsun/china_labor_union:goods .
docker push registry.cn-shenzhen.aliyuncs.com/satsun/china_labor_union:goods