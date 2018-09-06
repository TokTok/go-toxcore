package tox

import (
	"encoding/hex"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

// `go test -v -run Covers` will show untested functions
// TODO boundary value testing

var bsnodes = []string{
	"tox.ngc.zone", "33445", "15E9C309CFCB79FDDF0EBA057DABB49FE15F3803B1BFF06536AE2E5BA5E4690E",
	"node.tox.biribiri.org", "33445", "F404ABAA1C99A9D37D61AB54898F56793E1DEF8BD46B1038B9D822E8460FAB67",
	"178.62.250.138", "33445", "788236D34978D1D5BD822F0A5BEBD2C53C64CC31CD3149350EE27D4D9A2F9B6B",
	"198.98.51.198", "33445", "1D5A5F2F5D6233058BF0259B09622FB40B482E4FA0931EB8FD3AB8E7BF7DAF6F",
}

func init() {
	log.SetFlags(log.Flags() | log.Lshortfile)
}

func TestCreate(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		_t := NewTox(nil)
		if _t == nil {
			t.Error("nil")
		}
		_t.Kill()
	})
	t.Run("default options", func(t *testing.T) {
		opts := NewToxOptions()
		_t := NewTox(opts)
		if _t == nil {
			t.Error("nil")
		}
		_t.Kill()
	})
	t.Run("tcp options", func(t *testing.T) {
		opts := NewToxOptions()
		opts.Tcp_port = 44577
		_t := NewTox(opts)
		if _t == nil {
			t.Error("nil")
		}
		_t.Kill()
	})
	t.Run("tcp conflict", func(t *testing.T) {
		opts := NewToxOptions()
		opts.Tcp_port = 44587
		_t, _t2 := NewTox(opts), NewTox(opts)
		if _t == nil || _t2 != nil {
			t.Error("should non-nil/nil", _t, _t2)
		}
		_t.Kill()
		_t2.Kill()
	})
	t.Run("save profile", func(t *testing.T) {
		_t := NewTox(nil)
		sz := _t.GetSavedataSize()
		dat := _t.GetSavedata()
		if sz <= 0 || dat == nil || len(dat) != int(sz) {
			t.Error("cannot zero")
		}
		_t.Kill()
	})
	t.Run("load profile", func(t *testing.T) {
		_t := NewTox(nil)
		dat := _t.GetSavedata()
		_t.Kill()

		opts := NewToxOptions()
		opts.Savedata_data = dat
		opts.Savedata_type = SAVEDATA_TYPE_TOX_SAVE
		_t2 := NewTox(opts)
		dat2 := _t2.GetSavedata()
		if len(dat2) != len(dat) || string(dat2) != string(dat) {
			t.Error("must ==")
		}
		_t2.Kill()
	})
	t.Run("load error profile", func(t *testing.T) {
		_t := NewTox(nil)
		dat := _t.GetSavedata()
		addr := _t.SelfGetAddress()
		_t.Kill()

		opts := NewToxOptions()
		opts.Savedata_data = append([]byte("set-broken"), dat...)
		opts.Savedata_type = SAVEDATA_TYPE_TOX_SAVE
		_t2 := NewTox(opts)
		if _t2 == nil {
			t.Error("must non-nil")
		}
		if addr == _t2.SelfGetAddress() {
			t.Error("must !=", addr, _t2.SelfGetAddress())
		}
		_t2.Kill()
	})
	t.Run("load seckey", func(t *testing.T) {
		_t := NewTox(nil)
		addr := _t.SelfGetAddress()
		seckey := _t.SelfGetSecretKey()
		_t.Kill()

		opts := NewToxOptions()
		opts.Savedata_type = SAVEDATA_TYPE_SECRET_KEY
		binsk, _ := hex.DecodeString(seckey)
		opts.Savedata_data = binsk
		_t2 := NewTox(opts)
		if _t2.SelfGetSecretKey() != seckey {
			t.Error("must =")
		}
		if _t2.SelfGetAddress()[0:PUBLIC_KEY_SIZE*2] != addr[0:PUBLIC_KEY_SIZE*2] {
			t.Error("must =", _t2.SelfGetAddress(), addr)
		}
	})
	t.Run("destroy", func(t *testing.T) {
		_t := NewTox(nil)
		_t.Kill()
		if _t.toxcore != nil {
			t.Error("must nil")
		}
	})
}

