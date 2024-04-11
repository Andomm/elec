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
	Short: "Get prices for electricity",
	Long:  "Get prices for electricity for a given date and time",
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&Date, "date", "d", "", "Define the date to get prices for")
	viper.BindPFlag("date", rootCmd.PersistentFlags().Lookup("date"))
	// rootCmd.PersistentFlags().StringVarP(&Time, "time", "t", time.Now().Format("15:04"), "Define the time to get prices for")
	// viper.BindPFlag("time", rootCmd.PersistentFlags().Lookup("time"))
	rootCmd.PersistentFlags().IntVar(&Hours, "hours", 24, "Define the number of hours to get prices for")
	rootCmd.MarkFlagRequired("hours")
	viper.BindPFlag("hours", rootCmd.PersistentFlags().Lookup("hours"))
}
