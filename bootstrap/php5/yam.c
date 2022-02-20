#ifdef HAVE_CONFIG_H
#include "config.h"
#endif

#include "php.h"
#include "php_ini.h"
#include "ext/standard/info.h"
#include "php_yam.h"

ZEND_DECLARE_MODULE_GLOBALS(yam)

static int le_yam;

static int library_loaded = 0;
static void *library_handle = NULL;

ZEND_API void (*origin_execute_ex)(zend_execute_data *execute_data TSRMLS_DC);
ZEND_API void (*origin_execute_internal)(zend_execute_data *execute_data, struct _zend_fcall_info *fci, int return_value_used TSRMLS_DC);

long (*ModuleStartup)(void);
long (*ModuleShutdown)(void);
long (*RequestStartup)(void);
long (*RequestShutdown)(void);
void (*BeforeExecuteEx)(void *, void *);
void (*AfterExecuteEx)(void *, void *);
void (*BeforeExecuteInternal)(void *, void *);
void (*AfterExecuteInternal)(void *, void *);

const zend_function_entry yam_functions[] = {
    PHP_FE(confirm_yam_compiled, NULL) /* For testing, remove later. */
    PHP_FE_END                         /* Must be the last line in yam_functions[] */
};

zend_module_entry yam_module_entry = {
#if ZEND_MODULE_API_NO >= 20010901
    STANDARD_MODULE_HEADER,
#endif
    "yam",
    yam_functions,
    PHP_MINIT(yam),
    PHP_MSHUTDOWN(yam),
    PHP_RINIT(yam),     /* Replace with NULL if there's nothing to do at request start */
    PHP_RSHUTDOWN(yam), /* Replace with NULL if there's nothing to do at request end */
    PHP_MINFO(yam),
#if ZEND_MODULE_API_NO >= 20010901
    PHP_YAM_VERSION,
#endif
    STANDARD_MODULE_PROPERTIES};

#ifdef COMPILE_DL_YAM
ZEND_GET_MODULE(yam)
#endif

PHP_INI_BEGIN()
STD_PHP_INI_ENTRY("yam.library", "", PHP_INI_ALL, OnUpdateString, library, zend_yam_globals, yam_globals)
PHP_INI_END()

static void php_yam_init_globals(zend_yam_globals *yam_globals)
{
    yam_globals->library = NULL;
}

PHP_MINIT_FUNCTION(yam)
{
    REGISTER_INI_ENTRIES();
    origin_execute_ex = zend_execute_ex;
    zend_execute_ex = hook_execute_ex;
    origin_execute_internal = zend_execute_internal;
    zend_execute_internal = hook_execute_internal;
    return SUCCESS;
}

PHP_MSHUTDOWN_FUNCTION(yam)
{
    yam_unload_library();
    zend_execute_ex = origin_execute_ex;
    zend_execute_internal = origin_execute_internal;
    UNREGISTER_INI_ENTRIES();
    return SUCCESS;
}

PHP_RINIT_FUNCTION(yam)
{
    yam_load_library();
    if (library_loaded)
    {
        RequestStartup();
    }
    return SUCCESS;
}

PHP_RSHUTDOWN_FUNCTION(yam)
{
    if (library_loaded)
    {
        RequestShutdown();
    }
    return SUCCESS;
}

PHP_MINFO_FUNCTION(yam)
{
    php_info_print_table_start();
    php_info_print_table_header(2, "yam support", "enabled");
    php_info_print_table_end();
    DISPLAY_INI_ENTRIES();
}

/* Remove the following function when you have successfully modified config.m4
   so that your module can be compiled into PHP, it exists only for testing
   purposes. */

/* Every user-visible function in PHP should document itself in the source */
/* {{{ proto string confirm_yam_compiled(string arg)
   Return a string to confirm that the module is compiled in */
