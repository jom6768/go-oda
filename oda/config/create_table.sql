-- PostgreSQL command --
\l; -- List all databases
\c [database_name]; -- Connect database
\du; -- List all users
\du+; -- List all users with detailed information
\dt; -- List all tables
\dt+; -- List all tables with detailed information
exit -- Close connection

-- Connect database
% "/Applications/Postgres.app/Contents/Versions/16/bin/psql" -U myuser -d go_oda

SELECT current_user, session_user; -- Check user
SET ROLE [username]; -- Change user login

-- Create database
CREATE DATABASE go_oda;
CREATE USER myuser WITH ENCRYPTED PASSWORD 'mypass';
GRANT ALL PRIVILEGES ON DATABASE go_oda TO myuser;
GRANT ALL ON SCHEMA public TO myuser;

-- TMF632_Party Management API
CREATE TABLE customer (
    id VARCHAR(50) PRIMARY KEY NOT NULL,
    gender VARCHAR(20),
    countryOfBirth VARCHAR(50),
    nationality VARCHAR(50),
    maritalStatus VARCHAR(20),
    birthDate TIMESTAMP,
    givenName VARCHAR(50),
    preferredGivenName VARCHAR(50),
    familyName VARCHAR(50),
    legalName VARCHAR(50), -- UNIQUE
    middleName VARCHAR(50),
    fullName VARCHAR(100),
    formattedName VARCHAR(100),
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO customer (id, gender, countryOfBirth, nationality, maritalStatus, birthDate, givenName, preferredGivenName, familyName, legalName, middleName, fullName, formattedName, status)
VALUES ('42', 'female', 'United States', 'American', 'married', '1967-09-26T05:00:00.246Z', 'Jane', 'Lamborgizzia', 'Lamborgizzia', 'Smith', 'JL', 'Jane Smith ep Lamborgizzia', 'Jane Smith ep Lamborgizzia', 'validated') RETURNING id;

SELECT id, gender, countryOfBirth, nationality, maritalStatus, birthDate, givenName, preferredGivenName, familyName, legalName, middleName, fullName, formattedName, status FROM customer;

CREATE TABLE externalReference (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    externalIdentifierType VARCHAR(20),
    type VARCHAR(20),
    customer_id VARCHAR(50) REFERENCES customer (id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO externalReference (name, externalIdentifierType, type, customer_id) VALUES ('http://facebook.com/17263635', 'facebookId', 'ExternalIdentifier', '42');
INSERT INTO externalReference (name, externalIdentifierType, type, customer_id) VALUES ('http://google.com/17263635', 'googleId', 'ExternalIdentifier', '42');
INSERT INTO externalReference (name, externalIdentifierType, type, customer_id) VALUES ('http://facebook.com/17263635', 'facebookId', 'ExternalIdentifier', '44');

SELECT name, externalIdentifierType, type FROM externalReference WHERE customer_id = '42';

{
  "id": "42",
  "href": "https://serverRoot/tmf-api/party/v5/individual/42",
  "@type": "Individual",
  "@baseType": "Party",
  "gender": "female",
  "countryOfBirth": "United States",
  "nationality": "American",
  "maritalStatus": "married",
  "birthDate": "1967-09-26T05:00:00.246Z",
  "givenName": "Jane",
  "preferredGivenName": "Lamborgizzia",
  "familyName": "Lamborgizzia",
  "legalName": "Smith",
  "middleName": "JL",
  "fullName": "Jane Smith ep Lamborgizzia",
  "formattedName": "Jane Smith ep Lamborgizzia",
  "status": "validated",
  "externalReference": [
    {
      "name": "http://facebook.com/17263635",
      "externalIdentifierType": "facebookId",
      "@type": "ExternalIdentifier"
    }
  ]
}

-- Drop Table
DROP TABLE customer;

-- Drop database
DROP DATABASE go_oda WITH (FORCE);
