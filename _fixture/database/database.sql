DROP DATABASE IF EXISTS echo;
CREATE DATABASE echo DEFAULT CHARACTER SET utf8;

USE echo;

CREATE TABLE visitor
(
  id   int(11) not null AUTO_INCREMENT,
  name char(64),
  PRIMARY KEY (id)
) ENGINE = InnoDB;

INSERT INTO visitor (id, name)
VALUES (1, 'Tom'),
       (2, 'Bob'),
       (3, 'Jack');