-- PostgreSQL command --
\l; -- List all databases
\c [database_name]; -- Connect database
\du; -- List all users
\du+; -- List all users with detailed information
\dt; -- List all tables
\dt+; -- List all tables with detailed information
exit -- Close connection

-- Connect database --
"/Applications/Postgres.app/Contents/Versions/16/bin/psql" -U myuser -d go_oda

SELECT current_user,session_user; -- Check user
SET ROLE [username]; -- Change user login

-- Create database --
CREATE DATABASE go_oda;
CREATE USER myuser WITH ENCRYPTED PASSWORD 'mypass';
GRANT ALL PRIVILEGES ON DATABASE go_oda TO myuser;
GRANT ALL ON SCHEMA public TO myuser;

-- Drop database
DROP DATABASE go_oda WITH (FORCE);

-- Drop Table
DROP TABLE individual;

-- SELECT - IFNULL --
COALESCE(par.href,'')
COALESCE(ind.birthDate::TEXT,'')
TO_CHAR(current_timestamp,'DD/MM/YYYY HH24:MI:SS')
('2018-01-15T08:54:45.000Z'::TIMESTAMP)::DATE AS date
('2018-01-15T08:54:45.000Z'::TIMESTAMP)::TIME AS time


--// TMF629_Customer Management API \\--
-- Customer - SELECT --
SELECT par.id,par.href,cus.type,par.name,par.description,par.role,par.status,par.statusReason,par.startDateTime,par.endDateTime
FROM partyRole par INNER JOIN customer cus ON par.id=cus.partyRole_id;

-- Customer - INSERT --
WITH partyins AS (
    INSERT INTO partyRole (id,href,name,description,role,status,statusReason,startDateTime,endDateTime,type) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,'PartyRole') RETURNING id
), customerins AS (
    INSERT INTO customer (type,partyRole_id) VALUES ('Invididual',$1) RETURNING id
)
SELECT id FROM partyins;

-- Customer - DELETE --
WITH customerdel AS (
    DELETE FROM customer WHERE partyRole_id = $1
)
    DELETE FROM partyRole WHERE id = $1;

