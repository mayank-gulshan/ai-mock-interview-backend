// Package apperrors mirrors the custom exception types and ErrorResponse
// from Security/Exception.kt. The @RestControllerAdvice GlobalExceptionHandler
// itself is translated as RespondError below (see also
// internal/security/middleware.go for the AuthenticationException /
// CustomAuthEntryPoint equivalent), since Go has no global exception
// dispatch mechanism — each handler explicitly maps the returned error to
// the same HTTP status/body Spring would have produced.
package apperrors

import "net/http"

type UserAlreadyExistsError struct{ Msg string }

func (e *UserAlreadyExistsError) Error() string { return e.Msg }

type InvalidCredentialsError struct{ Msg string }

func (e *InvalidCredentialsError) Error() string { return e.Msg }

type InvalidTokenError struct{ Msg string }

func (e *InvalidTokenError) Error() string { return e.Msg }

// ErrorResponse mirrors Security/Exception.kt's ErrorResponse(val message: String).
type ErrorResponse struct {
	Message string `json:"message"`
}

// StatusFor mirrors the @ExceptionHandler status-code mapping in
// GlobalExceptionHandler:
//
//	UserAlreadyExistsException  -> 409 CONFLICT
//	InvalidCredentialsException -> 401 UNAUTHORIZED
//	InvalidTokenException       -> 401 UNAUTHORIZED
//	AuthenticationException     -> 401 UNAUTHORIZED ("Unauthorized")
//	(anything else)             -> 500 INTERNAL_SERVER_ERROR (Spring default)
func StatusFor(err error) (int, ErrorResponse) {
	switch e := err.(type) {
	case *UserAlreadyExistsError:
		return http.StatusConflict, ErrorResponse{Message: e.Msg}
	case *InvalidCredentialsError:
		return http.StatusUnauthorized, ErrorResponse{Message: e.Msg}
	case *InvalidTokenError:
		return http.StatusUnauthorized, ErrorResponse{Message: e.Msg}
	default:
		return http.StatusInternalServerError, ErrorResponse{Message: "Internal server error"}
	}
}
