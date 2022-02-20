#include "php.h"
#include "zend.h"
#include "zend_execute.h"
#include "zend_types.h"
#include "php_globals.h"
#include "ext/standard/php_smart_str.h"
#include "ext/standard/info.h"
#include "ext/standard/file.h"
#include "ext/standard/url.h"
#include "ext/curl/php_curl.h"
#include "_cgo_export.h"

#include <stdio.h>
#include <curl/curl.h>
#include <curl/easy.h>

void yam_globals(char *typ, int sl, unsigned long long ufn)
{
    zval **container;
    zend_bool jit_initialization = PG(auto_globals_jit);
    if (jit_initialization)
        zend_is_auto_global(ZEND_STRL(typ) TSRMLS_CC);
    if (zend_hash_find(&EG(symbol_table), typ, sl + 1, (void **)&container) == SUCCESS)
    {
        void *fn = (void *)ufn;
        globalsCallback(container, fn);
    }
}

void yam_global_by_key(char *typ, int sl, char *key, int kl, unsigned long long ufn)
{
    zval **container;
    zval **value;
    char *res;
    zend_bool jit_initialization = PG(auto_globals_jit);
    if (jit_initialization)
        zend_is_auto_global(ZEND_STRL(typ) TSRMLS_CC);
    if (zend_hash_find(&EG(symbol_table), typ, sl + 1, (void **)&container) == SUCCESS)
    {
        if (zend_hash_find(Z_ARRVAL_PP(container), key, kl + 1, (void **)&value) == SUCCESS)
        {
            res = Z_STRVAL_PP(value);
            void *fn = (void *)ufn;
            globalByKeyCallback(res, fn);
        }
    }
}

void yam_curl_getinfo(void *ex, unsigned long long ufn)
{
    zend_execute_data *execute_data = (zend_execute_data *)ex;
    void **p = execute_data->function_state.arguments;
    int arg_count = (int)(zend_uintptr_t)*p;
    zval *ch = NULL;
    int i = 0;
    ch = *(p - (arg_count - i));
    if (ch == NULL)
        return;
    zval func, ret;
    zval *args[1];
    args[0] = ch;
    ZVAL_STRING(&func, "curl_getinfo", 0);
    call_user_function(EG(function_table), NULL, &func, &ret, 1, args TSRMLS_CC);
    void *fn = (void *)ufn;
    zval *pr = &ret;
    curlGetInfoCallback(&pr, fn);
}

void yam_curl_inject_trace_header(int id, struct curl_slist *sl)
{
    int actual_resource_type;
    php_curl *ch = zend_list_find(id, &actual_resource_type);
    if (!ch)
        return;
    zend_hash_index_update(ch->to_free->slist, CURLOPT_HTTPHEADER, &sl, sizeof(struct curl_slist *), NULL);
    curl_easy_setopt(ch->cp, CURLOPT_HTTPHEADER, sl);
}

//void *yam_read_property(void *pobj, const char *property)
//{
//    if (obj == NULL)
//        return NULL;
//    zval *obj = (zval *)pobj;
//    if (Z_TYPE_P(obj) != IS_OBJECT)
//        return NULL;
//    zend_class_entry *ce = zend_get_class_entry(obj);
//    return (void *)zend_read_property(ce, obj, property, strlen(property), 0);
//}

const char *yam_redis_version()
{
    zend_module_entry *module;
    if (zend_hash_find(&module_registry, "redis", sizeof("redis"), (void **)&module) == FAILURE) {
		return NULL;
	}
    return module->version;
}

void *yam_get_redis_socket_gte200_lte228(void *pz)
{
    zval *id = (zval *)pz;
	zval **socket = NULL;
	int resource_type;
	if (Z_TYPE_P(id) != IS_OBJECT || zend_hash_find(Z_OBJPROP_P(id), "socket", sizeof("socket"), (void **)&socket) == FAILURE)
	    return NULL;
	return zend_list_find(Z_LVAL_PP(socket), &resource_type);
}

void *yam_get_redis_socket_gte400_lte430(void *pz)
{
    zval *id = (zval *)pz;
	yam_redis_object *redis = NULL;
	int resource_type;
	if (Z_TYPE_P(id) != IS_OBJECT)
	    return NULL;
	redis = (yam_redis_object *)zend_objects_get_address(id TSRMLS_CC);
	if (redis->sock) {
        return redis->sock;
    }
	return NULL;
}
