package gql

import (
	"context"
)

type key int

const headerKey key = iota

func WithHeader(ctx context.Context, key, value string) context.Context {
	return WithHeaders(ctx, map[string]string{key: value})
}

func WithHeaders(ctx context.Context, hs map[string]string) context.Context {
	headers := GetHeadersFromContext(ctx)
	for k, v := range hs {
		headers[k] = v
	}
	return context.WithValue(ctx, headerKey, headers)
}

func GetHeadersFromContext(ctx context.Context) map[string]string {
	h := ctx.Value(headerKey)
	var headers map[string]string
	if h == nil {
		headers = map[string]string{}
	} else {
		headers = h.(map[string]string)
	}
	return headers
}
