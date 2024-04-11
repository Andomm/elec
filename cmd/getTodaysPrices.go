/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	// "time"
)

type Price []struct {
	Hour    string `json:"aikaleima_suomi"`
	UtcHour string `json:"aikaleima_utc"`
	Price   string `json:"hinta"`
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type Styles struct {
	Header   lipgloss.Style
	Cell     lipgloss.Style
	Selected lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Selected: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")),
		Header:   lipgloss.NewStyle().Bold(true).Padding(0, 1),
		Cell:     lipgloss.NewStyle().Padding(0, 1),
	}
}

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

var Date string
var Time string
var Hours int

// getTodaysPricesCmd represents the getTodaysPrices command
var getTodaysPricesCmd = &cobra.Command{
	Use:   "prices",
	Short: "Get prices for electricity for today",
	Long: `Get prices for electricity for today.
		You can get prices for different times of the day`,
	Run: func(cmd *cobra.Command, args []string) {
		getTodaysPrices(Date, Hours)
	},
}

func getTodaysPrices(date string, hours int) {
	base, err := url.Parse("https://www.sahkohinta-api.fi/api/v1/halpa")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	params := url.Values{}
	if date != "" {
		params.Add("aikaraja", date)
	}
	params.Add("tunnit", fmt.Sprintf("%d", hours))
	params.Add("tulos", "sarja")
	base.RawQuery = params.Encode()

	fmt.Println("URL:", base.String())
	resp, err := http.Get(base.String())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	var result Price
	switch resp.StatusCode {
	case 200:
		body, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Error parsing data:", err)
			return
		}
	case 204:
		fmt.Println("No prices found for that day")
	case 400:
		fmt.Println("Invalid date")
	case 500:
		fmt.Println("Something unexpected:", resp.Status)
	}

	columns := []table.Column{
		{Title: "Hour", Width: 20},
		{Title: "Price", Width: 10},
	}

	rows := []table.Row{}

	for _, price := range result {
		row := table.Row{price.Hour, price.Price}
		rows = append(rows, row)
		fmt.Println(price.Hour, price.UtcHour, price.Price)
	}
	table := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	m := model{table}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error starting program:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(getTodaysPricesCmd)
}
