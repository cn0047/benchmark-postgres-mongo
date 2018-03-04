<?php

$dbh = new PDO('pgsql:host=xpostgres;port=5432;dbname=test;user=dbu;password=dbp');

$f = fopen(__DIR__ . '/file_storage.json', 'r');
while (($line = fgets($f)) !== false) {
    $d = json_decode($line, true);
    if (empty($d)) {
        continue;
    }
    printf("\rInported id: %s", $d['_id']);

    $s = $dbh->prepare('INSERT INTO storage (id, sha1, count) VALUES (:s1, :s2, :s3)');
    $s->bindValue(':s1', $d['_id'], PDO::PARAM_INT);
    $s->bindValue(':s2', $d['sha1'], PDO::PARAM_STR);
    $s->bindValue(':s3', $d['count'], PDO::PARAM_INT);
    $s->execute();

    foreach ($d['files'] as $el) {
        $s = $dbh->prepare('INSERT INTO file (id, storage_id, name) VALUES (default, :s2, :s3)');
        $s->bindValue(':s2', $d['_id'], PDO::PARAM_INT);
        $s->bindValue(':s3', $el['name'], PDO::PARAM_STR);
        $s->execute();
    }
}
