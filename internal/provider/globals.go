package provider

const (
	EmptyPath = ""
)

const (
	JobStatusWaiting = iota
	JobStatusQueued
	JobStatusSkipped
	JobStatusFailed
	JobStatusDone
)
