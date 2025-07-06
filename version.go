package hidapi

// #include <hidapi/hidapi.h>
import "C"

const VersionMajor = C.HID_API_VERSION_MAJOR
const VersionMinor = C.HID_API_VERSION_MINOR
const VersionPatch = C.HID_API_VERSION_PATCH
const Version = C.HID_API_VERSION
const VersionString = C.HID_API_VERSION_STR

func MakeVersion(major, minor, patch int) int {
	return (((major) << 24) | ((minor) << 8) | (patch))
}
