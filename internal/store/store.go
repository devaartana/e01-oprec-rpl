package store

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	DB                   = "link-shortener"
	Collection           = "data"
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicateUsername = errors.New("username already exists")
	ErrNotFound          = errors.New("user not found")
	ErrDuplicateSlug     = errors.New("slug already exists")
	QueryTimeoutDuration = 5 * time.Second
)

type Storage struct {
	Users interface {
		GetAllUsers(ctx context.Context) ([]User, error)
		Create(ctx context.Context, user *User) error
		Update(ctx context.Context, user *User) error
		GetByEmail(ctx context.Context, email string) (*User, error)
		DeleteByEmail(ctx context.Context, email string) error 
	}

	Links interface {
		Create(ctx context.Context, email string, link *Link) error
		GetBySlug(ctx context.Context, slug string) (*Link, error)
		GetAll(ctx context.Context, email string) ([]Link, error)
		DeleteBySlug(ctx context.Context, email string, slug string) error
		UpdateBySlug(ctx context.Context, email string, link *Link) error
	}
}

func NewStorage(db *mongo.Client) Storage {
	return Storage{
		Users: &UserStore{db},
		Links: &LinkStore{db},
	}
}

