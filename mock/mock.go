package mock

import (
	"testing"

	"go.nhat.io/matcher/v3"
)

// Mocker is Matcher mocker.
type Mocker func(tb testing.TB) *Matcher

// Nop is no mock Matcher.
var Nop = Mock()

var _ matcher.Matcher = (*Matcher)(nil)

// Mock creates Matcher mock with cleanup to ensure all the expectations are met.
func Mock(mocks ...func(m *Matcher)) Mocker {
	return func(tb testing.TB) *Matcher {
		tb.Helper()

		result := NewMatcher(tb)

		for _, m := range mocks {
			m(result)
		}

		return result
	}
}
