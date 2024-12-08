package domain

type VerifyEmail struct {
	Email      string `json:"email" gorm:"index;unique"` // The email of the user, unique.
	SecretCode int    `json:"secret_code"`               // An optional verification or identification code.
}
