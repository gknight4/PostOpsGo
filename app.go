package main 

import (
	"net/http"
  "os"
//  "path/filepath"
//  "io/ioutil"
//	"time"
	"encoding/json"
	"log"
	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"
	jwt "github.com/dgrijalva/jwt-go"
)

/*
install:
go get github.com/dgrijalva/jwt-go
go get github.com/gorilla/mux
go get golang.org/x/crypto/bcrypt
go get gopkg.in/mgo.v2

 * */

func addCors (w http.ResponseWriter){
	w.Header().Set("Access-Control-Allow-Origin", "*")	
	w.Header().Set("Access-Control-Allow-Headers", 
		"Accept, Accept-Encoding, Accept-Language, Authorization, Cache-Control, Connection, Content-Length, Content-Type, Host, Origin, Pragma, Referer, User-Agent")	
	w.Header().Set("Access-Control-Expose-Headers", 
		"Accept, Accept-Encoding, Accept-Language, Authorization, Cache-Control, Connection, Content-Length, Content-Type, Host, Origin, Pragma, Referer, User-Agent")	
	w.Header().Set("Access-Control-Allow-Methods", 
		"GET, POST, PUT, DELETE, PATCH")	
}

func corsOptions (w http.ResponseWriter, r *http.Request){
//	lo("options")
	addCors(w)
//	w.Header().Set("Access-Control-Allow-Origin", "*")	
//	w.Header().Set("Access-Control-Allow-Headers", 
//		"Accept, Accept-Encoding, Accept-Language, Authorization, Cache-Control, Connection, Content-Length, Content-Type, Host, Origin, Pragma, Referer, User-Agent")	
	ReturnJson(w, http.StatusOK, Msr(bToOk(false)))
}

func rootResponse (w http.ResponseWriter, r *http.Request){
  lo("root response")
	defer r.Body.Close()
	addCors(w)
	ReturnJson(w, http.StatusOK, Msr("root response"))
}

func checkUser (w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	lo("check user")
	addCors(w)
	vars := mux.Vars(r)
        got := checkUserName(vars["username"])
	ReturnJson(w, http.StatusOK, Msr(bToOk(got)))
}

func checkPass (w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	addCors(w)
	vars := mux.Vars(r)
	got := checkPassword(vars["hash"])
	lo("pass")
	lo(got)
	ReturnJson(w, http.StatusOK, Msr(bToOk(got)))
}

func checkIp (w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	vars := mux.Vars(r)
	got := checkIpAddress(vars["ip"])
	ReturnJson(w, http.StatusOK, Msr(bToOk(got)))
}

func newUser (w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	addCors(w)
	var user User
	lo("new user")
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
    lo("error")
    lo(err)
//		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if _, err := getUser(user.Useremail); err == nil {// if user *is* found
    lo("user exists")
    ReturnJson(w, http.StatusOK, Msr("userexists"))
    return
  } else {
    lo("new user")
    user.Flags = 1 // standard user
    user.Password, _ = HashPassword (user.Password)
    insertUser(user)
    ReturnJson(w, http.StatusOK, Msr("ok"))
  }
}

func checkAuth (w http.ResponseWriter, r *http.Request){
  lo("check auth")
	addCors(w)
	ReturnJson(w, http.StatusOK, Msr(bToOk(true)))
}
/* there's two kinds of users: parents and children
parents login with parentname / password
children need a childname, too
*or* parents have username / pass
and children have
parentname, username / pass
parent registers with parentname, parentemail, password


*/

func authUserPass (w http.ResponseWriter, r *http.Request){
// username and password (hash) are passed in the url	
	defer r.Body.Close()
  lo("auth");
	addCors(w)
	vars := mux.Vars(r)
//	lo("user: " + vars["useremail"])
//	lo("pass: " + vars["password"])
	user, _ := getUser(vars["useremail"])
//  lo("found user: " + user.Useremail)
//  lo(user.Password)
	res := CheckPasswordHash(vars["password"], user.Password)
//  lo(res)
	if res {
    lo("OK")
		token := GetJwtToken(user, 1)
		w.Header().Set("Authorization", "Bearer " + token)
    lo(token);
//		w.Header().Set("Authorization", "Bearer " + token)
    ReturnJson(w, http.StatusOK, Msr(bToOk(res)))
	} else {
    ReturnJson(w, http.StatusOK, Msr("nosuchaccount"))
  }
//	lo("returning")
}

