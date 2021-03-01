package main

import (
	"fmt"
	"userdb"
//	"encoding/json"
)

func main() {

	reqCreate := userdb.RequestCreateUser{ "JoeSmith", "abcd1234", "jsmith@bbb.com", "1 Broadway"}

	uid, err := userdb.CreateUser( &reqCreate )
	if err != nil {
		fmt.Println("error in user creation: ", err.Error() )		
	} else {
		fmt.Println("All good: ", uid)
	}
}