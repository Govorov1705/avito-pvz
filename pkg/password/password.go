package password

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	passwordBytes := []byte(password)

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, 12)
	if err != nil {
		return "", err
	}

	return string(hashedPasswordBytes), nil
}
