package constants

const (
	ErrUnhandled    = -999
	ErrUnauthorized = -1
	ErrValidation   = 1001
	ErrNotImage     = 1002
	ErrUploadFailed = 1003
	ErrFileExists   = 1004
	ErrFileNotFound = 1005
	ErrDatabase     = 1006
	ErrIO           = 1007
)

var ErrorMessages = map[int]string{
	ErrUnhandled:    "Unhandled Error",
	ErrUnauthorized: "Unauthorized",
	ErrValidation:   "Validation error",
	ErrNotImage:     "The file is not an image",
	ErrUploadFailed: "Failed to upload image",
	ErrFileExists:   "File exists",
	ErrFileNotFound: "File not found",
	ErrDatabase:     "Database error",
	ErrIO:           "I/O error",
}
