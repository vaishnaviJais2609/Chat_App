package ws

import (
	"testing"
	"time"
)

func newTestClient(username string) *Client {
	return &Client{
		username: username,
		send:     make(chan Message, 16),
	}
}

func TestRegisterRules(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	if !hub.Register(newTestClient("alice")) {
		t.Fatal("first user should be accepted")
	}

	if hub.Register(newTestClient("alice")) {
		t.Fatal("duplicate username should be rejected")
	}

	if !hub.Register(newTestClient("bob")) {
		t.Fatal("second user should be accepted")
	}

	if hub.Register(newTestClient("carol")) {
		t.Fatal("third user should be rejected: chat is full")
	}
}

func TestBroadcastSkipsSender(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	alice := newTestClient("alice")
	bob := newTestClient("bob")

	if !hub.Register(alice) || !hub.Register(bob) {
		t.Fatal("both users should register")
	}

	hub.broadcast <- Message{User: "alice", Message: "hi"}

	select {
	case got := <-bob.send:
		if got.User != "alice" || got.Message != "hi" {
			t.Fatalf("unexpected message: %+v", got)
		}
	case <-time.After(time.Second):
		t.Fatal("bob did not receive the message")
	}

	select {
	case <-alice.send:
		t.Fatal("sender should not receive their own message")
	case <-time.After(50 * time.Millisecond):
	}
}

func TestSlowClientDoesNotBlockHub(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	alice := newTestClient("alice")
	bob := newTestClient("bob")

	if !hub.Register(alice) || !hub.Register(bob) {
		t.Fatal("both users should register")
	}

	// Bob never reads. Overfill his buffer; the hub must keep running.
	done := make(chan struct{})
	go func() {
		for i := 0; i < 100; i++ {
			hub.broadcast <- Message{User: "alice", Message: "spam"}
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("hub blocked on a slow client")
	}
}
