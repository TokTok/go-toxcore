package tox

/*
#include <stdlib.h>
#include <string.h>
#include <vpx/vpx_image.h>
#include <tox/tox.h>
#include <tox/toxav.h>

void callbackCallWrapperForC(ToxAV *toxAV, uint32_t friend_number, bool audio_enabled,
                           bool video_enabled, void *user_data);
typedef void (*cb_call_ftype)(ToxAV *toxAV, uint32_t friend_number, bool audio_enabled,
                           bool video_enabled, void *user_data);
static void cb_call_wrapper_for_go(ToxAV *m, cb_call_ftype fn, void *userdata)
{ toxav_callback_call(m, fn, userdata); }

void callbackCallStateWrapperForC(ToxAV *toxAV, uint32_t friendNumber, uint32_t state, void* user_data);
typedef void (*cb_call_state_ftype)(ToxAV *toxAV, uint32_t friendNumber, uint32_t state, void *user_data);
static void cb_call_state_wrapper_for_go(ToxAV *m, cb_call_state_ftype fn, void *userdata)
{ toxav_callback_call_state(m, fn, userdata); }

void callbackAudioBitRateWrapperForC(ToxAV *toxAV, uint32_t friendNumber, uint32_t audioBitRate, void* user_data);
typedef void (*cb_audio_bit_rate_ftype)(ToxAV *toxAV, uint32_t friendNumber, uint32_t audioBitRate, void *user_data);
static void cb_audio_bit_rate_wrapper_for_go(ToxAV *m, cb_audio_bit_rate_ftype fn, void *userdata)
{ toxav_callback_audio_bit_rate(m, fn, userdata); }

void callbackVideoBitRateWrapperForC(ToxAV *toxAV, uint32_t friendNumber, uint32_t videoBitRate, void* user_data);
typedef void (*cb_video_bit_rate_ftype)(ToxAV *toxAV, uint32_t friendNumber, uint32_t videoBitRate, void *user_data);
static void cb_video_bit_rate_wrapper_for_go(ToxAV *m, cb_video_bit_rate_ftype fn, void *userdata)
{ toxav_callback_video_bit_rate(m, fn, userdata); }

void callbackAudioReceiveFrameWrapperForC(ToxAV *toxAV, uint32_t friendNumber, int16_t *pcm, size_t sample_count, uint8_t channels, uint32_t sampling_rate, void* user_data);
typedef void (*cb_audio_receive_frame_ftype)(ToxAV *toxAV, uint32_t friendNumber, const int16_t *pcm, size_t sample_count, uint8_t channels, uint32_t sampling_rate, void *user_data);
static void cb_audio_receive_frame_wrapper_for_go(ToxAV *m, cb_audio_receive_frame_ftype fn, void *userdata)
{ toxav_callback_audio_receive_frame(m, fn, userdata); }

void callbackVideoReceiveFrameWrapperForC(ToxAV *toxAV, uint32_t friendNumber, uint16_t width,
        uint16_t height, uint8_t *y, uint8_t *u, uint8_t *v,
        int32_t ystride, int32_t ustride, int32_t vstride, void* user_data);
typedef void (*cb_video_receive_frame_ftype)(ToxAV *toxAV, uint32_t friendNumber, uint16_t width,
        uint16_t height, const uint8_t *y, const uint8_t *u, const uint8_t *v,
        int32_t ystride, int32_t ustride, int32_t vstride, void *user_data);
static void cb_video_receive_frame_wrapper_for_go(ToxAV *m, cb_video_receive_frame_ftype fn, void *userdata)
{ toxav_callback_video_receive_frame(m, fn, userdata); }

extern void i420_to_rgb(int width, int height, const uint8_t *y, const uint8_t *u, const uint8_t *v,
            int ystride, int ustride, int vstride, unsigned char *out);
extern void rgb_to_i420(unsigned char* rgb, vpx_image_t *img);


// fix nouse compile warning
static inline __attribute__((__unused__)) void fixnousetoxav() {
    cb_call_wrapper_for_go(NULL, NULL, NULL);
    cb_call_state_wrapper_for_go(NULL, NULL, NULL);
    cb_audio_bit_rate_wrapper_for_go(NULL, NULL, NULL);
    cb_video_bit_rate_wrapper_for_go(NULL, NULL, NULL);
    cb_audio_receive_frame_wrapper_for_go(NULL, NULL, NULL);
    cb_video_receive_frame_wrapper_for_go(NULL, NULL, NULL);
}

*/
import "C"
import (
	"encoding/hex"
	"errors"
	"unsafe"
)

