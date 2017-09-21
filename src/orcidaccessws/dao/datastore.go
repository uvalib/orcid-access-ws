package dao

import (
   "database/sql"
   "fmt"
   _ "github.com/go-sql-driver/mysql"
   "orcidaccessws/api"
   "orcidaccessws/logger"
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
// get all by ID (should only be 1)
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
// set orcid by ID
//
func (db *DB) SetOrcidAttributesByCid(id string, orcid string) error {

   attributes, err := getOrcidAttributesByCid(db, id)
   if err != nil {
      return err
   }

   // if we did not find a record, create a new one
   if len(attributes) == 0 {

      stmt, err := db.Prepare("INSERT INTO orcid_attributes( cid, orcid ) VALUES( ?,? )")
      if err != nil {
         return err
      }

      _, err = stmt.Exec(id, orcid)
   } else {

      // we already have a record; do we actually need to do the update
      if attributes[0].Orcid != orcid {
         stmt, err := db.Prepare("UPDATE orcid_attributes SET orcid = ?, updated_at = NOW( ) WHERE cid = ? LIMIT 1")
         if err != nil {
            return err
         }
         _, err = stmt.Exec(orcid, id)
      }
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
      reg.Uri = fmt.Sprintf("http://orcid.org/%s", reg.Orcid)

      results = append(results, reg)
   }
   if err := rows.Err(); err != nil {
      return nil, err
   }

   logger.Log(fmt.Sprintf("OrcidAttributes request returns %d row(s)", len(results)))
   return results, nil
}
