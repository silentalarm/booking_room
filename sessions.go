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
	user, err := getUser(w, r)
	if err != nil || !user.Authenticated {
		http.Error(w, "User not found", http.StatusForbidden)
		return
	}
	fmt.Printf("id: %s name: %s campus: %s auth: %t",
		user.ID, user.Name, user.Campus, user.Authenticated)
	//fmt.Printf("%s", name)
	return
}

func userLogin(w http.ResponseWriter, r *http.Request, user *AuthUser) {
	session, _ := store.Get(r, "auth-session")
	session.Values["id"] = user.ID
	session.Values["name"] = user.Name
	session.Values["campus"] = user.Campus
	session.Options.MaxAge = 60 * 15 //15 minutes life cookie
	session.Save(r, w)
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
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func getUser(w http.ResponseWriter, r *http.Request) (*User, error) {
	session, err := store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, nil
	}
	user := User{}
	user.ID = session.Values["id"].(string)
	user.Name = session.Values["name"].(string)
	user.Campus = session.Values["campus"].(string)

	return &user, nil
}
