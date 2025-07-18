package hidapi

// #include <hidapi/hidapi.h>
import "C"
import (
	"math"
	"sync"
	"sync/atomic"
)

var (
	hidRefMu  sync.Mutex
	hidRefCnt = 0
)

type hidRef struct {
	ok uint32
}

func hidAcquire() (hidRef, error) {
	hidRefMu.Lock()
	defer hidRefMu.Unlock()

	if hidRefCnt < 0 {
		panic("hidapi: hidCount < 0")
	}
	if hidRefCnt >= math.MaxInt {
		panic("hidapi: hidCount >= math.MaxInt")
	}
	if hidRefCnt == 0 {
		err := hidInit()
		if err != nil {
			return hidRef{}, err
		}
	}
	hidRefCnt++

	ref := hidRef{ok: 1}
	return ref, nil
}

func (ref *hidRef) Close() error {
	if atomic.SwapUint32(&ref.ok, 0) == 0 {
		return nil
	}

	hidRefMu.Lock()
	defer hidRefMu.Unlock()

	if hidRefCnt <= 0 {
		panic("hidapi: hidCount <= 0")
	}
	hidRefCnt--
	if hidRefCnt == 0 {
		_ = hidExit()
	}

	return nil
}
