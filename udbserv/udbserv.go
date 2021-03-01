package main

import (
	"userdb"
	"encoding/json"
	"net/http"
	"io"
	"log"
	"strings"
	"strconv"
)

func handleCreateUser(w http.ResponseWriter, req *http.Request) {

	dec := json.NewDecoder( req.Body )
	var uids []int
	var expl string

	defer req.Body.Close()

	for {
		var cr userdb.RequestCreateUser
		if err := dec.Decode( &cr ) ; err == io.EOF {
			break
		} else if err != nil {
			// signal format error and stop parsing
			uids = append( uids, -1)
			expl = err.Error()
			break
		}

		uid, err := userdb.CreateUser( &cr )
		if err != nil {
			// signal format error and stop parsing
			uids = append( uids, -1)
			expl = err.Error()
			break
		} else {
			uids = append( uids, uid)
		}
	}

	enc := json.NewEncoder( w )

	if ( len( uids ) == 1 ) {
		enc.Encode( uids[0] )
	} else if uids == nil {
		expl = "Empty input"
	} else {
		enc.Encode( uids ) 
	}
	if ( expl != "" ) {
		enc.Encode( expl ) 
	}
}

func handleUpdateUser(w http.ResponseWriter, req *http.Request){
}

func handleGetUser(w http.ResponseWriter, req *http.Request) {
	arg := strings.TrimPrefix( req.URL.Path, "/get/" )

	var expl string

	enc := json.NewEncoder( w )

	uid, err := strconv.Atoi( arg )
	if err != nil {
		expl = "Non-numerical UID"
	} else {
		grsp, err := userdb.GetUser( uid )
		if err != nil {
			expl = err.Error()
		} else {
			enc.Encode( grsp )
		}
	}

	if expl != "" {
		enc.Encode( expl )
	}
}

func main() {

	// create a fictitous user for testing purposes
	reqCreate := userdb.RequestCreateUser{ "JoeSmith", "1234abcd", "jsmith@mmm.com", "2 Brodway"}

	_, err := userdb.CreateUser( &reqCreate )
	if err != nil  {
		log.Fatal( err.Error() )
	}

	http.HandleFunc("/create/", handleCreateUser)
	http.HandleFunc("/update/", handleUpdateUser)
	http.HandleFunc("/get/", handleGetUser)

	log.Fatal( http.ListenAndServe(":8080", nil) )
}
