### SP Toy Key value store

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

## SP Toy key value store

Toy key value store key as string and value as text and has an expired at.
After key is expired, `get` operation won't get the value.
Keys are soft-deleted by setting the expired as the midnight
Unix epoch is TIMESTAMP '1970-01-01 00:00:00'.

## set key value expired_at

To set the `key`,`value` and `expired_at` use the below command or update the json at `set_method.sh`.

Below are sample curl

```shell
curl -X POST "http://localhost:8080/set" \
     -H "Content-Type: application/json" \
     -d '{"key":"sp", "value":"tree", "expiredAt":"2024-12-31T23:59:59Z"}'
```

Using the set method

```shell
./set_method.sh
```

## get value for specific key

To get the `value` and `expired_at` use below curl or update the key on the `get_method.sh`.

Below are sample curl

```shell
curl "http://localhost:8080/get?key=sp"
```

Using the get method

```shell
./get_method.sh
```

## delete value for specific key

To delete a `key`, use the below curl set the key in the `delete_method.sh`

```shell
curl "http://localhost:8080/delete?key=sp"
```

using the delete method

```shell
./delete_method.sh
```

![toy store](screenshots/sp_toy_key_value_store.jpg)

## Database Schema

### Table: `toy_dynamo`

| Column Name  | Data Type   | Constraints      | Description                                   |
| ------------ | ----------- | ---------------- | --------------------------------------------- |
| `id`         | `SERIAL`    | Primary Key      | Unique identifier for each record.            |
| `key`        | `VARCHAR`   | Not Null, Unique | The key associated with the value.            |
| `value`      | `TEXT`      | Not Null         | The data stored against the key.              |
| `expired_at` | `TIMESTAMP` | Not Null         | Timestamp indicating when the record expires. |

- **`id`**: This is an auto-incrementing integer used as the primary key for the table.
- **`key`**: A unique identifier for each record. This field is used to retrieve or modify the value.
- **`value`**: The content or information associated with the key.
- **`expired_at`**: The timestamp when the record is considered expired. This is used to manage the validity of the data. Setting this to `TIMESTAMP '1970-01-01 00:00:00'` indicates that the record has been logically deleted.

### Indexes

- **`idx_toy_dynamo_expired_at`**: An index on the `expired_at` column to improve performance of queries filtering by this field.
