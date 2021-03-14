package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wdwiramadhan/bookhub-api/domain"
)

// mysqlAuthorRepository represent the connection database struct
type mysqlAuthorRepository struct {
	Conn *sql.DB
}

// NewMysqlAuthorRepository will create an object that represent the author.Repository interface
func NewMysqlAuthorRepository(Conn *sql.DB) domain.AuthorRepository {
	return &mysqlAuthorRepository{Conn: Conn}
}

func (m *mysqlAuthorRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Author, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()
	result = make([]domain.Author, 0)
	for rows.Next() {
		t := domain.Author{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.DateOfBirth,
			&t.UpdatedAt,
			&t.CreatedAt,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (m *mysqlAuthorRepository) Fetch(ctx context.Context) (res []domain.Author, err error) {
	query := `SELECT * FROM author`
	res, err = m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}
	return
}

func (m *mysqlAuthorRepository) Store(ctx context.Context, dataAuthor *domain.Author) (err error) {
	query := `INSERT INTO author (name, date_of_birth, updated_at, created_at) VALUES(?,?,?,?)`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	fmt.Println(err)
	if err != nil {
		return err
	}
	dateOfBirth, _ := time.Parse("2006-01-02", dataAuthor.DateOfBirth)
	_, err = stmt.ExecContext(ctx, dataAuthor.Name, dateOfBirth, time.Now(), time.Now())
	if err != nil {
		return
	}
	return
}

func (m *mysqlAuthorRepository) GetAuthorById(ctx context.Context, authorId int) (res domain.Author, err error) {
	query := `SELECT * FROM author where id=?`
	author, err := m.fetch(ctx, query, authorId)
	if err != nil {
		return
	}
	if len(author) > 0 {
		res = author[0]
	} else {
		return res, domain.ErrNotFound
	}
	return
}

func (m *mysqlAuthorRepository) UpdateAuthorById(ctx context.Context, authorId int, dataAuthor *domain.Author) (err error) {
	query := `UPDATE author SET name=?, date_of_birth=?, updated_at=? WHERE id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	_, err = stmt.ExecContext(ctx, dataAuthor.Name, dataAuthor.DateOfBirth, time.Now(), authorId)
	if err != nil {
		return
	}
	return
}

func (m *mysqlAuthorRepository) DeleteAuthorById(ctx context.Context, authorId int) (err error) {
	query := `DElETE FROM author WHERE id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	_, err = stmt.ExecContext(ctx, authorId)
	if err != nil {
		return
	}
	return
}
