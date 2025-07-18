//go:build windows

package hidapi

// #cgo pkg-config: hidapi
import "C"

const hidBackendInitThread = false
