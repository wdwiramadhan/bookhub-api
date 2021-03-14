package domain

import (
	"context"
	"time"
)

// Author ...
type Author struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DateOfBirth string    `json:"date_of_birth"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// AuthorUsecase represent the author's usecases
type AuthorUsecase interface {
	Fetch(c context.Context) ([]Author, error)
	Store(c context.Context, dataAuthor *Author) error
	GetAuthorById(c context.Context, authorId int) (Author, error)
	UpdateAuthorById(c context.Context, authorId int, dataAuthor *Author) error
	DeleteAuthorById(c context.Context, authorId int) error
}

// AuthorRepository represent the author's repository contract
type AuthorRepository interface {
	Fetch(ctx context.Context) ([]Author, error)
	Store(ctx context.Context, dataAuthor *Author) error
	GetAuthorById(ctx context.Context, authorId int) (Author, error)
	UpdateAuthorById(ctx context.Context, authorId int, dataAuthor *Author) error
	DeleteAuthorById(ctx context.Context, authorId int) error
}
