package jwt

import "fmt"

type ErrUnexpectedSigningMethod struct {
	Algorithm any
}

func (err ErrUnexpectedSigningMethod) Error() string {
	return fmt.Sprintf("unexpected signing method: %v", err.Algorithm)
}

type ErrFailedToParseToken struct {
	Reason string
}

func (err ErrFailedToParseToken) Error() string {
	return fmt.Sprintf("failed to parse token: %s", err.Reason)
}

type ErrInvalidToken struct {
	Reason string
}

func (err ErrInvalidToken) Error() string {
	return fmt.Sprintf("invalid token: %s", err.Reason)
}
