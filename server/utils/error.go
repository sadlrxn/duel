package utils

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// @External
// This is a helper function t make error.
func MakeError(caller string, category string, reason string, err error) error {
	return fmt.Errorf("%s:%s:%s:\n\r %v\n\r", caller, category, reason, err)
}

// @External
// Returns error object with error code.
func MakeErrorWithCode(
	caller string,
	category string,
	reason string,
	code string,
	err error,
) error {
	return fmt.Errorf(
		"%s:%s:%s:\n\r%s\n\r%v\n\r",
		caller,
		category,
		reason,
		codeInError(code),
		err,
	)
}

// @External
// Response an error
func RespondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}

// @External
// Returns whether the error is cuz from the error code
func IsErrorCode(err error, code string) bool {
	if err == nil || code == "" {
		return false
	}

	return strings.Contains(
		err.Error(),
		codeInError(code),
	)
}

// @Internal
// Builds error code part.
func codeInError(code string) string {
	return fmt.Sprintf(
		"DuelBackEndErrCode: %s",
		code,
	)
}
