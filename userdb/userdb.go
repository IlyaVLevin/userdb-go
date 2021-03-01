package userdb

import (
	"sync"
	"sync/atomic"
)

type userRecord struct {
	id int
	name string
	passwd string
	email string
	addr  string
}

var userMap map[int] atomic.Value           // UID -> userRecord
var name2idMap map[string] int              // name -> UID
var cm  sync.Mutex                          // protect record creation
var uidCounter int                          // global UID generator

type requestError struct {
	text string
}

func (e requestError) Error() string { return e.text }

type RequestCreateUser struct {
	Name string
	Passwd string
	Email string
	Addr  string
}

type RequestUpdateUser struct {
	Email string
	Addr  string
}

type ResponseGetUser struct {
	Name string
	Email string
	Addr  string
}

func Init() {
	userMap = make( map[int] atomic.Value )
	name2idMap = make( map[string] int)
	uidCounter = 1000
}

func CreateUser( r *RequestCreateUser ) (userId int, err error) {
	if ( r.Name == "" ) {
		err = requestError { "Empty name" }
		return
	}

	cm.Lock()
	defer cm.Unlock()
	
	userId = uidCounter
	uidCounter = uidCounter + 1
	
	_, ok := name2idMap[ r.Name ] 
	if ( ok ) {
		err = requestError { "Name already reserved" }
		return
	}
	newuser := new( userRecord )
	newuser.name = r.Name
	newuser.passwd = r.Passwd
	newuser.email = r.Email
	newuser.addr = r.Addr

	name2idMap[ r.Name ] = userId
	curu := userMap[userId]
	curu.Store( newuser )
	
	err = nil
	return
}

func UpdateUser( uid int, r* RequestUpdateUser ) error {
	curu, ok := userMap[ uid ] 
	if ( !ok ) {
		err := requestError { "UID not found" }
		return err
	}

	oldUser := curu.Load().(*userRecord)

	newUser := new( userRecord )
	*newUser = *oldUser

	// now update only requested fields
	if  r.Email != ""  {
		newUser.email = r.Email
	}
	if r.Addr != "" {
		newUser.addr = r.Addr
	}

	curu.Store( newUser )    // oldUser to be garbage collected
	
	return nil
}

func GetUser( uid int) (resp ResponseGetUser, err error) {
	err = nil

	curu, ok := userMap[ uid ] 
	if ( !ok ) {
		err = requestError { "UID not found"}
		return
	}

	user := curu.Load().(userRecord)

	resp.Name = user.name
	resp.Email = user.email
	resp.Addr = user.addr

	return
}
