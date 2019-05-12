package user

import (
	"time"
)

// User represents the user entity.
type User struct {
	ID                  string     `json:"id"`
	Email               string     `json:"email,omitempty" validate:"omitempty,email"`
	EmailVerified       bool       `json:"email_verified,omitempty"`
	HashedPassword      string     `json:"-"`
	PhoneNumber         string     `json:"phone_number,omitempty"`
	PhoneNumberVerified bool       `json:"phone_number_verified"`
	Name                string     `json:"name,omitempty"`
	FamilyName          string     `json:"family_name,omitempty"`
	GivenName           string     `json:"given_name,omitempty"`
	MiddleName          string     `json:"middle_name,omitempty"`
	Nickname            string     `json:"nickname,omitempty"`
	PreferredUsername   string     `json:"preferred_username,omitempty"`
	Profile             string     `json:"profile,omitempty"`
	Picture             string     `json:"picture,omitempty"`
	Website             string     `json:"website,omitempty" validate:"omitempty,url"`
	Gender              string     `json:"gender,omitempty"`
	BirthDate           BirthDate  `json:"birth_date,omitempty"`
	ZoneInfo            string     `json:"zone_info,omitempty"`
	Locale              string     `json:"locale,omitempty"`
	StreetAddress       string     `json:"street_address,omitempty"`
	Locality            string     `json:"locality,omitempty"`
	Region              string     `json:"region,omitempty"`
	PostalCode          string     `json:"postal_code,omitempty"`
	Country             string     `json:"country,omitempty"`
	CreatedAt           *time.Time `json:"created_at,omitempty"`
	UpdatedAt           *time.Time `json:"updated_at,omitempty"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
}

func NewUser(id string) User {
	return User{ID: id}
}
