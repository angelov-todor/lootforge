package auth

import "context"

type MockTokenVerifier struct {
	User *AuthUser
	Err  error
}

func (m *MockTokenVerifier) VerifyToken(ctx context.Context, idToken string) (*AuthUser, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.User, nil
}
