#include <time.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <tox/tox.h>

typedef union {
    uint32_t uint32;
    uint16_t uint16[2];
    uint8_t uint8[4];
}
    IP4;

extern const IP4 IP4_LOOPBACK;
extern const IP4 IP4_BROADCAST;

typedef union {
    uint8_t uint8[16];
    uint16_t uint16[8];
    uint32_t uint32[4];
    uint64_t uint64[2];
}
    IP6;

extern const IP6 IP6_LOOPBACK;
extern const IP6 IP6_BROADCAST;

typedef struct {
    uint8_t family;
    /*GNU_EXTENSION*/ union {
        IP4 ip4;
        IP6 ip6;
    };
}
    IP;

typedef struct {
    IP ip;
    uint16_t port;
}
    IP_Port;

typedef struct {
    IP_Port ip_port;
    uint8_t proxy_type; // a value from TCP_PROXY_TYPE
} TCP_Proxy_Info;

typedef struct {
    uint8_t ipv6enabled;
    uint8_t udp_disabled;
    TCP_Proxy_Info proxy_info;
    uint16_t port_range[2];
    uint16_t tcp_server_port;

    uint8_t hole_punching_enabled;
    bool local_discovery_enabled;

    /*logger_cb*/void *log_callback;
    void *log_user_data;
} Messenger_Options_Fake;


#define MAX_RECEIVED_STORED 32
#define CRYPTO_PUBLIC_KEY_SIZE         32
#define CRYPTO_SHARED_KEY_SIZE         32
#define CRYPTO_SYMMETRIC_KEY_SIZE      CRYPTO_SHARED_KEY_SIZE

typedef struct {
    uint8_t     public_key[CRYPTO_PUBLIC_KEY_SIZE];
    IP_Port     ip_port;
} Node_format_Fake;

#define DESIRED_CLOSE_CONNECTIONS 4
#define MAX_GROUP_CONNECTIONS 16
#define GROUP_IDENTIFIER_LENGTH (1 + CRYPTO_SYMMETRIC_KEY_SIZE) /* type + CRYPTO_SYMMETRIC_KEY_SIZE so we can use new_symmetric_key(...) to fill it */

typedef struct {
    uint8_t status;

    /*Group_Peer*/void *group;
    uint32_t numpeers;
#define MAX_GROUP_CONNECTIONS 16
    struct {
        uint8_t type; /* GROUPCHAT_CLOSE_* */
        uint8_t closest;
        uint32_t number;
        uint16_t group_number;
    } close[MAX_GROUP_CONNECTIONS];

    uint8_t real_pk[CRYPTO_PUBLIC_KEY_SIZE];
    struct {
        uint8_t entry;
        uint8_t real_pk[CRYPTO_PUBLIC_KEY_SIZE];
        uint8_t temp_pk[CRYPTO_PUBLIC_KEY_SIZE];
    } closest_peers[DESIRED_CLOSE_CONNECTIONS];
    uint8_t changed;

    uint8_t identifier[GROUP_IDENTIFIER_LENGTH];

    uint8_t title[TOX_MAX_NAME_LENGTH];
    uint8_t title_len;

    uint32_t message_number;
    uint16_t lossy_message_number;
    uint16_t peer_number;

    uint64_t last_sent_ping;

    int number_joined; /* friendcon_id of person that invited us to the chat. (-1 means none) */

    void *object;

    void (*peer_on_join)(void *, uint32_t, uint32_t);
    void (*peer_on_leave)(void *, uint32_t, uint32_t, void *);
    void (*group_on_delete)(void *, uint32_t);
} Group_c_Fake;

typedef enum {
    USERSTATUS_NONE,
    USERSTATUS_AWAY,
    USERSTATUS_BUSY,
    USERSTATUS_INVALID
}
    USERSTATUS;


#define NUM_SAVED_TCP_RELAYS 8
#define MAX_STATUSMESSAGE_LENGTH 1007

struct Messenger {
    /*Logger*/ void *log;

    /*Networking_Core*/void *net;
    /*Net_Crypto*/void *net_crypto;
    /*DHT*/void *dht;

    /*Onion*/void *onion;
    /*Onion_Announce*/void *onion_a;
    /*Onion_Client*/void *onion_c;

    /*Friend_Connections*/void *fr_c;

    /*TCP_Server*/void *tcp_server;
    /*Friend_Requests **/void *fr;
    uint8_t name[TOX_MAX_NAME_LENGTH];
    uint16_t name_length;

    uint8_t statusmessage[MAX_STATUSMESSAGE_LENGTH];
    uint16_t statusmessage_length;

    USERSTATUS userstatus;

