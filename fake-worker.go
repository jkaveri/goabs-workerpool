package workerpool

import (
	"context"
)

func fakeWorker(context.Context, []byte) error {
	return nil
}