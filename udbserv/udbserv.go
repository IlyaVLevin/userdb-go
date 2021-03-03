package udbserv

import (
	"userdb"
	"encoding/json"
	"net/http"
	"io"
	"strings"
	"strconv"
)

type UidResponse struct {
	Uid int
}

type ErrorResponse struct {
	Error string
}

type StatusResponse struct {
	Status string
}

func handleCreateUser(w http.ResponseWriter, req *http.Request) {

	dec := json.NewDecoder( req.Body )
	var uids []UidResponse
	var expl string

	defer req.Body.Close()

	dec.DisallowUnknownFields()

	for {
		var cr userdb.RequestCreateUser
		if err := dec.Decode( &cr ) ; err == io.EOF {
			break
		} else if err != nil {
			// signal format error and stop parsing
			uids = append( uids, UidResponse{ -1 })
			expl = err.Error()
			break
		}

		uid, err := userdb.CreateUser( &cr )
		if err != nil {
			// signal format error and stop parsing
			uids = append( uids, UidResponse{ -1 })
			expl = err.Error()
			break
		} else {
			uids = append( uids, UidResponse{ uid })
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
		enc.Encode( ErrorResponse{ expl } )
	}
}

func handleUpdateUser(w http.ResponseWriter, req *http.Request){

	arg := strings.TrimPrefix( req.URL.Path, "/update/" )
	if arg == "" {
		helpUpdateUser( w, req )
		return
	}

	var expl string

	enc := json.NewEncoder( w )

	defer req.Body.Close()

	uid, err := strconv.Atoi( arg )
	if err != nil {
		expl = "Non-numerical UID"
	} else {

		dec := json.NewDecoder( req.Body )
		dec.DisallowUnknownFields()

		var ur userdb.RequestUpdateUser

		// no batch mode for updates
		errDeco := dec.Decode( &ur )
		if errDeco == nil {
			errDeco = dec.Decode( &ur )
			if errDeco == nil {
				expl = "Multiple update requests"
			} else if errDeco != io.EOF {
				expl = "Malformed request: " + errDeco.Error()
			}
		} else {
			expl = "Malformed request: " + errDeco.Error()
		}

		if expl == "" {
			err = userdb.UpdateUser( uid, &ur )
			if err != nil {
				expl = err.Error()
			} else {
				enc.Encode( StatusResponse{ "OK" } )
			}
		}
	}

	if expl != "" {
		enc.Encode( StatusResponse{ "ERROR: " + expl } )
	}
}

func handleGetUser(w http.ResponseWriter, req *http.Request) {
	arg := strings.TrimPrefix( req.URL.Path, "/get/" )
	if arg == "" {
		helpGetUser( w, req )
		return
	}

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
		enc.Encode( ErrorResponse{ expl } )
	}
}

func helpUpdateUser(w http.ResponseWriter, _ *http.Request) {
	expl := "Call:  /update/<UID>   Data: 'email:<email>' and/or 'addr: <address>'"
	enc := json.NewEncoder( w )
	enc.Encode( ErrorResponse{ expl } )
}

func helpGetUser(w http.ResponseWriter, _ *http.Request) {
	expl := "call:  /get/<UID>"
	enc := json.NewEncoder( w )
	enc.Encode( ErrorResponse{ expl } )
}

func SetHttpHandlerFuncs() {

	http.HandleFunc("/create/", handleCreateUser)
	http.HandleFunc("/create", handleCreateUser)

	http.HandleFunc("/update/", handleUpdateUser)
	http.HandleFunc("/get/", handleGetUser)

	http.HandleFunc("/update", helpUpdateUser)
	http.HandleFunc("/get", helpGetUser)
}