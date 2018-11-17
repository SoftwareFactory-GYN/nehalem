
# Setup Tutorial 

## Required software: 
- Docker


### 01. Start the cassandra cluster:
` docker run --name dev_cassandra -d -p 9042:9042 -p 7000:7000 -p 7001:7001 -p 7199:7199 -p 9160:9160 cassandra`

### Wait 30s for db to start 

### 02. Them start the cli session:
` docker run -it --link dev_cassandra:cassandra --rm cassandra cqlsh cassandra`

### 03. Run the following commands in the db cli session to create database and tables:
`CREATE KEYSPACE nehalem WITH replication = {'class':'SimpleStrategy', 'replication_factor' : 1};
 USE nehalem;
 CREATE TABLE users(
    id UUID PRIMARY KEY,
    username text,
    password text,
 );`
