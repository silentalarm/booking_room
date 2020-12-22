package sessions

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"net/http"
)

type AuthUser struct {
	ID     string
	Name   string
	Campus string
	Staff  bool
}

type User struct {
	ID            string
	Name          string
	Campus        string
	Staff         bool
	Authenticated bool
}

var (
	authKeyOne       = securecookie.GenerateRandomKey(64)
	encryptionKeyOne = securecookie.GenerateRandomKey(32)

	Store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)
)

func ProfileUser(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	if user.Authenticated == false {
		fmt.Printf("user: %s auth: %t", user.Name, user.Authenticated)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}

func Init(w http.ResponseWriter, r *http.Request, user *AuthUser) {
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["id"] = user.ID
	session.Values["name"] = user.Name
	session.Values["campus"] = user.Campus
	session.Values["staff"] = user.Staff
	session.Values["authenticated"] = true
	session.Options.MaxAge = 60 * 15 //15 minutes life cookie

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["id"] = ""
	session.Values["name"] = ""
	session.Values["campus"] = ""
	session.Values["staff"] = false
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func GetUser(session *sessions.Session) *User {
	if _, ok := session.Values["authenticated"].(bool); !ok {
		user := User{Name: "неавторизован", Authenticated: false}
		return &user
	}
	user := User{
		ID:            session.Values["id"].(string),
		Name:          session.Values["name"].(string),
		Campus:        session.Values["campus"].(string),
		Staff:         session.Values["staff"].(bool),
		Authenticated: session.Values["authenticated"].(bool),
	}
	return &user
}

//FIX ME
func IsAuthenticated(w http.ResponseWriter, r *http.Request) (bool, error) {
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false, err
	}
	user := GetUser(session)

	return user.Authenticated, nil
}
