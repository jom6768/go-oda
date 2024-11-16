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
SELECT par.id,COALESCE(par.href,''),COALESCE(ind.type,''),COALESCE(par.type,'') AS "baseType",COALESCE(ind.gender,''),COALESCE(ind.countryOfBirth,''),COALESCE(ind.nationality,''),COALESCE(ind.maritalStatus,''),COALESCE(ind.birthDate::TEXT,''),COALESCE(ind.givenName,''),COALESCE(ind.preferredGivenName,''),COALESCE(ind.familyName,''),COALESCE(ind.legalName,''),COALESCE(ind.middleName,''),ind.fullName,COALESCE(ind.formattedName,''),COALESCE(ind.status,'')
FROM party par INNER JOIN individual ind ON par.id=ind.party_id;

SELECT par.id,COALESCE(par.href,''),COALESCE(org.type,''),COALESCE(par.type,'') AS "baseType",COALESCE(org.isLegalEntity::TEXT,''),COALESCE(org.isHeadOffice::TEXT,''),COALESCE(org.organizationType,''),COALESCE(org.name,''),COALESCE(org.tradingName,''),COALESCE(org.nameType,''),COALESCE(org.status,'')
FROM party par INNER JOIN organization org ON par.id=org.party_id;

CREATE TABLE party (
    id VARCHAR(50) PRIMARY KEY NOT NULL,
    href VARCHAR(100),
    type VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO party (id, href, type) VALUES ('42', 'http://localhost:8081/tmf-api/party/v5/individual/42', 'Party') RETURNING id;
INSERT INTO party (id, href, type) VALUES ('43', 'http://localhost:8081/tmf-api/party/v5/individual/43', 'Party') RETURNING id;
INSERT INTO party (id, href, type) VALUES ('44', 'http://localhost:8081/tmf-api/party/v5/individual/44', 'Party') RETURNING id;
INSERT INTO party (id, href, type) VALUES ('128', 'http://localhost:8081/tmf-api/party/v5/individual/128', 'Party') RETURNING id;

SELECT id, href, type FROM party WHERE id = '42';

CREATE TABLE individual (
    id SERIAL PRIMARY KEY,
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
    type VARCHAR(20),
    party_id VARCHAR(50) REFERENCES party (id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO individual (gender, countryOfBirth, nationality, maritalStatus, birthDate, givenName, preferredGivenName, familyName, legalName, middleName, fullName, formattedName, status, type, party_id)
VALUES ('female', 'United States', 'American', 'married', '1967-09-26T05:00:00.246Z', 'Jane', 'Lamborgizzia', 'Lamborgizzia', 'Smith', 'JL', 'Jane Smith ep Lamborgizzia', 'Jane Smith ep Lamborgizzia', 'validated', 'Individual', '42') RETURNING id;
INSERT INTO individual (gender, countryOfBirth, nationality, maritalStatus, birthDate, givenName, preferredGivenName, familyName, legalName, middleName, fullName, formattedName, status, type, party_id)
VALUES ('male', 'Thailand', 'Thai', 'single', '1988-08-08T00:00:00.000Z', 'Jane', 'Lamborgizzia', 'Lamborgizzia', 'Smith', 'JL', 'Jane Smith ep Lamborgizzia', 'Jane Smith ep Lamborgizzia', 'active', 'Individual', '43') RETURNING id;
INSERT INTO individual (birthDate, fullName, type, party_id)
VALUES ('1988-08-08T00:00:00.000Z', 'Jane Smith ep Lamborgizzia', 'Individual', '44') RETURNING id;

UPDATE individual SET gender = 'male', status = 'active' WHERE party_id = '46';

SELECT id, gender, countryOfBirth, nationality, maritalStatus, birthDate, givenName, preferredGivenName, familyName, legalName, middleName, fullName, formattedName, status, type, party_id FROM individual WHERE party_id = '42';

CREATE TABLE organization (
    id SERIAL PRIMARY KEY,
    isLegalEntity BOOLEAN,
    isHeadOffice BOOLEAN,
    organizationType VARCHAR(20),
    name VARCHAR(100),
    tradingName VARCHAR(100),
    nameType VARCHAR(20),
    status VARCHAR(50),
    type VARCHAR(20),
    party_id VARCHAR(50) REFERENCES party (id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO organization (isLegalEntity, isHeadOffice, organizationType, name, tradingName, nameType, status, type, party_id)
VALUES (true, true, 'company', 'Coffee Do Brazil', 'Coffee Do Brazil Fair Trade', 'inc', 'validated', 'Individual', '128') RETURNING id;

SELECT id, isLegalEntity, isHeadOffice, organizationType, name, tradingName, nameType, status FROM organization WHERE party_id = '128';

CREATE TABLE externalReference (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    externalIdentifierType VARCHAR(20),
    type VARCHAR(20),
    party_id VARCHAR(50) REFERENCES party (id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO externalReference (name, externalIdentifierType, type, party_id) VALUES ('http://facebook.com/17263635', 'facebookId', 'ExternalIdentifier', '42') RETURNING id;
INSERT INTO externalReference (name, externalIdentifierType, type, party_id) VALUES ('http://google.com/17263635', 'googleId', 'ExternalIdentifier', '42') RETURNING id;
INSERT INTO externalReference (name, externalIdentifierType, type, party_id) VALUES ('http://facebook.com/17263636', 'facebookId', 'ExternalIdentifier', '44') RETURNING id;
INSERT INTO externalReference (name, externalIdentifierType, type, party_id) VALUES ('http://coffeedobrazil.com', 'internetSite', 'ExternalIdentifier', '128') RETURNING id;

SELECT id, name, externalIdentifierType, type FROM externalReference WHERE party_id = '42';

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
DROP TABLE individual;

-- Drop database
DROP DATABASE go_oda WITH (FORCE);
