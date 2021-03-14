package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wdwiramadhan/bookhub-api/domain"
)

// mysqlProductRepository represent the connection database struct
type mysqlProductRepository struct {
	Conn *sql.DB
}

// NewMysqlProductRepository will create an object that represent the product.Repository interface
func NewMysqlProductRepository(Conn *sql.DB) domain.ProductRepository {
	return &mysqlProductRepository{Conn: Conn}
}

func (m *mysqlProductRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Product, err error) {
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
	result = make([]domain.Product, 0)
	for rows.Next() {
		t := domain.Product{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Price,
			&t.AuthorID,
			&t.Description,
			&t.Image,
			&t.UpdatedAt,
			&t.CreatedAt,
			&t.Author.ID,
			&t.Author.Name,
			&t.Author.DateOfBirth,
			&t.Author.UpdatedAt,
			&t.Author.CreatedAt,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (m *mysqlProductRepository) Fetch(ctx context.Context) (res []domain.Product, err error) {
	query := `SELECT * FROM product JOIN author ON product.author_id = author.id`
	res, err = m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}
	return
}

func (m *mysqlProductRepository) Store(ctx context.Context, p *domain.Product) (err error) {
	query := `INSERT INTO product VALUES(?,?,?,?,?,?,?,?)`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	_, err = stmt.ExecContext(ctx, p.ID, p.Name, p.Price, p.AuthorID, p.Description, p.Image, time.Now(), time.Now())
	if err != nil {
		return
	}
	return
}

func (m *mysqlProductRepository) GetByID(ctx context.Context, id int) (res domain.Product, err error) {
	query := `SELECT * FROM product JOIN author ON product.author_id = author.id WHERE product.id=?`
	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}
	return
}

func (m *mysqlProductRepository) Update(ctx context.Context, p *domain.Product, id int) (err error) {
	query := `UPDATE product SET name=?, price=?, author_id=?, description=?, updated_at=? WHERE id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	_, err = stmt.ExecContext(ctx, p.Name, p.Price, p.AuthorID, p.Description, time.Now(), id)
	if err != nil {
		return
	}
	return
}

func (m *mysqlProductRepository) Delete(ctx context.Context, id int) (err error) {
	query := `DELETE FROM product WHERE id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}
	return
}
