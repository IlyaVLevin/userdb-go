package userdb

import ( "testing" )

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

	for i := 0; i < len( reqUpdates ); i++ {
		err = UpdateUser( uid, &reqUpdates[i] )
		if err != nil {
			t.Fatalf("Error updating user record #%v as %+v: %v", uid, reqUpdates[i], err.Error() )
		}
		usobj, err := GetUser( uid )
		if err != nil {
			t.Fatalf("Error obtaining user record for #%v: %v", uid, err.Error() )		
		}
		if usobj != expResponses[i] {
			t.Fatalf("User update %+v for user record #%v resulted in %+v instead of %+v", reqUpdates[i], uid, usobj, expResponses[i] )		
		}
	}
}