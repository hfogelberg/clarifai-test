package main

import "time"

type Tag struct {
	Name  string
	Value float64
}

type Decoded struct {
	Url  string
	Tags []Tag
}

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
