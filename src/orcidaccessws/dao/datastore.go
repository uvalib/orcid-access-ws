package dao

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
//    "strconv"
//    "fmt"
//    "orcidaccessws/api"
//    "orcidaccessws/logger"
    "orcidaccessws/api"
    "orcidaccessws/logger"
    "fmt"
)

type DB struct {
    *sql.DB
}

type Mapper struct {
    FieldClass  string
    FieldSource string
    FieldMapped string
}

var Database * DB

//
// create the DB singletomn
//
func NewDB( dataSourceName string ) error {
    db, err := sql.Open( "mysql", dataSourceName )
    if err != nil {
        return err
    }
    if err = db.Ping( ); err != nil {
        return err
    }
    Database = &DB{ db }
    return nil
}

//
// check our DB health
//
func ( db *DB ) Check( ) error {
    if err := db.Ping( ); err != nil {
        return err
    }
    return nil
}

//
// get all orcid records
//
func ( db *DB ) GetAllOrcid( ) ( [] * api.Orcid, error ) {

    rows, err := db.Query( "SELECT * FROM orcids ORDER BY id ASC" )
    if err != nil {
        return nil, err
    }
    defer rows.Close( )

    return orcidResults( rows )
}

//
// get all by ID (should only be 1)
//
func ( db *DB ) GetOrcidByCid( id string ) ( [] * api.Orcid, error ) {

    rows, err := db.Query( "SELECT * FROM orcids WHERE cid = ? LIMIT 1", id )
    if err != nil {
        return nil, err
    }
    defer rows.Close( )

    return orcidResults( rows )
}

/*

//
// Determine if the supplied deposit authorization already exists
//
func ( db *DB ) DepositAuthorizationExists( e api.Authorization ) ( bool, error ) {

    rows, err := db.Query( "SELECT COUNT(*) FROM depositauth WHERE computing_id = ? AND degree = ? AND plan = ? AND title = ?", e.ComputingId, e.Degree, e.Plan, e.Title )
    if err != nil {
        return false, err
    }
    defer rows.Close( )

    var count int
    for rows.Next() {
        err := rows.Scan( &count )
        if err != nil {
            return false, err
        }
    }

    if err := rows.Err( ); err != nil {
        return false, err
    }

    return count != 0, nil
}

//
// get all greater than a specified ID
//
func ( db *DB ) SearchDepositAuthorizationById( id string ) ( [] * api.Authorization, error ) {

    rows, err := db.Query( "SELECT * FROM depositauth WHERE id > ? ORDER BY id ASC", id )
    if err != nil {
        return nil, err
    }
    defer rows.Close( )

    return depositAuthorizationResults( rows )
}

//
// get all similar to the a specified computing ID
//
func ( db *DB ) SearchDepositAuthorizationByCid( cid string ) ( [] * api.Authorization, error ) {

    rows, err := db.Query( "SELECT * FROM depositauth WHERE computing_id LIKE ? ORDER BY id ASC", fmt.Sprintf( "%s%%", cid ) )
    if err != nil {
        return nil, err
    }
    defer rows.Close( )

    return depositAuthorizationResults( rows )
}

//
// create a new deposit authorization
//
func ( db *DB ) CreateDepositAuthorization( reg api.Authorization ) ( * api.Authorization, error ) {

    stmt, err := db.Prepare( "INSERT INTO depositauth( employee_id, computing_id, first_name, middle_name, last_name, career, program, plan, degree, title, doctype, approved_at ) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)" )
    if err != nil {
        return nil, err
    }

    res, err := stmt.Exec( reg.EmployeeId,
                           reg.ComputingId,
                           reg.FirstName,
                           reg.MiddleName,
                           reg.LastName,
                           reg.Career,
                           reg.Program,
                           reg.Plan,
                           reg.Degree,
                           reg.Title,
                           reg.DocType,
                           reg.ApprovedAt )
    if err != nil {
        return nil, err
    }

    lastId, err := res.LastInsertId( )
    if err != nil {
        return nil, err
    }

    reg.Id = strconv.FormatInt( lastId, 10 )
    return &reg, nil
}

//
// delete by ID
//
func ( db *DB ) DeleteDepositAuthorizationById( id string ) ( int64, error ) {

    stmt, err := db.Prepare( "DELETE FROM depositauth WHERE id = ? LIMIT 1" )
    if err != nil {
        return 0, err
    }

    res, err := stmt.Exec( id )
    if err != nil {
        return 0, err
    }

    rowCount, err := res.RowsAffected( )
    if err != nil {
        return 0, err
    }

    return rowCount, nil
}

//
// get all available for export
//
func ( db *DB ) GetDepositAuthorizationForExport( ) ( [] * api.Authorization, error ) {

    rows, err := db.Query( "SELECT * FROM depositauth WHERE accepted_at IS NOT NULL AND exported_at IS NULL ORDER BY id ASC" )
    if err != nil {
        return nil, err
    }
    defer rows.Close( )

    return depositAuthorizationResults( rows )
}

//
// update all export items with the time of export
//
func ( db *DB ) UpdateExportedDepositAuthorization( exports [] * api.Authorization ) error {

    stmt, err := db.Prepare( "UPDATE depositauth SET exported_at = NOW( ) WHERE id = ? LIMIT 1" )
    if err != nil {
        return err
    }

    for _, rec := range exports {
        _, err := stmt.Exec( rec.Id )
        if err != nil {
            return err
        }
    }

    return nil
}

//
// update an item that has been 'fulfilled'
//
func ( db *DB ) UpdateFulfilledDepositAuthorizationById( id string, did string ) error {

    stmt, err := db.Prepare( "UPDATE depositauth SET exported_at = NULL, accepted_at = NOW( ), status = ?, libra_id = ? WHERE id = ? LIMIT 1" )
    if err != nil {
        return err
    }

    _, err = stmt.Exec( "submitted", did, id )
    if err != nil {
        return err
    }

    return nil
}

func ( db *DB ) GetFieldMapperList( ) ( [] * Mapper, error ) {

    rows, err := db.Query( "SELECT field_class, field_name, field_value FROM fieldmapper" )
    if err != nil {
        return nil, err
    }
    defer rows.Close( )

    results := make([ ] * Mapper, 0 )

    for rows.Next() {
        mapping := new( Mapper )
        err := rows.Scan(
            &mapping.FieldClass,
            &mapping.FieldSource,
            &mapping.FieldMapped )
        if err != nil {
            return nil, err
        }

        results = append( results, mapping )
    }

    if err := rows.Err( ); err != nil {
        return nil, err
    }

    return results, nil
}

 */

//
// private implementation methods
//

func orcidResults( rows * sql.Rows ) ( [] * api.Orcid, error ) {

    var optionalUpdatedAt sql.NullString

    results := make([ ] * api.Orcid, 0 )
    for rows.Next() {
        reg := new( api.Orcid )
        err := rows.Scan( &reg.Id,
            &reg.Cid,
            &reg.Orcid,
            &reg.CreatedAt,
            &optionalUpdatedAt )
        if err != nil {
            return nil, err
        }

        if optionalUpdatedAt.Valid {
            reg.UpdatedAt = optionalUpdatedAt.String
        }

        results = append( results, reg )
    }
    if err := rows.Err( ); err != nil {
        return nil, err
    }

    logger.Log( fmt.Sprintf( "Orcid request returns %d row(s)", len( results ) ) )
    return results, nil
}

