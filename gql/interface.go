package gql

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// Client abstracts the interface provided by hasura/go-graphql-client so their implementation can be replaced
// by something else
type Client interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}, options ...graphql.Option) error
	NamedQuery(ctx context.Context, name string, q interface{}, variables map[string]interface{}, options ...graphql.Option) error
	NamedQueryRaw(ctx context.Context, name string, q interface{}, variables map[string]interface{}, options ...graphql.Option) ([]byte, error)
	Mutate(ctx context.Context, m interface{}, variables map[string]interface{}, options ...graphql.Option) error
	NamedMutate(ctx context.Context, name string, m interface{}, variables map[string]interface{}, options ...graphql.Option) error
	NamedMutateRaw(ctx context.Context, name string, m interface{}, variables map[string]interface{}, options ...graphql.Option) ([]byte, error)
}
