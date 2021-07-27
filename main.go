package main

import (
	"embed"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/afero"
)

//go:embed index.html template.svg
var content embed.FS

var port = os.Getenv("PORT")

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		text := r.URL.Query().Get("text")
		font := r.URL.Query().Get("font")

		if text == "" {
			body, _ := content.ReadFile("index.html")
			w.Write(body)
			return
		}

		dxf, err := createDXF(text, font)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		fmt.Println("Converted to DXF")

		w.Header().Add("Content-Type", "application/octet-stream")
		w.Header().Add("Content-Length", fmt.Sprint(len(dxf)))
		w.Header().Add("Content-Disposition", fmt.Sprintf("inline; filename=%s.dxf", randomString(7)))
		fmt.Fprintf(w, string(dxf))
	})

	fmt.Println(http.ListenAndServe(":"+port, nil))
}

func createDXF(text, font string) ([]byte, error) {
	fs := afero.NewOsFs()
	afs := &afero.Afero{Fs: fs}
	svgFile, _ := afs.TempFile(".", "*.svg")
	epsFile, _ := afs.TempFile(".", "*.eps")
	dxfFile, _ := afs.TempFile(".", "*.dxf")

	err := template.
		Must(
			template.
				New("template.svg").
				Funcs(template.FuncMap{
					"offset": func(i int) int {
						return i * 30
					},
					"font": func() string {
						return font
					},
				}).
				ParseFS(content, "template.svg")).
		Execute(svgFile, strings.Split(text, "\n"))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	cmd := exec.Command("inkscape", "-E", epsFile.Name(), svgFile.Name())
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return nil, err
	}
	cmd = exec.Command("pstoedit", "-q", "-dt", "-f", "dxf:-polyaslines -mm", epsFile.Name(), dxfFile.Name())
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return afs.ReadFile(dxfFile.Name())
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
