package arango

import (
	"fmt"
	"testing"
)

func TestFindDocument(t *testing.T) {
	id := "_users/2172869"
	var user hash
	if err := Find(id, &user); err != nil {
		t.Error(err)
	}
	fmt.Println(user)
}
