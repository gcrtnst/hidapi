package hidapi

// #include <hidapi/hidapi.h>
import "C"

const VersionMajor = C.HID_API_VERSION_MAJOR
const VersionMinor = C.HID_API_VERSION_MINOR
const VersionPatch = C.HID_API_VERSION_PATCH
const Version = C.HID_API_VERSION
const VersionString = C.HID_API_VERSION_STR

func MakeVersion(mj, mn, p int) int { return (((mj) << 24) | ((mn) << 8) | (p)) }
