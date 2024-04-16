/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	// "time"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "elec",
	Short: "Get cheapest prices for electricity",
	Long:  "Get cheapest prices for electricity for a given date and time",
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&Date, "date", "d", "", "Define the date to get prices for (defaults to today)")
	viper.BindPFlag("date", rootCmd.PersistentFlags().Lookup("date"))
	rootCmd.PersistentFlags().IntVar(&Hours, "hours", 24, "Define the number of hours to get prices for")
	rootCmd.MarkFlagRequired("hours")
	rootCmd.PersistentFlags().BoolVarP(&Expensive, "expensive", "e", false, "Boolean flag to get the most expensive prices")
	viper.BindPFlag("hours", rootCmd.PersistentFlags().Lookup("hours"))
}
