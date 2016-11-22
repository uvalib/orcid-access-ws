if [ -z "$GOPATH" ]; then
   echo "ERROR: GOPATH is not defined"
   exit 1
fi

res=0
if [ $res -eq 0 ]; then
  GOOS=darwin go build -a -o bin/orcid-access-ws.darwin orcidaccessws
  res=$?
fi

if [ $res -eq 0 ]; then
  CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/orcid-access-ws.linux orcidaccessws
  res=$?
fi

exit $res