CREATE TABLE partyRole (
    id VARCHAR(50) PRIMARY KEY NOT NULL,
    href VARCHAR(100),
    name VARCHAR(50),
    description VARCHAR(250),
    role VARCHAR(50),
    status VARCHAR(50),
    statusReason VARCHAR(100),
    startDateTime TIMESTAMP,
    endDateTime TIMESTAMP,
    type VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO partyRole (id,href,name,description,role,status,statusReason,startDateTime,endDateTime,type)
VALUES ('1140','http://localhost:8629/tmf-api/customerManagement/v5/customer/1140','Moon Football Club','Testing','Buyer','Approved','Account details checked','2018-06-12T00:00:00Z','2019-01-01T00:00:00Z','PartyRole') RETURNING id;

UPDATE partyRole SET description = '',role = '' WHERE partyRole_id = '1140';

SELECT id,href,name,description,role,status,statusReason,startDateTime,endDateTime,type FROM partyRole WHERE partyRole_id = '1140';

CREATE TABLE customer (
    id SERIAL PRIMARY KEY,
    type VARCHAR(50),
    partyRole_id VARCHAR(50) REFERENCES partyRole (id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO customer (type,partyRole_id) VALUES ('Customer','1140') RETURNING id;

SELECT type,partyRole_id FROM customer WHERE partyRole_id = '1140';

CREATE TABLE contactMedium (
    id SERIAL PRIMARY KEY,
    preferred BOOLEAN,
    contactType VARCHAR(50),
    phoneNumber VARCHAR(20),
    city VARCHAR(50),
    country VARCHAR(50),
    postCode VARCHAR(20),
    street1 VARCHAR(100),
    startDateTime TIMESTAMP,
    endDateTime TIMESTAMP,
    type VARCHAR(50),
    partyRole_id VARCHAR(50) REFERENCES partyRole (id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO contactMedium (preferred,contactType,phoneNumber,startDateTime,endDateTime,type,partyRole_id)
VALUES (true,'homePhone','01 09 75 83 51','2018-06-12T00:00:00Z','2019-01-01T00:00:00Z','PhoneContactMedium','1140') RETURNING id;
INSERT INTO contactMedium (preferred,contactType,city,country,postCode,street1,startDateTime,endDateTime,type,partyRole_id)
VALUES (false,'homeAddress','Paris','France','75014','15 Rue des Canards','2018-06-12T00:00:00Z','2019-01-01T00:00:00Z','GeographicAddressContactMedium','1140') RETURNING id;

SELECT preferred,contactType,phoneNumber,city,country,postCode,street1,startDateTime,endDateTime,type FROM contactMedium WHERE partyRole_id = '1140';


--// TMF632_Party Management API \\--
-- Individual - SELECT --
SELECT par.id,par.href,ind.type,par.type AS "baseType",ind.gender,ind.countryOfBirth,ind.nationality,ind.maritalStatus,ind.birthDate,ind.givenName,ind.preferredGivenName,ind.familyName,ind.legalName,ind.middleName,ind.fullName,ind.formattedName,ind.status
FROM party par INNER JOIN individual ind ON par.id=ind.party_id WHERE par.id = $1 LIMIT 1;

-- Individual - INSERT --
WITH partyins AS (
    INSERT INTO party (id,href,type) VALUES ($1,$2,'Party') RETURNING id
), individualins AS (
    INSERT INTO individual (gender,countryOfBirth,nationality,maritalStatus,birthDate,givenName,preferredGivenName,familyName,legalName,middleName,fullName,formattedName,status,type,party_id) VALUES ($3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,'Invididual',$1) RETURNING id
)
SELECT id FROM partyins;

-- Individual - DELETE --
WITH individualdel AS (
    DELETE FROM individual WHERE party_id = $1
)
    DELETE FROM party WHERE id = $1;

-- Organization - SELECT --
SELECT par.id,par.href,org.type,par.type AS "baseType",org.isLegalEntity,org.isHeadOffice,org.organizationType,org.name,org.tradingName,org.nameType,org.status
FROM party par INNER JOIN organization org ON par.id=org.party_id WHERE par.id = $1 LIMIT 1;

-- Organization - INSERT --
WITH partyins AS (
    INSERT INTO party (id,href,type) VALUES ($1,$2,'Party') RETURNING id
), organizationins AS (
    INSERT INTO organization (isLegalEntity,isHeadOffice,organizationType,name,tradingName,nameType,status,type,party_id) VALUES ($3,$4,$5,$6,$7,$8,$9,'Organization',$1) RETURNING id
)
SELECT id FROM partyins;

-- Organization - DELETE --
WITH organizationdel AS (
    DELETE FROM organization WHERE party_id = $1
)
    DELETE FROM party WHERE id = $1;

CREATE TABLE party (
    id VARCHAR(50) PRIMARY KEY NOT NULL,
    href VARCHAR(100),
    type VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO party (id,href,type) VALUES ('42','http://localhost:8632/tmf-api/partyManagement/v5/individual/42','Party') RETURNING id;
INSERT INTO party (id,href,type) VALUES ('43','http://localhost:8632/tmf-api/partyManagement/v5/individual/43','Party') RETURNING id;
INSERT INTO party (id,href,type) VALUES ('44','http://localhost:8632/tmf-api/partyManagement/v5/individual/44','Party') RETURNING id;
INSERT INTO party (id,href,type) VALUES ('128','http://localhost:8632/tmf-api/partyManagement/v5/individual/128','Party') RETURNING id;
INSERT INTO party (id,href,type) VALUES ('129','http://localhost:8632/tmf-api/partyManagement/v5/individual/129','Party') RETURNING id;

SELECT id,href,type FROM party WHERE id = '42';

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
    legalName VARCHAR(50),-- UNIQUE
    middleName VARCHAR(50),
    fullName VARCHAR(100),
    formattedName VARCHAR(100),
    status VARCHAR(50),
    type VARCHAR(50),
    party_id VARCHAR(50) REFERENCES party (id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO individual (gender,countryOfBirth,nationality,maritalStatus,birthDate,givenName,preferredGivenName,familyName,legalName,middleName,fullName,formattedName,status,type,party_id)
VALUES ('female','United States','American','married','1967-09-26T05:00:00.246Z','Jane','Lamborgizzia','Lamborgizzia','Smith','JL','Jane Smith ep Lamborgizzia','Jane Smith ep Lamborgizzia','validated','Individual','42') RETURNING id;
INSERT INTO individual (gender,countryOfBirth,nationality,maritalStatus,birthDate,givenName,preferredGivenName,familyName,legalName,middleName,fullName,formattedName,status,type,party_id)
VALUES ('male','Thailand','Thai','single','1988-08-08T00:00:00.000Z','Jane','Lamborgizzia','Lamborgizzia','Smith','JL','Jane Smith ep Lamborgizzia','Jane Smith ep Lamborgizzia','validated','Individual','43') RETURNING id;
INSERT INTO individual (birthDate,fullName,type,party_id)
VALUES ('1988-08-08T00:00:00.000Z','Jane Smith ep Lamborgizzia','Individual','44') RETURNING id;

UPDATE individual SET gender = 'male',status = 'active' WHERE party_id = '46';

SELECT id,gender,countryOfBirth,nationality,maritalStatus,birthDate,givenName,preferredGivenName,familyName,legalName,middleName,fullName,formattedName,status,type,party_id FROM individual WHERE party_id = '42';

CREATE TABLE organization (
    id SERIAL PRIMARY KEY,
    isLegalEntity BOOLEAN,
    isHeadOffice BOOLEAN,
    organizationType VARCHAR(20),
    name VARCHAR(100),
    tradingName VARCHAR(100),
    nameType VARCHAR(20),
    status VARCHAR(50),
    type VARCHAR(50),
    party_id VARCHAR(50) REFERENCES party (id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO organization (isLegalEntity,isHeadOffice,organizationType,name,tradingName,nameType,status,type,party_id)
VALUES (true,true,'company','Coffee Do Brazil','Coffee Do Brazil Fair Trade','inc','validated','Organization','128') RETURNING id;
INSERT INTO organization (organizationType,name,tradingName,nameType,status,type,party_id)
VALUES ('company','Coffee Do Brazil','Coffee Do Brazil Fair Trade','inc','validated','Organization','129') RETURNING id;
INSERT INTO organization (name,type,party_id)
VALUES ('Coffee Do Brazil','Organization','131') RETURNING id;

DELETE FROM organization WHERE party_id = '131';

SELECT id,isLegalEntity,isHeadOffice,organizationType,name,tradingName,nameType,status FROM organization WHERE party_id = '128';

CREATE TABLE externalReference (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    externalIdentifierType VARCHAR(20),
    type VARCHAR(50),
    party_id VARCHAR(50) REFERENCES party (id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO externalReference (name,externalIdentifierType,type,party_id) VALUES ('http://facebook.com/17263635','facebookId','ExternalIdentifier','42') RETURNING id;
INSERT INTO externalReference (name,externalIdentifierType,type,party_id) VALUES ('http://google.com/17263635','googleId','ExternalIdentifier','42') RETURNING id;
INSERT INTO externalReference (name,externalIdentifierType,type,party_id) VALUES ('http://facebook.com/17263636','facebookId','ExternalIdentifier','44') RETURNING id;
INSERT INTO externalReference (name,externalIdentifierType,type,party_id) VALUES ('http://coffeedobrazil.com','internetSite','ExternalIdentifier','128') RETURNING id;

DELETE FROM externalReference WHERE party_id = $1;

SELECT id,name,externalIdentifierType,type FROM externalReference WHERE party_id = '42';
