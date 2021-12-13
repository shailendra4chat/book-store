package helpers

import (
	"net/http"

	"github.com/shailendra4chat/book-store/books/config"
)

var uri string = "http://" + config.Conf("UAPP_HOST") + ":" + config.Conf("UAPP_PORT") + "/auth/token"
var client = &http.Client{}

// Check If user is Admin
func IsAdmin(userToken string) bool {

	url := uri + "/admin"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-access-token", userToken)

	res, _ := client.Do(req)

	return res.StatusCode == 200
}

// Check If user is authorised
func IsAuthorisedUser(userToken string) bool {

	url := uri
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-access-token", userToken)

	res, _ := client.Do(req)

	return res.StatusCode == 200
}
