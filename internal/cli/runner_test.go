package cli

import (
	"errors"
	"testing"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

type timeoutErr struct{}

func (timeoutErr) Error() string   { return "timeout" }
func (timeoutErr) Timeout() bool   { return true }
func (timeoutErr) Temporary() bool { return true }

func TestShouldRetryAPIError429(t *testing.T) {
	err := &openai.APIError{HTTPStatusCode: 429}
	if !shouldRetry(err) {
		t.Fatal("expected retry for 429")
	}
}

func TestShouldRetryTimeoutNetError(t *testing.T) {
	if !shouldRetry(timeoutErr{}) {
		t.Fatal("expected retry for timeout net error")
	}
}

func TestShouldRetryNonRetryable(t *testing.T) {
	if shouldRetry(errors.New("validation failed")) {
		t.Fatal("did not expect retry for non-retryable error")
	}
}

func TestRetryDelayForAttempt(t *testing.T) {
	base := time.Second
	if got := retryDelayForAttempt(base, 1); got != 1*time.Second {
		t.Fatalf("expected 1s, got %s", got)
	}
	if got := retryDelayForAttempt(base, 3); got != 4*time.Second {
		t.Fatalf("expected 4s, got %s", got)
	}
}

func TestRetryDelayCapped(t *testing.T) {
	base := 20 * time.Second
	if got := retryDelayForAttempt(base, 3); got != 30*time.Second {
		t.Fatalf("expected capped 30s, got %s", got)
	}
}
