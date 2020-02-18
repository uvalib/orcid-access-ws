# set blank options variables
DBSECURE_OPT=""
DBHOST_OPT=""
DBNAME_OPT=""
DBUSER_OPT=""
DBPASSWD_OPT=""
ORCID_PUBLIC_URL_OPT=""
ORCID_SECURE_URL_OPT=""
ORCID_OAUTH_URL_OPT=""
ORCID_CLIENT_ID_OPT=""
ORCID_CLIENT_SECRET_OPT=""
SECRET_OPT=""
TIMEOUT_OPT=""
DEBUG_OPT=""

# secure database access
if [ -n "$DBSECURE" ]; then
   DBSECURE_OPT="--dbsecure=$DBSECURE"
fi

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

# ORCID OAuth endpoint URL
if [ -n "$ORCID_OAUTH_URL" ]; then
   ORCID_OAUTH_URL_OPT="--orcidoauthurl $ORCID_OAUTH_URL"
fi

# ORCID client identifier
if [ -n "$ORCID_CLIENT_ID" ]; then
   ORCID_CLIENT_ID_OPT="--orcidclientid $ORCID_CLIENT_ID"
fi

# ORCID client secret
if [ -n "$ORCID_CLIENT_SECRET" ]; then
   ORCID_CLIENT_SECRET_OPT="--orcidclientsecret $ORCID_CLIENT_SECRET"
fi

# shared secret
if [ -n "$AUTH_SHARED_SECRET" ]; then
   SECRET_OPT="--secret $AUTH_SHARED_SECRET"
fi

# ORCID service timeout
if [ -n "$SERVICE_TIMEOUT" ]; then
   TIMEOUT_OPT="--timeout $SERVICE_TIMEOUT"
fi

# service debugging
if [ -n "$ORCIDACCESS_DEBUG" ]; then
   DEBUG_OPT="--debug=$ORCIDACCESS_DEBUG"
fi

bin/orcid-access-ws $DBSECURE_OPT $DBHOST_OPT $DBNAME_OPT $DBUSER_OPT $DBPASSWD_OPT $ORCID_PUBLIC_URL_OPT $ORCID_SECURE_URL_OPT $ORCID_OAUTH_URL_OPT $ORCID_CLIENT_ID_OPT $ORCID_CLIENT_SECRET_OPT $SECRET_OPT $TIMEOUT_OPT $DEBUG_OPT

#
# end of file
#