type cb_call_ftype func(this *ToxAV, friendNumber uint32, audioEnabled bool, videoEnabled bool, userData interface{})
type cb_call_state_ftype func(this *ToxAV, friendNumber uint32, state uint32, userData interface{})
type cb_audio_bit_rate_ftype func(this *ToxAV, friendNumber uint32, audioBitRate uint32, userData interface{})
type cb_video_bit_rate_ftype func(this *ToxAV, friendNumber uint32, videoBitRate uint32, userData interface{})
type cb_audio_receive_frame_ftype func(this *ToxAV, friendNumber uint32, pcm []byte, sampleCount int, channels int, samplingRate int, userData interface{})
type cb_video_receive_frame_ftype func(this *ToxAV, friendNumber uint32, width uint16, height uint16, data []byte, userData interface{})

type ToxAV struct {
	tox   *Tox
	toxav *C.ToxAV

	// session datas
	out_image  []byte
	out_width  C.uint16_t
	out_hegith C.uint16_t
	in_image   *C.vpx_image_t
	in_width   C.uint16_t
	in_height  C.uint16_t

	// callbacks
	cb_call                          cb_call_ftype
	cb_call_user_data                interface{}
	cb_call_state                    cb_call_state_ftype
	cb_call_state_user_data          interface{}
	cb_audio_bit_rate                cb_audio_bit_rate_ftype
	cb_audio_bit_rate_user_data      interface{}
	cb_video_bit_rate                cb_video_bit_rate_ftype
	cb_video_bit_rate_user_data      interface{}
	cb_audio_receive_frame           cb_audio_receive_frame_ftype
	cb_audio_receive_frame_user_data interface{}
	cb_video_receive_frame           cb_video_receive_frame_ftype
	cb_video_receive_frame_user_data interface{}
}

func NewToxAV(tox *Tox) (*ToxAV, error) {
	if tox == nil {
		return nil, toxerr("tox can not nil")
	}

	tav := new(ToxAV)
	tav.tox = tox

	var cerr C.TOXAV_ERR_NEW
	tav.toxav = C.toxav_new(tox.toxcore, &cerr)
	if cerr != 0 {
		return nil, toxerr(cerr)
	}

	cbAVUserDatas.set(tav.toxav, tav)
	return tav, nil
}

func (this *ToxAV) Kill() {
	C.toxav_kill(this.toxav)
}

func (this *ToxAV) GetTox() *Tox {
	return this.tox
}

func (this *ToxAV) IterationInterval() int {
	return int(C.toxav_iteration_interval(this.toxav))
}

func (this *ToxAV) Iterate() {
	C.toxav_iterate(this.toxav)
}

func (this *ToxAV) Call(friendNumber uint32, audioBitRate uint32, videoBitRate uint32) (bool, error) {
	var cerr C.TOXAV_ERR_CALL
	r := C.toxav_call(this.toxav, C.uint32_t(friendNumber), C.uint32_t(audioBitRate), C.uint32_t(videoBitRate), &cerr)
	if cerr != 0 {
		return bool(r), toxerr(cerr)
	}
	return bool(r), nil
}

var cbAVUserDatas = newUserDataAV()

