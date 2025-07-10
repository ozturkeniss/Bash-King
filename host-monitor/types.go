package main

// Ortak struct'lar

type ProcessInfo struct {
	PID     string
	User    string
	CPU     float64
	Memory  float64
	Command string
}

type ConnectionInfo struct {
	Protocol   string
	LocalAddr  string
	RemoteAddr string
	State      string
	PID        string
	Program    string
	Count      int    // Security için
	RemoteIP   string // Security için
}

type UserInfo struct {
	Username string
	Created  string
}
