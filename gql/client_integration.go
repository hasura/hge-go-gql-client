//go:build integration
// +build integration

package gql

import (
	untyped "github.com/shahidhk/gql"
)

// As allows to impersonate a certain actor in tests
func (c *ActorAwareClient) As(actor *Actor) *ActorAwareClient {
	client, err := c.sudoFunc(actor)
	if err != nil {
		panic(err)
	}
	return client
}

// Untyped is only present in tests, allows us to return an untyped client derived from an ActorAwareClient.
// An Untyped client is less safe, but also less verbose, so we can use it to perform queries and mutations in
// tests more easily.
func (c *ActorAwareClient) Untyped() *untyped.Client {
	return c.untypedFunc()
}
