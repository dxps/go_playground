package postgres

import (
	"fmt"

	goreddit "devisions.org/go-reddit"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CommentStore struct {
	*sqlx.DB
}

func (s *CommentStore) GetComment(id uuid.UUID) (goreddit.Comment, error) {
	var c goreddit.Comment
	if err := s.Get(&c, `SELECT * FROM comments WHERE id = $1`, id); err != nil {
		return goreddit.Comment{}, fmt.Errorf("error getting comment: %w", err)
	}
	return c, nil
}

func (s *CommentStore) GetCommentsByPost(postID uuid.UUID) ([]goreddit.Comment, error) {
	var cs []goreddit.Comment
	if err := s.Select(&cs, `SELECT * FROM comments WHERE post_id = $1 ORDER BY VOTES DESC`, postID); err != nil {
		return []goreddit.Comment{}, fmt.Errorf("error getting comments: %w", err)
	}
	return cs, nil
}

func (s *CommentStore) SaveComment(c *goreddit.Comment) error {
	if err := s.Get(c, `INSERT INTO comments VALUES ($1, $2, $3, $4) RETURNING *`,
		c.ID, c.PostID, c.Content, c.Votes); err != nil {
		return fmt.Errorf("error adding comment: %w", err)
	}
	return nil
}

func (s *CommentStore) UpdateComment(c *goreddit.Comment) error {
	fmt.Printf(">>> UpdateComment > %+v\n", c)
	if err := s.Get(c, `UPDATE comments SET post_id = $1, content = $2, votes = $3 WHERE id = $4 RETURNING *`,
		c.PostID, c.Content, c.Votes, c.ID); err != nil {
		return fmt.Errorf("error updating comment: %w", err)
	}
	return nil
}

func (s *CommentStore) DeleteComment(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM comments WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error deleting comment: %w", err)
	}
	return nil
}
