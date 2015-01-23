package main

import (
	"fmt"
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
	flatsite.InputDir = getConf("INPUT_DIR", "tmpl")

	// walk filesystem, load all templates
	flatsite.Templates = template.New("templates")
	flatsite.Templates.Funcs(funcs)
	filepath.Walk(flatsite.InputDir, func(name string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "file %q traversal err: %s\n", name, err)
			return nil
		}
		isHidden := []rune(path.Base(name))[0] == '.'
		if info.IsDir() {
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
					fmt.Fprintf(os.Stderr, "%s \"%s\": %s\n", "failed to read template file", name, err)
				}
			}
		}
		return nil
	})

	// generate and output templates
	fmt.Printf("Generating public templates:\n")
	for _, v := range flatsite.Templates.Templates() {
		pth := NewPath(v.Name())
		if pth.Chunks[0] != "output" {
			continue
		}
		fmt.Fprintf(os.Stdout, "\t%s : %#v\n", v.Name(), v)
		page := NewMap()
		page.Set("path", Path{Chunks: pth.Chunks[1:]})
		if err := flatsite.generateFile(v, page); err != nil {
			fmt.Fprintf(os.Stderr, "\t\t%s \"%s\": %s\n", "failed to generate from template", v.Name(), err)
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
	OutputDir string
	InputDir  string
	Templates *template.Template
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
		v[i] = Path{Chunks: pth.Chunks[0 : i+1]}
	}
	return v
}

func (pth Path) LastChunk() string {
	return pth.Chunks[len(pth.Chunks)-1]
}

type Map map[string]interface{}

func NewMap() Map {
	m := make(map[string]interface{})
	m[""] = ""
	return m
}

func (m Map) Set(key string, value interface{}) string {
	if _, ok := m[""]; len(m) == 1 && ok {
		delete(m, "")
	}
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
	"Set": func(m Map, key string, value interface{}) interface{} {
		m.Set(key, value)
		return value
	},
	"Nul": func(_ ...interface{}) string {
		return ""
	},
}

func (flatsite *FlatSite) generateFile(tmpl *template.Template, page Map) error {
	outputPath := filepath.Join(page.Get("path").(Path).Chunks...)
	outputPathFull := filepath.Join(flatsite.OutputDir, outputPath)
	os.MkdirAll(path.Dir(outputPathFull), 0755)
	w, err := os.Create(outputPathFull)
	if err != nil {
		return fmt.Errorf("error creating static file %q: %s", outputPath, err)
	}
	defer w.Close()

	return tmpl.Execute(w, &page)
}
