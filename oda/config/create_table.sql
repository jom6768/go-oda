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

-- TMF632_customer
CREATE TABLE customer (
    customer_id VARCHAR(50) PRIMARY KEY NOT NULL,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE,
    phone VARCHAR(20),
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Drop Table
DROP TABLE customer;

-- Drop database
DROP DATABASE go_oda WITH (FORCE);
