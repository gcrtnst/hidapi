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
	dev *C.hid_device
}

func (d *Device) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.dev != nil {
		C.hid_close(d.dev)
		d.dev = nil
	}

	d.cln.Stop()
	return nil
}
