package mysql

import (
	"time"
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
	"github.com/wdwiramadhan/bookhub-api/domain"
)

type mysqlProductRepository struct {
	Conn *sql.DB
}

func NewMysqlProductRepository(Conn *sql.DB) domain.ProductRepository {
	return &mysqlProductRepository{Conn}
}

func (m *mysqlProductRepository) fetch(ctx context.Context, query string)  (result []domain.Product, err error) {
	rows, err := m.Conn.QueryContext(ctx, query)
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
			&t.Id,
			&t.Name,
			&t.Price,
			&t.Author,
			&t.Description,
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

func (m *mysqlProductRepository) Fetch(ctx context.Context) (res []domain.Product, err error){
	query := `SELECT * FROM product`
	res, err = m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}
	return
}

func (m *mysqlProductRepository) Store(ctx context.Context, p *domain.Product) (err error) {
	query := `INSERT INTO product VALUES(?,?,?,?,?,?,?)`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	_,err = stmt.ExecContext(ctx, p.Id, p.Name, p.Price, p.Author, p.Description, time.Now(), time.Now())
	if err != nil {
		return
	}
	return
}

func (m *mysqlProductRepository) GetById(ctx context.Context, id string) (res domain.Product, err error){
	query := "SELECT * FROM product WHERE id="+id
	list, err := m.fetch(ctx, query)
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

func (m *mysqlProductRepository) Update(ctx context.Context, p *domain.Product, id string) (err error) {
	query := `UPDATE product SET name=?, price=?, author=?, description=?, updated_at=? WHERE id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	_,err = stmt.ExecContext(ctx, p.Name, p.Price, p.Author, p.Description, time.Now(), id)
	if err != nil {
		return
	}
	return
}

func (m *mysqlProductRepository) Delete(ctx context.Context, id string) (err error) {
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