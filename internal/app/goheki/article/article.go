package article

import (
	//"github.com/maguro-alternative/goheki/internal/app/goheki/utility"
	"html/template"
	//"log"
	"net/http"
	"strings"
)

type Article struct {
	Id    int
	Title string
	Body  string
}

var tmpl *template.Template

func init() {
	funcMap := template.FuncMap{
		"nl2br": func(text string) template.HTML {
			return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br />", -1))
		},
	}

	tmpl, _ = template.New("article").Funcs(funcMap).ParseGlob("web/template/*")
}

func Index(w http.ResponseWriter, r *http.Request) {

}

func Show(w http.ResponseWriter, r *http.Request) {

}

func Create(w http.ResponseWriter, r *http.Request) {

}

func Edit(w http.ResponseWriter, r *http.Request) {

}

func Delete(w http.ResponseWriter, r *http.Request) {

}
