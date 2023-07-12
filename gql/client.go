package gql

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	gogql "github.com/hasura/go-graphql-client"
	untyped "github.com/shahidhk/gql"
)

const (
	HasuraClientName  = "Hasura-Client-Name"
	DefaultClientName = "lux-hasura-api"

	XHasuraRole        = "x-hasura-role"
	XHasuraUserID      = "x-hasura-user-id"
	XHasuraUserEmail   = "x-hasura-user-email"
	XHasuraAdminSecret = "X-Hasura-Admin-Secret"
)

type options struct {
	timeout time.Duration
}

var defaultOptions = options{
	timeout: 30 * time.Second,
}

type Option func(*options)

func WithTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.timeout = timeout
	}
}

// sudoFunc is a function to return derive a new ActorAwareClient with admin privileges
// from an existing ActorAwareClient.
//
// it's declared as a type because we want to make this function a [strategy pattern](https://wiki.c2.com/?StrategyPattern)
// some clients will use it to effectively promote themselves to Admin, while others, non-promotable
// will return an error.
type sudoFunc func(actor *Actor) (*ActorAwareClient, error)
type untypedFunc func() *untyped.Client

// ActorAwareClient is a graphql client which requests are made on behalf of an Actor
type ActorAwareClient struct {
	Client
	Actor *Actor
	// A function to derive an admin (or impersonate another user) client from this one. It will return an error if a
	// client cannot be promoted to an admin, see NewPromotableClient
	sudoFunc sudoFunc
	// A function to derive an untyped client from this one. This field is unexported, and production code won't
	// have any method to get the function or invoke it. Only when the integration build tag is used to compile the
	// code, the code will compile a method to return the untyped client.
	untypedFunc untypedFunc
}

// NewAdminClientFromHost is a helper constructor, that will append to the host, the default endpoint path
// to build a new admin client.
func NewAdminClientFromHost(host string, adminSecret string, clientName string, options ...Option) *ActorAwareClient {
	endpoint := fmt.Sprintf("%s/v1/graphql", host)
	return NewAdminClient(endpoint, adminSecret, clientName, options...)
}

// NewAdminClient creates a new client to access the given graphql endpoint using the given admin
// secret.
//
// options is a set of functional options used, among other things, to override the default timeout for
// the client.
func NewAdminClient(endpoint, adminSecret, clientName string, options ...Option) *ActorAwareClient {
	client := NewPromotableClient(endpoint, adminSecret, NewAdminActor(), clientName, options...)
	return client
}

// NewPromotableClient creates a new client to access the given graphql endpoint as the given actor, but that has the option
// to act as an admin, by calling ForceAdmin() over it. This is useful when a certain action needs to request
// superpowers, think of it as a sudoer in Linux.
//
// adminSecret is used for client authentication not authorization. Authorization is done when HGE checks the role
// and user headers. Think of adminSecret as an API key to access the endpoint, that's why promoting a client to an
// admin, only means changing the role to admin.
func NewPromotableClient(endpoint string, adminSecret string, actor *Actor, clientName string, options ...Option) *ActorAwareClient {
	client := NewClient(endpoint, adminSecret, actor, clientName, options...)

	client.sudoFunc = func(impersonated *Actor) (*ActorAwareClient, error) {
		return NewPromotableClient(endpoint, adminSecret, impersonated, clientName, options...), nil
	}

	return client
}

// NewClient creates a new client to access the given graphql endpoint as a user, setting the given http headers
// in the underlying httpClient.
func NewClient(endpoint string, adminSecret string, actor *Actor, clientName string, options ...Option) *ActorAwareClient {
	opts := defaultOptions
	for _, apply := range options {
		apply(&opts)
	}

	headers := HeadersFor(actor, adminSecret, clientName)
	httpClient := buildClient(headers)
	httpClient.Timeout = opts.timeout

	return &ActorAwareClient{
		Client: gogql.NewClient(endpoint, httpClient),
		Actor:  actor,
		sudoFunc: func(actor *Actor) (*ActorAwareClient, error) {
			return nil, errors.New("by default an actor aware client cannot impersonate another user")
		},
		untypedFunc: func() *untyped.Client {
			return untyped.NewClient(endpoint, headers)
		},
	}
}

// ForceAdmin allows the client to act on behalf of an admin, this function panics if the client cannot
// be promoted to an Admin client. Prefer AsAdmin instead.
func (c *ActorAwareClient) ForceAdmin() *ActorAwareClient {
	admin, err := c.AsAdmin()
	if err != nil {
		log.Panicf("Client (role=%s) cannot be promoted to admin. %s", c.Actor.Role, err)
	}
	return admin
}

// AsAdmin allows the client to act on behalf of an admin, this function returns an error in case
// the client is not promotable
func (c *ActorAwareClient) AsAdmin() (*ActorAwareClient, error) {
	return c.sudoFunc(c.Actor.AsAdmin())
}

func HeadersFor(actor *Actor, adminSecret, clientName string) map[string]string {
	headers := make(map[string]string)
	if actor != nil {
		var userId string
		if actor.UserID != nil {
			userId = actor.UserID.String()
			headers[XHasuraUserID] = userId
		}
		headers[XHasuraRole] = actor.Role
		headers[XHasuraUserEmail] = actor.Email
	}
	headers[HasuraClientName] = clientName
	headers[XHasuraAdminSecret] = adminSecret
	return headers
}

// buildClient creates a new HTTP transport that uses a RoundTripper to set the headers
// on every request.
// Headers come in two forms:
// * From the actor information and admin secret provided in ActorAwareClient initialization
// * Previously set in the context object
func buildClient(headers map[string]string) *http.Client {
	return &http.Client{
		Transport: headerRoundTripper{
			setHeaders: func(req *http.Request) {
				// we set the headers the client was configured with
				for hn, hv := range headers {
					req.Header.Set(hn, hv)
				}
				// and then those that might have been set in the context
				for hn, hv := range GetHeadersFromContext(req.Context()) {
					req.Header.Set(hn, hv)
				}
			},
			rt: http.DefaultTransport},
	}
}

type headerRoundTripper struct {
	setHeaders func(req *http.Request)
	rt         http.RoundTripper
}

func (h headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	h.setHeaders(req)
	return h.rt.RoundTrip(req)
}
