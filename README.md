PostgreSQL vs MongoDB
-

#### Prerequisites

````sql
# Postgre `storage` table example
 id |                   sha1                   | count
----+------------------------------------------+-------
 17 | dc9793d8f2379b73c4932d33a14c75eefa849fda |    64

# Postgre `file` table example
 id  | storage_id |         name
-----+------------+-----------------------
 113 |         17 | OdioCondimentumId.avi
 114 |         17 | Venenatis.png

# Mongo document example
{
  "_id": 17,
  "sha1": "dc9793d8f2379b73c4932d33a14c75eefa849fda",
  "count": 64,
  "files": [
    {
      "id": 1,
      "name": "OdioCondimentumId.avi"
    },
    {
      "id": 2,
      "name": "Venenatis.png"
    }
  ]
}
````

#### Prepare

````
docker network create --driver bridge xnet
````

#### PostgreSQL

````sql
# Start postgres
docker run -it --rm -p 5432:5432 --net=xnet --name xpostgres --hostname xpostgres \
    -e POSTGRES_DB=test -e POSTGRES_USER=dbu -e POSTGRES_PASSWORD=dbp \
    postgres:10.0

# Init postgres tables
docker exec -ti xpostgres psql -h localhost -p 5432 -U dbu -d test -c '
    DROP TABLE IF EXISTS storage;
    CREATE TABLE storage (
        id SERIAL PRIMARY KEY,
        sha1 CHARACTER VARYING(40),
        count INTEGER
    );
    DROP TABLE IF EXISTS file;
    CREATE TABLE file (
        id SERIAL PRIMARY KEY,
        storage_id INTEGER,
        name CHARACTER VARYING(50)
    );
    CREATE INDEX s_id ON file (storage_id);
'

# Import data into postgres
docker run -it --rm --net=xnet -v $PWD/data:/app -w /app cn007b/php php importDataIntoPostgres.php
# I don't know why I've used php here 😀 but it works, and import data into db 😉

# Check data
docker exec -ti xpostgres psql -h localhost -p 5432 -U dbu -d test -c '
    SELECT COUNT(*) FROM storage
    UNION ALL
    SELECT COUNT(*) FROM file
    UNION ALL
    SELECT COUNT(DISTINCT storage_id) FROM file
    UNION ALL
    SELECT COUNT(*) FROM storage s JOIN file f ON s.id = f.storage_id
'
# Also, you can run any arbitrary query
# to check that databases contain same data.

# Get pq for golang
docker run -it --rm --net=xnet -v $PWD:/app -w /app -e GOPATH='/app' \
    golang:latest sh -c 'go get github.com/lib/pq'

# Run benchmarks
docker run -it --rm --net=xnet -v $PWD:/app -w /app -e GOPATH='/app' \
    golang:latest go run src/benchmark/postgres/query1.go
docker run -it --rm --net=xnet -v $PWD:/app -w /app -e GOPATH='/app' \
    golang:latest go run src/benchmark/postgres/query2.go
docker run -it --rm --net=xnet -v $PWD:/app -w /app -e GOPATH='/app' \
    golang:latest go run src/benchmark/postgres/query3.go
docker run -it --rm --net=xnet -v $PWD:/app -w /app -e GOPATH='/app' \
    golang:latest go run src/benchmark/postgres/query4.go

# ⚠️ Now stop postgres.
````

#### MongoDB

````sql
# Start mongo
docker run -it --rm -p 27017:27017 --net=xnet --name xmongo  --hostname xmongo \
    -v $PWD/data/file_storage.json:/tmp/d.json \
    mongo:3.4.9

# Import data into mongo
docker exec -it xmongo mongoimport --port 27017 --drop -d test -c file_storage /tmp/d.json

# Add mongo user
docker exec -it xmongo mongo --port 27017 test \
    --eval 'db.createUser({user: "dbu", pwd: "dbp", roles: ["readWrite", "dbAdmin"]})'

# Check data
docker exec -it xmongo mongo --port 27017 --eval '
    db.file_storage.count();
    db.file_storage.aggregate([
        {$project: {count: {$size: "$files"}}},
        {$group: {_id: null, 'total': {$sum: "$count"}}}
    ]);
'
# Also, you can run any arbitrary query
# to check that databases contain same data.

# Get mgo for golang
docker run -it --rm --net=xnet -v $PWD:/app -w /app -e GOPATH='/app' \
    golang:latest sh -c 'go get gopkg.in/mgo.v2'

# Run benchmarks
docker run -it --rm --net=xnet -v $PWD:/app -w /app -e GOPATH='/app' \
    golang:latest go run src/benchmark/mongo/query1.go
docker run -it --rm --net=xnet -v $PWD:/app -w /app -e GOPATH='/app' \
    golang:latest go run src/benchmark/mongo/query2.go
docker run -it --rm --net=xnet -v $PWD:/app -w /app -e GOPATH='/app' \
    golang:latest go run src/benchmark/mongo/query3.go
docker run -it --rm --net=xnet -v $PWD:/app -w /app -e GOPATH='/app' \
    golang:latest go run src/benchmark/mongo/query4.go

# ⚠️ Now stop mongo.
````

#### Result

````
+-------------+-------------------+--------------------+
| Benchmark # | PostgreSQL        | MongoDB            |
+-------------+-------------------+--------------------+
| query1      | 8349 microseconds | 21721 microseconds |
| query2      | 5243 microseconds | 12990 microseconds |
| query3      | 4442 microseconds | 11163 microseconds |
| query4      | 7570 microseconds | 13085 microseconds |
+-------------+-------------------+--------------------+
````

You can run all benchmarking queries on your computer and check all values.
