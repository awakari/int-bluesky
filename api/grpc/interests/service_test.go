package interests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
	"time"
)

func TestService_Search(t *testing.T) {
	//
	svc := NewService(NewClientMock())
	svc = NewLoggingMiddleware(svc, slog.Default())
	//
	cases := map[string]struct {
		cursor string
		ids    []string
		err    error
	}{
		"ok": {
			ids: []string{
				"sub0",
				"sub1",
			},
		},
		"fail": {
			cursor: "fail",
			err:    ErrInternal,
		},
		"fail auth": {
			cursor: "fail_auth",
			err:    ErrAuth,
		},
	}
	//
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			ids, err := svc.Search(ctx, "group0", "user0", &Query{}, &Cursor{
				Id: c.cursor,
			})
			assert.Equal(t, c.ids, ids)
			assert.ErrorIs(t, err, c.err)
		})
	}
}
