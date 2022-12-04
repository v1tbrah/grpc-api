package api

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/djherbis/times"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pbapi "grpc-api/pkg/api"
)

func (g *GRPCServer) SaveFile(ctx context.Context, req *pbapi.SaveFileRequest) (*pbapi.SaveFileResponse, error) {
	log.Debug().Msg("api.SaveFile START")
	defer log.Debug().Msg("api.SaveFile END")

	// special for testing limiter
	time.Sleep(time.Second * 1)
	// special for testing limiter

	fileName := req.GetName()
	if fileName == "" {
		return &pbapi.SaveFileResponse{}, status.Error(codes.InvalidArgument, ErrFileNameIsEmpty.Error())
	}

	if strings.Contains(fileName, "/") {
		return &pbapi.SaveFileResponse{}, status.Error(codes.InvalidArgument, ErrFileNameIsInvalid.Error())
	}

	err := os.WriteFile(g.dirWithFiles+"/"+req.GetName(), req.GetData(), os.ModePerm)
	if err != nil {
		return &pbapi.SaveFileResponse{}, status.Error(codes.Internal, "")
	}

	return &pbapi.SaveFileResponse{}, nil
}

func (g *GRPCServer) GetFilesInfo(ctx context.Context, req *pbapi.GetFilesInfoRequest) (*pbapi.GetFilesInfoResponse, error) {
	log.Debug().Msg("api.GetFilesInfo START")
	defer log.Debug().Msg("api.GetFilesInfo END")

	// special for testing limiter
	time.Sleep(time.Second * 1)
	// special for testing limiter

	dirData, err := os.ReadDir(g.dirWithFiles)
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}

	filesInfo := make([]*pbapi.FileInfo, 0)

	for _, val := range dirData {

		fStat, errLstat := times.Lstat(g.dirWithFiles + "/" + val.Name())
		if errLstat != nil {
			if errors.Is(errLstat, os.ErrNotExist) {
				continue
			}
			return nil, status.Error(codes.Internal, "")
		}

		fInfo := &pbapi.FileInfo{
			Name: val.Name(),
		}

		if fStat.HasBirthTime() {
			fInfo.HasBirthTime = true
			fInfo.BirthTime = timestamppb.New(fStat.BirthTime())
		}

		if fStat.HasChangeTime() {
			fInfo.HasChangeTime = true
			fInfo.ChangeTime = timestamppb.New(fStat.ChangeTime())
		}

		filesInfo = append(filesInfo, fInfo)

	}

	return &pbapi.GetFilesInfoResponse{FilesInfo: filesInfo}, nil
}

func (g *GRPCServer) GetFiles(ctx context.Context, req *pbapi.GetFilesRequest) (*pbapi.GetFilesResponse, error) {
	log.Debug().Msg("api.GetFilesInfo START")
	defer log.Debug().Msg("api.GetFilesInfo END")

	// special for testing limiter
	time.Sleep(time.Second * 1)
	// special for testing limiter

	dirData, err := os.ReadDir(g.dirWithFiles)
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}

	files := make([][]byte, 0)

	for _, val := range dirData {

		fileData, errReadFile := readFile(g.dirWithFiles + "/" + val.Name())
		if errReadFile != nil {
			if errors.Is(errReadFile, os.ErrNotExist) {
				continue
			}
			return nil, status.Error(codes.Internal, "")
		}

		files = append(files, fileData)

	}

	return &pbapi.GetFilesResponse{Files: files}, nil
}

func readFile(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fr := bufio.NewReader(f)
	data, err := io.ReadAll(fr)
	if err != nil {
		return nil, err
	}

	return data, err
}