/*
{
"amount": 124,
"date": 20180305T030500,
"name": "some",
"userid": 12445
}
*/

func newTrans (w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	var trans Transaction
	if err := json.NewDecoder(r.Body).Decode(&trans); err != nil {
		return
	}
	insertTransaction(trans)
}

func proxyRequest  (w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
  var req HttpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    lo("error")
    lo(err)
    ReturnJson(w, http.StatusOK, Msr("err"))
		return
	}
	respData := doProxyRequest(req)
//	lo(req.Headers)
//  lo("proxy 1")
	addCors(w)
	ReturnJson(w, http.StatusOK, respData)
}

func newStringStore (w http.ResponseWriter, r *http.Request){
/* now receiving an *array* of type: text: entries
 **/  
//lo("new string store")
	defer r.Body.Close()
//  stringType := mux.Vars(r)["type"]
	addCors(w)
	var strsto []StringStore
	if err := json.NewDecoder(r.Body).Decode(&strsto); err != nil {
//    lo("error")
//    lo(err)
    ReturnJson(w, http.StatusOK, Msr("err"))
		return
	}
	claims, _ := r.Context().Value("claims").(jwt.MapClaims)
  userId := bson.ObjectIdHex(claims["id"].(string)) ;
//  lo("saving")
  for _, v := range strsto {
    v.UserId = userId
    lo(v)
    insertStringStore(v)
  }

//  strsto.UserId = bson.ObjectIdHex(claims["id"].(string))
//  strsto.Type = stringType
  
//	insertStringStore(strsto)
	ReturnJson(w, http.StatusOK, Msr("ok"))
}

func allStringStores (w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
  lo("get all string stores");
//  stringType := mux.Vars(r)["type"]
	addCors(w)
	claims, _ := r.Context().Value("claims").(jwt.MapClaims)
	strstos, _ := getStringStores(claims["id"].(string))
	ReturnJson(w, http.StatusOK, strstos)
}

func setStringStores (w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
  lo("set strings")
//  stringType := mux.Vars(r)["type"]
//  lo(stringType)
	addCors(w)
	var strstos StringStoreAll
	if err := json.NewDecoder(r.Body).Decode(&strstos); err != nil {
    ReturnJson(w, http.StatusOK, Msr("err"))
		return
	}
	claims, _ := r.Context().Value("claims").(jwt.MapClaims)
	userId := bson.ObjectIdHex(claims["id"].(string))
	removeAllStringStores(claims["id"].(string))
  for _, v := range strstos.Strings {
    v.UserId = userId
//    v.Type = stringType
//    lo(v)
    insertStringStore(v)
  }
  
  
/*  for i := 0 ; i < len(heads) ; i++ {
    lo(heads[i].Text)
    heads[i].UserId = userId
  }*/

//	head := make([]Header, 0)


  /*var heads HeaderArray
//  body, _ := ioutil.ReadOne(r.Body)
body, err := ioutil.ReadAll(r.Body)
    var arr []string
err2 := json.NewDecoder(r.Body).Decode(&arr)
lo(err2)
//    _ = json.Unmarshal([]byte(string(body)), &arr)
    _ = json.Unmarshal(body, &arr)
    
lo(string(body))
lo(err)
lo(arr)
lo(arr[0])
	if err := json.NewDecoder(r.Body).Decode(&heads); err != nil {
    lo(err)
    ReturnJson(w, http.StatusOK, Msr("err"))
		return
	}
//	lo(heads.heads[0])
//	claims, _ := r.Context().Value("claims").(jwt.MapClaims)
//	head.UserId = bson.ObjectIdHex(claims["id"].(string))
//	insertHeader(head)*/
	ReturnJson(w, http.StatusOK, Msr("ok"))
}





func lo(obj interface{}){ // args ...
	log.Println(obj)
}

//func lo(str string, obj interface{}){ // args ...
//	log.Println(str, obj)
//}

//     func (t *Interface) Empty() bool

func updateTrans (w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	vars := mux.Vars(r)
	var trans Transaction
	if err := json.NewDecoder(r.Body).Decode(&trans); err != nil {
		return
	}
	trans.ID = bson.ObjectIdHex(vars["transactionid"])
	updateTransaction(trans)
}

func deleteTrans (w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	vars := mux.Vars(r)
//	var trans Transaction
//	if err := json.NewDecoder(r.Body).Decode(&trans); err != nil {
//		return
//	}
//	trans.ID = bson.ObjectIdHex(vars["transactionid"])
	deleteTransaction(vars["transactionid"])
}