//export callbackCallWrapperForC
func callbackCallWrapperForC(m *C.ToxAV, friendNumber C.uint32_t, audioEnabled C.bool, videoEnabled C.bool, a3 unsafe.Pointer) {
	var this = cbAVUserDatas.get(m)
	if this.cb_call != nil {
		this.cb_call(this, uint32(friendNumber), bool(audioEnabled), bool(videoEnabled), this.cb_call_user_data)
	}
}

func (this *ToxAV) CallbackCall(cbfn cb_call_ftype, userData interface{}) {
	this.cb_call = cbfn
	this.cb_call_user_data = userData

	var _cbfn = (C.cb_call_ftype)(C.callbackCallWrapperForC)
	var _userData = unsafe.Pointer(this)
	_userData = nil

	C.cb_call_wrapper_for_go(this.toxav, _cbfn, _userData)
}

func (this *ToxAV) Answer(friendNumber uint32, audioBitRate uint32, videoBitRate uint32) (bool, error) {
	var cerr C.TOXAV_ERR_ANSWER
	r := C.toxav_answer(this.toxav, C.uint32_t(friendNumber), C.uint32_t(audioBitRate), C.uint32_t(videoBitRate), &cerr)
	if cerr != C.TOXAV_ERR_ANSWER_OK {
		return false, toxerr(cerr)
	}

	return bool(r), nil
}

//export callbackCallStateWrapperForC
func callbackCallStateWrapperForC(m *C.ToxAV, friendNumber C.uint32_t, state C.uint32_t, a3 unsafe.Pointer) {
	var this = cbAVUserDatas.get(m)
	if this.cb_call_state != nil {
		this.cb_call_state(this, uint32(friendNumber), uint32(state), this.cb_call_state_user_data)
	}
}

func (this *ToxAV) CallbackCallState(cbfn cb_call_state_ftype, userData interface{}) {
	this.cb_call_state = cbfn
	this.cb_call_state_user_data = userData

	var _cbfn = (C.cb_call_state_ftype)(C.callbackCallStateWrapperForC)
	var _userData = unsafe.Pointer(this)
	_userData = nil

	C.cb_call_state_wrapper_for_go(this.toxav, _cbfn, _userData)
}

func (this *ToxAV) CallControl(friendNumber uint32, control int) (bool, error) {
	var cerr C.TOXAV_ERR_CALL_CONTROL
	r := C.toxav_call_control(this.toxav, C.uint32_t(friendNumber), C.TOXAV_CALL_CONTROL(control), &cerr)
	if cerr != C.TOXAV_ERR_CALL_CONTROL_OK {
		return bool(r), toxerr(cerr)
	}
	return bool(r), nil
}

func (this *ToxAV) AudioSetBitRate(friendNumber uint32, audioBitRate uint32) (bool, error) {
	var cerr C.TOXAV_ERR_BIT_RATE_SET
	r := C.toxav_audio_set_bit_rate(this.toxav, C.uint32_t(friendNumber), C.uint32_t(audioBitRate), &cerr)
	if cerr != C.TOXAV_ERR_BIT_RATE_SET_OK {
		return bool(r), toxerr(cerr)
	}
	return bool(r), nil
}

func (this *ToxAV) VideoSetBitRate(friendNumber uint32, videoBitRate uint32) (bool, error) {
	var cerr C.TOXAV_ERR_BIT_RATE_SET
	r := C.toxav_video_set_bit_rate(this.toxav, C.uint32_t(friendNumber), C.uint32_t(videoBitRate), &cerr)
	if cerr != C.TOXAV_ERR_BIT_RATE_SET_OK {
		return bool(r), toxerr(cerr)
	}
	return bool(r), nil
}

//export callbackAudioBitRateWrapperForC
func callbackAudioBitRateWrapperForC(m *C.ToxAV, friendNumber C.uint32_t, audioBitRate C.uint32_t, a3 unsafe.Pointer) {
	var this = cbAVUserDatas.get(m)
	if this.cb_audio_bit_rate != nil {
		this.cb_audio_bit_rate(this, uint32(friendNumber), uint32(audioBitRate), this.cb_call_state_user_data)
	}
}

