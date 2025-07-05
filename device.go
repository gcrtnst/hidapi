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
	dev *hidDevice
}

func (d *Device) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	var err error
	if d.dev != nil {
		err = d.dev.Close()
	}

	d.cln.Stop()
	return err
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

	return d.ref.Close()
}
