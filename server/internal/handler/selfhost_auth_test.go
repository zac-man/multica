package handler

import "testing"

func TestSelfhostPasswordEqual(t *testing.T) {
	if !selfhostPasswordEqual("secret", "secret") {
		t.Fatal("expected match")
	}
	if selfhostPasswordEqual("secret", "Secret") {
		t.Fatal("expected case-sensitive mismatch")
	}
	if selfhostPasswordEqual("a", "ab") {
		t.Fatal("expected length mismatch")
	}
	if !selfhostPasswordEqual("", "") {
		t.Fatal("empty vs empty")
	}
}
