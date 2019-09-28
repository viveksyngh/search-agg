CREATE ROLE searchdbuser WITH LOGIN PASSWORD '12345';

CREATE DATABASE searchdb OWNER searchdbuser;

CREATE TABLE searchquery (
  id              SERIAL PRIMARY KEY,
  status          VARCHAR(100) NOT NULL,
  query           VARCHAR(255) NOT NULL,
  created_on      TIMESTAMP NOT NULL 
);