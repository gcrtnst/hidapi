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

func RuntimeVersionMajor() int {
	v := C.hid_version()
	if v == nil {
		return 0
	}

	return int(v.major)
}

func RuntimeVersionMinor() int {
	v := C.hid_version()
	if v == nil {
		return 0
	}

	return int(v.minor)
}

func RuntimeVersionPatch() int {
	v := C.hid_version()
	if v == nil {
		return 0
	}

	return int(v.patch)
}

func RuntimeVersion() int {
	v := C.hid_version()
	if v == nil {
		return 0
	}

	major := int(v.major)
	minor := int(v.minor)
	patch := int(v.patch)
	return MakeVersion(major, minor, patch)
}

func RuntimeVersionString() string {
	c := C.hid_version_str()
	return C.GoString(c)
}
