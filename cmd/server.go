package cmd

import (
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/devinmcgloin/sail/pkg/library"
	"github.com/devinmcgloin/sail/pkg/renderer"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/cobra"
)

// serverCmd represents the info command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server spins up a webserver to generate images on the fly",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt64("port")
		//slog.SetLevel(slog.ERROR)
		log.Printf("Serving requests at: http://localhost:%d\n", port)
		router := httprouter.New()
		router.GET("/", index)
		router.GET("/render/:category/:sketch", render)
		router.GET("/render/:category/:sketch/:seed", render)
		fs := http.Dir("sketches")
		router.ServeFiles("/view/*filepath", fs)

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
	},
}

func render(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sketchID := fmt.Sprintf("%s/%s", ps.ByName("category"), ps.ByName("sketch"))
	if !library.Exists(sketchID) {
		fmt.Fprintf(w, "Unable to load sketch: %s\n", sketchID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	seedString := ps.ByName("seed")
	var seed int64
	if seedString == "" {
		seed = time.Now().Unix()
	} else {
		i, err := strconv.ParseInt(seedString, 0, 64)
		if err != nil {
			seed = hash(seedString)
		} else {
			seed = i
		}
	}
	log.Printf("Rendering %s with seed %d\n", sketchID, seed)
	bytes, err := renderer.Render(sketchID, false, seed)
	if err != nil {
		fmt.Fprintf(w, "An Error Occured: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bytes.Bytes())
}

func index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sketches := library.List(".*")
	fmt.Fprintf(w, "Sail has the following sketches:\n%s", strings.Join(sketches, "\n"))
}

func hash(s string) int64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	sum := h.Sum64()
	if int64(sum) <= 0 {
		return int64(sum) * -1
	}
	return int64(sum)
}

func init() {
	serverCmd.Flags().Int64P("port", "p", 8080, "port to bind server responses to")
	rootCmd.AddCommand(serverCmd)
}
