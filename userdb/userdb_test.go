package userdb

import (
	"testing"
	"sync"
	"math/rand"
	"strconv"
)

func TestCreateUser(t *testing.T) {
	ResetDb()
	
	reqEmpty := RequestCreateUser{ "", "abcd1234", "jsmith@bbb.com", "1 Broadway"}
	_, err := CreateUser( &reqEmpty )
	if err == nil || err.Error() != "Empty name" {
		t.Fatalf("improper handling of empty name")
	}

	reqCreate := RequestCreateUser{ "JoeSmith", "abcd1234", "jsmith@bbb.com", "1 Broadway"}

	uid, err := CreateUser( &reqCreate )
	if err != nil {
		t.Fatalf("error in user creation: %v", err.Error() )		
	}
	
	usobj, err := GetUser( uid)
	if err != nil {
		t.Fatalf("Error obtaining newly created user #%v: %v", uid, err.Error() )		
	}
	
	if usobj.Name != reqCreate.Name || usobj.Email != reqCreate.Email || usobj.Addr != reqCreate.Addr {
		t.Fatalf("unexpected user #%v: %+v (should be %+v)", uid, usobj, reqCreate )		
	}

	// trying creation a user with the same name
	reqCreate1 := RequestCreateUser{ "JoeSmith", "1234abcd", "jsmith@mmm.com", "2 Brodway"}

	_, err = CreateUser( &reqCreate1 )
	if err == nil || err.Error() != "Name already reserved" {
		t.Fatalf("improper handling of duplicate name")
	}
}

func TestUpdateUser(t *testing.T) {
	ResetDb()

	reqCreate := RequestCreateUser{ "JoeSmith", "abcd1234", "jsmith@bbb.com", "1 Broadway"}

	uid, err := CreateUser( &reqCreate )
	if err != nil {
		t.Fatalf("error in user creation: %v", err.Error() )
	}

	reqUpdates := []RequestUpdateUser{ { Addr : "2 Broadway" }, { Email: "smithj@nnn.com", Addr : "3 Broadway" }, {} }
	expResp1 := ResponseGetUser{Name: "JoeSmith", Email: "jsmith@bbb.com", Addr: "2 Broadway" }
	expResp2 := ResponseGetUser{ Name: "JoeSmith", Email: "smithj@nnn.com", Addr: "3 Broadway" }
	expResponses := []ResponseGetUser{ expResp1, expResp2, expResp2 }

	for i, rui := range reqUpdates {
		err = UpdateUser( uid, &rui )
		if err != nil {
			t.Fatalf("Error updating user record #%v as %+v: %v", uid, rui, err.Error() )
		}
		usobj, err := GetUser( uid )
		if err != nil {
			t.Fatalf("Error obtaining user record for #%v: %v", uid, err.Error() )		
		}
		if usobj != expResponses[i] {
			t.Fatalf("User update %+v for user record #%v resulted in %+v instead of %+v", rui, uid, usobj, expResponses[i] )
		}
	}
}

var rndMut sync.Mutex

func randNum(n int) int {
	rndMut.Lock()
	defer rndMut.Unlock()

	return rand.Intn( n )
}

var namMut sync.Mutex
var namNum uint64 = 987654321

func generateNextName() string {
	namMut.Lock()
	nn := namNum
	namNum += 2
	namMut.Unlock()

	return strconv.FormatUint( nn, 36 )
}


func TestStressTest(t *testing.T) {
	ResetDb()

	const QUERY_BUF_SIZE = 1000      // number of uids to query and update -- the smaller the number the higher probability of collision
	const GETTERS_NUM = 10000        // number of querying threads  -- the higher the more probable is the collision
	const UPDATERS_NUM = 10000		 // number of updating threads  -- the higher the more probable is the collision
	const CREATORS_NUM = 10000		 // number of threads continuing to create new users in parallel with getters in updaters
	const NUM_QUERIES = 1000			// each getter issues this many queries
	const NUM_UPDATES = 1000			// each updater makes this many updates
	const NUM_CREATES = 1000			// each creator creates this many records

	var queryBuf [QUERY_BUF_SIZE]int

	reqCreate := RequestCreateUser{ "", "abcd1234", "jsmith@bbb.com", "1 Broadway"}

	for i, _ := range queryBuf {
		reqCreate.Name = generateNextName()
		uid, err := CreateUser( &reqCreate )
		if err != nil {
			t.Fatalf("error in user creation: %v", err.Error() )
		}
		queryBuf[i] = uid
	}

	var wg sync.WaitGroup

	for i:=0; i<CREATORS_NUM; i++ {
		wg.Add(1)
		go func() {
			for j:=0; j < NUM_CREATES; j++ {
				rqc := reqCreate
				rqc.Name = generateNextName()
				_, err := CreateUser( &rqc )
				if err != nil {
					t.Fatalf("error in user %v creation: %v", rqc.Name, err.Error() )
				}
			}
			wg.Done()
		}()
	}

	rand.Seed( 1 )     // or time.Now()?
	reqUpdate := RequestUpdateUser{ Email: "smithj@nnn.com", Addr : "3 Broadway" }

	for i:=0; i<UPDATERS_NUM; i++ {
		wg.Add(1)
		go func() {
			for j:=0; j < NUM_UPDATES; j++ {
				ind := randNum( QUERY_BUF_SIZE )
				err := UpdateUser( queryBuf[ind], &reqUpdate )     // same update every time
				if err != nil {
					t.Fatalf("error in update #%v: %v", queryBuf[ind], err.Error() )
				}
			}
			wg.Done()
		}()
	}

	for i:=0; i<GETTERS_NUM; i++ {
		wg.Add(1)
		go func() {
			for j:=0; j < NUM_QUERIES; j++ {
				ind := randNum( QUERY_BUF_SIZE )
				_, err := GetUser( queryBuf[ind] )
				if err != nil {
					t.Fatalf("error in querying #%v: %v", queryBuf[ind], err.Error() )
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
