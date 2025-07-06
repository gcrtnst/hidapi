package hidapi

// #include <stdlib.h>
// #include <hidapi/hidapi.h>
import "C"
import (
	"runtime"
	"sync"
	"unsafe"
)

type Device struct {
	mu  sync.Mutex
	cln runtime.Cleanup
	in  *hidDevice
}

func OpenPath(path string) (*Device, error) {
	hidMutex.Lock()
	defer hidMutex.Unlock()

	ref, err := hidAcquire()
	if err != nil {
		return nil, err
	}
	defer ref.Close() // ignore error

	cstr := C.CString(path)
	defer C.free(unsafe.Pointer(cstr))

	cdev := C.hid_open_path(cstr)
	if cdev == nil {
		return nil, hidError(nil)
	}

	return newDevice(cdev)
}

func newDevice(cdev *C.hid_device) (*Device, error) {
	dev := &Device{}
	dev.in = &hidDevice{dev: cdev}

	cleanup := func(in *hidDevice) { go in.Close() } // ignore error
	dev.cln = runtime.AddCleanup(dev, cleanup, dev.in)

	ref, err := hidAcquire()
	if err != nil {
		return nil, err
	}
	dev.in.ref = ref

	return dev, nil
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
