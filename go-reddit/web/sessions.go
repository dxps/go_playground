package web

import (
	"database/sql"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
)

// NewSessionManager creates a SessionManager that uses
// the provided data source name to work with the `sessions` table.
func NewSessionManager(dataSourceName string) (*scs.SessionManager, error) {

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	sessions := scs.New()
	sessions.Store = postgresstore.New(db)
	return sessions, nil
}
