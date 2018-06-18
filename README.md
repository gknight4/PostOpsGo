These are the files to run the Go backend server for PostOps.

The easiest install is just:\
go get github.com/gknight4/postops-go

that should install these files, *and* the dependencies.

Or, you can go straight to git:\
git clone https://github.com/gknight4/postops-go

Then these packages need to be installed:

go get github.com/dgrijalva/jwt-go\
go get github.com/gorilla/mux\
go get golang.org/x/crypto/bcrypt\
go get gopkg.in/mgo.v2

Then run the server:\
go run *.go

The server is available at:\
http://localhost:6026


Changelog:

Version: 0.2.1

Implement proxy mode\
disable tls certificate checking as https client.


Version: 0.1.5

Added SSL startup if run on server, not dev

Version: 0.1.4

First deployed to PostOps.us