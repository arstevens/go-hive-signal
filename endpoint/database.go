package endpoint

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

var (
	Host     = "localhost"
	Port     = 5432
	User     = "postgres"
	Password = "postgres"
	DBName   = "p2p_cdn"
)

const (
	tableName       = "registered_origins"
	fieldName       = "origin_id"
	insertStatement = "INSERT INTO " + tableName + " (" + fieldName + ") VALUES ($1)"
	removeStatement = "DELETE FROM " + tableName + " WHERE " + fieldName + " = $1"
	readStatement   = "SELECT " + fieldName + " FROM " + tableName
)

//EndpointRegistrationDatabase is a register of all valid Origin IDs
type EndpointRegistrationDatabase struct {
	backupDB      *sql.DB
	databaseMutex *sync.Mutex
	originSet     map[string]bool
}

//New creates a new EndpointRegistrationDatabase
func New() (*EndpointRegistrationDatabase, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		Host, Port, User, Password, DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("Failed to open databse %s in New()", DBName)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Failed to open connection to databse %s in New()", DBName)
	}

	inMemory, err := loadDatabaseIntoMemory(db)
	if err != nil {
		return nil, fmt.Errorf("Failed to read database %s into memory in New()", DBName)
	}

	return &EndpointRegistrationDatabase{
		backupDB:      db,
		databaseMutex: &sync.Mutex{},
		originSet:     inMemory,
	}, nil
}

//IsRegistered checks if 'originID' is registered
func (ed *EndpointRegistrationDatabase) IsRegistered(originID string) bool {
	ed.databaseMutex.Lock()
	defer ed.databaseMutex.Unlock()
	return ed.originSet[originID]
}

//AddOrigin adds 'originID' to the set of registered IDs
func (ed *EndpointRegistrationDatabase) AddOrigin(originID string) error {
	ed.databaseMutex.Lock()
	defer ed.databaseMutex.Unlock()

	if ed.originSet[originID] {
		return fmt.Errorf("Failed to add origin in EndpointRegistrationDatabase.AddOrigin(): "+
			"Origin %s already exists", originID)
	}
	_, err := ed.backupDB.Exec(insertStatement, originID)
	if err != nil {
		return fmt.Errorf("Failed to add origin in EndpointRegistrationDatabase.AddOrigin(): "+
			"Failed to insert %s into postgres database: %v", originID, err)
	}

	ed.originSet[originID] = true
	return nil
}

//RemoveOrigin removes 'originID' from the set of registered IDs
func (ed *EndpointRegistrationDatabase) RemoveOrigin(originID string) error {
	ed.databaseMutex.Lock()
	defer ed.databaseMutex.Unlock()

	if !ed.originSet[originID] {
		return fmt.Errorf("Failed to remove origin in EndpointRegistrationDatabase.RemoveOrigin(): "+
			"Origin %s is not registered", originID)
	}
	_, err := ed.backupDB.Exec(removeStatement, originID)
	if err != nil {
		return fmt.Errorf("Failed to remove origin in EndpointRegistrationDatabase.RemoveOrigin(): "+
			"Failed to remove %s from postgres database: %v", originID, err)
	}

	delete(ed.originSet, originID)
	return nil
}

/*Close closes the EndpointRegistrationDatabase object. Behaviour of any
method calls after Close() is called are undefined*/
func (ed *EndpointRegistrationDatabase) Close() error {
	ed.databaseMutex.Lock()
	ed.backupDB.Close()
	ed.databaseMutex.Unlock()
	return nil
}
