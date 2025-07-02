package hidapi

// #include <hidapi/hidapi.h>
import "C"
import "sync"

var (
	hidMutex sync.Mutex
	hidCount = 0
)

func hidEnter() error {
	hidMutex.Lock()
	defer hidMutex.Unlock()

	if hidCount < 0 {
		panic("hidapi: hidCount < 0")
	}
	if hidCount == 0 {
		err := hidInit()
		if err != nil {
			return err
		}
	}
	hidCount++
	return nil
}

func hidLeave() {
	hidMutex.Lock()
	defer hidMutex.Unlock()

	hidCount--
	if hidCount < 0 {
		panic("hidapi: hidCount < 0")
	}
	if hidCount == 0 {
		_ = hidExit()
	}
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
