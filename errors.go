package minutes

import "errors"

var (
	// ErrNotFound is returned when a single resource was requested
	// but was not found
	ErrNotFound = errors.New("Not found")
	// ErrInternalServer is returned when the request failed with
	// anything that was out of the requester's control
	ErrInternalServer = errors.New("Internal server error")
	// ErrMissingShow is returned when trying to request, add, or modify
	// a resource that is related to a Show, but the Show does not exist
	ErrMissingShow = errors.New("Show does not exist")
	// ErrMissingSeason is returned when trying to request, add, or modify
	// a resource that is related to a Season, but the Season does not exist
	ErrMissingSeason = errors.New("Season does not exist")
	// ErrNotImplemented is returned when trying to add a new resource on a
	// read-only Library
	ErrNotImplemented = errors.New("Not implemented")
)
