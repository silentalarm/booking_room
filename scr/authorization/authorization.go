package authorization

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
	ses "github.com/silentalarm/booking_room/scr/sessions"
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

var (
	AuthConfig       *oauth2.Config
	oauthStateString = bytesToHex(16)
	codeVerifier, _  = cv.CreateCodeVerifier()
	codeChallenge    = codeVerifier.CodeChallengeS256()
)

func Login(w http.ResponseWriter, r *http.Request) {
	url := AuthConfig.AuthCodeURL(oauthStateString, oauth2.SetAuthURLParam("code_challenge", codeChallenge), oauth2.SetAuthURLParam("code_challenge_method", "S256"))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	content, err := authUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	userInf, _ := getUserFromCallback(content)
	ses.Init(w, r, userInf)

	//fmt.Fprintf(w,"lol: %s", userInf)
	fmt.Fprintf(w, "Instr user info: %s", content)
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

func getUserFromCallback(bytes []byte) (*ses.AuthUser, error) {
	var message Message
	err := json.Unmarshal(bytes, &message)
	if err != nil {
		return nil, fmt.Errorf("failed on convert from JSON: %s", err.Error())
	}

	user := ses.AuthUser{
		ID:     strconv.Itoa(message.ID),
		Name:   message.Name,
		Campus: message.Campus[0].CampusName,
	}

	return &user, nil
}

func generateBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func bytesToHex(n int) string {
	return hex.EncodeToString(generateBytes(n))
}
