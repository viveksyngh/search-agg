# search-Aggregator
Application written in golang which aggregates search results from Google, DuckDuckGo and Wikipedia

## Steps to run

### Prerequisite
You need to have `docker` and `docker-compose` installed on your system

1. [Installation on Mac](https://docs.docker.com/docker-for-mac/install/)
    
2. [Install Docker Compose ](https://docs.docker.com/compose/install/) 


### Get the code
Clone this repository and change to repository directory

### Bring up all Services
Run below command to bring up all services
```
docker-compose build --parallel && docker-compose up
```

### Setup Database
Once all services are up, open a terminal window and run below command to exec to database container
```
docker exec -it $(docker ps --filter=name=search-agg_db -q) bash
```

Login to database inside the container using below command
```
psql "sslmode=disable host=localhost port=5432 user=searchdbuser dbname=searchdb password=12345"
```

To create the required tables for the app, run below create tables commands which are also present in `db.sql` file

```sql
CREATE TABLE searchquery (
  id              SERIAL PRIMARY KEY,
  status          VARCHAR(100) NOT NULL,
  query           VARCHAR(255) NOT NULL,
  created_on      TIMESTAMP NOT NULL 
);

CREATE TABLE googlesearchresult (
    id              SERIAL PRIMARY KEY,
    title           VARCHAR (255) NOT NULL,
    searchquery_id  INTEGER REFERENCES searchquery(id),
    url             VARCHAR (255) NOT NULL         
);

CREATE TABLE duckduckgosearchresult (
    id              SERIAL PRIMARY KEY,
    title           VARCHAR (255) NOT NULL,
    searchquery_id  INTEGER REFERENCES searchquery(id),
    url             VARCHAR (255) NOT NULL            
);

CREATE TABLE wikipediasearchresult (
    id              SERIAL PRIMARY KEY,
    result           TEXT NOT NULL,
    searchquery_id  INTEGER REFERENCES searchquery(id)        
);
```

### Open the application

Go to browser and access the application at `http://127.0.0.1:8000`
