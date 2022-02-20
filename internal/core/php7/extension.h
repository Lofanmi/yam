#ifndef _YAM_EXTENSION_H_
#define _YAM_EXTENSION_H_

extern void *yam_get_module(char *name, char *version);
extern void yam_globals(char *typ, unsigned long long ufn);
extern void yam_curl_getinfo(void *ch, unsigned long long ufn);

#endif
