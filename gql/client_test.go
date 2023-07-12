package gql

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func TestHasuraHeaders(t *testing.T) {
	var h http.Header
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h = r.Header
	}))
	defer ts.Close()

	sampleUUID := uuid.New()

	var query struct {
		Thing struct {
			Field int `graphql:"field"`
		} `graphql:"thing"`
	}

	for _, tC := range []struct {
		desc            string
		cl              Client
		expectedHeaders map[string]string
	}{
		{
			"For admin clients",
			NewAdminClientFromHost(ts.URL, "admin-secret", "test-client"),
			map[string]string{
				"hasura-client-name":    "test-client",
				"x-hasura-admin-secret": "admin-secret",
			},
		},
		{
			"For promotable clients",
			NewPromotableClient(ts.URL, "admin-secret", NewUserActor(&sampleUUID, "foo@bar.baz"), "test-client"),
			map[string]string{
				"hasura-client-name":    "test-client",
				"x-hasura-user-id":      sampleUUID.String(),
				"x-hasura-role":         "user",
				"x-hasura-user-email":   "foo@bar.baz",
				"x-hasura-admin-secret": "admin-secret",
			},
		},
		{
			"For promotable clients, when they are promoted",
			NewPromotableClient(ts.URL, "admin-secret", NewUserActor(&sampleUUID, "foo@bar.baz"), "test-client").ForceAdmin(),
			map[string]string{
				"hasura-client-name":    "test-client",
				"x-hasura-user-id":      sampleUUID.String(),
				"x-hasura-role":         "admin",
				"x-hasura-user-email":   "foo@bar.baz",
				"x-hasura-admin-secret": "admin-secret",
			},
		},
	} {
		t.Run(tC.desc, func(t *testing.T) {
			_ = tC.cl.Query(context.TODO(), &query, nil)
			for hn, hv := range tC.expectedHeaders {
				assert.Equal(t, hv, h.Get(hn))

			}
		})
	}
}

func TestTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
	}))
	defer ts.Close()

	cl := NewAdminClientFromHost(ts.URL, "admin-secret", "test-client", WithTimeout(1*time.Millisecond))

	var query struct {
		Thing struct {
			Field int `graphql:"field"`
		} `graphql:"thing"`
	}
	t0 := time.Now()
	err := cl.Query(context.TODO(), &query, nil)
	dt := time.Since(t0)

	assert.Less(t, dt, 1*time.Second)
	// there's a variety of errors that occur; let's skip
	// the check if this stays unstable
	//   Post "<url>": context deadline exceeded
	//   Post "<url>": net/http: request canceled (Client.Timeout exceeded while awaiting headers)
	//   Post "<url>": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
	assert.Contains(t, err.Error(), "exceeded")
}

func TestContextHeaders(t *testing.T) {
	var h http.Header
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h = r.Header
	}))
	defer ts.Close()

	cl := NewAdminClientFromHost(ts.URL, "admin-secret", "test-client")

	var query struct {
		Thing struct {
			Field int `graphql:"field"`
		} `graphql:"thing"`
	}

	ctx := context.TODO()
	ctx = WithHeader(ctx, "single-header", "foo")
	ctx = WithHeaders(ctx, map[string]string{
		"extra-header-1": "bar",
		"extra-header-2": "baz",
	})
	ctx = WithHeader(ctx, "extra-header-1", "reset")
	cl.Query(ctx, &query, nil)

	assert.Equal(t, h.Get("single-header"), "foo")
	assert.Equal(t, h.Get("extra-header-1"), "reset")
	assert.Equal(t, h.Get("extra-header-2"), "baz")
}
