package userdb

import (
	"sync"
)

type userRecord struct {
	id int
	name string
	passwd string
	email string
	addr  string
	mu    sync.Mutex
}

const maxNumOfColumns = 1000                     // 10 for testing. 10K in real life?

type columnType [maxNumOfColumns]userRecord

type stableStorage struct {
	numRows_ int
	numCols_ int
	matrix_  []columnType
	sm       sync.Mutex
}

func (s *stableStorage) addCell() *userRecord {

	// (numRows_, numCols_) points to the cell to be assigned

	s.sm.Lock()
	row, col := s.numRows_, s.numCols_
	if col == 0 {
		s.matrix_ = append( s.matrix_, columnType{} )
	}
	if col + 1 ==  maxNumOfColumns {
		s.numRows_ += 1
		s.numCols_ = 0
	} else {
		s.numCols_ += 1
	}
	u := &s.matrix_[row][col]
	s.sm.Unlock()

	return u
}

var stb stableStorage
var userMap map[int] *userRecord	// UID -> record in storage
var name2idMap map[string] int 		// name -> UID
var cm  sync.Mutex					// protect record creation
var uidCounter int					// global UID generator

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

func init() {
	ResetDb()
}

func ResetDb() {
	cm.Lock()
	userMap = make( map[int] *userRecord )
	name2idMap = make( map[string] int)
	uidCounter = 1000
	cm.Unlock()
	// Note: no need to reset the storage
}

func CreateUser( r *RequestCreateUser ) (userId int, err error) {
	if ( r.Name == "" ) {
		err = requestError { "Empty name" }
		return
	}

	newuser := new( userRecord )
	newuser.name = r.Name
	newuser.passwd = r.Passwd
	newuser.email = r.Email
	newuser.addr = r.Addr

	cm.Lock()
	
	userId, ok := name2idMap[ r.Name ] 
	if ( ok ) {
		cm.Unlock()
		err = requestError { "Name already reserved" }
		return
	}

	userId = uidCounter
	uidCounter = uidCounter + 1

	name2idMap[ r.Name ] = userId

	ur := stb.addCell()
	userMap[userId] = ur
	ur.mu.Lock()		// the record just created, so there will be no wait
	cm.Unlock()         // unlock global mutex but hold the local one
	defer ur.mu.Unlock()

	ur.name = r.Name         // initialize the record
	ur.passwd = r.Passwd
	ur.email = r.Email
	ur.addr = r.Addr

	err = nil
	return
}

func UpdateUser( uid int, r* RequestUpdateUser ) error {
	cm.Lock()
	user, ok := userMap[ uid ]
	cm.Unlock()

	if ( !ok ) {
		return requestError { "UID not found"}
	}

	user.mu.Lock()
	defer user.mu.Unlock()

	// now update only requested fields
	if  r.Email != ""  {
		user.email = r.Email
	}
	if r.Addr != "" {
		user.addr = r.Addr
	}

	return nil
}

func GetUser( uid int) (resp ResponseGetUser, err error) {
	err = nil

	cm.Lock()
	user, ok := userMap[ uid ]
	cm.Unlock()

	if ( !ok ) {
		err = requestError { "UID not found"}
		return
	}

	user.mu.Lock()
	defer user.mu.Unlock()

	resp.Name = user.name
	resp.Email = user.email
	resp.Addr = user.addr

	return
}
