package main

import (
	"fmt"
	"userdb"
	"encoding/json"
)

func main() {
	
	reqCreate := RequestCreateUser{ "JoeSmith", "abcd1234", "jsmith@bbb.com", "1 Broadway"}
	
	uid, err := CreateUser( reqCreate )
	if err != nil {
		fmt.Printf("error in user creation: %v", err.Error() )		
	}
	else 
		fmt.Println("All good")
}