package hidapi

// #include <hidapi/hidapi.h>
import "C"

const (
	CompileTimeVersionMajor  = C.HID_API_VERSION_MAJOR
	CompileTimeVersionMinor  = C.HID_API_VERSION_MINOR
	CompileTimeVersionPatch  = C.HID_API_VERSION_PATCH
	CompileTimeVersion       = C.HID_API_VERSION
	CompileTimeVersionString = C.HID_API_VERSION_STR
)

func MakeVersion(major, minor, patch int) int {
	return (((major) << 24) | ((minor) << 8) | (patch))
}
