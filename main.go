package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/schema"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gorilla/mux"

	gologin "github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/google"
	"golang.org/x/oauth2"
	googleOAuth2 "golang.org/x/oauth2/google"

	"github.com/joho/godotenv"

	"github.com/go-pg/pg/v10"
	"github.com/toferc/rq_web/database"
)

var (
	db         *pg.DB
	svc        *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	decoder    = schema.NewDecoder()
	s          = rand.NewSource(time.Now().Unix())
	randomSeed = rand.New(s)
)

// MaxMemory is the max upload size for images
const MaxMemory = 2 * 1024 * 1024

// DefaultCharacterPortrait is a base image used as a default
const DefaultCharacterPortrait = "/media/shadow.jpeg"

// DefaultCreaturePortrait is a base image used as a default
const DefaultCreaturePortrait = "/media/shadow.jpeg"

func init() {
	os.Setenv("DBUser", "chris")
	os.Setenv("DBPass", "12345")
	os.Setenv("DBName", "runequest")

	// Set markov chain for name generator
	generatedNameHomelands := []string{
		"sartar",
		"balazaring",
		"lunar",
		"grazelands",
	}

	// Generate model chains if needed
	for _, h := range generatedNameHomelands {

		if _, err := os.Stat(h + "MaleModel.json"); err == nil {
			continue
		} else {
			sc := buildModel(3, h+"_male_names.txt")
			saveModel(sc, h+"MaleModel.json")
		}

		if _, err := os.Stat(h + "FemaleModel.json"); err == nil {
			continue
		} else {
			sc := buildModel(3, h+"_female_names.txt")
			saveModel(sc, h+"FemaleModel.json")
		}
	}
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
		callback = "https://www.cradleofheroes.net/google/callback"

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

	database.CreateTSVColumn(db)
	database.CreateIndex(db)

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

	// Set Schema ignoreunknownkeys to true
	decoder.IgnoreUnknownKeys(true)

	// Miscellaneous actions
	// ResizeImages(db)

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notFound)

	fmt.Println("Starting Webserver at port " + port)

	r.HandleFunc("/about/", AboutHandler)

	r.HandleFunc("/add_to_user_roster/{id}", AddToUserRosterHandler)
	r.HandleFunc("/duplicate_character/{id}", DuplicateCharacterHandler)

	r.HandleFunc("/signup/", SignUpFunc)
	r.HandleFunc("/login/", LoginFunc)
	r.HandleFunc("/logout/", LogoutFunc)

	r.Handle("/google/login", google.StateHandler(stateConfig, google.LoginHandler(oauth2Config, nil)))
	r.Handle("/google/callback", google.StateHandler(stateConfig, google.CallbackHandler(oauth2Config, googleLoginFunc(), nil)))

	// Admin Views
	r.HandleFunc("/admin_view_users/{limit}/{offset}", UserIndexHandler)
	r.HandleFunc("/admin_view_users/", UserIndexHandler)
	r.HandleFunc("/admin_view_user_roster/{id}/{limit}/{offset}", AdminUserRosterViewHandler)
	r.HandleFunc("/admin_view_user_roster/{id}", AdminUserRosterViewHandler)

	// Character Views
	r.HandleFunc("/view_character/{id}", CharacterHandler)
	r.HandleFunc("/text_summary/{id}", TextSummaryHandler)
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

	// Note Handlers
	r.HandleFunc("/notes_index/{id}", NoteIndexHandler)
	r.HandleFunc("/view_note/{slug}", NoteHandler)
	r.HandleFunc("/add_note/{id}", AddNoteHandler)
	r.HandleFunc("/modify_note/{slug}", ModifyNoteHandler)
	r.HandleFunc("/delete_note/{id}", DeleteNoteHandler)

	// Random creation
	r.HandleFunc("/random_character/", RandomCharacterHandler)
	r.HandleFunc("/random_faction/", RandomFactionHandler)

	// Character Modification
	r.HandleFunc("/edit_magic/{id}", EditMagicHandler)
	r.HandleFunc("/bound_spirit/{id}", BoundSpiritHandler)

	// Homeland Handlers
	r.HandleFunc("/homeland_index/", HomelandIndexHandler)
	r.HandleFunc("/add_homeland/", AddHomelandHandler)
	r.HandleFunc("/view_homeland/{slug}", HomelandHandler)
	r.HandleFunc("/modify_homeland/{id}", ModifyHomelandHandler)
	r.HandleFunc("/delete_homeland/{id}", DeleteHomelandHandler)

	// Occupation Handlers
	r.HandleFunc("/occupation_index/", OccupationIndexHandler)
	r.HandleFunc("/add_occupation/", AddOccupationHandler)
	r.HandleFunc("/view_occupation/{slug}", OccupationHandler)
	r.HandleFunc("/modify_occupation/{id}", ModifyOccupationHandler)
	r.HandleFunc("/delete_occupation/{id}", DeleteOccupationHandler)

	// Cult Handlers
	r.HandleFunc("/cult_index/", CultIndexHandler)
	r.HandleFunc("/add_cult/", AddCultHandler)
	r.HandleFunc("/view_cult/{slug}", CultHandler)
	r.HandleFunc("/modify_cult/{id}", ModifyCultHandler)
	r.HandleFunc("/delete_cult/{id}", DeleteCultHandler)

	r.HandleFunc("/add_skills/{id}", AddSkillsHandler)
	r.HandleFunc("/add_passions/{id}", AddPassionsHandler)
	r.HandleFunc("/equip_weapons_armor/{id}", EquipWeaponsArmorHandler)

	// Faction Handlers
	r.HandleFunc("/faction_index/", FactionIndexHandler)
	r.HandleFunc("/user_faction_index/", UserFactionIndexHandler)
	r.HandleFunc("/add_faction/", AddFactionHandler)
	r.HandleFunc("/view_faction/{slug}", FactionHandler)
	r.HandleFunc("/modify_faction/{slug}", ModifyFactionHandler)
	r.HandleFunc("/delete_faction/{id}", DeleteFactionHandler)

	// Encounter Handlers
	r.HandleFunc("/encounter_index/", EncounterIndexHandler)
	r.HandleFunc("/user_encounter_index/", UserEncounterIndexHandler)
	r.HandleFunc("/add_encounter/", AddEncounterHandler)
	r.HandleFunc("/view_encounter/{slug}", EncounterHandler)
	r.HandleFunc("/modify_encounter/{slug}", ModifyEncounterHandler)
	r.HandleFunc("/delete_encounter/{slug}", DeleteEncounterHandler)

	r.HandleFunc("/add_creature/", NewCreatureHandler)
	r.HandleFunc("/modify_creature/{id}", ModifyCreatureHandler)

	// Admin handlers
	r.HandleFunc("/user_index/", UserIndexHandler)
	r.HandleFunc("/make_admin/{id}", MakeAdminHandler)
	r.HandleFunc("/delete_user/{id}", DeleteUserHandler)

	// API
	r.HandleFunc("/api/character", GetCraftedCharacterModels).Methods("GET")
	r.HandleFunc("/api/character/user/{id}", GetUserCharacterModels).Methods("GET")
	r.HandleFunc("/api/character/{id}", GetCharacterModel).Methods("GET")
	r.HandleFunc("/api/character/{id}", CreateCharacterModel).Methods("POST")
	r.HandleFunc("/api/character/{id}", DeleteCharacterModel).Methods("DELETE")
	r.HandleFunc("/api/character/{id}", UpdateCharacterModel).Methods("PUT")

	r.HandleFunc("/api/character_like/{id}", CharacterModelLikesHandler).Methods("PUT")

	// Search
	r.HandleFunc("/character_search_results/{query}", CharacterSearchHandler)

	// Index handlers

	r.HandleFunc("/random_roster/{limit}/{offset}", RandomCharacterIndexHandler)
	r.HandleFunc("/random_roster/", RandomCharacterIndexHandler)

	r.HandleFunc("/user_open_roster/{id}/{limit}/{offset}", UserOpenCharacterRosterHandler)
	r.HandleFunc("/user_open_roster/{id}", UserOpenCharacterRosterHandler)

	r.HandleFunc("/user_roster/{limit}/{offset}", UserCharacterRosterHandler)
	r.HandleFunc("/user_roster/", UserCharacterRosterHandler)

	r.HandleFunc("/all_characters/{limit}/{offset}", AllCharacterIndexHandler)
	r.HandleFunc("/all_characters/", AllCharacterIndexHandler)

	r.HandleFunc("/{limit}/{offset}", CraftedCharacterIndexHandler)
	r.HandleFunc("/", CraftedCharacterIndexHandler)

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":"+port, r))

}
