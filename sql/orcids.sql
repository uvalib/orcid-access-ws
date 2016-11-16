-- drop the table if it exists
DROP TABLE IF EXISTS orcids;

-- and create the new one
CREATE TABLE orcids(
   id          INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
   cid         VARCHAR( 32 ) NOT NULL DEFAULT '',
   orcid       VARCHAR( 32 ) NOT NULL DEFAULT '',
   created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at  TIMESTAMP NULL,
   UNIQUE KEY( cid )
) CHARACTER SET utf8 COLLATE utf8_bin;

-- set of degree mapping values (manually created)
INSERT INTO orcids( cid, orcid ) VALUES( "dpg3k", "0000-0002-0566-4186" );
INSERT INTO orcids( cid, orcid ) VALUES( "ecr2c", "0000-0003-4520-4923" );