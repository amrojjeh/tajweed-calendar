package assert

import "testing"

func Equal[T comparable](t *testing.T, actual, want T) {
	t.Helper()

	if actual != want {
		t.Errorf("got: %v; want: %v", actual, want)
	}
}
