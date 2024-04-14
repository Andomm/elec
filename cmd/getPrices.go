/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
	"time"
)

var Date string
var Time string
var Hours int
var Expensive bool

type Price []struct {
	Hour    string `json:"aikaleima_suomi"`
	UtcHour string `json:"aikaleima_utc"`
	Price   string `json:"hinta"`
}

var getPricesCmd = &cobra.Command{
	Use:   "prices",
	Short: "Get prices for electricity",
	Long:  `Get prices for electricity for a given date`,
	Run: func(cmd *cobra.Command, args []string) {
		getPrices(Date, Hours, Expensive)
	},
}

func handleNonOKStatus(statusCode int) error {
	switch statusCode {
	case 204:
		return errors.New("No prices found for that day")
	case 400:
		return errors.New("Invalid date")
	case 500:
		return errors.New("Something unexpected")
	default:
		return fmt.Errorf("Unexpected status code: %d", statusCode)

	}
}

func getPrices(date string, hours int, expensive bool) error {
	baseUrl := "https://www.sahkohinta-api.fi/api/v1/"
	if expensive {
		baseUrl += "kallis"
	} else {
		baseUrl += "halpa"
	}
	base, err := url.Parse(baseUrl)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	params := url.Values{}
	if date != "" {
		params.Add("aikaraja", date)
	} else {
		params.Add("aikaraja", time.Now().Format("2006-01-02"))
	}

	params.Add("tunnit", fmt.Sprintf("%d", hours))
	params.Add("tulos", "haja")
	base.RawQuery = params.Encode()

	fmt.Println("URL:", base.String())
	resp, err := http.Get(base.String())
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return handleNonOKStatus(resp.StatusCode)
	}

	var result Price

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}
	return CreateAndRunTable(result)
}

func init() {
	rootCmd.AddCommand(getPricesCmd)
}
