package util

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"strings"
)

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func GetErrorResponse(err error) gin.H {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ErrorMsg, len(ve))
		for i, fe := range ve {
			out[i] = ErrorMsg{Field: fe.Field(), Message: getErrorMsg(fe)}
		}

		return gin.H{"errors": out, "status": "error"}
	}

	var customErr JobTypeSchemaValidationErrors
	if errors.As(err, &customErr) {
		out := make([]ErrorMsg, len(customErr))
		for i, fe := range customErr {
			out[i] = ErrorMsg{Field: fe.FieldName, Message: getCustomErrorMsg(fe)}
		}
		return gin.H{"errors": out, "status": "error"}
	}

	return gin.H{"message": "Failed to parse required request body.", "status": "error"}
}

func GetErrorResponseWithMessage(message string) gin.H {
	return gin.H{"message": message, "status": "error"}
}

func GetErrorResponseAndMessage(err error, message string) gin.H {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ErrorMsg, len(ve))
		for i, fe := range ve {
			out[i] = ErrorMsg{Field: fe.Field(), Message: getErrorMsg(fe)}
		}

		return gin.H{"errors": out, "status": "error", "message": message}
	}

	var jtsve *JobTypeSchemaValidationErrors
	if errors.As(err, &jtsve) {
		out := make([]ErrorMsg, len(*jtsve))
		for i, fe := range *jtsve {
			out[i] = ErrorMsg{Field: fe.FieldName, Message: getCustomErrorMsg(fe)}
		}
		return gin.H{"errors": out, "status": "error", "message": message}
	}

	return gin.H{"message": message, "status": "error"}
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	}
	return "Unknown error"
}

func getCustomErrorMsg(fe JobTypeSchemaFieldError) string {
	switch fe.ValidationResult {
	case "missing":
		return "This field is required and missing"
	case "invalid":
		return "This field is invalid / has an invalid type"
	}
	return "Unknown error"
}

type CustomValidationErrorsTranslations map[string]string
type JobTypeSchemaValidationErrors []JobTypeSchemaFieldError

func (ce JobTypeSchemaValidationErrors) Error() string {

	buff := bytes.NewBufferString("JobTypeSchemaValidationErrors:\n")

	for i := 0; i < len(ce); i++ {

		buff.WriteString(ce[i].Error())
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

func (ce JobTypeSchemaValidationErrors) Translate(ut ut.Translator) CustomValidationErrorsTranslations {

	trans := make(CustomValidationErrorsTranslations)

	return trans
}

type JobTypeSchemaFieldError struct {
	FieldName        string
	FieldType        string
	ValidationResult string
}

func (e *JobTypeSchemaFieldError) Error() string {
	return fmt.Sprintf("missing field: %s of type %s", e.FieldName, e.FieldType)
}
