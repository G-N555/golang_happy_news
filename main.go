package main

import (
	"fmt"
	"golang_scraping/getArticlesFromYahoo"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
)

type Sentiment struct {
	Mixed float64
	Negative float64
	Neutral float64
	Positive float64
}

type Score struct{
		Title string
		Sentiment string
		SentimentScore Sentiment
}

func main(){
	links := getArticlesFromYahoo.GetArticlesFromYahoo()
	scores := []Score{}
	for i, v := range links{
		if v.Link != "" && v.Title != "" && i < 5{
			geziyor.NewGeziyor(&geziyor.Options{
				StartRequestsFunc: func(g *geziyor.Geziyor){
					g.GetRendered(v.Link, g.Opt.ParseFunc)
				},
				ParseFunc: func(g *geziyor.Geziyor, r *client.Response){
					r.HTMLDoc.Find("div > p.highLightSearchTarget").Each(func(i int, s *goquery.Selection){
						res, _ := getArticlesFromYahoo.DetectSentiment(s.Text())
						sScore := Score{
							Title: v.Title,
							Sentiment: *res.Sentiment,
							SentimentScore: Sentiment{
								Mixed: *res.SentimentScore.Mixed,
								Negative: *res.SentimentScore.Negative,
								Neutral: *res.SentimentScore.Neutral,
								Positive: *res.SentimentScore.Positive,
							},
						}

						scores = append(scores, sScore)
					})
				},
			}).Start()
		}
	}
	
	fmt.Println("score", scores)
}

