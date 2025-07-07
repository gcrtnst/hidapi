package hidapi

// #include <hidapi/hidapi.h>
import "C"

const (
	CompileTimeVersionMajor  = C.HID_API_VERSION_MAJOR
	CompileTimeVersionMinor  = C.HID_API_VERSION_MINOR
	CompileTimeVersionPatch  = C.HID_API_VERSION_PATCH
	CompileTimeVersionCode   = C.HID_API_VERSION
	CompileTimeVersionString = C.HID_API_VERSION_STR
)

func EncodeVersion(major, minor, patch int) int {
	return (((major) << 24) | ((minor) << 8) | (patch))
}

func RuntimeVersion() (int, int, int) {
	v := C.hid_version()
	if v == nil {
		return 0, 0, 0
	}

	major := int(v.major)
	minor := int(v.minor)
	patch := int(v.patch)
	return major, minor, patch
}

func RuntimeVersionString() string {
	c := C.hid_version_str()
	return C.GoString(c)
}
