package interests

import (
	"context"
)

type serviceMock struct {
}

func NewServiceMock() Service {
	return serviceMock{}
}

func (sm serviceMock) Search(ctx context.Context, groupId, userId string, q *Query, cursor *Cursor) (ids []string, err error) {
	switch cursor.Id {
	case "":
		ids = []string{
			"sub0",
			"sub1",
		}
	case "fail":
		err = ErrInternal
	case "fail_auth":
		err = ErrAuth
	}
	return
}
