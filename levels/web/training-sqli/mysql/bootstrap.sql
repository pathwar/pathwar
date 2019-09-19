USE training_sqli;

CREATE TABLE users (id INTEGER AUTO_INCREMENT PRIMARY KEY, username TEXT, password TEXT);
INSERT INTO users (username, password) VALUES ('admin', '__PASSPHRASE1__');

CREATE TABLE users2 (id INTEGER AUTO_INCREMENT PRIMARY KEY, username TEXT, password TEXT);
-- no entry in this table
