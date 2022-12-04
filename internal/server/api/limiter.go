package api

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type rateLimiter struct {
	saveFileLimiter     *methodLimiter
	getFilesInfoLimiter *methodLimiter
	getFilesLimiter     *methodLimiter
}

type methodLimiter struct {
	currLoading byte
	maxLoading  byte
	mu          sync.Mutex
}

func newRateLimiter() *rateLimiter {
	log.Debug().Msg("api.newRateLimiter START")
	defer log.Debug().Msg("api.newRateLimiter END")

	limiter := &rateLimiter{}

	limiter.saveFileLimiter = &methodLimiter{mu: sync.Mutex{}, maxLoading: 10}
	limiter.getFilesInfoLimiter = &methodLimiter{mu: sync.Mutex{}, maxLoading: 100}
	limiter.getFilesLimiter = &methodLimiter{mu: sync.Mutex{}, maxLoading: 10}

	return limiter

}

func (l *rateLimiter) interceptorLimiter(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	var mLimiter *methodLimiter

	switch info.FullMethod {
	case "/api.FileMngr/SaveFile":
		mLimiter = l.saveFileLimiter
	case "/api.FileMngr/GetFilesInfo":
		mLimiter = l.getFilesInfoLimiter
	case "/api.FileMngr/GetFiles":
		mLimiter = l.getFilesLimiter
	}

	if mLimiter == nil {
		return handler(ctx, req)
	}

	if mLimiter != nil {
		mLimiter.mu.Lock()
		newLoading := mLimiter.currLoading + 1
		if newLoading > mLimiter.maxLoading {
			mLimiter.mu.Unlock()
			return nil, status.Error(codes.ResourceExhausted, fmt.Sprintf("%s is rejected by grpc_ratelimit middleware, please retry later.", info.FullMethod))
		}

		mLimiter.currLoading = newLoading
		mLimiter.mu.Unlock()

		resp, err = handler(ctx, req)

		mLimiter.mu.Lock()
		mLimiter.currLoading--
		mLimiter.mu.Unlock()
	}

	return resp, err
}
