package user

import "fmt"

var (
	RoleAdmin = Role{"ADMIN"}
	RoleUser  = Role{"USER"}
)

var roles = map[string]Role{
	RoleAdmin.name: RoleAdmin,
	RoleUser.name:  RoleUser,
}

type Role struct {
	name string
}

func ParseRole(value string) (Role, error) {
	role, exists := roles[value]
	if !exists {
		return Role{}, fmt.Errorf("invalid role %q", value)
	}

	return role, nil
}

func MustParseRole(value string) Role {
	role, err := ParseRole(value)
	if err != nil {
		panic(err)
	}

	return role
}

func (r Role) Name() string {
	return r.name
}

func (r *Role) UnmarshalText(data []byte) error {
	role, err := ParseRole(string(data))
	if err != nil {
		return err
	}

	r.name = role.name
	return nil
}

func (r Role) MarshalText() ([]byte, error) {
	return []byte(r.name), nil
}

func (r Role) Equal(r2 Role) bool {
	return r.name == r2.name
}
