package hidapi

// #include <stdlib.h>
// #include <hidapi/hidapi.h>
import "C"
import (
	"runtime"
	"sync"
	"unsafe"
)

const MaxReportDescriptorSize = C.HID_API_MAX_REPORT_DESCRIPTOR_SIZE

type Device struct {
	mu  sync.Mutex
	cln runtime.Cleanup
	in  *hidDevice
}

func Open(vendorID uint16, productID uint16, serialNumber string) (*Device, error) {
	ref, err := hidAcquire()
	if err != nil {
		return nil, err
	}
	defer ref.Close()

	cvid := C.ushort(vendorID)
	cpid := C.ushort(productID)
	cser := convertStringToWCharPtr(serialNumber)

	cdev := C.hid_open(cvid, cpid, cser)
	if cdev == nil {
		return nil, hidError(nil)
	}

	return newDevice(cdev)
}

func OpenPath(path string) (*Device, error) {
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

	return d.ref.Close()
}
