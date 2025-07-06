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
