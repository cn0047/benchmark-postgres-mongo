PostgreSQL vs MongoDB
-

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

#### Prerequisites

Let's consider simple example,
suppose we have 2 entities (storage and file) with relationship `1 to many` (super simple case).
It will be 2 tables in `postgres`,
and in `mongo` will have document with embedded data.

Data will look like this:

Postgres tables examples:
````sql
# storage table example
 id |                   sha1                   | count
----+------------------------------------------+-------
 17 | dc9793d8f2379b73c4932d33a14c75eefa849fda |    64

# file table example
 id  | storage_id |         name
-----+------------+-----------------------
 113 |         17 | OdioCondimentumId.avi
 114 |         17 | Venenatis.png
````

Mongo document example:
````json
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

Before perform benchmarking let's import
**2500** entries as `storage` data and **24705** entries as `file` data into both databases.
<br>The target - is to have absolutely same data-sets in both databases!

#### Prepare

Please run next command:
````bash
docker network create --driver bridge xnet
````

#### PostgreSQL

Start postgres `docker` container:
````bash
docker run -it --rm -p 5432:5432 --net=xnet --name xpostgres --hostname xpostgres \
    -e POSTGRES_DB=test -e POSTGRES_USER=dbu -e POSTGRES_PASSWORD=dbp \
    postgres:10.5

# or (with data volume)
docker run -it --rm -p 5432:5432 --net=xnet --name xpostgres --hostname xpostgres \
    -v $PWD/.docker/data/postgresql/xpostgres:/var/lib/postgresql/data \
    -e POSTGRES_DB=test -e POSTGRES_USER=dbu -e POSTGRES_PASSWORD=dbp \
    postgres:10.5
````

Init postgres tables:
````bash
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
````

Import data into postgres
(I don't know why I've used `php` here 😀 but it works, and imports data into db 😉):
````bash
docker run -it --rm --net=xnet -v $PWD/data:/app -w /app cn007b/php php importDataIntoPostgres.php
````

Check data:
````sql
docker exec -ti xpostgres psql -h localhost -p 5432 -U dbu -d test -c '
    SELECT COUNT(*) FROM storage
    UNION ALL
    SELECT COUNT(*) FROM file
    UNION ALL
    SELECT COUNT(DISTINCT storage_id) FROM file
    UNION ALL
    SELECT COUNT(*) FROM storage s JOIN file f ON s.id = f.storage_id
'
````
Also, you can run any arbitrary query to check that databases contain same data.

Get `pq` for golang:
````bash
docker run -it --rm --net=xnet -v $PWD:/app -w /app -e GOPATH='/app' \
    cn007b/go:1.10 sh -c 'go get github.com/lib/pq'
````

**Run benchmarks:**

````bash
docker run -it --rm --net=xnet -v $PWD:/app -w /app -e GOPATH='/app' cn007b/go:1.10 sh -c '
    go run src/benchmark/main.go --v=v1 --db=postgres --q=query1
    go run src/benchmark/main.go --v=v1 --db=postgres --q=query2
    go run src/benchmark/main.go --v=v1 --db=postgres --q=query3
    go run src/benchmark/main.go --v=v1 --db=postgres --q=query4
'
docker run -it --rm --net=xnet -v $PWD:/app -w /app -e GOPATH='/app' cn007b/go:1.10 sh -c '
    go run src/benchmark/main.go --v=v2 --db=postgres --q=query1
    go run src/benchmark/main.go --v=v2 --db=postgres --q=query2
    go run src/benchmark/main.go --v=v2 --db=postgres --q=query3
    go run src/benchmark/main.go --v=v2 --db=postgres --q=query4
'
````

⚠️ Now stop postgres `docker` container.

#### MongoDB

Start mongo `docker` container:
````bash
docker run -it --rm -p 27017:27017 --net=xnet --name xmongo  --hostname xmongo \
    -v $PWD/:/tmp/data \
    mongo:4.0.1

# or (with data volume)
docker run -it --rm -p 27017:27017 --net=xnet --name xmongo  --hostname xmongo \
    -v $PWD/.docker/data/mongodb/xmongo:/data/db \
    -v $PWD/:/tmp/data \
    mongo:4.0.1
````

Import data into mongo:
````bash
docker exec -it xmongo mongo --port 27017 test --eval 'db.file_storage.drop()'

