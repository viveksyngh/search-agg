CREATE ROLE searchdbuser WITH LOGIN PASSWORD '12345';

CREATE DATABASE searchdb OWNER searchdbuser;

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