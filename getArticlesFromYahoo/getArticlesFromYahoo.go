package getArticlesFromYahoo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
)

type Score struct {
	Mixed float64
	Negative float64
	Neutral float64
	Positive float64
}

type SentimentResult struct {
	Sentiment string
	SentimentScore Score
}

type Article struct {
	Link string `json:"link"`
	Description string `json:"description"`
}

type Link struct {
	Title string
	Link string
}

func GetArticlesFromYahoo()(sentiment []Link){
	links := make([]Link, 0)
	// // provide empty array of Sentiment struct to be set result
	// allSe := make([]SentimentResult, 0)
	// allArticles := make([]Article, 0)
	// call scrape function by geziyor
	geziyor.NewGeziyor(&geziyor.Options{
		// access news site to start scraping
    StartRequestsFunc: func(g *geziyor.Geziyor) {
			// JS rendered request by GetRendered
			// geziyor
        g.GetRendered("https://news.yahoo.co.jp/", g.Opt.ParseFunc)
    },
		// parse function to scrape contents
		// g: type geziyor
		// r: type http client{ Body: response body, HTMLDoc: goquery document object, Request:meta data, render option }
    ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
				// scraping
				// r.HTMLDoc(goquery) is the one who is scraping
				// result come as array, so iterate it for getting SentimentResult(text recognize)
				sLink := Link{}
				r.HTMLDoc.Find("div.newsFeed > ul > li > div").Each(func(i int, s *goquery.Selection){
					sLink.Title = s.Find("div.newsFeed_item_text").Text()
					sLink.Link, _ = s.Find("a").Attr("href")
					links = append(links, sLink)
				})

    },
}).Start()
return links
}

// aws text recognition
// param: { text string }
// return { *comprehend.DetectSentimentOutput, error }
func DetectSentiment(text string)(*comprehend.DetectSentimentOutput, error){
	svc := comprehend.New(session.Must(session.NewSession()), &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})
	input := &comprehend.DetectSentimentInput{
		LanguageCode: aws.String("ja"),
		Text: aws.String(text),
	}
	res, err := svc.DetectSentiment(input)
	if err != nil {
		log.Println(err)
	}
	return res, nil
}

func writeJSON(data []Article){
	fmt.Println("data", data)
	// make a space 
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create JSON file")
	}

	_ = ioutil.WriteFile("facts.json", file, 0664)
}