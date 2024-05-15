package userdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/diegomagalhaes-dev/go-service/business/core/user"
	db "github.com/diegomagalhaes-dev/go-service/business/data/dbsql/pgx"
	"github.com/diegomagalhaes-dev/go-service/foundation/logger"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	log *logger.Logger
	db  *sqlx.DB
}

func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

func (s *Store) Create(ctx context.Context, usr user.User) error {
	const q = `
	INSERT INTO users
		(user_id, name, email, password_hash, roles, enabled, department, date_created, date_updated)
	VALUES
		(:user_id, :name, :email, :password_hash, :roles, :enabled, :department, :date_created, :date_updated)`

	if err := db.NamedExecContext(ctx, s.log, s.db, q, toDBUser(usr)); err != nil {
		if errors.Is(err, db.ErrDBDuplicatedEntry) {
			return fmt.Errorf("namedexeccontext: %w", user.ErrUniqueEmail)
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}
