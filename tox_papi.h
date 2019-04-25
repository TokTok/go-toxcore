#ifndef _TOX_PAPI_H_
#define _TOX_PAPI_H_

#include <stdint.h>
#include <tox/tox.h>

extern void tox_self_get_address_p(uintptr_t params);
extern void tox_friend_get_name_p(uintptr_t params);
extern void tox_friend_get_public_key_p(uintptr_t params);
extern void tox_friend_send_message_p(uintptr_t params);

#endif


