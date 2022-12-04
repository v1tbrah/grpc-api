package api

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	pbapi "grpc-api/pkg/api"
)

type GRPCServer struct {
	serv         *grpc.Server
	runAddress   string
	dirWithFiles string
	rateLimiter  *rateLimiter
	pbapi.UnimplementedFileMngrServer
}

func New(cfg Config) (grpcServer *GRPCServer, err error) {
	log.Debug().Str("cfg", cfg.String()).Msg("api.New START")
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("api.New END")
		} else {
			log.Debug().Msg("api.New END")
		}
	}()

	grpcServer = &GRPCServer{}

	limiter := newRateLimiter()
	grpcServer.rateLimiter = limiter

	serv := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logrus.New())),
				limiter.interceptorLimiter,
			),
		),
	)

	grpcServer.serv = serv

	runAddress := cfg.ServAPIAddr()
	grpcServer.runAddress = runAddress

	dirWithFiles := cfg.DirWithFiles()
	if isDir, err := isDirectory(dirWithFiles); !isDir || err != nil {
		return nil, fmt.Errorf("checking directory with files: dir name: %s: error: %w", dirWithFiles, err)
	}
	grpcServer.dirWithFiles = dirWithFiles

	return grpcServer, nil
}

func (g *GRPCServer) Run() {
	log.Debug().Msg("api.Run START")
	defer log.Debug().Msg("api.Run END")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	errG, ctx := errgroup.WithContext(context.Background())

	log.Info().Str("address", g.runAddress).Msg("grpc server starting")
	errG.Go(func() error {
		return g.startServing(ctx, shutdown)
	})

	if err := errG.Wait(); err != nil {
		log.Error().Err(err).Msg(err.Error())
		close(shutdown)
	}

	<-shutdown
	g.serv.GracefulStop()
	log.Info().Str("address", g.runAddress).Msg("grpc server gracefully stopped")

}

func (g *GRPCServer) startServing(ctx context.Context, shutdown chan os.Signal) (err error) {
	log.Debug().Msg("api.startServing START")
	defer log.Debug().Msg("api.startServing END")

	ended := make(chan struct{})

	go func() {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-shutdown:
			if ok {
				close(shutdown)
			}
			return
		default:
			pbapi.RegisterFileMngrServer(g.serv, g)

			l, errListen := net.Listen("tcp", g.runAddress)
			if errListen != nil {
				err = fmt.Errorf("net listen tcp %s server: %w", g.runAddress, errListen)
				ended <- struct{}{}
				return
			}

			if err = g.serv.Serve(l); err != nil {
				err = fmt.Errorf("serve %s server: %w", g.runAddress, err)
			}
			ended <- struct{}{}
		}
	}()

	select {
	case <-ctx.Done():
		return err
	case _, ok := <-shutdown:
		if ok {
			close(shutdown)
		}
		return err
	case <-ended:
		return err
	}

}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}
