package vercel_blob

import (
	"os"
)

// TokenProvider is a trait for providing a token to authenticate with the Vercel Blob Storage API.
//
// If your code is running inside a Vercel function then you will not need this.
//
// If your code is running outside of Vercel (e.g. a client application) then you will
// need to obtain a token from your Vercel application.  You can create a route
// to provide short-term tokens to authenticated users.  This trait allows you
// to connect to that route (or use some other method to obtain a token).
//
// The operation (e.g. list, put, download) and pathname (e.g. foo/bar.txt) are
// provided in case fine-grained authorization is required.  For operations that
// use the full URL (download / del) the pathname will be the URL.
type TokenProvider interface {
	GetToken(operation string, pathname string) (string, error)
}

func GetToken(provider TokenProvider, operation, pathname string) (string, error) {
	if provider != nil {
		return provider.GetToken(operation, pathname)
	}
	token := os.Getenv("BLOB_READ_WRITE_TOKEN")
	if token == "" {
		return "", ErrNotAuthenticated
	}
	return token, nil
}

// EnvTokenProvider is a token provider that reads the token from an environment variable.
//
// This is useful for testing but should not be used for real applications.
type EnvTokenProvider struct {
	token string
}

func (p *EnvTokenProvider) GetToken(operation, pathname string) (string, error) {
	if p.token != "" {
		return p.token, nil
	} else {
		envToken := os.Getenv("BLOB_READ_WRITE_TOKEN")
		if envToken != "" {
			return envToken, nil
		} else {
			return "", ErrNotAuthenticated
		}
	}
}

func NewEnvTokenProvider(envVar string) (*EnvTokenProvider, error) {
	token, exists := os.LookupEnv(envVar)
	if !exists {
		return nil, ErrNotAuthenticated
	}
	return &EnvTokenProvider{token}, nil
}
