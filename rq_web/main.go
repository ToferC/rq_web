package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gorilla/mux"

	"github.com/go-pg/pg"
	"github.com/toferc/ore_web_roller/database"
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
	os.Setenv("DBName", "ore_engine")
}

func main() {

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

	} else {
		// Not production
		db = pg.Connect(&pg.Options{
			User:     os.Getenv("DBUser"),
			Password: os.Getenv("DBPass"),
			Database: os.Getenv("DBName"),
		})
		os.Setenv("CookieSecret", "kimchee-typhoon")
		os.Setenv("BUCKET", "rqweb")
		os.Setenv("AWS_REGION", "us-east-1")
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

	// Check for terminal flag
	t := flag.Bool("t", false, "Activate the local terminal")

	flag.Parse()

	if *t == true {
		Terminal(db)
	} else {

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

		r.HandleFunc("/users/", UserIndexHandler)

		r.Path("/roll/{id}").HandlerFunc(RollHandler)
		r.Path("/roll/{id}").Queries(
			"ac", "",
			"d", "",
			"hd", "",
			"wd", "",
			"gf", "",
			"sp", "",
			"nr", "",
			"ed", "").HandlerFunc(RollHandler).Name("RollHandler")
		r.HandleFunc("/opposed/", OpposeHandler)
		r.Path("/opposed/").Queries(
			"name1", "",
			"ac", "",
			"nd", "",
			"hd", "",
			"wd", "",
			"gf", "",
			"sp", "",
			"name2", "",
			"ac2", "",
			"nd2", "",
			"hd2", "",
			"wd2", "",
			"gf2", "",
			"sp2", "",
		).HandlerFunc(OpposeHandler).Name("OpposeHandler")

		r.HandleFunc("/view_character/{id}", CharacterHandler)
		r.HandleFunc("/new/{setting}", NewCharacterHandler)
		r.HandleFunc("/modify/{id}", ModifyCharacterHandler)
		r.HandleFunc("/delete/{id}", DeleteCharacterHandler)

		r.HandleFunc("/index_powers/", PowerIndexHandler)
		r.HandleFunc("/view_power/{id}", PowerHandler)

		r.HandleFunc("/add_power/{id}", AddPowerHandler)
		r.HandleFunc("/add_power_from_list/{id}", PowerListHandler)
		r.HandleFunc("/modify_power/{id}/{power}", ModifyPowerHandler)

		r.HandleFunc("/add_standalone_power/", AddStandalonePowerHandler)
		r.HandleFunc("/modify_standalone_power/{id}", ModifyStandalonePowerHandler)
		r.HandleFunc("/delete_power/{id}/{power}", DeletePowerHandler)
		r.HandleFunc("/delete_standalone_power/{id}", DeleteStandalonePowerHandler)

		r.HandleFunc("/add_hyperstat/{id}/{stat}", AddHyperStatHandler)
		r.HandleFunc("/modify_hyperstat/{id}/{stat}", ModifyHyperStatHandler)
		r.HandleFunc("/delete_hyperstat/{id}/{stat}", DeleteHyperStatHandler)

		r.HandleFunc("/add_hyperskill/{id}/{skill}", AddHyperSkillHandler)
		r.HandleFunc("/modify_hyperskill/{id}/{skill}", ModifyHyperSkillHandler)
		r.HandleFunc("/delete_hyperskill/{id}/{skill}", DeleteHyperSkillHandler)

		r.HandleFunc("/add_skill/{id}/{stat}", AddSkillHandler)
		r.HandleFunc("/add_advantages/{id}", ModifyAdvantageHandler)

		r.HandleFunc("/user_index/", UserIndexHandler)

		http.Handle("/", r)

		log.Fatal(http.ListenAndServe(":"+port, r))
	}

}
