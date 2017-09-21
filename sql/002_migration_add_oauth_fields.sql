-- add OAUTH fields to orcid_attributes table
ALTER TABLE orcid_attributes ADD COLUMN oauth_access VARCHAR( 64 ) NOT NULL DEFAULT '' AFTER orcid;
ALTER TABLE orcid_attributes ADD COLUMN oauth_refresh VARCHAR( 64 ) NOT NULL DEFAULT '' AFTER oauth_access;
ALTER TABLE orcid_attributes ADD COLUMN oauth_scope VARCHAR( 64 ) NOT NULL DEFAULT '' AFTER oauth_refresh;