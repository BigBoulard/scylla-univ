# Preparing the Scylla DB

## Context

<https://university.scylladb.com/courses/using-scylla-drivers/lessons/golang-and-scylla-part-1/>
<https://university.scylladb.com/courses/using-scylla-drivers/lessons/golang-and-scylla-part-2-data-types/>
<https://university.scylladb.com/courses/using-scylla-drivers/lessons/golang-and-scylla-part-3-gocqlx/>

## Building Scylla containers

```sh
docker-compose up -d
[+] Running 3/3
 ✔ Container scylla-node1  Started                                                                                                                             0.1s
 ✔ Container scylla-node2  Started                                                                                                                             0.1s
 ✔ Container scylla-node3  Started                                                                                                                             0.0s
```

## Checking node status

```sh
$ docker exec -it scylla-node1 nodetool status
Datacenter: DC1
===============
Status=Up/Down
|/ State=Normal/Leaving/Joining/Moving
--  Address     Load       Tokens       Owns    Host ID                               Rack
UN  172.29.0.2  ?          256          ?       7ca9462a-819b-43f6-85cf-0a6857107377  Rack1

Note: Non-system keyspaces don't have the same replication settings, effective ownership information is meaningless
```

## Adding data into to Scylla

```sh
$ docker exec -it scylla-node1 cqlsh
Connected to  at 172.29.0.3:9042.
[cqlsh 5.0.1 | Cassandra 3.0.8 | CQL spec 3.3.1 | Native protocol v4]
Use HELP for help.
cqlsh> CREATE KEYSPACE catalog WITH REPLICATION = { 'class' : 'NetworkTopologyStrategy','DC1' : 3};
cqlsh> use catalog;
cqlsh:catalog>
cqlsh:catalog> CREATE TABLE mutant_data (
           ...    first_name text,
           ...    last_name text,
           ...    address text,
           ...    picture_location text,
           ...    PRIMARY KEY((first_name, last_name)));
cqlsh:catalog> insert into mutant_data ("first_name","last_name","address","picture_location") VALUES ('Bob','Loblaw','1313 Mockingbird Lane', 'http://www.facebook.com/bobloblaw');
cqlsh:catalog> insert into mutant_data ("first_name","last_name","address","picture_location") VALUES ('Bob','Zemuda','1202 Coffman Lane', 'http://www.facebook.com/bzemuda');
cqlsh:catalog> insert into mutant_data ("first_name","last_name","address","picture_location") VALUES ('Jim','Jeffries','1211 Hollywood Lane', 'http://www.facebook.com/jeffries');
cqlsh:catalog> exit
```
