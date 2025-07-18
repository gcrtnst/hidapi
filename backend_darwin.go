//go:build darwin

package hidapi

// #cgo pkg-config: hidapi
import "C"

// On macOS, hid_init and hid_exit must be called on the same thread.
// https://github.com/libusb/hidapi/wiki/Multi%E2%80%90threading-Notes
const hidBackendInitThread = true
