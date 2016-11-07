#
#
#

# ensure we have an endpoint
if [ -z "$ORCID_ACCESS_URL" ]; then
   echo "ERROR: ORCID_ACCESS_URL is not defined"
   exit 1
fi

# issue the command
echo "$ORCID_ACCESS_URL"
curl $ORCID_ACCESS_URL/version

exit 0

#
# end of file
#
