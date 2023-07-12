package wrap

import (
	"context"

	"github.com/hasura/go-graphql-client"
	"github.com/hasura/hge-go-gql-client/gql"
)

// WithContextHook wraps a given gql.Client, returning a gql.Client
// that will pass the context through given context hook before each
// query/mutation. A typical use would be to set certain headers
// using gql.WithHeaders.
func WithContextHook(cl gql.Client, hook func(context.Context) (context.Context, error)) gql.Client {
	return &client{
		hook: hook,
		cl:   cl,
	}
}

type client struct {
	hook func(context.Context) (context.Context, error)
	cl   gql.Client
}

func (c *client) Query(ctx context.Context, q interface{}, variables map[string]interface{}, options ...graphql.Option) error {
	headerCtx, err := c.hook(ctx)
	if err != nil {
		return err
	}
	return c.cl.Query(headerCtx, q, variables)
}

func (c *client) NamedQuery(ctx context.Context, name string, q interface{}, variables map[string]interface{}, options ...graphql.Option) error {
	headerCtx, err := c.hook(ctx)
	if err != nil {
		return err
	}
	return c.cl.NamedQuery(headerCtx, name, q, variables)
}

func (c *client) NamedQueryRaw(ctx context.Context, name string, q interface{}, variables map[string]interface{}, options ...graphql.Option) ([]byte, error) {
	headerCtx, err := c.hook(ctx)
	if err != nil {
		return nil, err
	}
	return c.cl.NamedQueryRaw(headerCtx, name, q, variables)
}

func (c *client) Mutate(ctx context.Context, m interface{}, variables map[string]interface{}, options ...graphql.Option) error {
	headerCtx, err := c.hook(ctx)
	if err != nil {
		return err
	}
	return c.cl.Mutate(headerCtx, m, variables)
}

func (c *client) NamedMutate(ctx context.Context, name string, m interface{}, variables map[string]interface{}, options ...graphql.Option) error {
	headerCtx, err := c.hook(ctx)
	if err != nil {
		return err
	}
	return c.cl.NamedMutate(headerCtx, name, m, variables)
}

func (c *client) NamedMutateRaw(ctx context.Context, name string, m interface{}, variables map[string]interface{}, options ...graphql.Option) ([]byte, error) {
	headerCtx, err := c.hook(ctx)
	if err != nil {
		return nil, err
	}
	return c.cl.NamedMutateRaw(headerCtx, name, m, variables)
}

// assert that *client implements gql.Client
var _ gql.Client = &client{}
