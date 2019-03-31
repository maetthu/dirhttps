package cmd

import (
	"fmt"
	"github.com/maetthu/dirhttps/internal/lib/version"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	CertFilename = "cert.pem"
	KeyFilename = "key.pem"
)

var rootCmd = &cobra.Command{
	Use:   "dirhttps",
	Short: "Serving contents of current directory by HTTPS.",
	Args: cobra.NoArgs,
	Version: fmt.Sprintf("%s -- %s", version.Version, version.Commit),
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()

		if err != nil {
			log.Fatalf("Error determining home directory: %s", err)
		}

		configDir := fmt.Sprintf("%s/.config/dirhttps", homeDir)

		certFile := filepath.Join(configDir, CertFilename)
		keyFile := filepath.Join(configDir, KeyFilename)
		certAvailable := true

		if _, err := os.Stat(certFile); err != nil  {
			certAvailable = false
		}

		if _, err := os.Stat(keyFile); err != nil  {
			certAvailable = false
		}

		if !certAvailable {
			log.Printf("Necessary certificate and/or key file not found.")
			log.Fatalf(
				"Place a certificate \"%s\" and a key file \"%s\" to %s",
				CertFilename,
				KeyFilename,
				configDir,
				)
		}

		dir, err := os.Getwd()

		if err != nil {
			log.Fatalf("Error determining current directory: %s", err)
		}

		addr, err := cmd.Flags().GetString("listen")

		if err != nil {
			log.Fatalf("Error parsing listen flage: %s", err)
		}

		log.Printf("Listening for HTTPS connections on %s", addr)
		log.Printf("Serving from directory %s", dir)

		// permit some CORS stuff
		corsHandler := cors.New(cors.Options{
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {return true},
		})

		handler := corsHandler.Handler(logger(http.FileServer(http.Dir(dir))))
		log.Fatal(http.ListenAndServeTLS(addr, certFile, keyFile, handler))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init(){
	rootCmd.Flags().StringP("listen", "l", ":8443", "Listen address")
	//rootCmd.Flags().Bool("version", false, "Print version and exit")
}

func logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf( "%s \"%s %s %s\"\n", r.RemoteAddr, r.Method, r.URL, r.Proto)
		handler.ServeHTTP(w, r)
	})
}