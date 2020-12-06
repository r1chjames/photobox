CREATE DATABASE photobox;

USE photobox;
CREATE TABLE albums (
    id                      VARCHAR(36) PRIMARY KEY,
    name                    VARCHAR(255) NOT NULL,
    description             VARCHAR(255),
    tags                    VARCHAR(255),
    metadata                JSON
);

CREATE TABLE photos (
    id                      VARCHAR(100) PRIMARY KEY,
    name                    VARCHAR(255) NOT NULL,
    filesystemPath          VARCHAR(255) NOT NULL,
    albumId                 VARCHAR(36) NOT NULL,
    tags                    VARCHAR(255),
    metadata                JSON,

    FOREIGN KEY (albumId) REFERENCES albums(id)
);

CREATE TABLE settings (
    setting                 VARCHAR(255) NOT NULL,
    value                   VARCHAR(255) NOT NULL
);

CREATE TABLE jobs (
    job                     VARCHAR(255) NOT NULL,
    running                 BOOLEAN,
    lastRun                 VARCHAR(50)
);

# Create base data
INSERT INTO jobs VALUES ('Photo_index', false, null);

CREATE USER 'photobox'@'%' IDENTIFIED BY 'photobox';
GRANT USAGE ON *.* TO 'photobox'@'%';
GRANT ALL privileges ON photobox.* TO 'photobox'@'%';
