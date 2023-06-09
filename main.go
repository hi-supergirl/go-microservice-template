package main

import (
	"fmt"
	"os"

	"github.com/hi-supergirl/go-microservice-template/server"
	"github.com/spf13/cobra"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "go-microservice-template.exe",
	Short: "A template for gin based micro-service project",
	Long: `
	This template will use the following tools:
	- gin : Http web framework
	- fx : A dependency injection system
	- gorm : Database access solution
	- koanf : a simple, extremely lightweight, extensible, configuration management library
	- cobra : A CLI application
	- zap :  A fast, structured, leveled logging
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Start web server from rootCmd ...")
		fmt.Println("configFile =", configFile)
		server.StartApplication(configFile)
	},
}

var subCommand = &cobra.Command{
	Use:     "server",
	Short:   "start web server",
	Aliases: []string{"s"},
	Args:    cobra.ExactArgs(0), // only 0 parameter for command1
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Start web server from subCommand ...")
		fmt.Println("configFile =", configFile)
		server.StartApplication(configFile)
	},
}

func init() {
	rootCmd.AddCommand(subCommand)
	rootCmd.PersistentFlags().StringVarP(&configFile, "conf", "c", "./config/properties.json", "config file path")
}

// Take the code references in the following order if you want to implement the same template by your own from scratch,
// 1. https://github.com/hi-supergirl/go-micro-service-example/tree/master/dive-in-cobra
// 2. https://github.com/hi-supergirl/go-learning-fx/tree/master/IntegrateFxWithGin-3
// 3. https://github.com/hi-supergirl/go-learning-fx/tree/master/IntegrateFxWithZapLogger4
// 5. https://github.com/hi-supergirl/go-micro-service-example/tree/master/dive-zapSugaredLogger
// 6. https://github.com/hi-supergirl/go-practices/tree/master/syncDemo
// 7. https://github.com/hi-supergirl/go-learning-fx/tree/master/callOrder
// 8. https://github.com/hi-supergirl/go-learning-gin/tree/master/loginlogout/jwtDemo
// 9. https://github.com/hi-supergirl/go-learning-gorm
// 10. https://github.com/hi-supergirl/go-micro-service-example/tree/master/dive-koanf
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
