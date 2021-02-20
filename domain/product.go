package domain

import (
	"context"
	"time"
)

// Product ...
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Price       int64     `json:"price"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// ProductUseCase represent the product's usecases
type ProductUseCase interface {
	Fetch(ctx context.Context) ([]Product, error)
	Store(context.Context, *Product) error
	GetByID(ctx context.Context, id string) (Product, error)
	Update(ctx context.Context, ar *Product, id string) error
	Delete(ctx context.Context, id string) error
}

// ProductRepository represent the product's repository contract
type ProductRepository interface {
	Fetch(ctx context.Context) ([]Product, error)
	Store(ctx context.Context, a *Product) error
	GetByID(ctx context.Context, id string) (Product, error)
	Update(ctx context.Context, ar *Product, id string) error
	Delete(ctx context.Context, id string) error
}
