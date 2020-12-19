package main

import (
	"fmt"
	auth "github.com/silentalarm/booking_room/scr/authorization"
	bk "github.com/silentalarm/booking_room/scr/booking"
	pg "github.com/silentalarm/booking_room/scr/page"
	ses "github.com/silentalarm/booking_room/scr/sessions"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

func init() {
	auth.AuthConfig = &oauth2.Config{
		RedirectURL:  "https://booking21.herokuapp.com/callback",
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

func main() {
	fs := http.FileServer(http.Dir("static/clubRegistration.html"))

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/webfonts/", http.StripPrefix("/webfonts/", http.FileServer(http.Dir("webfonts"))))

	port := os.Getenv("PORT")
	//port := "8185"

	http.HandleFunc("/", pg.Index)
	http.HandleFunc("/login", auth.Login)
	http.HandleFunc("/callback", auth.CallbackHandler)
	http.HandleFunc("/logout", ses.Delete)
	http.HandleFunc("/savereserve", bk.SaveReserve)
	http.HandleFunc("/delreserve", bk.DeleteReserveFromUser)
	http.Handle("/clubregistration", fs)
	http.HandleFunc("/saveclub", pg.InsertNewClub)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":"+port, nil)

}

//func main()  {
//	router := mux.NewRouter()
//	router.PathPrefix("/").Handler(http.FileServer(rice.MustFindBox("website").HTTPBox()))
//	http.ListenAndServe(":8082", router)
//}
