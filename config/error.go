package config

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorMessage(message string, messageCode codes.Code) error {
	return status.Errorf(
		messageCode,
		message,
	)
}