func TestBase(t *testing.T) {
	_t := NewTox(nil)
	defer _t.Kill()

	t.Run("name", func(t *testing.T) {
		if _t.SelfGetName() != "" {
			t.Error("must empty")
		}
		if _t.SelfGetNameSize() != 0 {
			t.Error("must zero")
		}
		tname := "test name"
		if err := _t.SelfSetName(tname); err != nil {
			t.Error(err)
		}
		if size := _t.SelfGetNameSize(); size != len(tname) {
			t.Error("must =", size, len(tname))
		}
		tname = strings.Repeat("n", MAX_NAME_LENGTH)
		if err := _t.SelfSetName(tname); err != nil {
			t.Error(err)
		}
		tname = strings.Repeat("n", MAX_NAME_LENGTH+1)
		if err := _t.SelfSetName(tname); err == nil {
			t.Error("must failed", err)
		}
	})
	t.Run("local status", func(t *testing.T) {
		if _t.SelfGetStatusMessageSize() != 0 {
			t.Error("must zero")
		}
		if stm, err := _t.SelfGetStatusMessage(); err != nil || len(stm) != 0 {
			t.Error("must empty", stm, err)
		}
		tmsg := "test status msg"
		if ok, err := _t.SelfSetStatusMessage(tmsg); !ok || err != nil {
			t.Error("must ok", err)
		}
		if stm, err := _t.SelfGetStatusMessage(); err != nil || stm != tmsg {
			t.Error("must =", stm, err)
		}
		tmsg = strings.Repeat("s", MAX_STATUS_MESSAGE_LENGTH)
		if ok, err := _t.SelfSetStatusMessage(tmsg); !ok || err != nil {
			t.Error("must ok", err)
		}
		tmsg = strings.Repeat("s", MAX_STATUS_MESSAGE_LENGTH+1)
		if ok, err := _t.SelfSetStatusMessage(tmsg); ok || err == nil {
			t.Error("must failed", err)
		}
		if _t.SelfGetConnectionStatus() != CONNECTION_NONE {
			t.Error("must none")
		}
	})
	t.Run("address/pubkey", func(t *testing.T) {
		addr := _t.SelfGetAddress()
		if len(addr) != ADDRESS_SIZE*2 {
			t.Error("size")
		}
		pubkey := _t.SelfGetPublicKey()
		if len(pubkey) != PUBLIC_KEY_SIZE*2 {
			t.Error("size")
		}
		if addr[0:len(pubkey)] != pubkey {
			t.Error(addr)
		}
	})
	t.Run("seckey", func(t *testing.T) {
		seckey := _t.SelfGetSecretKey()
		if len(seckey) != SECRET_KEY_SIZE*2 {
			t.Error("size")
		}
	})
	t.Run("nospam", func(t *testing.T) {
	})
}

func TestBootstrap(t *testing.T) {
	_t := NewTox(nil)
	defer _t.Kill()
	port, _ := strconv.Atoi(bsnodes[1])

	t.Run("success", func(t *testing.T) {
		if ok, err := _t.Bootstrap(bsnodes[0], uint16(port), bsnodes[2]); !ok || err != nil {
			t.Error("must ok", ok, err)
		}
	})
	t.Run("failed", func(t *testing.T) {
		brkey := bsnodes[2]
		brkey = "XYZAB" + bsnodes[2][3:]
		if ok, err := _t.Bootstrap(bsnodes[0], uint16(port), brkey); ok || err == nil {
			t.Error("must failed", ok, err)
		}
		if ok, err := _t.Bootstrap("a.b.c.d", uint16(port), bsnodes[2]); ok || err == nil {
			t.Error("must failed", ok, err)
		}
	})
	t.Run("relay", func(t *testing.T) {
		if ok, err := _t.AddTcpRelay(bsnodes[0], uint16(port), bsnodes[2]); !ok || err != nil {
			t.Error("must ok", ok, err)
		}
		if ok, err := _t.AddTcpRelay("a.b.c.d", uint16(port), bsnodes[2]); ok || err == nil {
			t.Error("must failed", ok, err)
		}
	})
}

type MiniTox struct {
	t      *Tox
	stopch chan struct{}
}

func NewMiniTox() *MiniTox {
	minitox := &MiniTox{}
	opts := NewToxOptions()
	opts.Local_discovery_enabled = false
	minitox.t = NewTox(opts)
	minitox.stopch = make(chan struct{}, 0)
	return minitox
}

func (minitox *MiniTox) Iterate() {
	tickch := time.Tick(100 * time.Millisecond)
	for {
		select {
		case <-tickch:
			minitox.t.Iterate()
		case <-minitox.stopch:
			return
		}
	}
}

func (minitox *MiniTox) bootstrap() {
	for idx := 0; idx < len(bsnodes)/3; idx++ {
		port, err := strconv.Atoi(bsnodes[1+idx*3])
		_, err = minitox.t.Bootstrap(bsnodes[0+idx*3], uint16(port), bsnodes[2+idx*3])
		if err != nil {
		}
		_, err = minitox.t.AddTcpRelay(bsnodes[0+idx*3], uint16(port), bsnodes[2+idx*3])
		if err != nil {
		}
	}
}

func (minitox *MiniTox) stop() {
	minitox.stopch <- struct{}{}
}

var err error

