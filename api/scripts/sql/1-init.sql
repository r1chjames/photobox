CREATE DATABASE photobox;

CREATE USER 'photobox'@'%' IDENTIFIED BY 'photobox';
GRANT ALL privileges ON photobox.* TO 'photobox'@'%';
