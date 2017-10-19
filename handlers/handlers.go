package handlers

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/haisum/recaptcha"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
)

type messageModel struct {
	Title   string
	Message string
}

var templates map[string]*template.Template
var reCaptchaSiteKey = os.Getenv("RECAPTCHA_SITE_KEY")
var reCaptcha = recaptcha.R{
	Secret: os.Getenv("RECAPTCHA_SECRET_KEY"),
}

var m = minify.New()

func compileTemplates(templatePaths ...string) (*template.Template, error) {
	var tmpl *template.Template

	for _, templatePath := range templatePaths {
		name := filepath.Base(templatePath)

		if tmpl == nil {
			tmpl = template.New(name)
		} else {
			tmpl = tmpl.New(name)
		}

		b, err := ioutil.ReadFile(templatePath)
		if err != nil {
			return nil, err
		}

		mb, err := m.Bytes("text/html", b)
		if err != nil {
			return nil, err
		}

		tmpl.Parse(string(mb))
	}

	return tmpl, nil
}

func minifyCSSFiles(templateDirectory string) {
	cssFileDirectory := filepath.Join(templateDirectory, "../assets/css/")
	cssFilePaths, _ := filepath.Glob(filepath.Join(cssFileDirectory, "*.css"))

	for _, cssFilePath := range cssFilePaths {
		if strings.HasSuffix(cssFilePath, ".min.css") {
			continue
		}

		cssFile, err := ioutil.ReadFile(cssFilePath)
		if err != nil {
			panic(err)
		}

		cssFile, err = m.Bytes("text/css", cssFile)
		if err != nil {
			panic(err)
		}

		minCSSFilePath := strings.Replace(cssFilePath, ".css", ".min.css", 1)
		err = ioutil.WriteFile(minCSSFilePath, cssFile, 0666)
		if err != nil {
			panic(err)
		}
	}
}

func renderTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {

	var user *User

	if userData := r.Context().Value("user"); userData != nil {
		user = userData.(*User)
	} else {
		user = &User { LoggedIn: false}
	}

	obj := make(map[string]interface{})

	obj["User"] = user
	obj["Context"] = data
	//
	//obj["User"] = user
	//for key, value := range data.(map[string]interface{}) {
	//	obj[key] = value
	//}

	templates[name].ExecuteTemplate(w, "layout", obj)
}
// Execute loads templates from the specified directory and configures routes.
func Execute(templateDirectory string) error {
	if _, err := os.Stat(templateDirectory); err != nil {
		return fmt.Errorf("Could not find template directory '%s'", templateDirectory)
	}

	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)
	minifyCSSFiles(templateDirectory)

	// Load template paths.
	templatePaths, _ := filepath.Glob(filepath.Join(templateDirectory, "*.tmpl"))
	sharedPaths, _ := filepath.Glob(filepath.Join(templateDirectory, "shared/*.tmpl"))

	// Load the templates.
	templates = make(map[string]*template.Template)
	for _, templatePath := range templatePaths {
		tmpl := template.Must(compileTemplates(append(sharedPaths, templatePath)...))

		name := strings.Split(filepath.Base(templatePath), ".")[0]
		templates[name] = tmpl
	}

	// Configure the routes.

	handleFunc := func(pattern string, handler http.HandlerFunc) {
		http.HandleFunc(pattern, authMiddleware(handler))
	}
	handleFunc("/", index)
	handleFunc("/robots.txt", robots)
	handleFunc("/events", events)
	handleFunc("/team", validate(team))
	handleFunc("/gallery", gallery)
	handleFunc("/partners", partners)
	handleFunc("/sign-up", signUp)
	handleFunc("/login", login)
	handleFunc("/contact", contact)
	handleFunc("/unsubscribe", unsubscribe)

	return nil
}
