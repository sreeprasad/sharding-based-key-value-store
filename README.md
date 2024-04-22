### Sharding based SP Toy Key value store

## How to run

Give permission to postgres-init to execute script to create users

```shell
chmod +x ./postgres-init/*.sh
chmod +x ./start_method.sh
chmod +x ./get_method.sh
chmod +x ./delete_method.sh
```

start the docker to run the the postgres database and pgbouncer

```shell
docker-compose up --build
```

start the application

```shell
go run application.go
```

## Sharding based SP Toy key value store

This is same as https://github.com/sreeprasad/toy-key-value-store
except that now keys are stored in 2 database shards based on
hash of key mod 2

![sharded toy store](screenshots/sharded_toy_store.mp4)