func allTrans (w http.ResponseWriter, r *http.Request){
	/* not quite right, yet
	this needs to search for the transactions associated with
	the userid that's making the request
also, there's a problem that this should retrieve all of the
parent's transactions, not just one child
	*/
	defer r.Body.Close()
	vars := mux.Vars(r)
//	userid := vars["userid"]
//	log.Println(userid)
	
//	var trans []Transaction
	trans, _ := getTransactions(vars["userid"])
	ReturnJson(w, http.StatusOK, trans)
}

func unAuthHTTPReturn (w http.ResponseWriter, r *http.Request){
	addCors(w)
	ReturnJson(w, http.StatusUnauthorized, Msr(bToOk(false)))
}


func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", corsOptions).Methods("OPTIONS")// 
	r.HandleFunc("/", rootResponse).Methods("GET")
	r.HandleFunc("/auth/stringstore", corsOptions).Methods("OPTIONS")
	r.HandleFunc("/auth/all/stringstore", corsOptions).Methods("OPTIONS")
	r.HandleFunc("/auth/proxy", corsOptions).Methods("OPTIONS")
  
	r.HandleFunc("/open/users/{username}", checkUser).Methods("GET")
	r.HandleFunc("/open/users/{username}", corsOptions).Methods("OPTIONS")
	r.HandleFunc("/open/commonpassword/{hash}", corsOptions).Methods("OPTIONS")
	r.HandleFunc("/open/users", corsOptions).Methods("OPTIONS")
//	r.Handle("/auth/users", AddContext(http.HandlerFunc(corsOptions))).Methods("OPTIONS")
	r.HandleFunc("/open/authenticate/{useremail}/{password}", corsOptions).Methods("OPTIONS")
	r.PathPrefix("/open/authenticate").Handler(http.HandlerFunc(corsOptions)).Methods("OPTIONS")
//	r.Handle("/auth/check", AddContext(http.HandlerFunc(corsOptions))).Methods("OPTIONS")
	r.HandleFunc("/auth/check", corsOptions).Methods("OPTIONS")
	r.HandleFunc("/auth/users", corsOptions).Methods("OPTIONS")
/***********note that "options" seem to not work as authorized, must be in the clear************/        

	r.HandleFunc("/open/commonpassword/{hash}", checkPass).Methods("GET")
	r.HandleFunc("/open/ips/{ip}", checkIp).Methods("GET")
	r.HandleFunc("/open/users", newUser).Methods("POST")
	r.HandleFunc("/open/authenticate/{useremail}/{password}", authUserPass).Methods("GET")
	r.Handle("/auth/check", AddContext(http.HandlerFunc(checkAuth))).Methods("GET")
	r.Handle("/auth/stringstore", AddContext(http.HandlerFunc(newStringStore))).Methods("POST")
	r.Handle("/auth/stringstore", AddContext(http.HandlerFunc(allStringStores))).Methods("GET")
	r.Handle("/auth/all/stringstore", AddContext(http.HandlerFunc(setStringStores))).Methods("PUT")
	r.Handle("/auth/proxy", AddContext(http.HandlerFunc(proxyRequest))).Methods("POST")

  //	r.HandleFunc("/auth/transactions", newTrans).Methods("POST")
	r.Handle("/auth/transactions", AddContext(http.HandlerFunc(newTrans))).Methods("POST")
	r.Handle("/auth/transactions/{userid}", AddContext(http.HandlerFunc(allTrans))).Methods("GET")
	r.Handle("/auth/transactions/{transactionid}", AddContext(http.HandlerFunc(updateTrans))).Methods("PUT")
	r.Handle("/auth/transactions/{transactionid}", AddContext(http.HandlerFunc(deleteTrans))).Methods("DELETE")
//	r.Handle("/auth/users", AddContext(http.HandlerFunc(newChild))).Methods("POST")
	r.PathPrefix("/open/authenticate").Handler(http.HandlerFunc(corsOptions)).Methods("GET")
	
	initMongo()
//	if err := http.ListenAndServe(":6026", r); err != nil {
//		log.Fatal(err)
//	}
// func ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error	
  name, _ := os.Hostname()
  if name == "genes" {
    if err := http.ListenAndServe(":6026", r); err != nil {
      log.Fatal(err)
    }
  } else {
    err := http.ListenAndServeTLS(":6026", "./ssl/cert.pem", "./ssl/privkey.pem", r)  
    if err != nil {
            log.Fatal(err)
    }	
  }
}

