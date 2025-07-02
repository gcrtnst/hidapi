package hidapi

// #include <hidapi/hidapi.h>
import "C"
import (
	"strconv"
	"unicode/utf16"
	"unsafe"
)

type Error struct {
	Text string
}

func (err *Error) Error() string {
	return err.Text
}

func hidError(dev *Device) error {
	cdev := (*C.hid_device)(nil)
	if dev != nil {
		dev.mu.Lock()
		defer dev.mu.Unlock()
		cdev = dev.dev
	}

	cstr := C.hid_error(cdev)
	gstr := convertWCharPtrToString(cstr)

	text := gstr
	if text == "hid_error is not implemented yet" {
		text = "an error occurred with hidapi"
	}

	return &Error{Text: text}
}

func convertWCharPtrToString(p *C.wchar_t) string {
	if C.sizeof_wchar_t == unsafe.Sizeof(uint16(0)) { // UTF-16
		l := convertWCharPtrToStringLen(p)
		u := (*uint16)(unsafe.Pointer(p))
		s := unsafe.Slice(u, l)
		d := utf16.Decode(s)
		return string(d)
	}

	if C.sizeof_wchar_t == unsafe.Sizeof(rune(0)) { // UTF-32
		l := convertWCharPtrToStringLen(p)
		r := (*rune)(unsafe.Pointer(p))
		s := unsafe.Slice(r, l)
		return string(s)
	}

	panic("unexpected size of wchar_t: " + strconv.Itoa(C.sizeof_wchar_t))
}

func convertWCharPtrToStringLen(p *C.wchar_t) int {
	l := 0
	for {
		var u unsafe.Pointer
		var c C.wchar_t
		u = unsafe.Pointer(p)
		u = unsafe.Pointer(uintptr(u) + uintptr(l*C.sizeof_wchar_t))
		c = *(*C.wchar_t)(u)

		if c == 0 {
			break
		}
		l++
	}
	return l
}
