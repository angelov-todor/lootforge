package auth

import (
	"context"

	firebaseAuth "firebase.google.com/go/v4/auth"
)

type AuthUser struct {
	UID         string
	Email       string
	DisplayName string
	PhotoURL    string
}

type TokenVerifier interface {
	VerifyToken(ctx context.Context, idToken string) (*AuthUser, error)
}

type FirebaseTokenVerifier struct {
	client *firebaseAuth.Client
}

func NewFirebaseTokenVerifier(client *firebaseAuth.Client) *FirebaseTokenVerifier {
	return &FirebaseTokenVerifier{client: client}
}

func (v *FirebaseTokenVerifier) VerifyToken(ctx context.Context, idToken string) (*AuthUser, error) {
	token, err := v.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}
	return &AuthUser{
		UID:         token.UID,
		Email:       stringClaim(token.Claims, "email"),
		DisplayName: stringClaim(token.Claims, "name"),
		PhotoURL:    stringClaim(token.Claims, "picture"),
	}, nil
}

func stringClaim(claims map[string]interface{}, key string) string {
	if v, ok := claims[key].(string); ok {
		return v
	}
	return ""
}