docker exec -it xmongo sh -c '
    for dumpFile in $( find /tmp/data/ -iname file_storage.* -type f ); do
        mongoimport --port 27017 -d test -c file_storage $dumpFile;
    done
'
````

Add mongo user:
````bash
docker exec -it xmongo mongo --port 27017 test \
    --eval 'db.createUser({user: "dbu", pwd: "dbp", roles: ["readWrite", "dbAdmin"]})'
````

Check data
````bash
docker exec -it xmongo mongo --port 27017 --eval '
    db.file_storage.count();
    db.file_storage.aggregate([
        {$project: {count: {$size: "$files"}}},
        {$group: {_id: null, 'total': {$sum: "$count"}}}
    ]);
'
````
Also, you can run any arbitrary query to check that databases contain same data.

Get `mgo` for golang:
````bash
docker run -it --rm --net=xnet -v $PWD:/app -w /app -e GOPATH='/app' \
    cn007b/go:1.10 sh -c 'go get gopkg.in/mgo.v2'
````

**Run benchmarks:**

````bash
docker run -it --rm --net=xnet -v $PWD:/app -w /app -e GOPATH='/app' cn007b/go:1.10 sh -c '
    go run src/benchmark/main.go --v=v1 --db=mongo --q=query1
    go run src/benchmark/main.go --v=v1 --db=mongo --q=query2
    go run src/benchmark/main.go --v=v1 --db=mongo --q=query3
    go run src/benchmark/main.go --v=v1 --db=mongo --q=query4
'
docker run -it --rm --net=xnet -v $PWD:/app -w /app -e GOPATH='/app' cn007b/go:1.10 sh -c '
    go run src/benchmark/main.go --v=v2 --db=mongo --q=query1
    go run src/benchmark/main.go --v=v2 --db=mongo --q=query2
    go run src/benchmark/main.go --v=v2 --db=mongo --q=query3
    go run src/benchmark/main.go --v=v2 --db=mongo --q=query4
'
````

⚠️ Now stop mongo `docker` container.

#### Result (2018-03-04)

**500** entries in both dbs.

````
+-------------+-------------------+--------------------+
| Benchmark # | PostgreSQL 10.0   | MongoDB 3.4.9      |
+-------------+-------------------+--------------------+
| query1      | 8349 microseconds | 21721 microseconds |
| query2      | 5243 microseconds | 12990 microseconds |
| query3      | 4442 microseconds | 11163 microseconds |
| query4      | 7570 microseconds | 13085 microseconds |
+-------------+-------------------+--------------------+
````

#### Result (2018-08-14)

**2500** entries in both dbs.

`v1` - steps: connect into db, get data, print data. 
`v2` - steps: get data, print data.

````
+-------------+--------------------+--------------------+
| Benchmark # | PostgreSQL 10.5    | MongoDB 4.0.1      |
+-------------+--------------------+--------------------+
| v1                                                    |
+-------------+--------------------+--------------------+
| v1.query1   | 26023 microseconds | 53008 microseconds |
| v1.query2   |  7658 microseconds | 14987 microseconds |
| v1.query3   |  4891 microseconds | 13856 microseconds |
| v1.query4   |  6385 microseconds | 16003 microseconds |
+-------------+--------------------+--------------------+
| v2                                                    |
+-------------+--------------------+--------------------+
| v2.query1   | 23289 microseconds | 42244 microseconds |
| v2.query2   |  5128 microseconds |  1629 microseconds |
| v2.query3   |  3996 microseconds |   522 microseconds |
| v2.query4   |  5040 microseconds |  1681 microseconds |
+-------------+--------------------+--------------------+
| v1 with data volume                                   |
+-------------+--------------------+--------------------+
| v1.query1   | 34137 microseconds | 52020 microseconds |
| v1.query2   | 20691 microseconds | 14950 microseconds |
| v1.query3   | 15952 microseconds | 14410 microseconds |
| v1.query4   | 18025 microseconds | 15389 microseconds |
+-------------+--------------------+--------------------+
| v2 with data volume                                   |
+-------------+--------------------+--------------------+
| v2.query1   | 33596 microseconds | 43352 microseconds |
| v2.query2   | 20072 microseconds | 1711 microseconds  |
| v2.query3   | 14750 microseconds | 573 microseconds   |
| v2.query4   | 17727 microseconds | 1678 microseconds  |
+-------------+--------------------+--------------------+
````

You can run all benchmarking queries on your computer and check all values.
