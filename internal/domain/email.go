package domain

type VerifyEmail struct {
	Email      string `json:"email" gorm:"index;unique"`
	SecretCode int    `json:"secret_code"`
}
