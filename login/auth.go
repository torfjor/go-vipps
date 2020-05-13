package login

import (
	"context"
	"errors"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

const IssuerURLTesting = "https://apitest.vipps.no/access-management-1.0/access/"
const IssuerURL = "https://api.vipps.no/access-management-1.0/access/"

// Provider is a convenience wrapper around oidc.Provider tailored to the Vipps
// Login API
type Provider struct {
	provider    *oidc.Provider
	oauthConfig oauth2.Config
	verifier    *oidc.IDTokenVerifier
}

// ClientConfig represents a configuration for a Provider
type ClientConfig struct {
	ClientID     string
	ClientSecret string
	IssuerURL    string
	RedirectURL  string
	Scopes       []string
}

// Claims represents the claims contained in Vipps ID tokens.
type Claims struct {
	Address []struct {
		Country   string `json:"country"`
		Street    string `json:"street_address"`
		Type      string `json:"address_type"`
		Formatted string `json:"formatted"`
		Zip       string `json:"postal_code"`
		Region    string `json:"region"`
	} `json:"address"`
	Phone      string `json:"phone_number"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
	UserID     string `json:"sub"`
}

// NewLoginClient returns a configured Vipps Login Provider.
func NewLoginClient(ctx context.Context, config *ClientConfig) (*Provider, error) {
	if config.IssuerURL == "" {
		config.IssuerURL = IssuerURLTesting
	}
	provider, err := oidc.NewProvider(ctx, config.IssuerURL)
	if err != nil {
		return nil, err
	}

	oauthConfig := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  config.RedirectURL,
		Scopes:       append([]string{oidc.ScopeOpenID}, config.Scopes...),
	}

	return &Provider{
		provider:    provider,
		oauthConfig: oauthConfig,
		verifier: provider.Verifier(&oidc.Config{
			ClientID: config.ClientID,
		}),
	}, err
}

// AuthCodeURL returns a URL to OAuth 2.0 provider's consent page that asks for
// permissions for the configured scopes explicitly.
func (c *Provider) AuthCodeURL(state string) string {
	return c.oauthConfig.AuthCodeURL(state)
}

// ExchangeCodeForClaims takes an oauth2 authorization code, exchanges it for a
// token, and returns the contained ID token's claims, if any
func (c *Provider) ExchangeCodeForClaims(ctx context.Context, code string) (*Claims, error) {
	token, err := c.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err

	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("oauth2: nb id_token in response")
	}

	idToken, err := c.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}

	claims := Claims{
		UserID: idToken.Subject,
	}

	if err = idToken.Claims(&claims); err != nil {
		return nil, err
	}

	return &claims, nil
}
