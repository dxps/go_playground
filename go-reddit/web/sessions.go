package web

import (
	"context"
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

type SessionData struct {
	FlashMessage string
	Form         interface{} // a 'generic' type to allow working with any form type
	// UserID uuid.UUID
}

func GetSessionData(session *scs.SessionManager, ctx context.Context) SessionData {

	var data SessionData
	data.FlashMessage = session.PopString(ctx, "flash")
	data.Form = session.Pop(ctx, "form")
	// data.UserID, _ = session.Get(ctx, "user_id").(uuid.UUID)
	if data.Form == nil {
		data.Form = map[string]string{}
	}
	return data
}
