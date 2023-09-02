package passwordx

import "golang.org/x/crypto/bcrypt"

type Password interface {
	HashPassword(password string) (string, error)
	CheckPassword(password, hashedPassword string) error
}

type PasswordImpl struct {
	cost int
}

func NewPassword(
	cost int,
) Password {
	return &PasswordImpl{
		cost: cost,
	}
}

func (p *PasswordImpl) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (p *PasswordImpl) CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
