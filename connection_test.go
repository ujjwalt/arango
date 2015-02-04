package arango

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	NewConnection("http://localhost:8529")
	os.Exit(m.Run())
}

func TestCurrentDatabase(t *testing.T) {
	d, err := CurrentDatabase()
	if err != nil || d == nil {
		t.Error(err)
	}
}
