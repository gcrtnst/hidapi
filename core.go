package hidapi

// #include <hidapi/hidapi.h>
import "C"

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
