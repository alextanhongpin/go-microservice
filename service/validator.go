package service

import validator "gopkg.in/go-playground/validator.v9"

// Validate performs validation scoped for the service package. Validation is
// part of the business rules, hence there is an exception to make it global.
var Validate func(interface{}) error

func init() {
	Validate = validator.New().Struct
}
