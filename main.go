package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gorilla/mux"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/google"
	"golang.org/x/oauth2"
	googleOAuth2 "golang.org/x/oauth2/google"

	"github.com/joho/godotenv"

	"github.com/go-pg/pg"
	"github.com/toferc/rq_web/database"
)

var (
	db       *pg.DB
	svc      *s3.S3
	uploader *s3manager.Uploader
)

// MaxMemory is the max upload size for images
const MaxMemory = 2 * 1024 * 1024

// DefaultCharacterPortrait is a base image used as a default
const DefaultCharacterPortrait = "/media/shadow.jpeg"

func init() {
	os.Setenv("DBUser", "chris")
	os.Setenv("DBPass", "12345")
	os.Setenv("DBName", "runequest")
}

func main() {

	var callback string

	if os.Getenv("ENVIRONMENT") == "production" {
		// Production system
		url, ok := os.LookupEnv("DATABASE_URL")

		if !ok {
			log.Fatalln("$DATABASE_URL is required")
		}

		options, err := pg.ParseURL(url)

		if err != nil {
			log.Fatalf("Connection error: %s", err.Error())
		}

		db = pg.Connect(options)

		// Set Google Oauth Callback
		callback = "https://rq-web.herokuapp.com/google/callback"

	} else {
		// Not production
		db = pg.Connect(&pg.Options{
			User:     os.Getenv("DBUser"),
			Password: os.Getenv("DBPass"),
			Database: os.Getenv("DBName"),
		})
		os.Setenv("CookieSecret", "kimchee-typhoon")
		os.Setenv("BUCKET", "runequeset")
		os.Setenv("AWS_REGION", "us-east-1")

		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		configGoogleOAUTH()

		callback = "http://localhost:8080/google/callback"
	}

	defer db.Close()

	err := database.InitDB(db)
	if err != nil {
		panic(err)
	}

	fmt.Println(db)

	// Create AWS session using local default config
	// or Env Variables if on Heroku
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		log.Fatal(err)
		fmt.Println("Error creating session", err)
		os.Exit(1)
	}

	fmt.Println("Session created ", sess)
	svc := s3.New(sess)

	fmt.Println(svc)

	uploader = s3manager.NewUploader(sess)

	// Config Google Oauth
	config := &Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	}

	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint:     googleOAuth2.Endpoint,
		RedirectURL:  callback,
		Scopes:       []string{"profile", "email"},
	}
	stateConfig := gologin.DebugOnlyCookieConfig

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()

	fmt.Println("Starting Webserver at port " + port)
	r.HandleFunc("/", CharacterIndexHandler)
	r.HandleFunc("/about/", AboutHandler)
	r.HandleFunc("/user_roster/", UserCharacterRosterHandler)
	r.HandleFunc("/add_to_user_roster/{id}", AddToUserRosterHandler)

	r.HandleFunc("/signup/", SignUpFunc)
	r.HandleFunc("/login/", LoginFunc)
	r.HandleFunc("/logout/", LogoutFunc)

	r.Handle("/google/login", google.StateHandler(stateConfig, google.LoginHandler(oauth2Config, nil)))
	r.Handle("/google/callback", google.StateHandler(stateConfig, google.CallbackHandler(oauth2Config, googleLoginFunc(), nil)))

	r.HandleFunc("/users/", UserIndexHandler)

	r.HandleFunc("/view_character/{id}", CharacterHandler)
	r.HandleFunc("/new/", NewCharacterHandler)
	r.HandleFunc("/modify/{id}", ModifyCharacterHandler)
	r.HandleFunc("/delete/{id}", DeleteCharacterHandler)

	// Character Creation Handlers
	r.HandleFunc("/cc1_choose_homeland/", ChooseHomelandHandler)
	r.HandleFunc("/cc12_personal_history/{id}", PersonalHistoryHandler)
	r.HandleFunc("/cc2_choose_runes/{id}", ChooseRunesHandler)
	r.HandleFunc("/cc3_roll_stats/{id}", RollStatisticsHandler)
	r.HandleFunc("/cc4_apply_homeland/{id}", ApplyHomelandHandler)
	r.HandleFunc("/cc5_apply_occupation/{id}", ApplyOccupationHandler)
	r.HandleFunc("/cc6_apply_cult/{id}", ApplyCultHandler)
	r.HandleFunc("/cc7_personal_skills/{id}", PersonalSkillsHandler)
	r.HandleFunc("/cc8_finishing_touches/{id}", FinishingTouchesHandler)

	// Homeland Handlers
	r.HandleFunc("/homeland_index/", HomelandIndexHandler)
	r.HandleFunc("/add_homeland/", AddHomelandHandler)
	r.HandleFunc("/view_homeland/{id}", HomelandHandler)
	r.HandleFunc("/modify_homeland/{id}", ModifyHomelandHandler)
	r.HandleFunc("/delete_homeland/{id}", DeleteHomelandHandler)

	// Occupation Handlers
	r.HandleFunc("/occupation_index/", OccupationIndexHandler)
	r.HandleFunc("/add_occupation/", AddOccupationHandler)
	r.HandleFunc("/view_occupation/{id}", OccupationHandler)
	r.HandleFunc("/modify_occupation/{id}", ModifyOccupationHandler)
	r.HandleFunc("/delete_occupation/{id}", DeleteOccupationHandler)

	// Cult Handlers
	r.HandleFunc("/cult_index/", CultIndexHandler)
	r.HandleFunc("/add_cult/", AddCultHandler)
	r.HandleFunc("/view_cult/{id}", CultHandler)
	r.HandleFunc("/modify_cult/{id}", ModifyCultHandler)
	r.HandleFunc("/delete_cult/{id}", DeleteCultHandler)

	r.HandleFunc("/add_character_content/{id}", AddCharacterContentHandler)
	r.HandleFunc("/equip_weapons_armor/{id}", EquipWeaponsArmorHandler)
	//r.HandleFunc("/add_advantages/{id}", ModifyAdvantageHandler)

	r.HandleFunc("/user_index/", UserIndexHandler)

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":"+port, r))

}
