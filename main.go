package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/hfogelberg/toogo"
	"github.com/kyokomi/cloudinary"
	"github.com/urfave/negroni"
)

const CloudinaryRoot = "http://res.cloudinary.com/golizzard/image/upload/h_200,c_scale/v1507814747/"

func main() {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/upload", uploadHandler)
	r.HandleFunc("/result/{id}", resultHandler)

	static := http.StripPrefix("/public/", http.FileServer(http.Dir("public")))
	r.PathPrefix("/public").Handler(static)

	mux := http.NewServeMux()
	mux.Handle("/", r)

	port := toogo.Getenv("PORT", ":80")
	n := negroni.Classic()
	n.UseHandler(mux)
	http.ListenAndServe(port, n)

}

func recognizeImage(img string) ([]Tag, error) {
	var tags []Tag
	var image *Image

	url := "https://api.clarifai.com/v2/models/aaa03c23b3724a16a56b629203edc62c/outputs"
	jsonInput := fmt.Sprintf(`{"inputs": [{"data": {"image": {"url": "%s"}}}]}`, img)
	jsonStr := []byte(jsonInput)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Printf("Error creating request %s\n", err.Error())
		return tags, err
	}

	key := fmt.Sprintf("Key %s", toogo.Getenv("CLARIFAI_TEST_APP", ""))
	req.Header.Set("Authorization", key)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error posting to Clarifai %s\n", err.Error())
		return tags, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &image)
	if err != nil {
		fmt.Printf("Error unmarshaling image data %s\n", err.Error())
		return tags, err
	}

	concepts := image.Outputs[0].Data.Concepts
	for _, c := range concepts {
		tag := Tag{
			Name:  c.Name,
			Value: c.Value,
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.New("").ParseFiles("templates/index.html", "templates/layout.html")
	err = tpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Fatalln("Error serving index template ", err.Error())
		return
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("image")
	if err != nil {
		log.Printf("Error reading form image %s\n", err.Error())
		return
	}
	defer file.Close()

	root := "./public/temp/"
	os.Mkdir(root, 0700)

	// Generate a file name
	id := fmt.Sprint(time.Now().Unix())
	path := root + "/" + id + ".jpg"

	out, err := os.Create(path)
	if err != nil {
		log.Printf("Error creating file in public/temp %s", err.Error())
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Printf("Error writing file to public/tmp %s", err.Error())
		return
	}

	if err := cloudinaryUpload(path, id); err != nil {
		log.Printf("Error uploading to Cloudinary %s\n", err.Error())
		return
	}

	fmt.Println("Done!")
	http.Redirect(w, r, "/result/"+id, http.StatusPermanentRedirect)

}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	u := strings.Split(url, "/")
	id := u[len(u)-1]

	img := fmt.Sprintf("%s%s.jpg", CloudinaryRoot, id)

	tags, err := recognizeImage(img)
	if err != nil {
		return
	}

	decoded := Decoded{
		Url:  img,
		Tags: tags,
	}

	tpl, err := template.New("").ParseFiles("templates/result.html", "templates/layout.html")
	err = tpl.ExecuteTemplate(w, "layout", decoded)
	if err != nil {
		log.Fatalln("Error serving result template ", err.Error())
		return
	}
}

func cloudinaryUpload(src string, fileName string) error {
	ctx := context.Background()

	key := toogo.Getenv("CLOUDINARY_API_KEY", "925374862654622")
	secret := toogo.Getenv("CLOUDINARY_API_SECRET", "doHBawwQUw7L2vYVKq5Dl9wbdUE")
	cloud := toogo.Getenv("CLOUDINARY_CLOUD_NAME", "golizzard")

	con := fmt.Sprintf("cloudinary://%s:%s@%s", key, secret, cloud)
	ctx = cloudinary.NewContext(ctx, con)

	data, _ := ioutil.ReadFile(src)

	if err := cloudinary.UploadStaticImage(ctx, fileName, bytes.NewBuffer(data)); err != nil {
		log.Println("Error uploading image to cloudinary")
		return err
	}

	_ = os.Remove(src)

	return nil
}
