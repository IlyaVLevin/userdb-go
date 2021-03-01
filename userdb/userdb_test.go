package userdb

import ( "testing" )

func TestCreatUser(t *testing.T) {
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
		t.Fatalf("unexpected user #%v: %v (should be %v)", uid, usobj, reqCreate )		
	}

	// 
	reqCreate1 := RequestCreateUser{ "JoeSmith", "1234abcd", "jsmith@mmm.com", "2 Brodway"}

	_, err = CreateUser( &reqCreate1 )
	if err == nil || err.Error() != "Name already reserved" {
		t.Fatalf("improper handling of duplicate name")
	}
}