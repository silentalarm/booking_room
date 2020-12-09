package main

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"net/http"
)

var (
	authKeyOne       = securecookie.GenerateRandomKey(64)
	encryptionKeyOne = securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)
)

func profileUser(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth-session")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	name := session.Values["name"]
	fmt.Printf("%s", name)
}

func userLogin(w http.ResponseWriter, r *http.Request, user *User) {
	session, _ := store.Get(r, "auth-session")
	session.Values["authenticated"] = true
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

	session.Values["authenticated"] = false
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
