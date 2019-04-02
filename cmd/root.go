package cmd

import (
	"fmt"
	"github.com/maetthu/dirhttps/internal/lib/version"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
)

const (
	certFilename = "cert.pem"
	keyFilename  = "key.pem"
)

var rootCmd = &cobra.Command{
	Use:   "dirhttps",
	Short: "Serving contents of current directory by HTTPS.",
	Args: cobra.NoArgs,
	Version: fmt.Sprintf("%s -- %s", version.Version, version.Commit),
	Run: func(cmd *cobra.Command, args []string) {

		certFile, err := cmd.Flags().GetString("cert")

		if err != nil {
			log.Fatalf("Error parsing cert flag: %s", err)
		}

		keyFile, err := cmd.Flags().GetString("key")

		if err != nil {
			log.Fatalf("Error parsing key flag: %s", err)
		}

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
				"Store a certificate to \"%s\" and a key file to \"%s\" or provide the --cert and --key flags",
				certFile,
				keyFile,
				)
		}

		dir, err := os.Getwd()

		if err != nil {
			log.Fatalf("Error determining current directory: %s", err)
		}

		addr, err := cmd.Flags().GetString("listen")

		if err != nil {
			log.Fatalf("Error parsing listen flag: %s", err)
		}

		log.Printf("Listening for HTTPS connections on %s", addr)
		log.Printf("Serving from directory %s", dir)

		handler := logger(http.FileServer(http.Dir(dir)))

		// permit some CORS stuff
		if disableCORS, _ := cmd.Flags().GetBool("no-cors"); !disableCORS {
			handler = cors.New(cors.Options{
				AllowCredentials: true,
				AllowOriginFunc: func(origin string) bool {return true},
			}).Handler(handler)
		}

		// Disable browser cache?
		if cache, _ := cmd.Flags().GetBool("cache"); !cache {
			handler = nocache(handler)
		}

		// Dump request headers?
		if d, _ := cmd.Flags().GetBool("dump"); d {
			handler = dump(handler)
		}

		log.Fatal(http.ListenAndServeTLS(addr, certFile, keyFile, handler))
	},
}

// Execute runs root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init(){
	home, err := homedir.Dir()

	if err != nil {
		log.Fatalf("Error determining home directory: %s", err)
	}

	configDir := fmt.Sprintf("%s/.config/dirhttps", home)

	certFile := filepath.Join(configDir, certFilename)
	keyFile := filepath.Join(configDir, keyFilename)

	rootCmd.Flags().StringP("cert", "c", certFile, "Certificate file")
	rootCmd.Flags().StringP("key", "k", keyFile, "Key file")
	rootCmd.Flags().StringP("listen", "l", ":8443", "Listen address")

	rootCmd.Flags().Bool("cache", false, "Enable client side caching")
	rootCmd.Flags().BoolP("dump", "d",false, "Dump client request headers to STDOUT")
	rootCmd.Flags().Bool("no-cors", false, "Disable CORS handling")
}

func logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf( "%s \"%s %s %s\"\n", r.RemoteAddr, r.Method, r.URL, r.Proto)
		handler.ServeHTTP(w, r)
	})
}

func nocache(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-store")
		handler.ServeHTTP(w, r)
	})
}

func dump(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if dump, err := httputil.DumpRequest(r, true); err == nil {
			fmt.Printf("---\n%s---\n", dump)
		}

		handler.ServeHTTP(w, r)
	})
}
