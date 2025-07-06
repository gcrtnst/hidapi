package hidapi

// #include <hidapi/hidapi.h>
import "C"
import (
	"strconv"
	"unicode/utf16"
	"unsafe"
)

func convertWCharPtrToString(p *C.wchar_t) string {
	if p == nil {
		return ""
	}

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
	if p == nil {
		return 0
	}

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

func convertStringToWCharPtr(s string) *C.wchar_t {
	if s == "" {
		return nil
	}

	if C.sizeof_wchar_t == unsafe.Sizeof(uint16(0)) {
		r := []rune(s)
		t := utf16.Encode(r)
		a := append(t, 0)
		return (*C.wchar_t)(unsafe.Pointer(&a[0]))
	}

	if C.sizeof_wchar_t == unsafe.Sizeof(rune(0)) {
		t := []rune(s)
		a := append(t, 0)
		return (*C.wchar_t)(unsafe.Pointer(&a[0]))
	}

	panic("unexpected size of wchar_t: " + strconv.Itoa(C.sizeof_wchar_t))
}
