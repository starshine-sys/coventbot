package levels

import "testing"

func TestHumanizeInt64(t *testing.T) {
	tests := []struct {
		Input  int64
		Expect string
	}{{1, "1"}, {0, "0"}, {-100, "-100"}, {1234, "1.2k"}, {22953, "23.0k"}, {22949, "22.9k"}}

	for _, tt := range tests {
		if out := HumanizeInt64(tt.Input); out != tt.Expect {
			t.Errorf("%d: expect %q, got %q", tt.Input, tt.Expect, out)
		}
	}
}