    /*Friend*/void *friendlist;
    uint32_t numfriends;

    time_t lastdump;

    uint8_t has_added_relays; // If the first connection has occurred in do_messenger
    Node_format_Fake loaded_relays[NUM_SAVED_TCP_RELAYS]; // Relays loaded from config

    void (*friend_message)(struct Messenger *m, uint32_t, unsigned int, const uint8_t *, size_t, void *);
    void (*friend_namechange)(struct Messenger *m, uint32_t, const uint8_t *, size_t, void *);
    void (*friend_statusmessagechange)(struct Messenger *m, uint32_t, const uint8_t *, size_t, void *);
    void (*friend_userstatuschange)(struct Messenger *m, uint32_t, unsigned int, void *);
    void (*friend_typingchange)(struct Messenger *m, uint32_t, bool, void *);
    void (*read_receipt)(struct Messenger *m, uint32_t, uint32_t, void *);
    void (*friend_connectionstatuschange)(struct Messenger *m, uint32_t, unsigned int, void *);
    void (*friend_connectionstatuschange_internal)(struct Messenger *m, uint32_t, uint8_t, void *);
    void *friend_connectionstatuschange_internal_userdata;

    void *conferences_object; /* Set by new_groupchats()*/
    void (*conference_invite)(struct Messenger *m, uint32_t, const uint8_t *, uint16_t, void *);

    void (*file_sendrequest)(struct Messenger *m, uint32_t, uint32_t, uint32_t, uint64_t, const uint8_t *, size_t,
                             void *);
    void (*file_filecontrol)(struct Messenger *m, uint32_t, uint32_t, unsigned int, void *);
    void (*file_filedata)(struct Messenger *m, uint32_t, uint32_t, uint64_t, const uint8_t *, size_t, void *);
    void (*file_reqchunk)(struct Messenger *m, uint32_t, uint32_t, uint64_t, size_t, void *);

    void (*msi_packet)(struct Messenger *m, uint32_t, const uint8_t *, uint16_t, void *);
    void *msi_packet_userdata;

    void (*lossy_packethandler)(struct Messenger *m, uint32_t, const uint8_t *, size_t, void *);
    void (*lossless_packethandler)(struct Messenger *m, uint32_t, const uint8_t *, size_t, void *);

    void (*core_connection_change)(struct Messenger *m, unsigned int, void *);
    unsigned int last_connection_status;

    Messenger_Options_Fake options;
};

typedef struct {
    /*Messenger*/void *m;
    /*Friend_Connections*/void *fr_c;

    Group_c_Fake *chats;
    uint32_t num_chats;
} Group_Chats;

extern void *group_get_object(/*const Group_Chats*/ void *g_c, uint32_t groupnumber);

enum {
    GROUPCHAT_STATUS_NONE,
    GROUPCHAT_STATUS_VALID,
    GROUPCHAT_STATUS_CONNECTED
};

/* return 1 if the groupnumber is not valid.
 * return 0 if the groupnumber is valid.
 */
static uint8_t groupnumber_not_valid(const Group_Chats *g_c, int groupnumber)
{
    if ((unsigned int)groupnumber >= g_c->num_chats) {
        return 1;
    }

    if (g_c->chats == NULL) {
        return 1;
    }

    if (g_c->chats[groupnumber].status == GROUPCHAT_STATUS_NONE) {
        return 1;
    }

    return 0;
}

// padding cause offset not really size, can not use offset adder
static Group_c_Fake *get_group_c(Tox *tox, int groupnumber)
{
    int fos = offsetof(struct Messenger, conferences_object);
    // int fos2 = offsetof(Group_c_Fake, identifier);
    int conferences_object_offset = fos;

    char **p = (char**)(&((char*)tox)[0] + conferences_object_offset);
    // void *p2 = ((struct Messenger*)tox)->conferences_object;
    Group_Chats* grpchats = (Group_Chats*)(*p);
    if (groupnumber_not_valid(grpchats, groupnumber)) {
        return 0;
    }
    Group_c_Fake *g = &grpchats->chats[groupnumber];
    return g;
}

void tox_conference_get_identifier(Tox *tox, uint32_t conference_number, void *idbuf) {
    Group_c_Fake *g = get_group_c(tox, conference_number);
    if (g) {
        memcpy(idbuf, g->identifier, GROUP_IDENTIFIER_LENGTH);
    }
}
void tox_conference_get_pubkey(Tox *tox, uint32_t conference_number, void *pkbuf) {
    Group_c_Fake *g = get_group_c(tox, conference_number);
    if (g) {
        memcpy(pkbuf, g->real_pk, CRYPTO_PUBLIC_KEY_SIZE);
    }
}