func (this *ToxAV) CallbackAudioBitRate(cbfn cb_audio_bit_rate_ftype, userData interface{}) {
	this.cb_audio_bit_rate = cbfn
	this.cb_audio_bit_rate_user_data = userData

	var _cbfn = (C.cb_audio_bit_rate_ftype)(C.callbackAudioBitRateWrapperForC)
	var _userData = unsafe.Pointer(this)
	_userData = nil

	C.cb_audio_bit_rate_wrapper_for_go(this.toxav, _cbfn, _userData)
}

//export callbackVideoBitRateWrapperForC
func callbackVideoBitRateWrapperForC(m *C.ToxAV, friendNumber C.uint32_t, videoBitRate C.uint32_t, a3 unsafe.Pointer) {
	var this = cbAVUserDatas.get(m)
	if this.cb_video_bit_rate != nil {
		this.cb_video_bit_rate(this, uint32(friendNumber), uint32(videoBitRate), this.cb_call_state_user_data)
	}
}

func (this *ToxAV) CallbackVideoBitRate(cbfn cb_video_bit_rate_ftype, userData interface{}) {
	this.cb_video_bit_rate = cbfn
	this.cb_video_bit_rate_user_data = userData

	var _cbfn = (C.cb_video_bit_rate_ftype)(C.callbackVideoBitRateWrapperForC)
	var _userData = unsafe.Pointer(this)
	_userData = nil

	C.cb_video_bit_rate_wrapper_for_go(this.toxav, _cbfn, _userData)
}

func (this *ToxAV) AudioSendFrame(friendNumber uint32, pcm []byte, sampleCount int, channels int, samplingRate int) (bool, error) {
	pcm_ := (*C.int16_t)(unsafe.Pointer(&pcm[0]))
	var cerr C.TOXAV_ERR_SEND_FRAME
	r := C.toxav_audio_send_frame(this.toxav, C.uint32_t(friendNumber), pcm_, C.size_t(sampleCount), C.uint8_t(channels), C.uint32_t(samplingRate), &cerr)
	if cerr != C.TOXAV_ERR_SEND_FRAME_OK {
		return false, toxerr(cerr)
	}
	return bool(r), nil
}

func (this *ToxAV) VideoSendFrame(friendNumber uint32, width uint16, height uint16, data []byte) (bool, error) {
	if this.in_image != nil && (uint16(this.in_width) != width || uint16(this.in_height) != height) {
		C.vpx_img_free(this.in_image)
		this.in_image = nil
	}

	if this.in_image == nil {
		this.in_width = C.uint16_t(width)
		this.in_height = C.uint16_t(height)
		this.in_image = C.vpx_img_alloc(nil, C.VPX_IMG_FMT_I420, C.uint(this.in_width), C.uint(this.in_height), 1)
	}

	C.rgb_to_i420((*C.uchar)(unsafe.Pointer(&data[0])), this.in_image)

	var cerr C.TOXAV_ERR_SEND_FRAME
	r := C.toxav_video_send_frame(this.toxav, C.uint32_t(friendNumber), C.uint16_t(width), C.uint16_t(height),
		(*C.uint8_t)(this.in_image.planes[0]),
		(*C.uint8_t)(this.in_image.planes[1]),
		(*C.uint8_t)(this.in_image.planes[2]),
		&cerr)
	if cerr != C.TOXAV_ERR_SEND_FRAME_OK {
		return false, toxerr(cerr)
	}
	return bool(r), nil
}

