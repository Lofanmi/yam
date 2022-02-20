#ifndef _YAM_EXTENSION_H_
#define _YAM_EXTENSION_H_

#include "zend.h"

struct curl_slist;

typedef struct {
    zend_object std;
    void *sock;
} yam_redis_object;

extern void yam_globals(char *typ, int sl, unsigned long long ufn);
extern void yam_global_by_key(char *typ, int sl, char *key, int kl, unsigned long long ufn);
extern void yam_curl_getinfo(void *ex, unsigned long long ufn);
extern void yam_curl_inject_trace_header(int id, struct curl_slist *sl);
extern const char *yam_redis_version();
void *yam_get_redis_socket_gte200_lte228(void *pz);
void *yam_get_redis_socket_gte400_lte430(void *pz);

#endif
