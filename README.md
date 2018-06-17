These are the files to run the backend server for PostOps.

These packages need to be installed:

go get github.com/dgrijalva/jwt-go\
go get github.com/gorilla/mux\
go get golang.org/x/crypto/bcrypt\
go get gopkg.in/mgo.v2\

Then run the server:\
go run *.go

The server is available at:\
http://localhost:6026


Changelog:

Version: 0.2.1

Implement proxy mode\
disable tls certificate checking as https client.\


Version: 0.1.5

Added SSL startup if run on server, not dev

Version: 0.1.4

First deployed to PostOps.us