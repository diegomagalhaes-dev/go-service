package user

import (
	"context"
	"net/mail"

	"github.com/diegomagalhaes-dev/go-service/business/data/order"
	"github.com/diegomagalhaes-dev/go-service/foundation/logger"
	"github.com/google/uuid"
)

type Storer interface {
	Create(ctx context.Context, usr User) error
	Update(ctx context.Context, usr User) error
	Delete(ctx context.Context, usr User) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]User, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	QueryByIDs(ctx context.Context, userID []uuid.UUID) ([]User, error)
	QueryByEmail(ctx context.Context, email mail.Address) (User, error)
}

type Core struct {
	storer Storer
	log    *logger.Logger
}

func NewCore(log *logger.Logger, storer Storer) *Core {
	return &Core{
		storer: storer,
		log:    log,
	}
}

func (c *Core) Create(ctx context.Context, nu NewUser) (User, error) {

	return User{}, nil
}
