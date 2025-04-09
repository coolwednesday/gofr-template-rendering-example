package main

import (
	"encoding/json"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http/response"
	"strconv"
)

type Article struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

type PageData struct {
	Articles []Article
	Tag      string
	Page     int
}

func FetchArticles(ctx *gofr.Context) (any, error) {
	tag := ctx.Param("tag")
	if tag == "" {
		tag = "go"
	}

	pageStr := ctx.Param("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	service := ctx.GetHTTPService("dev-articles")

	queryParams := map[string]interface{}{
		"tag":      tag,
		"page":     page,
		"per_page": 4,
	}

	resp, err := service.Get(ctx, "/articles", queryParams)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var articles []Article

	if err := json.NewDecoder(resp.Body).Decode(&articles); err != nil {
		return nil, err
	}

	return response.Template{
		Data: PageData{Articles: articles, Tag: tag, Page: page},
		Name: "devTo.html",
	}, nil
}

func main() {
	app := gofr.New()

	app.AddHTTPService("dev-articles", "https://dev.to/api")
	app.GET("/dev-articles", FetchArticles)

	app.AddStaticFiles("/", "./static")

	app.Run()
}
