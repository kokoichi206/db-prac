# db-prac

## setup
``` sh
# https://hub.docker.com/_/postgres/
docker run -it --rm --network some-network postgres psql -h some-postgres -U postgres
```

``` sh
docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=test -d postgres:12-alpine
docker exec -it postgres12 createdb --username=postgres --owner=postgres test_db
docker exec -it postgres12 psql -h localhost -p 5432 -U postgres
docker exec -it postgres12 dropdb test_db

# persist data (relative path)
docker run --name postgres12 -v $PWD/data:/var/lib/postgresql/data -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=test -d postgres:12-alpine
```

## postgreSQL
``` postgresql
\l
\c test_db
```
