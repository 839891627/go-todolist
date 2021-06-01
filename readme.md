1. `docker run --name mongo -p 27017:27017 mongo`
1. `docker run -p 4222:4222 -ti nats:latest` 消息中间件
   > `export MICRO_BROKER=nats`  
   > `export MICRO_BROKER_ADDRESS=127.0.0.1:4222`
1. 因为使用了 etcd 来注册发现。所以需要填配置环境变量。两个方式都可以
   1. 环境变量方式(我们从用这个方式)
    - `export  MICRO_REGISTRY=etcd`
    - `export MICRO_REGISTRY_ADDRESS=127.0.0.1:2379`
   2.  命令参数方式
    - micro --registry=etcd --registry_address=127.0.0.1:2379 api --handler=http
1. task-srv: srv 服务
   `go run main.go`
    > 如果编译的话：
    > 1. protoc --proto_path=. --micro_out=${MODIFY}:. --go_out=${MODIFY}:. proto/task/task.proto
    > 2. protoc-go-inject-tag -input=proto/task/task.pb.go
    > 3. go build -o task-srv main.go

1. achievement-srv: srv 消费者（事件监听）  
    `go run main.go`
1. task-api: 提供 restful api 服务。端口号 8888
   `go run main.go`
1. 启动反向代理。类似 nginx。端口号 8080
   > 因为使用了环境变量使用 etcd 来发现
    
    `micro api --handler=http`
1. `micro web` 启用 web 管理界面工具
## etcd
```bash
rm -rf /tmp/etcd-data.tmp && mkdir -p /tmp/etcd-data.tmp && \
  docker run \
  -p 2379:2379 \
  -p 2380:2380 \
  --mount type=bind,source=/tmp/etcd-data.tmp,destination=/etcd-data \
  --name etcd-gcr-v3.4.13 \
  -d \
  quay.io/coreos/etcd:v3.4.13 \
  /usr/local/bin/etcd \
  --name s1 \
  --data-dir /etcd-data \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls http://0.0.0.0:2379 \
  --listen-peer-urls http://0.0.0.0:2380 \
  --initial-advertise-peer-urls http://0.0.0.0:2380 \
  --initial-cluster s1=http://0.0.0.0:2380 \
  --initial-cluster-token tkn \
  --initial-cluster-state new \
  --log-level info \
  --logger zap \
  --log-outputs stderr
```