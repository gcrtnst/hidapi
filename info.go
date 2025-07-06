package hidapi

// #include <hidapi/hidapi.h>
import "C"

type BusType int

const (
	BusUnknown   BusType = C.HID_API_BUS_UNKNOWN
	BusUSB       BusType = C.HID_API_BUS_USB
	BusBluetooth BusType = C.HID_API_BUS_BLUETOOTH
	BusI2C       BusType = C.HID_API_BUS_I2C
	BusSPI       BusType = C.HID_API_BUS_SPI
)

type DeviceInfo struct {
	Path               string
	VendorID           uint16
	ProductID          uint16
	SerialNumber       string
	ReleaseNumber      uint16
	ManufacturerString string
	ProductString      string
	UsagePage          uint16
	Usage              uint16
	InterfaceNumber    int
	BusType            BusType
}

func Enumerate(vendorID uint16, productID uint16) ([]DeviceInfo, error) {
	hidMutex.Lock()
	defer hidMutex.Unlock()

	ref, err := hidAcquire()
	if err != nil {
		return nil, err
	}
	defer ref.Close() // ignore error

	cvid := C.ushort(vendorID)
	cpid := C.ushort(productID)

	cdevs := C.hid_enumerate(cvid, cpid)
	if cdevs == nil {
		return nil, hidError(nil)
	}
	defer C.hid_free_enumeration(cdevs)

	devs := convertDeviceInfoListFromC(cdevs)
	return devs, nil
}

func convertDeviceInfoListFromC(c *C.struct_hid_device_info) []DeviceInfo {
	devsCap := 0
	for p := c; p != nil; p = p.next {
		devsCap++
	}

	devs := make([]DeviceInfo, 0, devsCap)
	for p := c; p != nil; p = p.next {
		devi := convertDeviceInfoFromC(p)
		devs = append(devs, devi)
	}
	return devs
}

func convertDeviceInfoFromC(c *C.struct_hid_device_info) DeviceInfo {
	return DeviceInfo{
		Path:               C.GoString(c.path),
		VendorID:           uint16(c.vendor_id),
		ProductID:          uint16(c.product_id),
		SerialNumber:       convertWCharPtrToString(c.serial_number),
		ReleaseNumber:      uint16(c.release_number),
		ManufacturerString: convertWCharPtrToString(c.manufacturer_string),
		ProductString:      convertWCharPtrToString(c.product_string),
		UsagePage:          uint16(c.usage_page),
		Usage:              uint16(c.usage),
		InterfaceNumber:    int(c.interface_number),
		BusType:            BusType(c.bus_type),
	}
}