func waitcond(cond func() bool, timeout int) {
	// TODO might infinite loop
	btime := time.Now()
	cnter := 0
	for {
		if cond() {
			// print("\n")
			return
		}

		etime := time.Now()
		dtime := etime.Sub(btime)
		if timeout > 0 && int(dtime.Seconds()) > timeout {
			return // timeout
		}

		if cnter%15 == 0 {
			// print(".")
		}
		cnter += 1
		time.Sleep(51 * time.Millisecond)
	}
}

func link(a *MiniTox, b *MiniTox) error {
	port, err := a.t.SelfGetUdpPort()
	if err != nil {
		return err
	}
	if ok, err := b.t.Bootstrap("localhost", port, a.t.SelfGetDhtId()); !ok || err != nil {
		return err
	}
	return nil
}

// login udp / login tcp
func TestCommunication(t *testing.T) {
	t.Run("Login/connect", func(t *testing.T) {
		t.Parallel()

		_t := NewMiniTox()
		defer _t.t.Kill()
		_t.bootstrap()
		waitcond(func() bool {
			if _t.t.IterationInterval() == 0 {
				t.Error("why")
			}
			_t.t.Iterate()
			if _t.t.SelfGetConnectionStatus() > CONNECTION_NONE {
				return true
			}
			return false
		}, 60)
		if _t.t.SelfGetConnectionStatus() == CONNECTION_NONE {
			t.Error("maybe iterate not use")
		}
	})

	t.Run("Friend/add friend", func(t *testing.T) {
		t.Parallel()

		t1 := NewMiniTox()
		t2 := NewMiniTox()
		defer t1.t.Kill()
		defer t2.t.Kill()

		if err := link(t1, t2); err != nil {
			t.Error("must ok", err)
		}

		t1.t.CallbackFriendRequest(func(_ *Tox, friendId, msg string, d interface{}) {
			_, err := t1.t.FriendAddNorequest(friendId)
			if err != nil {
				t.Fail()
			}
		}, nil)

		go t1.Iterate()
		go t2.Iterate()
		defer t1.stop()
		defer t2.stop()

		waitcond(func() bool {
			return t1.t.SelfGetConnectionStatus() != CONNECTION_NONE && t2.t.SelfGetConnectionStatus() != CONNECTION_NONE
		}, 100)
		friendNumber, err := t2.t.FriendAdd(t1.t.SelfGetAddress(), "hoho")
		if err != nil {
			t.Error(err, friendNumber)
		}
		_, err = t2.t.FriendAdd(t1.t.SelfGetAddress(), "hehe")
		if err == nil {
			t.Error(err)
		}
		if t2.t.SelfGetFriendListSize() != 1 {
			t.Error("friend size not match")
		}
		lst := t2.t.SelfGetFriendList()
		if len(lst) != 1 {
			t.Error("friend list not match")
		}

		friendNumber2, err := t2.t.FriendByPublicKey(t1.t.SelfGetAddress())
		if err != nil {
			t.Error(err)
		}
		if friendNumber2 != friendNumber {
			t.Error("friend number not match")
		}
		pubkey, err := t2.t.FriendGetPublicKey(friendNumber)
		if err != nil {
			t.Error(err, pubkey)
		}
		if pubkey != t1.t.SelfGetPublicKey() {
			t.Error("friend pubkey not match")
		}
		if !t2.t.FriendExists(friendNumber) {
			t.Error("added friend not exists")
		}
	})

	t.Run("Friend/friend status", func(t *testing.T) {
		t.Parallel()

		t1 := NewMiniTox()
		t2 := NewMiniTox()
		defer t1.t.Kill()
		defer t2.t.Kill()

		if err := link(t1, t2); err != nil {
			t.Error("must ok", err)
		}

		t1.t.CallbackFriendRequest(func(_ *Tox, friendId, msg string, d interface{}) {
			t1.t.FriendAddNorequest(friendId)
		}, nil)

		// testing
		t1.t.CallbackFriendConnectionStatus(func(_ *Tox, friendNumber uint32, status int,
			d interface{}) {
		}, nil)
		t1nameChanged := false
		t2.t.CallbackFriendName(func(_ *Tox, friendNumber uint32, name string, d interface{}) {
			if len(name) > 0 {
				t1nameChanged = true
			}
		}, nil)
		t1statusMessageChanged := false
		t2.t.CallbackFriendStatusMessage(func(_ *Tox, friendNumber uint32, stmsg string, d interface{}) {
			if len(stmsg) > 0 {
				t1statusMessageChanged = true
			}
		}, nil)

		go t1.Iterate()
		go t2.Iterate()
		defer t1.stop()
		defer t2.stop()

		waitcond(func() bool {
			return t1.t.SelfGetConnectionStatus() != CONNECTION_NONE && t2.t.SelfGetConnectionStatus() != CONNECTION_NONE
		}, 100)
		friendNumber, _ := t2.t.FriendAdd(t1.t.SelfGetAddress(), "hoho")

		waitcond(func() bool {
			return 1 == t1.t.SelfGetFriendListSize()
		}, 100)
		waitcond(func() bool {
			status, err := t2.t.FriendGetConnectionStatus(friendNumber)
			if err != nil {
				t.Error(err, status)
				return false
			}
			return status > CONNECTION_NONE
		}, 100)
		if status, err := t2.t.FriendGetConnectionStatus(friendNumber); err != nil || status == CONNECTION_NONE {
			t.Error(err, status)
		}

		err = t1.t.SelfSetName("t1")
		if err != nil {
			t.Error(err)
		}
		waitcond(func() bool { return t1nameChanged }, 100)
		t1name, err := t2.t.FriendGetName(friendNumber)
		t1size, err := t2.t.FriendGetNameSize(friendNumber)
		if err != nil {
			t.Error(err)
		}
		if t1name != "t1" {
			t.Error(t1name)
		}
		if t1size != len(t1name) {
			t.Error(t1size, t1name)
		}
		_, err = t1.t.SelfSetStatusMessage("t1status")
		if err != nil {
			t.Error(err)
		}
		waitcond(func() bool { return t1statusMessageChanged }, 100)
		t1stmsg, err := t2.t.FriendGetStatusMessage(friendNumber)
		if err != nil {
			t.Error(err)
		}
		if t1stmsg != "t1status" {
			t.Error(t1stmsg, t1stmsg != "t1status")
		}
		t1stmsgsz, err := t2.t.FriendGetStatusMessageSize(friendNumber)
		if err != nil {
			t.Error(err)
		}
		if t1stmsgsz != len("t1status") {
			t.Error(t1stmsgsz, len("t1status"))
		}

		t1st, err := t2.t.FriendGetStatus(friendNumber)
		if err != nil {
			t.Error(err)
		}
		if t1st != USER_STATUS_NONE {
			t.Error(t1st)
		}
	})

	t.Run("Friend/friend message", func(t *testing.T) {
		t.Parallel()

		t1 := NewMiniTox()
		t2 := NewMiniTox()
		defer t1.t.Kill()
		defer t2.t.Kill()

		if err := link(t1, t2); err != nil {
			t.Error("must ok", err)
		}

		t1.t.CallbackFriendRequest(func(_ *Tox, friendId, msg string, d interface{}) {
			t1.t.FriendAddNorequest(friendId)
		}, nil)
		recvmsg := ""
		t1.t.CallbackFriendMessage(func(_ *Tox, friendNumber uint32, msg string, d interface{}) {
			recvmsg = msg
		}, nil)

		go t1.Iterate()
		go t2.Iterate()
		defer t1.stop()
		defer t2.stop()

		waitcond(func() bool {
			return t1.t.SelfGetConnectionStatus() != CONNECTION_NONE && t2.t.SelfGetConnectionStatus() != CONNECTION_NONE
		}, 100)
		friendNumber, _ := t2.t.FriendAdd(t1.t.SelfGetAddress(), "hoho")
		waitcond(func() bool {
			return 1 == t1.t.SelfGetFriendListSize()
		}, 100)
		waitcond(func() bool {
			status, _ := t2.t.FriendGetConnectionStatus(friendNumber)
			return status > CONNECTION_NONE
		}, 100)
		_, err := t2.t.FriendSendMessage(friendNumber, "hohoo")
		if err != nil {
			t.Error(err)
		}
		waitcond(func() bool {
			return len(recvmsg) > 0
		}, 100)
		if recvmsg != "hohoo" {
			t.Error("send/recv message failed")
		}
		_, err = t2.t.FriendSendAction(friendNumber, "actfoo")
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Friend/friend delete", func(t *testing.T) {
		t.Parallel()

		t1 := NewMiniTox()
		t2 := NewMiniTox()
		defer t1.t.Kill()
		defer t2.t.Kill()

		if err := link(t1, t2); err != nil {
			t.Error("must ok", err)
		}

		t1.t.CallbackFriendRequest(func(_ *Tox, friendId, msg string, d interface{}) {
			t1.t.FriendAddNorequest(friendId)
		}, nil)

		go t1.Iterate()
		go t2.Iterate()
		defer t1.stop()
		defer t2.stop()

		waitcond(func() bool {
			return t1.t.SelfGetConnectionStatus() != CONNECTION_NONE && t2.t.SelfGetConnectionStatus() != CONNECTION_NONE
		}, 100)
		friendNumber, _ := t2.t.FriendAdd(t1.t.SelfGetAddress(), "hoho")
		waitcond(func() bool {
			return 1 == t1.t.SelfGetFriendListSize()
		}, 100)
		_, err = t2.t.FriendDelete(friendNumber)
		if err != nil {
			t.Error(err)
		}
		if t2.t.FriendExists(friendNumber) {
			t.Error("deleted friend appearence")
		}
		_, err = t2.t.FriendDelete(friendNumber)
		if err == nil {
			t.Error("delete deleted friend should failed")
		}
	})

	t.Run("Group/add del", func(t *testing.T) {
		t.Parallel()

		t1 := NewMiniTox()
		t2 := NewMiniTox()
		defer t1.t.Kill()
		defer t2.t.Kill()

		if err := link(t1, t2); err != nil {
			t.Error("must ok", err)
		}

		t1.t.CallbackFriendRequest(func(_ *Tox, friendId, msg string, d interface{}) {
			t1.t.FriendAddNorequest(friendId)
		}, nil)

		go t1.Iterate()
		go t2.Iterate()
		defer t1.stop()
		defer t2.stop()

		waitcond(func() bool {
			return t1.t.SelfGetConnectionStatus() != CONNECTION_NONE
		}, 100)
		gn, err := t1.t.AddGroupChat()
		if err != nil || gn != 0 {
			t.Error(err)
		}
		_, err = t1.t.DelGroupChat(gn)
		if err != nil {
			t.Error(err)
		}
		if n := t1.t.CountChatList(); n != 0 {
			t.Error(n)
		}
		if len(t1.t.GetChatList()) != 0 {
			t.Error("should 0")
		}
		var gcnt = 5
		for idx := 0; idx < gcnt; idx++ {
			gn, err = t1.t.AddGroupChat()
			if gn != idx {
				t.Error(gn, idx)
			}
			title := fmt.Sprintf("group%d", idx)
			_, err = t1.t.GroupSetTitle(gn, title)
			if err != nil {
				t.Error(err)
			}
			ntitle, err := t1.t.GroupGetTitle(gn)
			if err != nil {
				t.Error(err)
			}
			if ntitle != title {
				t.Error(ntitle, title)
			}
			names := t1.t.GroupGetNames(gn)
			if len(names) != 1 {
				t.Error(len(names), 1)
			}
			pubkeys := t1.t.GroupGetPeerPubkeys(gn)
			if len(pubkeys) != 1 {
				t.Error(len(names), 1)
			}
			gtype, err := t1.t.GroupGetType(uint32(gn))
			if err != nil {
				t.Error(err)
			}
			if uint8(gtype) != CONFERENCE_TYPE_TEXT {
				t.Error(gtype, CONFERENCE_TYPE_TEXT)
			}
			if t1.t.GroupNumberPeers(gn) != 1 {
				t.Error(1)
			}
			pname, err := t1.t.GroupPeerName(gn, 0)
			if err != nil {
				t.Error(err)
			}
			if len(pname) != 0 {
				t.Error(pname)
			}
			pubkey, err := t1.t.GroupPeerPubkey(gn, 0)
			if err != nil {
				t.Error(err)
			}
			if !strings.HasPrefix(t1.t.SelfGetAddress(), pubkey) {
				t.Error("get peer pubkey")
			}
			if !t1.t.GroupPeerNumberIsOurs(gn, 0) {
				t.Error("ours")
			}
			if t1.t.GroupPeerNumberIsOurs(gn, 789) {
				t.Error("not ours")
			}
			_, err = t1.t.GroupActionSend(gn, "abc")
			if err == nil {
				t.Error("should not nil")
			}
			_, err = t1.t.GroupMessageSend(gn, "abc")
			if err == nil {
				t.Error("should not nil")
			}
			peers := t1.t.GroupGetPeers(gn)
			if len(peers) != 1 {
				t.Error("should 1")
			}
			if _, err = t1.t.JoinGroupChat(5, ""); err == nil {
				t.Error("should not nil")
			}
			if _, err = t1.t.InviteFriend(123, gn); err == nil {
				t.Error("should nil")
			}
			if cnt := t1.t.CountChatList(); int(cnt) != idx+1 {
				t.Error(cnt, idx+1)
			}
			if grps := t1.t.GetChatList(); len(grps) != idx+1 {
				t.Error(len(grps), idx+1)
			}
		}
		grps := t1.t.GetChatList()
		if len(grps) != gcnt {
			t.Error(len(grps), gcnt)
		}
		if t1.t.CountChatList() != uint32(gcnt) {
			t.Error(t1.t.CountChatList(), gcnt)
		}
	})

	t.Run("Group/group invite", func(t *testing.T) {
		t.Parallel()

		t1 := NewMiniTox()
		t2 := NewMiniTox()
		defer t1.t.Kill()
		defer t2.t.Kill()

		if err := link(t1, t2); err != nil {
			t.Error("must ok", err)
		}

		t1.t.CallbackFriendRequest(func(_ *Tox, friendId, msg string, d interface{}) {
			t1.t.FriendAddNorequest(friendId)
		}, nil)

		t1.t.CallbackConferenceInvite(func(_ *Tox, friendNumber uint32, itype uint8, data string, ud interface{}) {
			switch itype {
			case CONFERENCE_TYPE_TEXT:
				_, err := t1.t.JoinGroupChat(friendNumber, data)
				if err != nil {
					t.Error(err)
				}
			case CONFERENCE_TYPE_AV:
				_, err := t1.t.JoinAVGroupChat(friendNumber, data, nil)
				if err != nil {
					t.Error(err)
				}
			}
		}, nil)

		groupNameChangeTimes := 0
		t2.t.CallbackConferencePeerListChanged(func(_ *Tox, groupNumber uint32, ud interface{}) {
			groupNameChangeTimes += 1
		}, nil)

		go t1.Iterate()
		go t2.Iterate()
		defer t1.stop()
		defer t2.stop()

		waitcond(func() bool {
			return t1.t.SelfGetConnectionStatus() != CONNECTION_NONE && t2.t.SelfGetConnectionStatus() != CONNECTION_NONE
		}, 100)

		t2.t.FriendAdd(t1.t.SelfGetAddress(), "autotests")
		waitcond(func() bool {
			return t1.t.SelfGetFriendListSize() == 1
		}, 100)

		fn, _ := t2.t.FriendByPublicKey(t1.t.SelfGetPublicKey())
		gn, _ := t2.t.AddGroupChat()

		// must wait friend online and can call InviteFriend
		waitcond(func() bool {
			st, _ := t2.t.FriendGetConnectionStatus(fn)
			return st > CONNECTION_NONE
		}, 100)

		_, err = t2.t.InviteFriend(fn, gn)
		if err != nil {
			t.Error(err)
		}
		if err != nil {
			t.Error(err)
		}
		waitcond(func() bool {
			return t1.t.CountChatList() == 1
		}, 100)
		if t1.t.CountChatList() != 1 {
			t.Error("must 1 chat", t1.t.CountChatList())
		}
		if t2.t.CountChatList() != 1 {
			t.Error("must 1 chat", t2.t.CountChatList())
		}
		waitcond(func() bool {
			return t1.t.GroupNumberPeers(gn) > 0
		}, 100)

		if _, err := t1.t.DelGroupChat(gn); err != nil {
			t.Error(err)
		}
		if _, err := t2.t.DelGroupChat(gn); err != nil {
			t.Error(err)
		}

		if groupNameChangeTimes == 0 {
			t.Error("must > 0")
		}
	})

	t.Run("Group/group message", func(t *testing.T) {
		t.Parallel()

		t1 := NewMiniTox()
		t2 := NewMiniTox()
		defer t1.t.Kill()
		defer t2.t.Kill()

		if err := link(t1, t2); err != nil {
			t.Error("must ok", err)
		}

		t1.t.CallbackFriendRequest(func(_ *Tox, friendId, msg string, d interface{}) {
			t1.t.FriendAddNorequest(friendId)
		}, nil)

		t1.t.CallbackConferenceInvite(func(_ *Tox, friendNumber uint32, itype uint8, data string, ud interface{}) {
			switch itype {
			case CONFERENCE_TYPE_TEXT:
				t1.t.JoinGroupChat(friendNumber, data)
			case CONFERENCE_TYPE_AV:
				t1.t.JoinAVGroupChat(friendNumber, data, nil)
			}
		}, nil)

		recved_act := ""
		recved_msg := ""
		t1.t.CallbackGroupMessage(func(_ *Tox, groupNumber, peerNumber int, msg string, ud interface{}) {
			recved_msg = msg
		}, nil)
		t1.t.CallbackGroupAction(func(_ *Tox, groupNumber, peerNumber int, msg string, ud interface{}) {
			recved_act = msg
		}, nil)

		go t1.Iterate()
		go t2.Iterate()
		defer t1.stop()
		defer t2.stop()

		waitcond(func() bool {
			return t1.t.SelfGetConnectionStatus() != CONNECTION_NONE && t2.t.SelfGetConnectionStatus() != CONNECTION_NONE
		}, 100)

		t2.t.FriendAdd(t1.t.SelfGetAddress(), "autotests")
		waitcond(func() bool {
			return t1.t.SelfGetFriendListSize() == 1
		}, 100)

		fn, _ := t2.t.FriendByPublicKey(t1.t.SelfGetPublicKey())
		gn, _ := t2.t.AddGroupChat()

		// must wait friend online and can call InviteFriend
		waitcond(func() bool {
			st, _ := t2.t.FriendGetConnectionStatus(fn)
			return st > CONNECTION_NONE
		}, 100)

		_, err = t2.t.InviteFriend(fn, gn)
		if err != nil {
			t.Error(err)
		}
		waitcond(func() bool {
			return t1.t.CountChatList() == 1
		}, 100)

		// must wait peer join
		waitcond(func() bool {
			return t2.t.GroupNumberPeers(gn) == 2
		}, 10)

		if _, err := t2.t.GroupMessageSend(gn, "foo123"); err != nil {
			t.Error(err)
		}
		if _, err := t2.t.GroupActionSend(gn, "bar123"); err != nil {
			t.Error(err)
		}
		waitcond(func() bool {
			return len(recved_msg) > 0 && len(recved_act) > 0
		}, 10)
		if recved_msg != "foo123" || recved_act != "bar123" {
			t.Errorf("Received msg='%s', act='%s', but wanted '%s' and '%s'",
				recved_msg, recved_act, "foo123", "bar123")
		}
	})

	t.Run("Group/issue 6", func(t *testing.T) {
		t.Parallel()

		opts := NewToxOptions()
		opts.ThreadSafe = true
		opts.Tcp_port = 34567
		_t1 := NewTox(opts)
		log.Println(_t1)
		go func() {
			for {
				_t1.Iterate()
				time.Sleep(100 * time.Millisecond)
			}
		}()

		opts2 := NewToxOptions()
		opts2.ThreadSafe = true
		opts2.Tcp_port = 34568
		_t2 := NewTox(opts2)
		log.Println(_t2)
		_t2.CallbackGroupInviteAdd(func(_ *Tox, friendNumber uint32, itype uint8, data string, userData interface{}) {
			log.Println(friendNumber, itype)
		}, nil)
		go func() {
			for {
				_t2.Iterate()
				time.Sleep(100 * time.Millisecond)
			}
		}()

		waitcond(func() bool { return _t1.IsConnected() > 0 }, 100)
		waitcond(func() bool { return _t2.IsConnected() > 0 }, 100)
		log.Println("both connected")

		gid := _t1.AddAVGroupChat(nil)
		// ok, err := _t1.DelGroupChat(gid)
		// log.Println(ok, err)
		log.Println(gid)
	})

	t.Run("File/send", func(t *testing.T) {
		t.Parallel()

		t1 := NewMiniTox()
		t2 := NewMiniTox()
		defer t1.t.Kill()
		defer t2.t.Kill()

		if err := link(t1, t2); err != nil {
			t.Error("must ok", err)
		}

		t1.t.CallbackFriendRequest(func(_ *Tox, friendId, msg string, d interface{}) {
			t1.t.FriendAddNorequest(friendId)
		}, nil)

		t1.t.CallbackFileRecv(func(_ *Tox, friendNumber uint32, fileNumber uint32,
			kind uint32, fileSize uint64, fileName string, d interface{}) {
			log.Println(fileNumber, fileSize, fileName)
			_, err := t1.t.FileSeek(friendNumber, fileNumber, 15)
			if err != nil {
				t.Error(err)
			}
			_, err = t1.t.FileControl(friendNumber, fileNumber, FILE_CONTROL_RESUME)
			if err != nil {
				t.Error(err)
			}
		}, nil)
		recvData := ""
		t1.t.CallbackFileRecvChunk(func(_ *Tox, friendNumber uint32, fileNumber uint32,
			position uint64, data []byte, d interface{}) {
			// log.Println(fileNumber, position, len(data))
			recvData += string(data)
		}, nil)
		t1.t.CallbackFileRecvControl(func(_ *Tox, friendNumber uint32, fileNumber uint32,
			control int, ud interface{}) {
			// log.Println(fileNumber, control)
		}, nil)

		t2.t.CallbackFileChunkRequest(func(_ *Tox, friend_number uint32, file_number uint32,
			position uint64, length int, d interface{}) {
			// log.Println(file_number, position, length)
			if length == 0 {
				return
			}
			s := strings.Repeat("T", length)
			_, err := t2.t.FileSendChunk(friend_number, file_number, position, []byte(s))
			if err != nil {
				t.Error(err)
			}

		}, nil)
		sendRecvDone := false
		t2.t.CallbackFileRecvControl(func(_ *Tox, friendNumber uint32, fileNumber uint32,
			control int, ud interface{}) {
			// log.Println(fileNumber, control)
			if control == FILE_CONTROL_CANCEL {
				sendRecvDone = true
			}
		}, nil)

		go t1.Iterate()
		go t2.Iterate()
		defer t1.stop()
		defer t2.stop()

		waitcond(func() bool {
			return t1.t.SelfGetConnectionStatus() != CONNECTION_NONE && t2.t.SelfGetConnectionStatus() != CONNECTION_NONE
		}, 100)

		t2.t.FriendAdd(t1.t.SelfGetAddress(), "autotests")
		waitcond(func() bool {
			return t1.t.SelfGetFriendListSize() == 1
		}, 100)

		fn, _ := t2.t.FriendByPublicKey(t1.t.SelfGetPublicKey())
		// must wait friend online and can call InviteFriend
		waitcond(func() bool {
			st, _ := t2.t.FriendGetConnectionStatus(fn)
			return st > CONNECTION_NONE
		}, 100)

		fh, err := t2.t.FileSend(fn, FILE_KIND_DATA, 12345, "123456", "testfile.txt")
		if err != nil {
			t.Error(err, fh)
		}
		fid, err := t2.t.FileGetFileId(fn, fh)
		if len(fid) != FILE_ID_LENGTH*2 {
			t.Error("file id length not match:", len(fid), FILE_ID_LENGTH*2)
		}

		waitcond(func() bool {
			return len(recvData) > 0 && sendRecvDone
		}, 10)
		if len(recvData) != 12345-15 {
			t.Error("recv size not match")
		}

		// select {}
	})
}

