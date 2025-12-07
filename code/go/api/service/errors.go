package service

import (
	"github.com/pkg/errors"
)

var (
	// ErrNoChanges is returned when there is no changes for stored entity.
	ErrNoChanges = errors.New("no changes")

	// ErrNoPermissions is returned when there is no permission for entity.
	ErrNoPermissions = errors.New("no permissions")
)
