package mysql

import (
	"database/sql"
	"github.com/vision8tech/snippetbox/pkg/models"
)

type SnippetsStore struct {
	DB *sql.DB
}

func (store *SnippetsStore) Insert(title, content, expires string) (int, error) {

	insertStmt := `INSERT INTO snippets (title, content, created, expires)
		VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := store.DB.Exec(insertStmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil

}

func (store *SnippetsStore) Get(id int) (*models.Snippet, error) {

	s := &models.Snippet{}
	stmt := `SELECT id, title, content, created, expires FROM snippets
			 WHERE expires > UTC_TIMESTAMP() AND id = ?`
	err := store.DB.QueryRow(stmt, id).
		Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil

}

func (store *SnippetsStore) Latest() ([]*models.Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
		WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	rows, err := store.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// create an empty slice for any found rows
	var snippets []*models.Snippet
	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil

}
