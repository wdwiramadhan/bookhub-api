package domain

import (
	"context"
	"time"
)

// Product ...
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Price       int64     `json:"price"`
	AuthorID    int       `json:"author_id"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
	Author      Author    `json:"author"`
}

// ProductUseCase represent the product's usecases
type ProductUseCase interface {
	Fetch(ctx context.Context) ([]Product, error)
	Store(context.Context, *Product) error
	GetByID(ctx context.Context, id int) (Product, error)
	Update(ctx context.Context, ar *Product, id int) error
	Delete(ctx context.Context, id int) error
}

// ProductRepository represent the product's repository contract
type ProductRepository interface {
	Fetch(ctx context.Context) ([]Product, error)
	Store(ctx context.Context, a *Product) error
	GetByID(ctx context.Context, id int) (Product, error)
	Update(ctx context.Context, ar *Product, id int) error
	Delete(ctx context.Context, id int) error
}