func TestAV(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		if tv1, err := NewToxAV(nil); tv1 != nil {
			t.Error("must nil", err)
		}
		t1 := NewMiniTox()
		tv1, err := NewToxAV(t1.t)
		if err != nil {
			t.Error(err, tv1)
		}
	})
}

// go test -v -run Covers
func TestCovers(t *testing.T) {
	t1 := NewMiniTox()
	defer t1.t.Kill()

	tv := reflect.ValueOf(t1.t)
	mnum := tv.NumMethod()
	if false {
		t.Log(mnum)
	}

	mths := make(map[string]bool)
	for i := 0; i < mnum; i++ {
		mth := tv.Type().Method(i)
		// t.Log(i, mth.Name)
		mths[mth.Name] = true
	}

	//
	_, file, _, _ := runtime.Caller(0)
	t.Log(file)

	fset := token.NewFileSet() // positions are relative to fset

	// Parse the file containing this very example
	// but stop after processing the imports.
	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		t.Log(err)
		return
	}

	t.Log("walking ast...")
	v := &callVisitor{t: t}
	v.fns = make(map[string]bool)
	ast.Walk(v, f)
	// t.Log(v.fns)

	notins := make(map[string]bool)
	for mn, _ := range mths {
		if _, ok := v.fns[mn]; !ok {
			t.Log("not tested:", mn)
			notins[mn] = false
		}
	}

	t.Log("test covers:", mnum-len(notins), mnum)
}

