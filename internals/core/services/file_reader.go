package services

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/jocbarbosa/viswals-backend/internals/core/model"
	"github.com/jocbarbosa/viswals-backend/internals/core/port"
)

type FileReaderService struct {
	logger    port.Logger
	messaging port.Messaging
	filepath  string
}

func NewFileReader(logger port.Logger, m port.Messaging, filepath string) *FileReaderService {
	return &FileReaderService{
		logger:    logger,
		filepath:  filepath,
		messaging: m,
	}
}

func (f *FileReaderService) ReadFile() error {
	f.logger.Info("Reading file")

	file, err := os.Open(f.filepath)
	if err != nil {
		f.logger.Error("error opening file", err)
		return err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))

	if _, err := reader.Read(); err != nil {
		f.logger.Error("error reading file", err)
		return err
	}

	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			f.logger.Error("error reading file", err)
			return err
		}

		id, _ := strconv.Atoi(record[0])
		createdAt, _ := strconv.ParseInt(record[4], 10, 64)
		mergedAt, _ := strconv.ParseInt(record[6], 10, 64)
		parentUserID, _ := strconv.Atoi(record[7])

		user := model.User{
			ID:           id,
			FirstName:    record[1],
			LastName:     record[2],
			Email:        record[3],
			CreatedAt:    time.Unix(createdAt/1000, 0),
			MergedAt:     time.Unix(mergedAt/1000, 0),
			ParentUserID: parentUserID,
		}

		userBytes, err := json.Marshal(user)
		if err != nil {
			f.logger.Error("error marshalling user", err)
			return err
		}

		f.logger.Info("sending user to queue", user.ID)
		err = f.messaging.Write(port.Message{
			Value: userBytes,
		})
		if err != nil {
			f.logger.Error("error writing message to queue", err)
			return err
		}
	}

	return nil
}
