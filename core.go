package hidapi

// #include <hidapi/hidapi.h>
import "C"
import (
	"math"
	"strconv"
	"sync"
	"sync/atomic"
	"unicode/utf16"
	"unsafe"
)

var (
	hidMutex sync.Mutex
	hidCount = 0
)

type Error struct {
	Text string
}

func (err *Error) Error() string {
	return err.Text
}

type hidRef struct {
	OK atomic.Bool
}

func hidAcquire() (*hidRef, error) {
	hidMutex.Lock()
	defer hidMutex.Unlock()

	if hidCount < 0 {
		panic("hidapi: hidCount < 0")
	}
	if hidCount >= math.MaxInt {
		panic("hidapi: hidCount >= math.MaxInt")
	}
	if hidCount == 0 {
		err := hidInit()
		if err != nil {
			return nil, err
		}
	}
	hidCount++

	ref := &hidRef{}
	ref.OK.Store(true)
	return ref, nil
}

func (ref *hidRef) Close() error {
	if !ref.OK.Swap(false) {
		return nil
	}

	hidMutex.Lock()
	defer hidMutex.Unlock()

	if hidCount <= 0 {
		panic("hidapi: hidCount <= 0")
	}
	hidCount--
	if hidCount == 0 {
		_ = hidExit()
	}

	return nil
}

func hidInit() error {
	ret := C.hid_init()
	if ret != 0 {
		return hidError(nil)
	}
	return nil
}

func hidExit() bool {
	ret := C.hid_exit()
	return ret == 0
}

func hidError(dev *Device) error {
	cdev := (*C.hid_device)(nil)
	if dev != nil {
		dev.mu.Lock()
		defer dev.mu.Unlock()
		cdev = dev.cptr()
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
