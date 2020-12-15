package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Message struct {
	ID     int    `json:"id"`
	Name   string `json:"login"`
	Campus []CampusData
}

type CampusData struct {
	CampusName string `json:"name"`
}

type AuthUser struct {
	ID     string
	Name   string
	Campus string
}

var (
	AuthConfig       *oauth2.Config
	oauthStateString = Hex(16)
	codeVerifier, _  = cv.CreateCodeVerifier()
	codeChallenge    = codeVerifier.CodeChallengeS256()
)

func init() {
	AuthConfig = &oauth2.Config{
		RedirectURL: "http://localhost:8185/callback",

		ClientID:     "c7a7c50ad67f03a72f23c77545b25ac48d616bc1e5daef046d956ed55acf95fd",
		ClientSecret: "157505de170d0b275ab4e10041d4dba1f4f90e21bd1ab5567fc9694b1f040716",
		Scopes:       []string{"public"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://api.intra.42.fr/oauth/authorize",
			TokenURL:  "https://api.intra.42.fr/oauth/token",
			AuthStyle: oauth2.AuthStyleInHeader,
		},
	}
}

func authLogin(w http.ResponseWriter, r *http.Request) {
	url := AuthConfig.AuthCodeURL(oauthStateString, oauth2.SetAuthURLParam("code_challenge", codeChallenge), oauth2.SetAuthURLParam("code_challenge_method", "S256"))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	content, err := authUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	userInf, _ := getUserFromCallback(content)
	userLogin(w, r, userInf)

	//fmt.Fprintf(w,"lol: %s", userInf)
	//fmt.Fprintf(w, "Instr user info: %s", content)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func authUserInfo(state string, code string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := AuthConfig.Exchange(oauth2.NoContext, code, oauth2.SetAuthURLParam("code_verifier", codeVerifier.String()))
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	url := "https://api.intra.42.fr/v2/me"
	var bearer = "Bearer " + token.AccessToken
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, nil
}

func getUserFromCallback(bytes []byte) (*AuthUser, error) {
	var message Message
	err := json.Unmarshal(bytes, &message)
	if err != nil {
		return nil, fmt.Errorf("failed on convert from JSON: %s", err.Error())
	}

	user := AuthUser{
		ID:     strconv.Itoa(message.ID),
		Name:   message.Name,
		Campus: message.Campus[0].CampusName,
	}

	return &user, nil
}

func Bytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func Hex(n int) string {
	return hex.EncodeToString(Bytes(n))
}
