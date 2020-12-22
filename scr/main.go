package main

import (
	"fmt"
	auth "github.com/silentalarm/booking_room/scr/authorization"
	bk "github.com/silentalarm/booking_room/scr/booking"
	pg "github.com/silentalarm/booking_room/scr/clubs"
	ses "github.com/silentalarm/booking_room/scr/sessions"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

func init() {
	auth.AuthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		ClientID:     os.Getenv("INTRA_CLIENT_ID"),
		ClientSecret: os.Getenv("INTRA_CLIENT_SECRET"),
		Scopes:       []string{"public"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://api.intra.42.fr/oauth/authorize",
			TokenURL:  "https://api.intra.42.fr/oauth/token",
			AuthStyle: oauth2.AuthStyleInHeader,
		},
	}
}

func main() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/webfonts/", http.StripPrefix("/webfonts/", http.FileServer(http.Dir("webfonts"))))

	port := os.Getenv("PORT")
	//port := "8185"

	http.HandleFunc("/", bk.Index)
	http.HandleFunc("/login", auth.Login)
	http.HandleFunc("/callback", auth.CallbackHandler)
	http.HandleFunc("/logout", ses.Delete)
	http.HandleFunc("/savereserve", bk.SaveReserve)
	http.HandleFunc("/delreserve", bk.DeleteReserveFromUser)
	http.HandleFunc("/clubregistration", pg.InsertNewClub)
	http.HandleFunc("/clubs", pg.ClubsTable)
	http.HandleFunc("/clubstoapproved", pg.ClubsToApprovedTable)
	//http.HandleFunc("/saveclub", pg.InsertNewClub)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":"+port, nil)

}

//func main()  {
//	router := mux.NewRouter()
//	router.PathPrefix("/").Handler(http.FileServer(rice.MustFindBox("website").HTTPBox()))
//	http.ListenAndServe(":8082", router)
//}
