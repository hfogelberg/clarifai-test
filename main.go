package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/hfogelberg/toogo"
)

type Image struct {
	Status struct {
		Code        int    `db:"code" json:"code"`
		Description string `db:"description" json:"description"`
	} `db:"status" json:"status"`
	Outputs []struct {
		ID     string `db:"id" json:"id"`
		Status struct {
			Code        int    `db:"code" json:"code"`
			Description string `db:"description" json:"description"`
		} `db:"status" json:"status"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
		Model     struct {
			ID         string    `db:"id" json:"id"`
			Name       string    `db:"name" json:"name"`
			CreatedAt  time.Time `db:"created_at" json:"created_at"`
			AppID      string    `db:"app_id" json:"app_id"`
			OutputInfo struct {
				Message string `db:"message" json:"message"`
				Type    string `db:"type" json:"type"`
				TypeExt string `db:"type_ext" json:"type_ext"`
			} `db:"output_info" json:"output_info"`
			ModelVersion struct {
				ID        string    `db:"id" json:"id"`
				CreatedAt time.Time `db:"created_at" json:"created_at"`
				Status    struct {
					Code        int    `db:"code" json:"code"`
					Description string `db:"description" json:"description"`
				} `db:"status" json:"status"`
			} `db:"model_version" json:"model_version"`
			DisplayName string `db:"display_name" json:"display_name"`
		} `db:"model" json:"model"`
		Input struct {
			ID   string `db:"id" json:"id"`
			Data struct {
				Image struct {
					URL string `db:"url" json:"url"`
				} `db:"image" json:"image"`
			} `db:"data" json:"data"`
		} `db:"input" json:"input"`
		Data struct {
			Concepts []struct {
				ID    string  `db:"id" json:"id"`
				Name  string  `db:"name" json:"name"`
				Value float64 `db:"value" json:"value"`
				AppID string  `db:"app_id" json:"app_id"`
			} `db:"concepts" json:"concepts"`
		} `db:"data" json:"data"`
	} `db:"outputs" json:"outputs"`
}

func main() {
	url := "https://api.clarifai.com/v2/models/aaa03c23b3724a16a56b629203edc62c/outputs"

	img := "https://puppaint.com/img/portfolio/lucky.jpg"
	jsonInput := fmt.Sprintf(`{"inputs": [{"data": {"image": {"url": "%s"}}}]}`, img)
	jsonStr := []byte(jsonInput)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Printf("Error creating request %s\n", err.Error())
		return
	}

	key := fmt.Sprintf("Key %s", toogo.Getenv("CLARIFAI_TEST_APP", ""))
	req.Header.Set("Authorization", key)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error posting to Clarifai %s\n", err.Error())
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Response status %s\n", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)

	var image *Image
	err = json.Unmarshal(body, &image)
	if err != nil {
		fmt.Printf("Error unmarshaling image data %s\n", err.Error())
		return
	}

	concepts := image.Outputs[0].Data.Concepts
	for _, c := range concepts {
		fmt.Printf("%s: %0.3f\n", c.Name, c.Value)
	}

}
