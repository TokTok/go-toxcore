

#include <stdlib.h>
#include <stdio.h>
#include <assert.h>
#include <tox/tox.h>


#define PAPI_PARAM_HEADER uintptr_t gosz; Tox* tox;
#define PAPI_GETCHK_PARAM(stname) stname* prm=(stname*)(void*)params; \
    if (prm->gosz != sizeof(stname)) {                                  \
        fprintf(stderr, "gostsz=%ld, cstsz=%ld\n", prm->gosz, sizeof(stname)); \
        assert(prm->gosz == sizeof(stname)); }

typedef struct {
    PAPI_PARAM_HEADER;
    uint8_t* addr;
} self_get_address_param;

void tox_self_get_address_p(uintptr_t params) {
    PAPI_GETCHK_PARAM(self_get_address_param);
    tox_self_get_address(prm->tox, prm->addr);
}

typedef  struct {
    PAPI_PARAM_HEADER;
    bool retval;
    uint32_t friendNumber;
    uint8_t* name     ;
    long error        ;
    long namelen      ;
}friend_get_name_param;

void tox_friend_get_name_p(uintptr_t params) {
    PAPI_GETCHK_PARAM(friend_get_name_param);
    prm->namelen = tox_friend_get_name_size(prm->tox, prm->friendNumber, (TOX_ERR_FRIEND_QUERY*)&prm->error);
    prm->namelen = 100;
    prm->retval = tox_friend_get_name(prm->tox, prm->friendNumber, prm->name, (TOX_ERR_FRIEND_QUERY*)&prm->error);
}

typedef  struct {
    PAPI_PARAM_HEADER;
    bool retval;
    uint32_t friendNumber;
    uint8_t* pubkey         ;
    long error        ;
}friend_get_public_key_param;

void tox_friend_get_public_key_p(uintptr_t params) {
    PAPI_GETCHK_PARAM(friend_get_public_key_param);
    prm->retval = tox_friend_get_public_key(prm->tox, prm->friendNumber, prm->pubkey,
                               (TOX_ERR_FRIEND_GET_PUBLIC_KEY*)&prm->error);
}

typedef struct {
    PAPI_PARAM_HEADER;
    uint32_t retval;
    uint32_t friend_number;
    TOX_MESSAGE_TYPE type;
    uint8_t* message;
    size_t length;
    TOX_ERR_FRIEND_SEND_MESSAGE error;
} send_message_param;

void tox_friend_send_message_p(uintptr_t params) {
    PAPI_GETCHK_PARAM(send_message_param);
    prm->retval = tox_friend_send_message(prm->tox, prm->friend_number, prm->type,
                                           prm->message, prm->length, &prm->error);
}