//export callbackAudioReceiveFrameWrapperForC
func callbackAudioReceiveFrameWrapperForC(m *C.ToxAV, friendNumber C.uint32_t, pcm *C.int16_t, sampleCount C.size_t, channels C.uint8_t, samplingRate C.uint32_t, a3 unsafe.Pointer) {
	var this = cbAVUserDatas.get(m)
	if this.cb_audio_receive_frame != nil {
		length := sampleCount * C.size_t(channels) * 2
		pcm_p := unsafe.Pointer(pcm)
		pcm_b := C.GoBytes(pcm_p, C.int(length))
		this.cb_audio_receive_frame(this, uint32(friendNumber), pcm_b, int(sampleCount), int(channels), int(samplingRate), this.cb_audio_receive_frame_user_data)
	}
}

func (this *ToxAV) CallbackAudioReceiveFrame(cbfn cb_audio_receive_frame_ftype, userData interface{}) {
	this.cb_audio_receive_frame = cbfn
	this.cb_audio_receive_frame_user_data = userData

	var _cbfn = (C.cb_audio_receive_frame_ftype)(C.callbackAudioReceiveFrameWrapperForC)
	var _userData = unsafe.Pointer(this)
	_userData = nil

	C.cb_audio_receive_frame_wrapper_for_go(this.toxav, _cbfn, _userData)
}

//export callbackVideoReceiveFrameWrapperForC
func callbackVideoReceiveFrameWrapperForC(m *C.ToxAV, friendNumber C.uint32_t, width C.uint16_t, height C.uint16_t, y *C.uint8_t, u *C.uint8_t, v *C.uint8_t, ystride C.int32_t, ustride C.int32_t, vstride C.int32_t, a3 unsafe.Pointer) {
	var this = cbAVUserDatas.get(m)
	if this.cb_video_receive_frame != nil {

		if this.out_image != nil && (this.out_width != width || this.out_hegith != height) {
			this.out_image = nil
		}

		var buf_size int = int(width) * int(height) * 3

		if this.out_image == nil {
			this.out_width = width
			this.out_hegith = height
			this.out_image = make([]byte, buf_size, buf_size)
		}

		out := unsafe.Pointer(&(this.out_image[0]))
		C.i420_to_rgb(C.int(width), C.int(height), y, u, v, C.int(ystride), C.int(ustride), C.int(vstride), (*C.uchar)(out))

		this.cb_video_receive_frame(this, uint32(friendNumber), uint16(width), uint16(height), this.out_image, this.cb_video_receive_frame_user_data)

	}
}

func (this *ToxAV) CallbackVideoReceiveFrame(cbfn cb_video_receive_frame_ftype, userData interface{}) {
	this.cb_video_receive_frame = cbfn
	this.cb_video_receive_frame_user_data = userData

	var _cbfn = (C.cb_video_receive_frame_ftype)(C.callbackVideoReceiveFrameWrapperForC)
	var _userData = unsafe.Pointer(this)
	_userData = nil

	C.cb_video_receive_frame_wrapper_for_go(this.toxav, _cbfn, _userData)
}

// TODO
// toxav_add_av_groupchat
// toxav_join_av_groupchat
// toxav_group_send_audio

func (this *Tox) AddAVGroupChat() int {
	r := C.toxav_add_av_groupchat(this.toxcore, nil, nil)
	return int(r)
}

func (this *Tox) JoinAVGroupChat(friendNumber uint32, cookie string) (int, error) {
	data, err := hex.DecodeString(cookie)
	if err != nil {
		return 0, errors.New("Invalid cookie:" + cookie)
	}
	var _fn = C.uint32_t(friendNumber)
	var _data = (*C.uint8_t)((unsafe.Pointer)(&data[0]))
	var length = len(data)
	var _length = C.uint16_t(length)

	// TODO nil => real
	r := C.toxav_join_av_groupchat(this.toxcore, _fn, _data, _length, nil, nil)
	if int(r) == -1 {
		return int(r), errors.New("Join av group chat failed")
	}

	if this.hooks.ConferenceJoin != nil {
		this.hooks.ConferenceJoin(friendNumber, uint32(r), cookie)
	}
	return int(r), nil
}
