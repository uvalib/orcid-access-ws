# set blank options variables
DBHOST_OPT=""
DBNAME_OPT=""
DBUSER_OPT=""
DBPASSWD_OPT=""
ORCIDURL_OPT=""
ORCIDUSER_OPT=""
ORCIDPASSWD_OPT=""
TOKENURL_OPT=""
TIMEOUT_OPT=""
DEBUG_OPT=""

# database host
if [ -n "$DBHOST" ]; then
   DBHOST_OPT="--dbhost $DBHOST"
fi

# database name
if [ -n "$DBNAME" ]; then
   DBNAME_OPT="--dbname $DBNAME"
fi

# database user
if [ -n "$DBUSER" ]; then
   DBUSER_OPT="--dbuser $DBUSER"
fi

# database password
if [ -n "$DBPASSWD" ]; then
   DBPASSWD_OPT="--dbpassword $DBPASSWD"
fi

# ORCID endpoint URL
if [ -n "$ORCID_URL" ]; then
   ORCIDURL_OPT="--orcidurl $ORCID_URL"
fi

# ORCID user name
#if [ -n "$ORCID_USER" ]; then
#   ORCIDUSER_OPT="--orciduser $ORCID_USER"
#fi

# ORCID password
#if [ -n "$ORCID_PASSWD" ]; then
#   ORCIDPASSWD_OPT="--orcidpassword $ORCID_PASSWD"
#fi

# token authentication service URL
if [ -n "$TOKENAUTH_URL" ]; then
   TOKENURL_OPT="--tokenauth $TOKENAUTH_URL"
fi

# ORCID service timeout
if [ -n "$ORCID_TIMEOUT" ]; then
   TIMEOUT_OPT="--timeout $ORCID_TIMEOUT"
fi

# service debugging
if [ -n "$ORCIDACCESS_DEBUG" ]; then
   DEBUG_OPT="--debug"
fi

bin/orcid-access-ws $DBHOST_OPT $DBNAME_OPT $DBUSER_OPT $DBPASSWD_OPT $ORCIDURL_OPT $ORCIDUSER_OPT $ORCIDPASSWD_OPT $TOKENURL_OPT $TIMEOUT_OPT $DEBUG_OPT

#
# end of file
#
