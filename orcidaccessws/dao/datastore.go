package dao

import (
	// needed
	_ "github.com/go-sql-driver/mysql"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/api"
)

// our storage interface
type Storage interface {
	Check() error
	GetAllOrcidAttributes() ([]*api.OrcidAttributes, error)
	GetOrcidAttributesByCid(id string) ([]*api.OrcidAttributes, error)
	DelOrcidAttributesByCid(id string) error
	SetOrcidAttributesByCid(id string, attributes api.OrcidAttributes) error
	//Destroy() error
}

// our singleton store
var Store Storage

// our factory
func NewDatastore() error {
	var err error
	// mock implementation here
	Store, err = newDBStore()
	return err
}

//
// end of file
//
