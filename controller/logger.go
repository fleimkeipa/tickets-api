package controller

import (
	"bytes"
	"encoding/json"

	"github.com/fleimkeipa/tickets-api/models"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Logger is a wrapper around zap.SugaredLogger
type Logger struct {
	logger *zap.SugaredLogger
}

// NewLogger initializes a new Logger instance
func NewLogger(logger *zap.SugaredLogger) *Logger {
	return &Logger{logger: logger}
}

// responseWriter wraps echo.Response to capture the response body
type responseWriter struct {
	echo.Response
	body *bytes.Buffer
}

// Write captures the response body while continuing to write to the original response
func (rc *responseWriter) Write(b []byte) (int, error) {
	rc.body.Write(b)            // Buffer the response body
	return rc.Response.Write(b) // Write the response to the client
}

// WriteHeader captures the status code
func (rc *responseWriter) WriteHeader(statusCode int) {
	rc.Response.WriteHeader(statusCode)
}

// LoggerMiddleware intercepts the response to log any errors present in the response body
func (rc *Logger) LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Wrap the original response writer to intercept the response body
		var res = c.Response()

		var rawPath = c.Path()
		if rawPath == "/swagger/*" {
			return next(c)
		}

		var bodyBuffer = new(bytes.Buffer)
		var writer = &responseWriter{
			Response: *res,
			body:     bodyBuffer,
		}
		c.Response().Writer = writer

		// Proceed with request handling
		var err = next(c)

		// If the response status code indicates a success (100-399),
		// pass the request to the next handler without logging.
		if res.Status > 99 && res.Status < 400 {
			return err
		}

		// After the handler, check if there was an error in the response
		var failureResponse models.FailureResponse
		if json.Unmarshal(writer.body.Bytes(), &failureResponse) == nil {
			// If the response contains an error, log it
			if failureResponse.Error != "" {
				rc.logger.Errorf("Error logged: %v", failureResponse.Error)
			}
		}

		return err
	}
}
