package hidapi

// #include <hidapi/hidapi.h>
import "C"
import (
	"runtime"
	"sync"
)

type Device struct {
	mu  sync.Mutex
	cln runtime.Cleanup
	in  *hidDevice
}

func (d *Device) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	var err error
	if d.in != nil {
		err = d.in.Close()
	}

	d.cln.Stop()
	return err
}

func (d *Device) cptr() *C.hid_device {
	if d.in == nil {
		return nil
	}
	return d.in.dev
}

type hidDevice struct {
	ref hidRef
	dev *C.hid_device
}

func (d *hidDevice) Close() error {
	if d.dev != nil {
		C.hid_close(d.dev)
		d.dev = nil
	}

	hidMutex.Lock()
	defer hidMutex.Unlock()
	return d.ref.Close()
}
