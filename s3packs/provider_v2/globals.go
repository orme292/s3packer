package provider_v2

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
