package arpcmocks

import (
	"context"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

type IDTokenStub struct {
	oauth2.TokenSource
}

var DefaultToken = &oauth2.Token{
	AccessToken:  "secret-access-token",
	TokenType:    "custom-token-type",
	RefreshToken: "refresh-token",
	Expiry:       time.Now().Add(time.Hour),
	ExpiresIn:    999999999,
}

func (i *IDTokenStub) Token() (*oauth2.Token, error) {
	return DefaultToken, nil
}

func TokenSource(
	token oauth2.TokenSource,
) func(ctx context.Context, audience string, opts ...idtoken.ClientOption) (oauth2.TokenSource, error) {
	return func(_ context.Context, _ string, _ ...idtoken.ClientOption) (oauth2.TokenSource, error) {
		return token, nil
	}
}
