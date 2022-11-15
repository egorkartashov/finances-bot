package utils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const codeUnknown = "unknown"

func GetResponseCode(err error) string {
	if err == nil {
		return codes.OK.String()
	}
	s, ok := status.FromError(err)
	code := codeUnknown
	if ok {
		code = s.Code().String()
	}
	return code
}
