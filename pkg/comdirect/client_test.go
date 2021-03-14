package comdirect

import (
	"log"
	"testing"
)

func TestGenerateSessionId(t *testing.T) {
	sessionId, err := generateSessionId()

	if err != nil {
		t.Errorf("Error generating session id")
	}

	if len(sessionId) != 32 {
		t.Errorf("Length of session id not equal to 32: %d", len(sessionId))
	}

	log.Println(sessionId)
}

func TestGenerateRequestId(t *testing.T) {
	generateRequestId()
}
