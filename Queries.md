Queries
-

````
# postgres 1
EXPLAIN
SELECT s.id, s.sha1, f.name
FROM storage s
JOIN file f ON s.id = f.storage_id
WHERE s.count > 0
ORDER by id DESC, name ASC
OFFSET 1000 LIMIT 10
;

# postgres 2
EXPLAIN
SELECT s.id, s.sha1, f.name
FROM storage s
JOIN file f ON s.id = f.storage_id
WHERE s.sha1 = '806b9a087e6822c1548c606e8e6348b7f08b62ff' AND f.name = 'Sit.avi'
;

# postgres 3
EXPLAIN
SELECT s.id, s.sha1
FROM storage s
WHERE s.id IN (171, 352)
;

# mongo 1
db.file_storage.explain().aggregate([
    {$match: {"count": {$gt: 0}}},
    {$unwind: "$files"},
    {$project: {_id: 1, sha1: 1, "name": "$files.name"}},
    {$sort: {_id: -1, "name": 1}},
    {$skip : 1000},
    {$limit : 10},
]);

# mongo 2
db.file_storage.explain().find(
    {sha1: "806b9a087e6822c1548c606e8e6348b7f08b62ff"},
    {_id: 1, sha1: 1, "files": {"$elemMatch": {"name": "Sit.avi"}}}
);

# mongo 3
db.file_storage.explain().find(
    {_id: {$in: [171, 352]}},
    {_id: 1, sha1: 1}
);
````
