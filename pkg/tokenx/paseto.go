package tokenx

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
)

type PasetoImpl struct {
	paseto *paseto.V2
	symmetricKey string
	duration time.Duration
}

func NewPasetoImpl(
	symmetricKey string,
	duration time.Duration,
) (Token, error) {
	if len(symmetricKey) != 32 {
		return nil, fmt.Errorf("symmetric key length must be 32, got %d", len(symmetricKey))
	}
	return &PasetoImpl{
		paseto: paseto.NewV2(),
		symmetricKey: symmetricKey,
		duration: duration,
	}, nil
}

func (p *PasetoImpl) CreateToken(username string) (string, error) {
	payload, err := NewPayload(username, p.duration)
	if err != nil {
		return "", err
	}

	// Encrypt data
	return p.paseto.Encrypt([]byte(p.symmetricKey), payload, nil)
}

func (p *PasetoImpl) VerifyToken(tokenString string) (*Payload, error) {
	var newPayload Payload
	err := p.paseto.Decrypt(tokenString, []byte(p.symmetricKey), &newPayload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return &newPayload, nil
}