type callVisitor struct {
	t   *testing.T
	fns map[string]bool
}

func (v *callVisitor) Visit(node ast.Node) (w ast.Visitor) {
	t := v.t
	if false {
		nt := reflect.TypeOf(node)
		switch nt.Kind() {
		case reflect.Ptr:
			t.Log(nt.Elem().Kind(), nt.Elem().Name())
		default:
			t.Log(nt.Kind())
		}
	}

	switch ty := node.(type) {
	case *ast.File:
		for _, d := range ty.Decls {
			v.Visit(d)
		}
	case *ast.FuncDecl:
		v.Visit(ty.Body)
	case *ast.GenDecl:
		for _, d := range ty.Specs {
			v.Visit(d)
		}
	case *ast.BlockStmt:
		for _, s := range ty.List {
			v.Visit(s)
		}
	case *ast.ExprStmt:
		v.Visit(ty.X)
	case *ast.CallExpr:
		// t.Logf("%+v\n", ty)
		v.Visit(ty.Fun)
		for _, a := range ty.Args {
			v.Visit(a)
		}
	case *ast.FuncLit:
		v.Visit(ty.Body)
	case *ast.IfStmt:
		v.Visit(ty.Body)
		v.Visit(ty.Cond)
		if ty.Init != nil {
			v.Visit(ty.Init)
		}
		if ty.Else != nil {
			v.Visit(ty.Else)
		}
	case *ast.AssignStmt:
		for _, s := range ty.Rhs {
			v.Visit(s)
		}
	case *ast.ForStmt:
		if ty.Cond != nil {
			v.Visit(ty.Cond)
		}
		v.Visit(ty.Body)
		if ty.Init != nil {
			v.Visit(ty.Init)
		}
		if ty.Post != nil {
			v.Visit(ty.Post)
		}
	case *ast.ReturnStmt:
		for _, s := range ty.Results {
			v.Visit(s)
		}
	case *ast.SwitchStmt:
		if ty.Init != nil {
			v.Visit(ty.Init)
		}
		v.Visit(ty.Body)
	case *ast.GoStmt:
		v.Visit(ty.Call)
	case *ast.SelectStmt:
		v.Visit(ty.Body)
	case *ast.SelectorExpr:
		if ty.Sel.IsExported() {
			// t.Log(ty.Sel.String(), ty.Sel.Name, ty.X)
			v.fns[ty.Sel.Name] = true
		}
		v.Visit(ty.X)
	case *ast.BinaryExpr:
		v.Visit(ty.X)
		v.Visit(ty.Y)
	case *ast.UnaryExpr:
		v.Visit(ty.X)
	case *ast.ValueSpec:
		for _, val := range ty.Values {
			v.Visit(val)
		}
	case *ast.CaseClause:
		for _, b := range ty.Body {
			v.Visit(b)
		}
		for _, l := range ty.List {
			v.Visit(l)
		}
	default:
		if false {
			t.Logf("%+v, %+v ===\n", ty, node)
		}
	}
	// t.Log(node)
	return nil
}
