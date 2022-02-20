PHP_ARG_ENABLE(yam, whether to enable yam support,
Make sure that the comment is aligned:
[  --enable-yam           Enable yam support])

if test "$PHP_YAM" != "no"; then
  PHP_SUBST(YAM_SHARED_LIBADD)
  PHP_NEW_EXTENSION(yam, yam.c, $ext_shared)
fi
