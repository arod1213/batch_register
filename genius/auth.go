// https://api.genius.com/oauth/authorize?
// client_id=YOUR_CLIENT_ID&
// redirect_uri=YOUR_REDIRECT_URI&
// scope=REQUESTED_SCOPE&
// state=SOME_STATE_VALUE&
// response_type=code
package genius

// import (
// 	"encoding/json"
// 	"log"
// 	"net/url"
// 	"os"

// 	"github.com/arod1213/auto_ingestion/utils"
// )

// type Auth struct {
// 	Code         string `json:"code"`
// 	ClientID     string `json:"client_id"`
// 	ClientSecret string `json:"client_secret"`
// 	RedirectURI  string `json:"redirect_uri"`
// 	ResponseType string `json:"response_type"`
// 	GrantType    string `json:"grant_type"`
// }

// func GetAuth() (*Auth, error) {
// 	params := url.Values{}
// 	params.Set("client_id", os.Getenv("GENIUS_CLIENT_ID"))
// 	params.Set("redirect_uri", "http://genius.com/callback")
// 	params.Set("scope", "me")
// 	// params.Set("state", "SOME_STATE_VALUE")
// 	params.Set("response_type", "code")
// 	body, err := utils.FetchBodyNoAuth("https://api.genius.com/oauth/authorize", "GET", params)
// 	if err != nil {
// 		log.Println("error fetching body: ", err)
// 		return nil, err
// 	}
// 	var auth Auth
// 	err = json.Unmarshal(body, &auth)
// 	if err != nil {
// 		log.Println("error unmarshalling body: ", err)
// 		return nil, err
// 	}
// 	return &auth, nil
// }
