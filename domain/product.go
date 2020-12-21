package domain

import (
	"context"
	"time"
)

type Product struct{
	Id string `json:"id"`
	Name string `json:"name"`
	Price int64 `json:"price"`
	Author string `json:"author"`
	Description string `json:"description"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type ProductUseCase interface{
	Fetch(ctx context.Context) ([]Product, error)
	Store(context.Context, *Product) error
	GetById(ctx context.Context, id string) (Product, error)
}

type ProductRepository interface{
	Fetch(ctx context.Context) (res []Product, err error)
	Store(ctx context.Context, a *Product) error
	GetById(ctx context.Context, id string) (Product, error)
}