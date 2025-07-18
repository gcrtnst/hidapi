package hidapi

// #include <hidapi/hidapi.h>
import "C"
import (
	"errors"
	"runtime"
)

var hidInitReqCh chan<- *hidInitRequest

func init() {
	if hidBackendInitThread {
		t := &hidInitThread{}
		reqCh := make(chan *hidInitRequest)
		go t.run(reqCh)
		hidInitReqCh = reqCh
	}
}

func hidInit() error {
	if hidInitReqCh == nil {
		return hidInitRaw()
	}

	errCh := make(chan error)
	req := &hidInitRequest{
		target: true,
		errCh:  errCh,
	}
	hidInitReqCh <- req
	return <-errCh
}

func hidExit() error {
	if hidInitReqCh == nil {
		return hidExitRaw()
	}

	errCh := make(chan error)
	req := &hidInitRequest{
		target: false,
		errCh:  errCh,
	}
	hidInitReqCh <- req
	return <-errCh
}

type hidInitRequest struct {
	target bool
	errCh  chan<- error
}

type hidInitThread struct {
	thread bool
	state  bool
}

func (t *hidInitThread) run(reqCh <-chan *hidInitRequest) {
	defer t.hidExit()

	for req := range reqCh {
		t.handleRequest(req)
	}
}

func (t *hidInitThread) handleRequest(req *hidInitRequest) {
	if req == nil {
		return
	}

	var err error
	if req.target {
		err = t.hidInit()
	} else {
		err = t.hidExit()
	}

	req.errCh <- err
	close(req.errCh)
}

func (t *hidInitThread) hidInit() (err error) {
	if t.state {
		err = nil
		return
	}

	t.threadLock()
	err = hidInitRaw()
	t.state = err == nil
	return
}

func (t *hidInitThread) hidExit() error {
	if !t.state {
		return nil
	}

	if !t.thread {
		panic("hidapi: thread not locked when hidExit called")
	}

	err := hidExitRaw()
	if err != nil {
		return err
	}

	t.state = false
	return nil
}

func (t *hidInitThread) threadLock() {
	if !t.thread {
		runtime.LockOSThread()
		t.thread = true
	}
}

func hidInitRaw() error {
	ret := C.hid_init()
	if ret != 0 {
		return hidError(nil)
	}
	return nil
}

func hidExitRaw() error {
	ret := C.hid_exit()
	if ret != 0 {
		return errors.New("hid_exit failed")
	}
	return nil
}
