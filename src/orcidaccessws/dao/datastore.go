package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"orcidaccessws/api"
	"orcidaccessws/logger"
    "orcidaccessws/config"
)

type DB struct {
	*sql.DB
}

type Mapper struct {
	FieldClass  string
	FieldSource string
	FieldMapped string
}

var Database *DB

//
// create the DB singletomn
//
func NewDB(dataSourceName string) error {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	Database = &DB{db}
	return nil
}

//
// check our DB health
//
func (db *DB) Check() error {
	if err := db.Ping(); err != nil {
		return err
	}
	return nil
}

//
// get all orcid records
//
func (db *DB) GetAllOrcidAttributes() ([]*api.OrcidAttributes, error) {

	rows, err := db.Query("SELECT * FROM orcid_attributes ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return orcidResults(rows)
}

//
// get all by ID (should only be 1)
//
func (db *DB) GetOrcidAttributesByCid(id string) ([]*api.OrcidAttributes, error) {
	return (getOrcidAttributesByCid(db, id))
}

//
// delete by ID (should only be 1)
//
func (db *DB) DelOrcidAttributesByCid(id string) error {

	stmt, err := db.Prepare("DELETE FROM orcid_attributes WHERE cid = ? LIMIT 1")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)

	return err
}

//
// set orcid attributes by ID
//
func (db *DB) SetOrcidAttributesByCid(id string, attributes api.OrcidAttributes) error {

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

func getOrcidAttributesByCid(db *DB, id string) ([]*api.OrcidAttributes, error) {

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
		err := rows.Scan(&reg.Id,
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
		reg.Uri = fmt.Sprintf("%s/%s", config.Configuration.OrcidOauthUrl, reg.Orcid)

		results = append(results, reg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	logger.Log(fmt.Sprintf("OrcidAttributes request returns %d row(s)", len(results)))
	return results, nil
}