PHP_FUNCTION(confirm_yam_compiled)
{
    char *arg = NULL;
    int arg_len, len;
    char *strg;

    if (zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "s", &arg, &arg_len) == FAILURE)
    {
        return;
    }

    len = spprintf(&strg, 0, "Congratulations! You have successfully modified ext/%.78s/config.m4. Module %.78s is now compiled into PHP.", "yam", arg);
    RETURN_STRINGL(strg, len, 0);
}
/* }}} */
/* The previous line is meant for vim and emacs, so it can correctly fold and 
   unfold functions in source code. See the corresponding marks just before 
   function definition, where the functions purpose is also documented. Please 
   follow this convention for the convenience of others editing your code.
*/

int yam_load_library()
{
    if (library_loaded)
    {
        return 0;
    }
    char *path = YAM_G(library);
    if (path == "")
    {
        zend_error(E_WARNING, "yam_load_library() yam.library path is empty [skip]");
        return 1;
    }
    library_handle = DL_LOAD(path);
    if (!library_handle)
    {
        zend_error(E_WARNING, "yam_load_library() failed: %d %s", -1, path);
        return -1;
    }
    ModuleStartup = (long (*)(void))DL_FETCH_SYMBOL(library_handle, "ModuleStartup");
    if (!ModuleStartup)
    {
        DL_UNLOAD(library_handle);
        return -2;
    }
    ModuleShutdown = (long (*)(void))DL_FETCH_SYMBOL(library_handle, "ModuleShutdown");
    if (!ModuleShutdown)
    {
        DL_UNLOAD(library_handle);
        return -2;
    }
    RequestStartup = (long (*)(void))DL_FETCH_SYMBOL(library_handle, "RequestStartup");
    if (!RequestStartup)
    {
        DL_UNLOAD(library_handle);
        return -2;
    }
    RequestShutdown = (long (*)(void))DL_FETCH_SYMBOL(library_handle, "RequestShutdown");
    if (!RequestShutdown)
    {
        DL_UNLOAD(library_handle);
        return -2;
    }
    BeforeExecuteEx = (void (*)(void *, void *))DL_FETCH_SYMBOL(library_handle, "BeforeExecuteEx");
    if (!BeforeExecuteEx)
    {
        DL_UNLOAD(library_handle);
        return -2;
    }
    AfterExecuteEx = (void (*)(void *, void *))DL_FETCH_SYMBOL(library_handle, "AfterExecuteEx");
    if (!AfterExecuteEx)
    {
        DL_UNLOAD(library_handle);
        return -2;
    }
    BeforeExecuteInternal = (void (*)(void *, void *))DL_FETCH_SYMBOL(library_handle, "BeforeExecuteInternal");
    if (!BeforeExecuteInternal)
    {
        DL_UNLOAD(library_handle);
        return -2;
    }
    AfterExecuteInternal = (void (*)(void *, void *))DL_FETCH_SYMBOL(library_handle, "AfterExecuteInternal");
    if (!AfterExecuteInternal)
    {
        DL_UNLOAD(library_handle);
        return -2;
    }
    library_loaded = 1;
    ModuleStartup();
    return 0;
}

int yam_unload_library()
{
    if (library_loaded)
    {
        library_loaded = 0;
        ModuleShutdown();
        DL_UNLOAD(library_handle);
        library_handle = NULL;
    }
    return 0;
}

void hook_execute_ex(zend_execute_data *execute_data TSRMLS_DC)
{
    if (!library_loaded)
    {
        origin_execute_ex(execute_data);
        return;
    }
    BeforeExecuteEx(NULL, execute_data);
    origin_execute_ex(execute_data);
    AfterExecuteEx(NULL, execute_data);
}

void hook_execute_internal(zend_execute_data *execute_data, struct _zend_fcall_info *fci, int return_value_used TSRMLS_DC)
{
    if (!library_loaded)
    {
        if (origin_execute_internal)
            origin_execute_internal(execute_data, fci, return_value_used);
        return;
    }
    zval **return_value = &EX_TMP_VAR(execute_data, execute_data->opline->result.var)->var.ptr;
    BeforeExecuteInternal(execute_data, return_value);
    if (origin_execute_internal)
        origin_execute_internal(execute_data, fci, return_value_used);
    else
        execute_internal(execute_data, fci, return_value_used);
    AfterExecuteInternal(execute_data, return_value);
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: noet sw=4 ts=4 fdm=marker
 * vim<600: noet sw=4 ts=4
 */
