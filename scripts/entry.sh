# set blank options variables
DBHOST_OPT=""
DBNAME_OPT=""
DBUSER_OPT=""
DBPASSWD_OPT=""
ORCID_PUBLIC_URL_OPT=""
ORCID_SECURE_URL_OPT=""
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

# ORCID public endpoint URL
if [ -n "$ORCID_PUBLIC_API_URL" ]; then
   ORCID_PUBLIC_URL_OPT="--orcidpublicurl $ORCID_PUBLIC_API_URL"
fi

# ORCID secure endpoint URL
if [ -n "$ORCID_SECURE_API_URL" ]; then
   ORCID_SECURE_URL_OPT="--orcidsecureurl $ORCID_SECURE_API_URL"
fi

# token authentication service URL
if [ -n "$TOKENAUTH_URL" ]; then
   TOKENURL_OPT="--tokenauth $TOKENAUTH_URL"
fi

# ORCID service timeout
if [ -n "$SERVICE_TIMEOUT" ]; then
   TIMEOUT_OPT="--timeout $SERVICE_TIMEOUT"
fi

# service debugging
if [ -n "$ORCIDACCESS_DEBUG" ]; then
   DEBUG_OPT="--debug"
fi

bin/orcid-access-ws $DBHOST_OPT $DBNAME_OPT $DBUSER_OPT $DBPASSWD_OPT $ORCID_PUBLIC_URL_OPT $ORCID_SECURE_URL_OPT $TOKENURL_OPT $TIMEOUT_OPT $DEBUG_OPT

#
# end of file
#
