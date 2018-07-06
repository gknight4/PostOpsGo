package main 

/*
import into mongo:
mongoimport -d users_db -c cp --type csv --file hash2.txt --headerline
*/

import (
	"log"
//	"github.com/BurntSushi/toml"
	"gopkg.in/mgo.v2/bson"
	mgo "gopkg.in/mgo.v2"
	"time"
)

const (
	USERCOLLECTION = "users"
	TRANSACTIONCOLLECTION = "transaction"
	HEADERCOLLECTION = "header"
	STRINGSTORECOLLECTION = "stringstore"
	CPCOLLECTION = "cp"
)

type User struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Useremail  string        `bson:"useremail" json:"useremail"`
	Password string        `bson:"password" json:"password"`
	Flags int        `bson:"flags" json:"flags"`
}

type StringStore struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Text        string        `bson:"text" json:"text"`
	Type        string        `bson:"type" json:"type"`
	UserId      bson.ObjectId `bson:"userid" json:"userid"`
}

type StringStoreAll struct {
  Type  string  `json:"type"`
  Strings []StringStore `json:strings`
}

type Url struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Text        string        `bson:"text" json:"text"`
	UserId      bson.ObjectId `bson:"userid" json:"userid"`
}

type Transaction struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Amount        int        `bson:"amount" json:"amount"`
	Date  time.Time        `bson:"date" json:"date"`
	Name string        `bson:"name" json:"name"`
	ParentId bson.ObjectId        `bson:"parentid" json:"parentid"`
	UserId bson.ObjectId        `bson:"userid" json:"userid"`
}

type Cp struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Hash        string        `bson:"hash" json:"hash"`
}

type UsersDAO struct {
	Server string
	Database string
	Username string
	Password string
	db *mgo.Database
}

var userdb = UsersDAO{}

func (u *UsersDAO) Connect() {
  info := &mgo.DialInfo{
    Addrs:[]string{u.Server},
    Timeout: 60*time.Second,
    Database:"admin",
    Username: u.Username,
    Password: u.Password}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		log.Fatal(err)
	}
	u.db = session.DB(u.Database)
}

func (u *UsersDAO) ConnectO() {
	session, err := mgo.Dial(u.Server)
	if err != nil {
		log.Fatal(err)
	}
	u.db = session.DB(u.Database)
}

func (u *UsersDAO) Insert(user User) error {
//	log.Println("inserting")
//	log.Println(user)
	err := u.db.C(USERCOLLECTION).Insert(&user)
	return err
}

func insertTransaction(trans Transaction) error {
	log.Println(trans)
	trans.ID = bson.NewObjectId()
	err := userdb.db.C(TRANSACTIONCOLLECTION).Insert(&trans)
	return err
}

//func insertFieldType() error {
//}

func insertStringStore(obj StringStore) error {
	obj.ID = bson.NewObjectId()
	return userdb.db.C(STRINGSTORECOLLECTION).
    Insert(&obj)
}

func getStringStores(userid string)([]StringStore, error){
	var arr []StringStore
	return arr, userdb.db.C(STRINGSTORECOLLECTION).
		Find(bson.M{"userid": bson.ObjectIdHex(userid)}).All(&arr)
}

func removeAllStringStores(userid string) error {
	_, err := userdb.db.C(STRINGSTORECOLLECTION).
		RemoveAll(bson.M{"userid": bson.ObjectIdHex(userid)})
	return err
}





func getTransactions(userid string)([]Transaction, error){
	var trans []Transaction
	log.Println(userid)
	err := userdb.db.C(TRANSACTIONCOLLECTION).
		Find(bson.M{"userid": bson.ObjectIdHex(userid)}).All(&trans)
	return trans, err
}

func updateTransaction(trans Transaction) error {
	err := userdb.db.C(TRANSACTIONCOLLECTION).
		UpdateId(trans.ID, &trans)
	return err
}

func deleteTransaction(transid string) error {
//	trans Transaction
//	trans.ID = bson.ObjectIdHex(transid)
lo("delete")
	err := userdb.db.C(TRANSACTIONCOLLECTION).
		RemoveId(bson.ObjectIdHex(transid))
	return err
}


func testTrans(){
	var trans Transaction
	trans.Amount = 123
	trans.Date = time.Now();
	trans.Name = "name"
	trans.UserId = bson.NewObjectId()
	insertTransaction(trans)
}

func insertUser(user User) error {
	user.ID = bson.NewObjectId()
	return userdb.Insert(user)
}

func (u *UsersDAO) FindByUserId(id string) (User, error) {
    var user User
    err := u.db.C(USERCOLLECTION).FindId(bson.ObjectIdHex(id)).One(&user)
    return user, err
}

func (u *UsersDAO) FindByUserName(username string) (User, error) {
	var user User
//	c := u.db.C(USERCOLLECTION)
	err := u.db.C(USERCOLLECTION).Find(bson.M{"username": username}).One(&user)
	return user, err
}

func (u *UsersDAO) FindByUserEmail(useremail string) (User, error) {
	var user User
	err := u.db.C(USERCOLLECTION).Find(bson.M{"useremail": useremail}).One(&user)
	return user, err
}

func (u *UsersDAO) FindCp(hash string) (Cp, error) {
	var cp Cp
//	c := u.db.C(USERCOLLECTION)
	err := u.db.C(CPCOLLECTION).Find(bson.M{"hash": hash}).One(&cp)
	return cp, err
}

func getUser (useremail string) (User, error) {
	return userdb.FindByUserEmail(useremail)
}

func getUserFromId (id string) (User, error) {
    return userdb.FindByUserId(id)
}

func checkUserName (useremail string) bool {
// return true if the username is found
  lo(useremail);
	_, err := userdb.FindByUserEmail(useremail)
	return err == nil
}

func checkPassword (hash string) bool {
// return true if the username is found	
	_, err := userdb.FindCp(hash)
	return err == nil
}


/* find by username:
db.inventory.find( { status: "A", qty: { $lt: 30 } } )
*/

/*func initUsers(){
	var user = User{}
	user.ID = bson.NewObjectId()
	user.Parentemail = "a@b.c"
	user.Parentname = "gknight4"
	user.Password = "pass"
	user.Username = "gknight4"
	userdb.Insert(user)
}

func testUsers(){
//	var user User
	user, err := userdb.FindByUserName("gknight4")
	log.Println(user)
	log.Println(err)
}

func testCp(){
// fea60547de42ef75d775f7ad5c3c8b60486380bbb06850943b5358d357e4fb46	
	hash, err := userdb.FindCp("fea60547de42ef75d775f7ad5c3c8b60486380bbb06850943b5358d357e4fb46")
	log.Println(hash)
	log.Println(err)
}*/

func initMongo(){
//	var userdb UsersDAO
	userdb.Server = "localhost"
	userdb.Database = "users_db"
//  lo(userdb) ;
  userdb.Username = "golang" ;
  userdb.Password = "WDsdC6uM"
	userdb.Connect()
//	testCp()
	
	
	
}
