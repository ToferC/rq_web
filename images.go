package main

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-pg/pg"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/disintegration/imaging"

	"github.com/toferc/rq_web/database"
	"github.com/toferc/rq_web/models"
	"github.com/toferc/runequest"
)

// ProcessImage resizes and uploads an image to AWS
func ProcessImage(h *multipart.FileHeader, file multipart.File, cm *models.CharacterModel) error {
	// Upload image to s3

	fNameSplit := strings.Split(h.Filename, ".")
	fName := runequest.ToSnakeCase(fNameSplit[0]) + ".jpeg"

	// example path media/Major/TestImage/Jason_White.jpg
	path := fmt.Sprintf("/media/%s/%s/%s",
		cm.Author.UserName,
		runequest.ToSnakeCase(cm.Character.Name),
		fName,
	)

	fmt.Println(path)

	img, err := imaging.Decode(file)
	if err != nil {
		fmt.Print("Imaging Open error")
		log.Print("Error decoding", err)
		return err
	}

	newImage := imaging.Resize(img, 350, 0, imaging.Lanczos)

	buf := new(bytes.Buffer)
	err = imaging.Encode(buf, newImage, imaging.JPEG)
	if err != nil {
		log.Printf("JPEG encoding error: %v", err)
	}

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("BUCKET")),
		Key:    aws.String(path),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		log.Panic(err)
		fmt.Println("Error uploading file ", err)
		return err
	}
	fmt.Printf("successfully uploaded %q to %q\n",
		h.Filename, os.Getenv("BUCKET"))

	if cm.Image == nil {
		cm.Image = new(models.Image)
	}
	cm.Image.Path = path

	fmt.Println(path)
	return nil
}

// ResizeImages loads all images from AWS S3 and resizes them before re-saving
func ResizeImages(db *pg.DB) error {

	cms, _ := database.ListAllCharacterModels(db)

	for _, cm := range cms {

		fmt.Println("Attempting to upload: ", cm.Character.Name)

		buff := &aws.WriteAtBuffer{}

		fmt.Println("Downloading image")

		if cm.Image != nil {
			downloader := s3manager.NewDownloader(session.New())
			_, err := downloader.Download(buff,
				&s3.GetObjectInput{
					Bucket: aws.String("runequeset"),
					Key:    aws.String(cm.Image.Path),
				})

			if err != nil {
				log.Println("Error - failed to download file", err)
				continue
			}

			fmt.Println("Decoding Image")

			img, err := imaging.Decode(bytes.NewReader(buff.Bytes()))
			if err != nil {
				fmt.Print("Imaging Open error")
				log.Print("Error decoding", err)
				continue
			}

			b := img.Bounds()

			if b.Max.X > 350 {

				fmt.Println("Resizing image")

				newImage := imaging.Resize(img, 350, 0, imaging.Lanczos)

				buf := new(bytes.Buffer)

				fmt.Println("Encoding image")

				err = imaging.Encode(buf, newImage, imaging.JPEG)
				if err != nil {
					log.Printf("JPEG encoding error: %v", err)
				}

				fmt.Println("Uploading image")

				_, err = uploader.Upload(&s3manager.UploadInput{
					Bucket: aws.String(os.Getenv("BUCKET")),
					Key:    aws.String(cm.Image.Path),
					Body:   bytes.NewReader(buf.Bytes()),
				})
				if err != nil {
					log.Panic(err)
					fmt.Println("Error uploading file ", err)
					return err
				}
				fmt.Printf("successfully uploaded image to %q\n",
					os.Getenv("BUCKET"))

			} else {
				fmt.Println("Image ok")
			}
		} else {
			log.Println("no file")
			cm.Image = new(models.Image)
			cm.Image.Path = DefaultCharacterPortrait

		}
	}
	return nil
}
