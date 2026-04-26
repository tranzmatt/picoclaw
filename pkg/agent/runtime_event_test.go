package agent

import (
	"testing"
	"time"

	runtimeevents "github.com/sipeed/picoclaw/pkg/events"
)

func subscribeRuntimeEventsForTest(
	t *testing.T,
	al *AgentLoop,
	buffer int,
	kinds ...runtimeevents.Kind,
) (<-chan runtimeevents.Event, func()) {
	t.Helper()

	if al == nil {
		t.Fatal("agent loop is nil")
	}
	channel := al.RuntimeEvents()
	if channel == nil {
		t.Fatal("runtime event channel is nil")
	}
	if len(kinds) > 0 {
		channel = channel.OfKind(kinds...)
	}
	sub, ch, err := channel.SubscribeChan(
		t.Context(),
		runtimeevents.SubscribeOptions{Name: "agent-runtime-test", Buffer: buffer},
	)
	if err != nil {
		t.Fatalf("SubscribeChan failed: %v", err)
	}
	return ch, func() {
		if err := sub.Close(); err != nil {
			t.Errorf("runtime subscription close failed: %v", err)
		}
	}
}

func waitForRuntimeEvent(
	t *testing.T,
	ch <-chan runtimeevents.Event,
	timeout time.Duration,
	match func(runtimeevents.Event) bool,
) runtimeevents.Event {
	t.Helper()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case evt, ok := <-ch:
			if !ok {
				t.Fatal("runtime event stream closed before expected event arrived")
			}
			if match(evt) {
				return evt
			}
		case <-timer.C:
			t.Fatal("timed out waiting for expected runtime event")
		}
	}
}
