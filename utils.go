package tox

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"unsafe"
)

func safeptr(b []byte) unsafe.Pointer {
	h := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return unsafe.Pointer(h.Data)
}

func toxerr(errno interface{}) error {
	return errors.New(fmt.Sprintf("toxcore error: %v", errno))
}

func toxerrf(f string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(f, args...))
}

var toxdebug = false

func SetDebug(debug bool) {
	toxdebug = debug
}

var loglevel = 0

func SetLogLevel(level int) {
	loglevel = level
}

func FileExist(fname string) bool {
	_, err := os.Stat(fname)
	if err != nil {
		return false
	}
	return true
}

// the go-toxcore-c has data lost problem
// we need first write tmp file, and if ok, then mv to real file
func (this *Tox) WriteSavedata(fname string) error {
	if !FileExist(fname) {
		err := ioutil.WriteFile(fname, this.GetSavedata(), 0755)
		if err != nil {
			return err
		}
	} else {
		data, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}
		liveData := this.GetSavedata()
		if bytes.Compare(data, liveData) != 0 {
			tfp, err := ioutil.TempFile(filepath.Dir(fname), "gotcb")
			if err != nil {
				return err
			}
			if _, err := tfp.Write(liveData); err != nil {
				return err
			}
			tfname := tfp.Name()
			if err := tfp.Close(); err != nil {
				return err
			}
			if err := os.Remove(fname); err != nil {
				return err
			}
			if err := os.Rename(tfname, fname); err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *Tox) LoadSavedata(fname string) ([]byte, error) {
	return ioutil.ReadFile(fname)
}

func LoadSavedata(fname string) ([]byte, error) {
	return ioutil.ReadFile(fname)
}

func ConnStatusString(status int) (s string) {
	switch status {
	case CONNECTION_NONE:
		s = "CONNECTION_NONE"
	case CONNECTION_TCP:
		s = "CONNECTION_TCP"
	case CONNECTION_UDP:
		s = "CONNECTION_UDP"
	}
	return
}
