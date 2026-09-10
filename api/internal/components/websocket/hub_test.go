package websocket

import (
	"context"
	"testing"
	"time"
)

// newTestClient registers a client in the given workspace without a real
// connection, so hub routing can be tested directly.
func newTestClient(h *Hub, workspaceID string) *Client {
	c := &Client{
		hub:         h,
		send:        make(chan []byte, 64),
		workspaceID: workspaceID,
	}
	h.Register(c)
	return c
}

// waitForRegistered polls until the hub has registered the expected number of
// clients (registration is asynchronous via the run loop).
func waitForRegistered(t *testing.T, h *Hub, want int) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if h.ClientCount() == want {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("timed out waiting for %d clients, have %d", want, h.ClientCount())
}

func recv(t *testing.T, c *Client) ([]byte, bool) {
	t.Helper()
	select {
	case msg, open := <-c.send:
		// A closed channel (client unregistered) must be reported as "no
		// message" rather than a zero-value message.
		return msg, open
	case <-time.After(500 * time.Millisecond):
		return nil, false
	}
}

// TestHub_WorkspaceEventReachesOnlyThatWorkspace is the core isolation
// property: an event scoped to one workspace must never reach clients in
// another (issue #74 — the pre-existing hub broadcast photo IDs to everyone).
func TestHub_WorkspaceEventReachesOnlyThatWorkspace(t *testing.T) {
	h := NewHub(context.Background())
	a := newTestClient(h, "ws-a")
	b := newTestClient(h, "ws-b")
	waitForRegistered(t, h, 2)

	h.BroadcastWorkspaceEvent("ws-a", Event{
		Type:    EventThumbnailReady,
		Payload: ThumbnailReadyPayload{PhotoID: "secret-photo", Size: "m"},
	})

	if _, ok := recv(t, a); !ok {
		t.Fatal("workspace A client should receive its own workspace's event")
	}
	if msg, ok := recv(t, b); ok {
		t.Fatalf("workspace B client must NOT receive workspace A's event, got %s", msg)
	}
}

// TestHub_GlobalEventReachesAllClients covers the documented exception:
// deployment-wide events (no tenant data) reach every client.
func TestHub_GlobalEventReachesAllClients(t *testing.T) {
	h := NewHub(context.Background())
	a := newTestClient(h, "ws-a")
	b := newTestClient(h, "ws-b")
	waitForRegistered(t, h, 2)

	h.BroadcastEvent(Event{Type: EventSystemHealth, Payload: SystemHealthPayload{Status: "ok"}})

	if _, ok := recv(t, a); !ok {
		t.Fatal("global event should reach workspace A")
	}
	if _, ok := recv(t, b); !ok {
		t.Fatal("global event should reach workspace B")
	}
}

// TestHub_WorkspaceEventWithEmptyWorkspace_IsDropped pins the fail-closed
// behaviour: a tenant event with no workspace must not be broadcast globally.
func TestHub_WorkspaceEventWithEmptyWorkspace_IsDropped(t *testing.T) {
	h := NewHub(context.Background())
	a := newTestClient(h, "ws-a")
	waitForRegistered(t, h, 1)

	h.BroadcastWorkspaceEvent("", Event{
		Type:    EventThumbnailReady,
		Payload: ThumbnailReadyPayload{PhotoID: "leak"},
	})

	if msg, ok := recv(t, a); ok {
		t.Fatalf("an unscoped tenant event must be dropped, got %s", msg)
	}
}

// TestHub_UnregisterStopsDelivery verifies a disconnected client receives
// nothing further.
func TestHub_UnregisterStopsDelivery(t *testing.T) {
	h := NewHub(context.Background())
	a := newTestClient(h, "ws-a")
	waitForRegistered(t, h, 1)

	h.unregister <- a
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) && h.ClientCount() != 0 {
		time.Sleep(time.Millisecond)
	}
	if h.ClientCount() != 0 {
		t.Fatal("client should have been unregistered")
	}

	h.BroadcastWorkspaceEvent("ws-a", Event{Type: EventThumbnailReady, Payload: ThumbnailReadyPayload{PhotoID: "x"}})
	if msg, ok := recv(t, a); ok {
		t.Fatalf("unregistered client must not receive messages, got %s", msg)
	}
}
