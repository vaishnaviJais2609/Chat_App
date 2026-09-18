package ws

import "testing"

func TestNewHub_AllowedUsers(t *testing.T) {
	h := NewHub("alice", "bob")

	if !h.AllowedUsers["alice"] || !h.AllowedUsers["bob"] {
		t.Fatalf("expected both alice and bob to be allowed users")
	}
	if h.AllowedUsers["charlie"] {
		t.Fatalf("expected charlie to NOT be an allowed user")
	}
}

func TestHub_IsFullAndUsernameTaken(t *testing.T) {
	h := NewHub("alice", "bob")

	if h.IsFull() {
		t.Fatalf("expected hub to not be full with no clients")
	}
	if h.IsUsernameTaken("alice") {
		t.Fatalf("expected alice to not be taken before registering")
	}

	h.clients["alice"] = &Client{Username: "alice"}
	if !h.IsUsernameTaken("alice") {
		t.Fatalf("expected alice to be taken after registering")
	}
	if h.IsFull() {
		t.Fatalf("expected hub to not be full with only one client")
	}

	h.clients["bob"] = &Client{Username: "bob"}
	if !h.IsFull() {
		t.Fatalf("expected hub to be full with two clients")
	}
}
