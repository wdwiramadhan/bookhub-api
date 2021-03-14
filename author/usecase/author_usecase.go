package usecase

import (
	"context"
	"time"

	"github.com/wdwiramadhan/bookhub-api/domain"
)

// AuthorUsecase represent the author use case struct
type AuthorUsecase struct {
	authorRepo     domain.AuthorRepository
	contextTimeout time.Duration
}

// NewAuthorUsecase will create new an author usecase object representation of domain.AuthorUsecase interface
func NewAuthorUsecase(a domain.AuthorRepository, timeout time.Duration) domain.AuthorUsecase {
	return &AuthorUsecase{
		authorRepo:     a,
		contextTimeout: timeout,
	}
}

// Fetch will get authors
func (a *AuthorUsecase) Fetch(c context.Context) (res []domain.Author, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err = a.authorRepo.Fetch(ctx)
	if err != nil {
		return nil, err
	}
	return
}

// Store will create new author
func (a *AuthorUsecase) Store(c context.Context, dataAuthor *domain.Author) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	err = a.authorRepo.Store(ctx, dataAuthor)
	if err != nil {
		return err
	}
	return
}

func (a *AuthorUsecase) GetAuthorById(c context.Context, authorId int) (res domain.Author, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err = a.authorRepo.GetAuthorById(ctx, authorId)
	if err != nil {
		return
	}
	return
}

func (a *AuthorUsecase) UpdateAuthorById(c context.Context, authorId int, dataAuthor *domain.Author) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	err = a.authorRepo.UpdateAuthorById(ctx, authorId, dataAuthor)
	if err != nil {
		return
	}
	return
}

func (a *AuthorUsecase) DeleteAuthorById(c context.Context, authorId int) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	err = a.authorRepo.DeleteAuthorById(ctx, authorId)
	if err != nil {
		return
	}
	return
}
