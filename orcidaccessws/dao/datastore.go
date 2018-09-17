package dao

import (
	"database/sql"
	"fmt"
	// needed
	_ "github.com/go-sql-driver/mysql"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/api"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/config"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/logger"
)

type dbStruct struct {
	*sql.DB
}

//
// DB -- the database instance
//
var DB *dbStruct

//
// NewDB -- create the database singletomn
//
func NewDB(dbHost string, dbSecure bool, dbName string, dbUser string, dbPassword string, dbTimeout string) error {

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?allowOldPasswords=1&tls=%t&sql_notes=false&timeout=%s&readTimeout=%s&writeTimeout=%s",
		dbUser,
		dbPassword,
		dbHost,
		dbName,
		dbSecure,
		dbTimeout,
		dbTimeout,
		dbTimeout)

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	DB = &dbStruct{db}
	return nil
}

//
// CheckDB -- check our database health
//
func (db *dbStruct) CheckDB() error {
	return db.Ping()
}

//
// GetAllOrcidAttributes -- get all orcid records
//
func (db *dbStruct) GetAllOrcidAttributes() ([]*api.OrcidAttributes, error) {

	rows, err := db.Query("SELECT * FROM orcid_attributes ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return orcidResults(rows)
}

//
// GetOrcidAttributesByCid -- get all by ID (should only be 1)
//
func (db *dbStruct) GetOrcidAttributesByCid(id string) ([]*api.OrcidAttributes, error) {
	return (getOrcidAttributesByCid(db, id))
}

//
// DelOrcidAttributesByCid -- delete by ID (should only be 1)
//
func (db *dbStruct) DelOrcidAttributesByCid(id string) error {

	stmt, err := db.Prepare("DELETE FROM orcid_attributes WHERE cid = ? LIMIT 1")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)

	return err
}

//
// SetOrcidAttributesByCid -- set orcid attributes by ID
//
func (db *dbStruct) SetOrcidAttributesByCid(id string, attributes api.OrcidAttributes) error {

	existing, err := getOrcidAttributesByCid(db, id)
	if err != nil {
		return err
	}

	// if we did not find a record, create a new one
	if len(existing) == 0 {

		stmt, err := db.Prepare("INSERT INTO orcid_attributes( cid, orcid, oauth_access, oauth_refresh, oauth_scope ) VALUES( ?,?,?,?,? )")
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

		stmt, err := db.Prepare("UPDATE orcid_attributes SET orcid = ?, oauth_access = ?, oauth_refresh = ?, oauth_scope = ?, updated_at = NOW( ) WHERE cid = ? LIMIT 1")
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

func getOrcidAttributesByCid(db *dbStruct, id string) ([]*api.OrcidAttributes, error) {

	rows, err := db.Query("SELECT * FROM orcid_attributes WHERE cid = ? LIMIT 1", id)
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
