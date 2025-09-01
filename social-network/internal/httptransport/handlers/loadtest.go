package handlers

import (
	"context"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
)

// PostLoadtestWrite запись в бд для нагрузочного тестирования. (POST /loadtest/write).
func (h *Handlers) PostLoadtestWrite(
	ctx context.Context, req serverhttp.PostLoadtestWriteRequestObject,
) (serverhttp.PostLoadtestWriteResponseObject, error) {
	return h.loadTest.Write(ctx, req)
}
