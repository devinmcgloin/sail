package cmd

import (
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/devinmcgloin/sail/pkg/renderer"
	"github.com/devinmcgloin/sail/pkg/slog"

	"github.com/spf13/cobra"
)

// serverCmd represents the info command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server spins up a webserver to generate images on the fly",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt64("port")
		//slog.SetLevel(slog.ERROR)
		http.HandleFunc("/", index)

		http.HandleFunc("/render", render)

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	},
}

func render(w http.ResponseWriter, r *http.Request) {
	sketchID := "delaunay/ring" //fmt.Sprintf("%s/%s", ps.ByName("category"), ps.ByName("sketch"))
	seedString := ""            //ps.ByName("seed")
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
	slog.InfoPrintf("Rendering %s with seed %d\n", sketchID, seed)
	bytes, err := renderer.Render(sketchID, false, seed)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "An Error Occured: %s\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bytes.Bytes())
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world!\n")
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
