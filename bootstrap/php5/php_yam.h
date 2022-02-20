/*
  +----------------------------------------------------------------------+
  | PHP Version 5                                                        |
  +----------------------------------------------------------------------+
  | Copyright (c) 1997-2014 The PHP Group                                |
  +----------------------------------------------------------------------+
  | This source file is subject to version 3.01 of the PHP license,      |
  | that is bundled with this package in the file LICENSE, and is        |
  | available through the world-wide-web at the following url:           |
  | http://www.php.net/license/3_01.txt                                  |
  | If you did not receive a copy of the PHP license and are unable to   |
  | obtain it through the world-wide-web, please send a note to          |
  | license@php.net so we can mail you a copy immediately.               |
  +----------------------------------------------------------------------+
  | Author:                                                              |
  +----------------------------------------------------------------------+
*/

/* $Id$ */

#ifndef PHP_YAM_H
#define PHP_YAM_H

extern zend_module_entry yam_module_entry;
#define phpext_yam_ptr &yam_module_entry

#define PHP_YAM_VERSION "0.1.0" /* Replace with version number for your extension */

#ifdef PHP_WIN32
#	define PHP_YAM_API __declspec(dllexport)
#elif defined(__GNUC__) && __GNUC__ >= 4
#	define PHP_YAM_API __attribute__ ((visibility("default")))
#else
#	define PHP_YAM_API
#endif

#ifdef ZTS
#include "TSRM.h"
#endif

PHP_MINIT_FUNCTION(yam);
PHP_MSHUTDOWN_FUNCTION(yam);
PHP_RINIT_FUNCTION(yam);
PHP_RSHUTDOWN_FUNCTION(yam);
PHP_MINFO_FUNCTION(yam);

PHP_FUNCTION(confirm_yam_compiled);	/* For testing, remove later. */

int yam_load_library(void);
int yam_unload_library(void);
void hook_execute_ex(zend_execute_data *execute_data TSRMLS_DC);
void hook_execute_internal(zend_execute_data *execute_data, struct _zend_fcall_info *fci, int return_value_used TSRMLS_DC);


ZEND_BEGIN_MODULE_GLOBALS(yam)
	char *library;
ZEND_END_MODULE_GLOBALS(yam)

/* In every utility function you add that needs to use variables 
   in php_yam_globals, call TSRMLS_FETCH(); after declaring other 
   variables used by that function, or better yet, pass in TSRMLS_CC
   after the last function argument and declare your utility function
   with TSRMLS_DC after the last declared argument.  Always refer to
   the globals in your function as YAM_G(variable).  You are 
   encouraged to rename these macros something shorter, see
   examples in any other php module directory.
*/

#ifdef ZTS
#define YAM_G(v) TSRMG(yam_globals_id, zend_yam_globals *, v)
#else
#define YAM_G(v) (yam_globals.v)
#endif

#endif	/* PHP_YAM_H */


/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: noet sw=4 ts=4 fdm=marker
 * vim<600: noet sw=4 ts=4
 */
