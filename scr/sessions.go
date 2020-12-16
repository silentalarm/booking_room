package main

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"net/http"
)

type User struct {
	ID            string
	Name          string
	Campus        string
	Authenticated bool
}

var (
	authKeyOne       = securecookie.GenerateRandomKey(64)
	encryptionKeyOne = securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)
)

func profileUser(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := getUser(session)

	if user.Authenticated == false {
		fmt.Printf("user: %s auth: %t", user.Name, user.Authenticated)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fmt.Fprintf(w, "id: %s name: %s campus: %s auth: %t",
		user.ID, user.Name, user.Campus, user.Authenticated)
}

func userLogin(w http.ResponseWriter, r *http.Request, user *AuthUser) {
	session, err := store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["id"] = user.ID
	session.Values["name"] = user.Name
	session.Values["campus"] = user.Campus
	session.Values["authenticated"] = true
	session.Options.MaxAge = 60 * 15 //15 minutes life cookie

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func userLogout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["id"] = ""
	session.Values["name"] = ""
	session.Values["campus"] = ""
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func getUser(session *sessions.Session) *User {
	if _, ok := session.Values["authenticated"].(bool); !ok {
		user := User{Name: "неавторизован", Authenticated: false}
		return &user
	}
	user := User{
		ID:            session.Values["id"].(string),
		Name:          session.Values["name"].(string),
		Campus:        session.Values["campus"].(string),
		Authenticated: session.Values["authenticated"].(bool),
	}
	return &user
}

//FIX ME
func isAuthenticated(w http.ResponseWriter, r *http.Request) (bool, error) {
	session, err := store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false, err
	}
	user := getUser(session)

	return user.Authenticated, nil
}
