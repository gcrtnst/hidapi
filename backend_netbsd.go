//go:build netbsd

package hidapi

// #cgo pkg-config: hidapi-netbsd
import "C"

const hidBackendInitThread = false
