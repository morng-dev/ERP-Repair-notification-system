package entities

type RegisterRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Firsname string `json:"first_name" validate:"required"`
	Lastname string `json:"last_name" validate:"required"`
	Avatar   string `json:"avatar,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user,omitempty"`
}

type ChangePasswordRequest struct {
	Oldpassword string `json:"old_password"`
	Newpassword string `json:"new_password"`
}

type RefreshTokenRequest struct {
	Refreshtoken string `json:"refresh_token" validate:"required"`
}
