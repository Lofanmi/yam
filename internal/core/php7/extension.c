#include <stdio.h>
#include <php.h>
#include <php_globals.h>
#include <zend_execute.h>
#include <zend_types.h>

#include "_cgo_export.h"

static zend_module_entry g_entry = {0};

void (*origin_execute_ex)(zend_execute_data *execute_data);
void (*origin_execute_internal)(zend_execute_data *execute_data, zval *return_value);

void hook_execute_ex(zend_execute_data *execute_data)
{
    BeforeExecuteEx(execute_data);
    origin_execute_ex(execute_data);
    AfterExecuteEx(execute_data);
}

void hook_execute_internal(zend_execute_data *execute_data, zval *return_value)
{
    BeforeExecuteInternal(execute_data, return_value);
    if (origin_execute_internal)
        origin_execute_internal(execute_data, return_value);
    else
        execute_internal(execute_data, return_value);
    AfterExecuteInternal(execute_data, return_value);
}

int yam_module_startup_func(int type, int module_number)
{
    ModuleStartup(type, module_number);
    origin_execute_ex = zend_execute_ex;
    zend_execute_ex = hook_execute_ex;
    origin_execute_internal = zend_execute_internal;
    zend_execute_internal = hook_execute_internal;
    return 0;
}

int yam_module_shutdown_func(int type, int module_number)
{
    ModuleShutdown(type, module_number);
    zend_execute_internal = origin_execute_internal;
    zend_execute_ex = origin_execute_ex;
    return 0;
}

int yam_request_startup_func(int type, int module_number)
{
    RequestStartup(type, module_number);
    return 0;
}

int yam_request_shutdown_func(int type, int module_number)
{
    RequestShutdown(type, module_number);
    return 0;
}

void *yam_get_module(char *name, char *version)
{
    zend_module_entry te = {
        STANDARD_MODULE_HEADER,
        name,
        NULL,
        yam_module_startup_func,
        yam_module_shutdown_func,
        yam_request_startup_func,
        yam_request_shutdown_func,
        NULL,
        version,
        STANDARD_MODULE_PROPERTIES};
    memcpy(&g_entry, &te, sizeof(zend_module_entry));
    return &g_entry;
};

void yam_globals(char *typ, unsigned long long ufn)
{
    zval *container;
    zend_bool jit_initialization = PG(auto_globals_jit);
    if (jit_initialization)
    {
        zend_is_auto_global_str(ZEND_STRL(typ));
    }
    container = zend_hash_str_find(&EG(symbol_table), ZEND_STRL(typ));
    void *p = NULL;
    if (UNEXPECTED(Z_TYPE_P(container) != IS_ARRAY))
    {
        if (Z_TYPE_P(container) != IS_REFERENCE || Z_TYPE_P(Z_REFVAL_P(container)) != IS_ARRAY)
        {
            return;
        }
        p = (void *)Z_REFVAL_P(container);
    }
    else
    {
        p = (void *)container;
    }
    void *fn = (void *)ufn;
    globalsCallback(p, fn);
}

void yam_curl_getinfo(void *ch, unsigned long long ufn)
{
    zval args[1];
    ZVAL_COPY(&args[0], (zval *)ch);
    zval func, ret;
    ZVAL_STRING(&func, "curl_getinfo");
    call_user_function(CG(function_table), NULL, &func, &ret, 1, args);
    void *fn = (void *)ufn;
    curlGetInfoCallback(&ret, fn);
    zval_ptr_dtor(&args[0]);
    zval_ptr_dtor(&func);
    zval_ptr_dtor(&ret);
}
