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

const maxNumOfColumns = 10

type columnType [maxNumOfColumns]userRecord

type stableStorage struct {
	numRows_ int
	numCols_ int
	matrix_  []columnType
	sm       sync.Mutex
}

type storageKey struct {
	row_, col_ int
}

func (s *stableStorage) getCell( k storageKey ) (*userRecord, error) {

	var prow * columnType
	curNumCols := maxNumOfColumns

	s.sm.Lock()         // otherwise it is not safe to read numRows_ and numCols_
	if k.row_ == s.numRows_ {
		curNumCols = s.numCols_
		if curNumCols > 0 {         // current row already started
			prow = &s.matrix_[k.row_]
		}
	} else if k.row_ < s.numRows_{
		prow = &s.matrix_[k.row_]
	}
	s.sm.Unlock()
	
	if prow == nil {
		return nil, requestError{ "row out of range" }
	}
	
	if k.col_ >= curNumCols {
		return nil, requestError{ "column out of range" }
	}

	// however accessing the row is safe - it never gets relocated
	return &prow[k.col_], nil
}

func (s *stableStorage) addCell() (storageKey, *userRecord) {

	s.sm.Lock()
	k := storageKey{s.numRows_, s.numCols_}
	if k.col_ == 0 {
		s.matrix_ = append( s.matrix_, columnType{} )
	}
	if k.col_ + 1 ==  maxNumOfColumns {
		s.numRows_ += 1
		s.numCols_ = 0
	} else {
		s.numCols_ += 1
	}
	u := &s.matrix_[k.row_][k.col_]
	s.sm.Unlock()

	return k, u
}


var stb stableStorage
var userMap map[int] storageKey		// UID -> storage
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
	userMap = make( map[int] storageKey )
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

	k, ur := stb.addCell()
	userMap[userId] = k
	ur.mu.Lock()
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
	uk, ok := userMap[ uid ]
	cm.Unlock()

	if ( !ok ) {
		return requestError { "UID not found"}
	}

	user, err := stb.getCell( uk )
	if err != nil {
		return requestError { "user record not found - internal error"}
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
	uk, ok := userMap[ uid ]
	cm.Unlock()

	if ( !ok ) {
		err = requestError { "UID not found"}
		return
	}

	user, err := stb.getCell( uk )
	if err != nil {
		err = requestError { "user record not found - internal error"}
		return
	}

	user.mu.Lock()
	defer user.mu.Unlock()

	resp.Name = user.name
	resp.Email = user.email
	resp.Addr = user.addr

	return
}
