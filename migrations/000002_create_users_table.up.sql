CREATE TABLE users (
    id SERIAL NOT NULL,
    email varchar(100) NOT NULL,
    password_hash varchar(100) NOT NULL,
    PRIMARY KEY (id)
);