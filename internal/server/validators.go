package server

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"

	notesv1 "github.com/DeathPolarizer/notes/gen/notes/v1"
)

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
)

func decodePageToken(pageToken string) (int32, error) {
	offsetBytes, err := base64.StdEncoding.DecodeString(pageToken)
	if err != nil {
		return 0, err
	}
	offset, err := strconv.Atoi(string(offsetBytes))
	if err != nil {
		return 0, err
	}
	return int32(offset), nil
}

func validateListNotesRequest(req *notesv1.ListNotesRequest) (int32, int32, error) {
	pageSize := req.GetPageSize()
	switch {
	case pageSize == 0:
		pageSize = DefaultPageSize
	case pageSize > MaxPageSize || pageSize < 0:
		return 0, 0, fmt.Errorf("%w: pageSize", ErrInvalidArgument)
	}
	var offset int32
	var err error
	pageToken := req.GetPageToken()
	if pageToken != "" {
		offset, err = decodePageToken(pageToken)
	}
	if err != nil {
		return 0, 0, fmt.Errorf("%w: pageToken", ErrInvalidArgument)
	}
	if offset < 0 {
		return 0, 0, fmt.Errorf("%w: pageToken", ErrInvalidArgument)
	}
	return pageSize, offset, nil
}
