package controllers

import (
	"context"
	"fmt"

	"github.com/jocbarbosa/viswals-backend/internals/core/port"
	"github.com/jocbarbosa/viswals-backend/internals/core/services"
)

type FileReaderController struct {
	logger    port.Logger
	messaging port.Messaging
	filepath  string
}

func NewFileReader(l port.Logger, m port.Messaging, f string) *FileReaderController {
	fmt.Println("Creating new file reader")
	return &FileReaderController{
		logger:    l,
		messaging: m,
		filepath:  f,
	}
}

func (f *FileReaderController) ReadFile(ctx context.Context) error {
	fmt.Println("Reading file from", f.filepath)

	svc := services.NewFileReader(f.logger, f.messaging, f.filepath)

	err := svc.ReadFile()
	if err != nil {
		f.logger.Error("Error reading file", err)
		return err
	}

	return nil
}
