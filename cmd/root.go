package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "search-filter",
	Short: "Search Filter Service CLI",
	Long:  `CLI для управления сервисом search-filter (запуск сервера, миграции базы данных и др).`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
