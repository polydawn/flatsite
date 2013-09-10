package main

import (
	. "fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	flatsite := &FlatSite{}
	flatsite.OutputDir = getConf("OUTPUT_DIR", "www")
	flatsite.InputDir  = getConf("INPUT_DIR",  "tmpl")

	// read config for site
	site := SiteData{
		Name: getConf("SITE_NAME", "SITE_NAME"),
		BaseUrl: getConf("BASE_URL", "http://localhost/"),
	}

	// walk filesystem, load all templates
	flatsite.Templates = template.New("templates")
	flatsite.Templates.Funcs(funcs)
	filepath.Walk(flatsite.InputDir, func(name string, info os.FileInfo, err error) error {
		isHidden := []rune(path.Base(name))[0] == '.'
		if (info.IsDir()) {
			if isHidden {
				return filepath.SkipDir
			} else {
				return nil
			}
		} else {
			if isHidden {
				return nil
			} else {
				if content, err := ioutil.ReadFile(name); err == nil {
					name, _ = filepath.Rel(flatsite.InputDir, name)
					flatsite.Templates.New(name).Parse(string(content))
				} else {
					Fprintf(os.Stderr, "%s \"%s\": %s\n", "failed to read template file", name, err)
				}
			}
		}
		return nil
	})

	// generate and output templates
	Printf("Generating public templates:\n")
	for _, v := range flatsite.Templates.Templates() {
		pth := NewPath(v.Name())
		if pth.Chunks[0] != "output" { continue; }
		Fprintf(os.Stdout, "\t%s : %#v\n", v.Name(), v)
		page := PageData{
			Site: site,
			Path: pth,
			Template: v,
		}
		if err := flatsite.generateFile(page); err != nil {
			Fprintf(os.Stderr, "\t\t%s \"%s\": %s\n", "failed to generate from template", v.Name(), err)
		}
	}
}

func getConf(key string, defalt string) string {
	env := os.Getenv(key)
	if env == "" {
		return defalt
	} else {
		return env
	}
}

type FlatSite struct {
	OutputDir  string
	InputDir   string
	Templates  *template.Template
}

type SiteData struct {
	Name     string
	BaseUrl  string
}

type PageData struct {
	Site      SiteData
	Path      Path
	Template  *template.Template
	Title     string
}

func NewPath(pth string) Path {
	return Path{
		Chunks: strings.Split(pth, string(os.PathSeparator)),
	}
}

type Path struct {
	Chunks []string
}

func (pth Path) String() string {
	return strings.Join(pth.Chunks, "/")
}

func (pth Path) Paths() []Path {
	n := len(pth.Chunks)
	v := make([]Path, n)
	for i := 0; i < n; i++ {
		v[i] = Path{Chunks: pth.Chunks[0:i+1]}
	}
	return v
}

func (pth Path) LastChunk() string {
	return pth.Chunks[len(pth.Chunks)-1]
}

func (page *PageData) SetTitle(title string) string {
	page.Title = title
	return ""
}

type Map map[string]interface{}

func NewMap() Map {
	return make(map[string]interface{})
}

func (m Map) Set(key string, value interface{}) string {
	m[key] = value
	return ""
}

func (m Map) Get(key string) interface{} {
	return m[key]
}

var funcs = template.FuncMap{
	"eq": func(a interface{}, b interface{}) bool {
		return a == b
	},
	"NewMap": func() Map {
		return NewMap()
	},
}

func (flatsite *FlatSite) generateFile(page PageData) error {
	outputPath := filepath.Join(page.Path.Chunks[1:]...)
	outputPathFull := filepath.Join(flatsite.OutputDir, outputPath)
	os.MkdirAll(path.Dir(outputPathFull), 0755)
	w, err := os.Create(outputPathFull)
	if err != nil {
		return Errorf("error creating static file %s: %s", outputPath, err)
	}
	defer w.Close()

	page.Path = Path{Chunks: page.Path.Chunks[1:]}
	return page.Template.Execute(w, &page)
}
