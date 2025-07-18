//go:build !windows && !darwin && !linux && !netbsd

package hidapi

// #cgo pkg-config: hidapi-libusb
import "C"
