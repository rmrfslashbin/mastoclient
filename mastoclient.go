package mastoclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/mattn/go-mastodon"
	"github.com/rs/zerolog"
)

// Options for the weather query
type Option func(c *Config)

// Config for the weather query
type Config struct {
	log          *zerolog.Logger
	instance     string
	clientKey    string
	clientSecret string
	accessToken  string
}

// NewConfig creates a new Config
func New(opts ...Option) (*Config, error) {
	c := &Config{}

	// apply the list of options to Config
	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// WithToken sets the token to use
func WithAccessToken(accessToken string) Option {
	return func(c *Config) {
		c.accessToken = accessToken
	}
}

// WithClientID sets the client ID to use
func WithClientkey(clientKey string) Option {
	return func(c *Config) {
		c.clientKey = clientKey
	}
}

// WithClientSecret sets the client secret to use
func WithClientSecret(clientSecret string) Option {
	return func(c *Config) {
		c.clientSecret = clientSecret
	}
}

// WithInstance sets the instance to use
func WithInstance(instance string) Option {
	return func(c *Config) {
		c.instance = instance
	}
}

// WithLogger sets the logger to use
func WithLogger(log *zerolog.Logger) Option {
	return func(c *Config) {
		c.log = log
	}
}

// SetAccessToken sets the access token
func (c *Config) SetAccessToken(accessToken string) {
	c.accessToken = accessToken
}

// SetClientKey sets the client key
func (c *Config) SetClientKey(clientKey string) {
	c.clientKey = clientKey
}

// SetClientSecret sets the client secret
func (c *Config) SetClientSecret(clientSecret string) {
	c.clientSecret = clientSecret
}

// SetInstance sets the instance
func (c *Config) SetInstance(instance string) {
	c.instance = instance
}

// SetLogger sets the logger
func (c *Config) SetLogger(log *zerolog.Logger) {
	c.log = log
}

// prefight checks if the config is set up correctly and returns a mastodon client
func (c *Config) preflight() (*mastodon.Client, error) {
	// Check set up
	if c.instance == "" {
		return nil, &NoInstanceError{}
	}

	if c.clientKey == "" {
		return nil, &NoClientKeyError{}
	}

	if c.clientSecret == "" {
		return nil, &NoClientSecretError{}
	}

	if c.accessToken == "" {
		return nil, &NoAccessTokenError{}
	}

	// Set up Mastodon client
	client := mastodon.NewClient(&mastodon.Config{
		Server:       c.instance,
		ClientID:     c.clientKey,
		ClientSecret: c.clientSecret,
		AccessToken:  c.accessToken,
	})

	return client, nil
}

// GetUserByID gets a user by ID
func (c *Config) GetUserByID(id string) (*mastodon.Account, error) {
	client, err := c.preflight()
	if err != nil {
		return nil, err
	}
	// Get user
	user, err := client.GetAccount(context.Background(), mastodon.ID(id))
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Me gets the current user
func (c *Config) Me() (*mastodon.Account, error) {
	client, err := c.preflight()
	if err != nil {
		return nil, err
	}
	// Get user
	user, err := client.GetAccountCurrentUser(context.Background())
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Post a toot
func (c *Config) Post(toot *mastodon.Toot) (*mastodon.ID, error) {
	client, err := c.preflight()
	if err != nil {
		return nil, err
	}

	// Post the toot
	if status, err := client.PostStatus(context.Background(), toot); err != nil {
		fmt.Println(err)
		return nil, err
	} else {
		return &status.ID, nil
	}
}

func RegisterApp(input *RegisterAppInput) (*mastodon.Application, error) {
	app, err := mastodon.RegisterApp(context.Background(), &mastodon.AppConfig{
		Server:       input.InstanceURL,
		ClientName:   input.ClientName,
		RedirectURIs: input.RedirectURI,
		Scopes:       strings.Join(input.Scopes, " "),
		Website:      input.Website,
	})
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (c *Config) GetAuthTokenFromCode(authCode *string, redirectURI *string) (*string, error) {
	client, err := c.preflight()
	if err != nil {
		return nil, err
	}

	if err = client.AuthenticateToken(context.Background(), *authCode, *redirectURI); err != nil {
		return nil, err
	}

	return &client.Config.AccessToken, nil
}
