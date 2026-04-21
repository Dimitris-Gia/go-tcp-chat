package parser

import (
	"strings"
	"testing"
)

func TestGetPortNumber_Default(t *testing.T) {
	port, err := GetPortNumber([]string{"./TCPChat"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if port != defaultPort {
		t.Fatalf("expected %d, got %d", defaultPort, port)
	}
}

func TestGetPortNumber_CustomPort(t *testing.T) {
	port, err := GetPortNumber([]string{"./TCPChat", "2525"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if port != 2525 {
		t.Fatalf("expected 2525, got %d", port)
	}
}

func TestGetPortNumber_TooManyArgs(t *testing.T) {
	_, err := GetPortNumber([]string{"./TCPChat", "2525", "localhost"})
	if err == nil {
		t.Fatal("expected error for too many args, got nil")
	}
	if !strings.Contains(err.Error(), "[USAGE]: ./TCPChat $port") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestGetPortNumber_NonNumericPort(t *testing.T) {
	_, err := GetPortNumber([]string{"./TCPChat", "abc"})
	if err == nil {
		t.Fatal("expected error for non-numeric port, got nil")
	}
}

func TestGetPortNumber_PortTooLow(t *testing.T) {
	_, err := GetPortNumber([]string{"./TCPChat", "0"})
	if err == nil {
		t.Fatal("expected error for port 0, got nil")
	}
}

func TestGetPortNumber_PortTooHigh(t *testing.T) {
	_, err := GetPortNumber([]string{"./TCPChat", "99999"})
	if err == nil {
		t.Fatal("expected error for port 99999, got nil")
	}
}

func TestGetPortNumber_BoundaryMin(t *testing.T) {
	port, err := GetPortNumber([]string{"./TCPChat", "1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if port != 1 {
		t.Fatalf("expected 1, got %d", port)
	}
}

func TestGetPortNumber_BoundaryMax(t *testing.T) {
	port, err := GetPortNumber([]string{"./TCPChat", "65535"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if port != 65535 {
		t.Fatalf("expected 65535, got %d", port)
	}
}
