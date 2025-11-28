package json

import (
	"encoding/json"
	"project/pkg/logger"
)

func JSONMarshal(data interface{}) []byte {
	bytes, err := json.Marshal(data)
	if err != nil {
		logger.Sugar.Error(err)
		return nil
	}

	return bytes
}
