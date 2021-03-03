package main

import (
	"udbserv"
	"userdb"
	"net/http"
	"log"
)

func main() {

	// create a fictitous user for testing purposes
	reqCreate := userdb.RequestCreateUser{ "JoeSmith", "1234abcd", "jsmith@mmm.com", "2 Brodway"}

	_, err := userdb.CreateUser( &reqCreate )
	if err != nil  {
		log.Fatal( err.Error() )
	}

	udbserv.SetHttpHandlerFuncs()

	log.Fatal( http.ListenAndServe(":8080", nil) )
}
