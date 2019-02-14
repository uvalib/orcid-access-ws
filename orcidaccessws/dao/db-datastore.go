package dao

import (
	"database/sql"
	"fmt"
	"time"

	// needed
	_ "github.com/go-sql-driver/mysql"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/api"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/config"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/logger"
)

// this is our DB implementation
type storage struct {
	*sql.DB
}

//
// newDBStore -- create a DB version of the storage singleton
//
func newDBStore() (Storage, error) {

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?allowOldPasswords=1&tls=%s&sql_notes=false&timeout=%s&readTimeout=%s&writeTimeout=%s",
		config.Configuration.DbUser,
		config.Configuration.DbPassphrase,
		config.Configuration.DbHost,
		config.Configuration.DbName,
		config.Configuration.DbSecure,
		config.Configuration.DbTimeout,
		config.Configuration.DbTimeout,
		config.Configuration.DbTimeout)

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	//taken from https://github.com/go-sql-driver/mysql/issues/461
	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(2)

	return &storage{db}, nil
}

//
// CheckDB -- check our database health
//
func (s *storage) Check() error {
	return s.Ping()
}

//
// GetAllOrcidAttributes -- get all orcid records
//
func (s *storage) GetAllOrcidAttributes() ([]*api.OrcidAttributes, error) {

	rows, err := s.Query("SELECT * FROM orcid_attributes ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return orcidResults(rows)
}

//
// GetOrcidAttributesByCid -- get all by ID (should only be 1)
//
func (s *storage) GetOrcidAttributesByCid(id string) ([]*api.OrcidAttributes, error) {
	return (getOrcidAttributesByCid(s, id))
}

//
// DelOrcidAttributesByCid -- delete by ID (should only be 1)
//
func (s *storage) DelOrcidAttributesByCid(id string) error {

	stmt, err := s.Prepare("DELETE FROM orcid_attributes WHERE cid = ? LIMIT 1")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)

	return err
}

//
// SetOrcidAttributesByCid -- set orcid attributes by ID
//
func (s *storage) SetOrcidAttributesByCid(id string, attributes api.OrcidAttributes) error {

	existing, err := getOrcidAttributesByCid(s, id)
	if err != nil {
		return err
	}

	// if we did not find a record, create a new one
	if len(existing) == 0 {

		stmt, err := s.Prepare("INSERT INTO orcid_attributes( cid, orcid, oauth_access, oauth_refresh, oauth_scope ) VALUES( ?,?,?,?,? )")
		if err != nil {
			return err
		}

		_, err = stmt.Exec(
			id,
			attributes.Orcid,
			attributes.OauthAccessToken,
			attributes.OauthRefreshToken,
			attributes.OauthScope)

	} else {

		// a special case where we preserve the existing ORCID if none provided
		newOrcid := existing[0].Orcid
		if len(attributes.Orcid) != 0 {
			newOrcid = attributes.Orcid
		}

		stmt, err := s.Prepare("UPDATE orcid_attributes SET orcid = ?, oauth_access = ?, oauth_refresh = ?, oauth_scope = ?, updated_at = NOW( ) WHERE cid = ? LIMIT 1")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(
			newOrcid,
			attributes.OauthAccessToken,
			attributes.OauthRefreshToken,
			attributes.OauthScope,
			id)
	}

	return err
}

//
// private implementation methods
//

func getOrcidAttributesByCid(s *storage, id string) ([]*api.OrcidAttributes, error) {

	rows, err := s.Query("SELECT * FROM orcid_attributes WHERE cid = ? LIMIT 1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return orcidResults(rows)
}

func orcidResults(rows *sql.Rows) ([]*api.OrcidAttributes, error) {

	var optionalUpdatedAt sql.NullString

	results := make([]*api.OrcidAttributes, 0)
	for rows.Next() {
		reg := new(api.OrcidAttributes)
		err := rows.Scan(&reg.ID,
			&reg.Cid,
			&reg.Orcid,
			&reg.OauthAccessToken,
			&reg.OauthRefreshToken,
			&reg.OauthScope,
			&reg.CreatedAt,
			&optionalUpdatedAt)
		if err != nil {
			return nil, err
		}

		if optionalUpdatedAt.Valid {
			reg.UpdatedAt = optionalUpdatedAt.String
		}

		// hack for now...
		reg.URI = fmt.Sprintf("%s/%s", config.Configuration.OrcidOauthURL, reg.Orcid)

		results = append(results, reg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	logger.Log(fmt.Sprintf("OrcidAttributes request returns %d row(s)", len(results)))
	return results, nil
}

//
// end of file
//
