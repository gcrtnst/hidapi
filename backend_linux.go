//go:build linux

package hidapi

// #cgo pkg-config: hidapi-hidraw
import "C"

const hidBackendInitThread = false
