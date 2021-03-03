package udbserv

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
//	"fmt"
	"io"
)

type UsrvHandler struct {
	t *testing.T
}

func (h* UsrvHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if  strings.HasPrefix(req.URL.Path, "/update") {
		handleUpdateUser(w, req)
	} else if strings.HasPrefix(req.URL.Path, "/get") {
		handleGetUser(w, req)
	} else if strings.HasPrefix(req.URL.Path, "/create") {
		handleCreateUser(w, req)
	} else {
		h.t.Fatalf("Unkonwn route: %v", req.URL.Path)
	}
}

func readAndPrintResponse( t *testing.T, resp *http.Response, err error, url string, expReply string) {
	if err != nil {
		t.Fatalf( err.Error() )
	}

	cont, err := io.ReadAll( resp.Body )
	resp.Body.Close()
	
	if err != nil {
		t.Fatalf( err.Error() )
	}

	realReply := strings.TrimSpace( string( cont ) )

	if realReply != expReply {
		t.Fatalf( "Reply '%s' to '%s' differs from expected '%s'", realReply, url, expReply )
	}
//	fmt.Println( realReply )
}

func checkGet( t *testing.T, url string, expReply string ) {
	resp, err := http.Get( url )
	readAndPrintResponse( t, resp, err, url, expReply)
}

func checkPost( t* testing.T, url string, body string, expReply string ) {
	resp, err := http.Post( url, "application/json", strings.NewReader( body ) );
	readAndPrintResponse( t, resp, err, url, expReply )
}

func TestGeneralTest(t *testing.T) {
	server := httptest.NewServer( &UsrvHandler{ t }  )
	defer server.Close()

	// Get for non-existing uid
	getUrl := server.URL + "/get/1000"
	expReply := "{\"Error\":\"UID not found\"}"
	checkGet( t, getUrl, expReply )

	// Normal Create
	createUrl := server.URL + "/create"
	createBody := "{\"name\":\"ilya\", \"passwd\":\"asdfgh\", \"email\":\"Bububu@hhh\"}"
	expReply = "{\"Uid\":1000}"
	checkPost( t, createUrl, createBody, expReply )

	// Normal Get
	expReply = "{\"Name\":\"ilya\",\"Email\":\"Bububu@hhh\",\"Addr\":\"\"}"
	checkGet( t, getUrl, expReply )

	// Create for existing username
	expReply = "{\"Uid\":-1}\n{\"Error\":\"Name already reserved\"}"
	checkPost( t, createUrl, createBody, expReply )

	// Malformed Create
	createBody = "{\"passwd\":\"asdfgh\", \"email\":\"Bububu@hhh\"}"
	expReply = "{\"Uid\":-1}\n{\"Error\":\"Empty name\"}"
	checkPost( t, createUrl, createBody, expReply )

	// Normal Update
	updateUrl := server.URL + "/update/1000"
	updateBody := "{\"addr\":\"123 Broadway\"}"
	expReply = "{\"Status\":\"OK\"}"
	checkPost( t, updateUrl, updateBody, expReply )

	// Update for non-existing uid
	updateUrlBad := server.URL + "/update/99"
	expReply = "{\"Status\":\"ERROR: UID not found\"}"
	checkPost( t, updateUrlBad, updateBody, expReply )

	// Update for malformed uid (get by username currently not supported)
	updateUrlBad = server.URL + "/update/ilya"
	expReply = "{\"Status\":\"ERROR: Non-numerical UID\"}"
	checkPost( t, updateUrlBad, updateBody, expReply )

	// Malformed Update request - unknown field
	updateBodyBad := "{\"address\":\"123 Broadway\"}"
	expReply = "{\"Status\":\"ERROR: Malformed request: json: unknown field \\\"address\\\"\"}"
	checkPost( t, updateUrl, updateBodyBad, expReply )

	// Normal Get for updated record
	expReply = "{\"Name\":\"ilya\",\"Email\":\"Bububu@hhh\",\"Addr\":\"123 Broadway\"}"
	checkGet( t, getUrl, expReply )
}
