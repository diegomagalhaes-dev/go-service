package user

import (
	"net/mail"
	"time"

	"github.com/diegomagalhaes-dev/go-service/business/core/event"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Name         string
	Email        mail.Address
	Roles        []Role
	PasswordHash []byte
	Department   string
	Enabled      bool
	DateCreated  time.Time
	DateUpdated  time.Time
}

type NewUser struct {
	Name            string
	Email           mail.Address
	Roles           []Role
	Department      string
	Password        string
	PasswordConfirm string
}

type UpdateUser struct {
	Name            *string
	Email           *mail.Address
	Roles           []Role
	Department      *string
	Password        *string
	PasswordConfirm *string
	Enabled         *bool
}

// UpdatedEvent constructs an event for when a user is updated.
func (uu UpdateUser) UpdatedEvent(userID uuid.UUID) event.Event {
	params := EventParamsUpdated{
		UserID: userID,
		UpdateUser: UpdateUser{
			Enabled: uu.Enabled,
		},
	}

	rawParams, err := params.Marshal()
	if err != nil {
		panic(err)
	}

	return event.Event{
		Source:    EventSource,
		Type:      EventUpdated,
		RawParams: rawParams,
	}
}
