// Package test provides a simple way to test the code using snapshots.
package test

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func Match(t *testing.T, values ...any) {
	t.Helper()

	snaps.
		WithConfig(snaps.Dir("testdata")).
		MatchSnapshot(t, values...)
}
