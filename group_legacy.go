package tox

/*
#include <stdlib.h>
#include <string.h>
#include <tox/tox.h>


*/
import "C"

// legacy group callback type

type cb_group_invite_ftype func(this *Tox, friendNumber uint32, itype uint8, cookie string, userData interface{})
type cb_group_message_ftype func(this *Tox, groupNumber int, peerNumber int, message string, userData interface{})

type cb_group_action_ftype func(this *Tox, groupNumber int, peerNumber int, action string, userData interface{})
type cb_group_title_ftype func(this *Tox, groupNumber int, peerNumber int, title string, userData interface{})
type cb_group_namelist_change_ftype func(this *Tox, groupNumber int, peerNumber int, change uint8, userData interface{})

// tox_callback_group_***

func (this *Tox) CallbackGroupInvite(cbfn cb_group_invite_ftype, userData interface{}) {
	this.CallbackGroupInviteAdd(cbfn, userData)
}
func (this *Tox) CallbackGroupInviteAdd(cbfn cb_group_invite_ftype, userData interface{}) {
	cbfn_ := func(this *Tox, friendNumber uint32, itype uint8, cookie string, userData interface{}) {
		cbfn(this, friendNumber, itype, cookie, userData)
	}
	this.CallbackConferenceInviteAdd(cbfn_, userData)
}

func (this *Tox) CallbackGroupMessage(cbfn cb_group_message_ftype, userData interface{}) {
	this.CallbackGroupMessageAdd(cbfn, userData)
}
func (this *Tox) CallbackGroupMessageAdd(cbfn cb_group_message_ftype, userData interface{}) {
	cbfn_ := func(this *Tox, groupNumber uint32, peerNumber uint32, message string, userData interface{}) {
		cbfn(this, int(groupNumber), int(peerNumber), message, userData)
	}
	this.CallbackConferenceMessageAdd(cbfn_, userData)
}

func (this *Tox) CallbackGroupAction(cbfn cb_group_action_ftype, userData interface{}) {
	this.CallbackGroupActionAdd(cbfn, userData)
}
func (this *Tox) CallbackGroupActionAdd(cbfn cb_group_action_ftype, userData interface{}) {
	cbfn_ := func(this *Tox, groupNumber uint32, peerNumber uint32, message string, userData interface{}) {
		cbfn(this, int(groupNumber), int(peerNumber), message, userData)
	}
	this.CallbackConferenceActionAdd(cbfn_, userData)
}

func (this *Tox) CallbackGroupTitle(cbfn cb_group_title_ftype, userData interface{}) {
	this.CallbackGroupTitleAdd(cbfn, userData)
}
func (this *Tox) CallbackGroupTitleAdd(cbfn cb_group_title_ftype, userData interface{}) {
	cbfn_ := func(this *Tox, groupNumber uint32, peerNumber uint32, title string, userData interface{}) {
		cbfn(this, int(groupNumber), int(peerNumber), title, userData)
	}
	this.CallbackConferenceTitleAdd(cbfn_, userData)
}

// methods
func (this *Tox) AddGroupChat() (int, error) {
	gn, err := this.ConferenceNew()
	return int(gn), err
}

func (this *Tox) DelGroupChat(groupNumber int) (int, error) {
	return this.ConferenceDelete(uint32(groupNumber))
}

func (this *Tox) GroupPeerName(groupNumber int, peerNumber int) (string, error) {
	return this.ConferencePeerGetName(uint32(groupNumber), uint32(peerNumber))
}

func (this *Tox) GroupPeerPubkey(groupNumber int, peerNumber int) (string, error) {
	return this.ConferencePeerGetPublicKey(uint32(groupNumber), uint32(peerNumber))
}

func (this *Tox) InviteFriend(friendNumber uint32, groupNumber int) (int, error) {
	return this.ConferenceInvite(friendNumber, uint32(groupNumber))
}

func (this *Tox) JoinGroupChat(friendNumber uint32, cookie string) (int, error) {
	groupNumber, err := this.ConferenceJoin(friendNumber, cookie)
	return int(groupNumber), err
}

func (this *Tox) GroupActionSend(groupNumber int, action string) (int, error) {
	return this.ConferenceSendMessage(uint32(groupNumber), MESSAGE_TYPE_ACTION, action)
}

func (this *Tox) GroupMessageSend(groupNumber int, message string) (int, error) {
	return this.ConferenceSendMessage(uint32(groupNumber), MESSAGE_TYPE_NORMAL, message)
}

func (this *Tox) GroupSetTitle(groupNumber int, title string) (int, error) {
	return this.ConferenceSetTitle(uint32(groupNumber), title)
}

func (this *Tox) GroupGetTitle(groupNumber int) (string, error) {
	return this.ConferenceGetTitle(uint32(groupNumber))
}

func (this *Tox) GroupPeerNumberIsOurs(groupNumber int, peerNumber uint32) bool {
	return this.ConferencePeerNumberIsOurs(uint32(groupNumber), peerNumber)
}

func (this *Tox) GroupNumberPeers(groupNumber int) int {
	cnt := this.ConferencePeerCount(uint32(groupNumber))
	return int(cnt)
}

// extra combined api
func (this *Tox) GroupGetNames(groupNumber int) []string {
	return this.ConferenceGetNames(uint32(groupNumber))
}

func (this *Tox) GroupGetPeerPubkeys(groupNumber int) []string {
	return this.ConferenceGetPeerPubkeys(uint32(groupNumber))
}

func (this *Tox) GroupGetPeers(groupNumber int) map[int]string {
	peers := this.ConferenceGetPeers(uint32(groupNumber))
	rets := make(map[int]string, 0)
	for pn, name := range peers {
		rets[int(pn)] = name
	}
	return rets
}

func (this *Tox) CountChatList() uint32 {
	return this.ConferenceGetChatlistSize()
}

func (this *Tox) GetChatList() []int32 {
	gns := this.ConferenceGetChatlist()
	var ret []int32
	for _, gn := range gns {
		ret = append(ret, int32(gn))
	}
	return ret
}

func (this *Tox) GroupGetType(groupNumber uint32) (int, error) {
	return this.ConferenceGetType(groupNumber)
}
