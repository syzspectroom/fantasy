package model

// LoginInput input struct
type LoginInput struct {
	Email    string `mod:"trim,lcase" json:"email"  validate:"required,email"`
	Password string `json:"password"  validate:"required"`
	Meta     *LoginMeta
}

// LoginResponse input struct
type LoginResponse struct {
	Token   string `json:"token"`
	Refresh string `json:"refresh_token"`
}

// LoginMeta input struct
type LoginMeta struct {
	IP        string
	UserAgent string
	UserID    string
}

// RefreshInputStruct Refresh query input
type RefreshInputStruct struct {
	RefreshToken string `json:"refresh_token"`
	Meta         *LoginMeta
}
