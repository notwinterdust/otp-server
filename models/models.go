package models

type User struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type OTPAccount struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Issuer    string  `json:"issuer"`
	Secret    string  `json:"secret"`
	Algorithm string  `json:"algorithm"`
	Digits    int     `json:"digits"`
	Period    int     `json:"period"`
	Type      string  `json:"type"`
	Counter   *int    `json:"counter,omitempty"`
	CreatedAt string  `json:"createdAt"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type SyncRequest struct {
	Accounts []OTPAccount `json:"accounts"`
}

type SyncResponse struct {
	Accounts []OTPAccount `json:"accounts"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}
