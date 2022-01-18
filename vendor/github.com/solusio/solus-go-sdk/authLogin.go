package solus

// AuthLoginRequest represents all required properties for authentication.
type AuthLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthLoginResponse represents an authentication response.
type AuthLoginResponse struct {
	Credentials Credentials `json:"credentials"`
}

// Credentials represents obtained credentials.
type Credentials struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   string `json:"expires_at"`
}
