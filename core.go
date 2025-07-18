package hidapi

// #include <hidapi/hidapi.h>
import "C"
import (
	"math"
	"sync"
	"sync/atomic"
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
	ok uint32
}

func hidAcquire() (hidRef, error) {
	if hidCount < 0 {
		panic("hidapi: hidCount < 0")
	}
	if hidCount >= math.MaxInt {
		panic("hidapi: hidCount >= math.MaxInt")
	}
	if hidCount == 0 {
		err := hidInit()
		if err != nil {
			return hidRef{}, err
		}
	}
	hidCount++

	ref := hidRef{ok: 1}
	return ref, nil
}

func (ref *hidRef) Close() error {
	if atomic.SwapUint32(&ref.ok, 0) == 0 {
		return nil
	}

	if hidCount <= 0 {
		panic("hidapi: hidCount <= 0")
	}
	hidCount--
	if hidCount == 0 {
		_ = hidExit()
	}

	return nil
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
