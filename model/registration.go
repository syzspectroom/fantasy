package model

// Registration input struct
type Registration struct {
	Email           string `mod:"trim,lcase" json:"email"  validate:"required,email"`
	Password        string `json:"password"  validate:"required,min=4,eqfield=ConfirmPassword"`
	ConfirmPassword string `json:"password_confirmation"`
}
