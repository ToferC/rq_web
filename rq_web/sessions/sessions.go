package sessions

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/toferc/ore_web_roller/models"
)

//Store the cookie store which is going to store session data in the ChooseHyperSkillDice
var Store = sessions.NewCookieStore([]byte(os.Getenv("CookieSecret")))

//IsLoggedIn will check if the user has an active session and return true
func IsLoggedIn(req *http.Request) bool {
	session, _ := Store.Get(req, "session")
	if session.Values["loggedin"] == "true" {
		return true
	}
	return false
}

//IsAuthor checks if the user matches the Character's author and return true
func IsAuthor(req *http.Request, cr models.CharacterModel) bool {
	session, _ := Store.Get(req, "session")
	if session.Values["username"] == cr.Author.UserName {
		return true
	}
	return false
}
