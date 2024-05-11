package user

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/diegomagalhaes-dev/go-service/foundation/validate"
	"github.com/google/uuid"
)

type QueryFilter struct {
	ID               *uuid.UUID    `validate:"omitempty"`
	Name             *string       `validate:"omitempty,min=3"`
	Email            *mail.Address `validate:"omitempty"`
	StartCreatedDate *time.Time    `validate:"omitempty"`
	EndCreatedDate   *time.Time    `validate:"omitempty"`
}

func (qf *QueryFilter) Validate() error {
	if err := validate.Check(qf); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

func (qf *QueryFilter) WithUserID(userID uuid.UUID) {
	qf.ID = &userID
}

func (qf *QueryFilter) WithName(name string) {
	qf.Name = &name
}

func (qf *QueryFilter) WithEmail(email mail.Address) {
	qf.Email = &email
}

func (qf *QueryFilter) WithStartDateCreated(startDate time.Time) {
	d := startDate.UTC()
	qf.StartCreatedDate = &d
}

func (qf *QueryFilter) WithEndCreatedDate(endDate time.Time) {
	d := endDate.UTC()
	qf.EndCreatedDate = &d
}
