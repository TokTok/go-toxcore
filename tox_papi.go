package tox

/*
#include <stdlib.h>
#include <string.h>

#include "tox_papi.h"

*/
import "C"
import (
	"encoding/hex"
	"strings"
	"unsafe"
)

type self_get_address_param struct {
	gosz   uintptr
	t      unsafe.Pointer
	retval bool
	addr   *byte
}

func new_self_get_address_param() *self_get_address_param {
	this := &self_get_address_param{}
	this.gosz = uintptr(unsafe.Sizeof(*this))
	return this
}

func (this *Tox) toxptr() unsafe.Pointer { return unsafe.Pointer(this.toxcore) }

func (this *Tox) SelfGetAddress_p() string {
	var addr [ADDRESS_SIZE]byte
	prm := new_self_get_address_param()
	prm.t = this.toxptr()
	prm.addr = &addr[0]
	// log.Println(prm.gosz)
	// os.Exit(0)

	C.tox_self_get_address_p(C.uintptr_t(uintptr(unsafe.Pointer(prm))))

	return strings.ToUpper(hex.EncodeToString(addr[:]))
}

type friend_get_name_param struct {
	gosz         uintptr
	t            unsafe.Pointer
	retval       bool
	friendNumber uint32
	name         *byte
	error        int
	namelen      int
}

func new_friend_get_name_param() *friend_get_name_param {
	this := &friend_get_name_param{}
	this.gosz = uintptr(unsafe.Sizeof(*this))
	return this
}

func (this *Tox) FriendGetName_p(friendNumber uint32) (string, error) {
	prm := new_friend_get_name_param()
	prm.t = this.toxptr()
	prm.friendNumber = friendNumber
	var name [MAX_NAME_LENGTH]byte
	prm.name = &name[0]

	C.tox_friend_get_name_p(C.uintptr_t(uintptr(unsafe.Pointer(prm))))
	if !prm.retval {
		// error
	}
	return string(name[:prm.namelen]), nil
}

type friend_get_public_key_param struct {
	gosz         uintptr
	t            unsafe.Pointer
	retval       bool
	friendNumber uint32
	pubkey       *byte
	error        int
}

func new_friend_get_public_key_param() *friend_get_public_key_param {
	this := &friend_get_public_key_param{}
	this.gosz = uintptr(unsafe.Sizeof(*this))
	return this
}

func (this *Tox) FriendGetPublicKey_p(friendNumber uint32) (string, error) {
	prm := new_friend_get_public_key_param()
	prm.t = this.toxptr()
	prm.friendNumber = friendNumber
	var pubkey [PUBLIC_KEY_SIZE]byte
	prm.pubkey = &pubkey[0]

	C.tox_friend_get_public_key_p(C.uintptr_t(uintptr(unsafe.Pointer(prm))))
	if !prm.retval {
		// error
	}
	return string(pubkey[:PUBLIC_KEY_SIZE]), nil
}
