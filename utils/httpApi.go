package utils

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// To be called in all http backend api handlers to log an error
func HttpError(response http.ResponseWriter, logger *logrus.Logger, message string, httpStatusCode int) {
	logger.Error(message)
	http.Error(response, message, httpStatusCode)
}
