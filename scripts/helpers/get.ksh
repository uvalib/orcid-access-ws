if [ -z "$GOPATH" ]; then
   echo "ERROR: $GOPATH not defined"
   exit 1
fi

cd $GOPATH/src
rm -fr vendor

go get -u github.com/FiloSottile/gvt

gvt fetch github.com/gorilla/mux
gvt fetch github.com/parnurzeal/gorequest
gvt fetch gopkg.in/xmlpath.v1
gvt fetch -tag v1.2 github.com/go-sql-driver/mysql

# for tests
gvt fetch gopkg.in/yaml.v2
