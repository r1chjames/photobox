CREATE DATABASE photobox;

CREATE USER 'photobox'@'%' IDENTIFIED BY 'photobox';
GRANT USAGE ON *.* TO '%'@'%';
GRANT ALL privileges ON photobox.* TO 'photobox'@'%';
