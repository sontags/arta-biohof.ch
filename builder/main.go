package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"text/template"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/mitchellh/go-homedir"
	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v2"
)

func main() {
	var configPath string
	var outputPath string
	flag.StringVar(&configPath, "config", "./config.yaml", "path to the configuration file")
	flag.StringVar(&outputPath, "out", "", "path to the output file; leave empty to print to STDOUT")
	flag.Parse()

	fs := osfs.New(".")
	c, err := NewConfig(configPath, fs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	out, err := RenderTemplate(c.Content, c.Template)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if outputPath != "" {
		outputPath, err := homedir.Expand(outputPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = os.WriteFile(outputPath, out.Bytes(), 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Println(out.String())
	}
}

type content struct {
	Name     string            `json:"name,omitempty" yaml:"name"`
	Path     string            `json:"path,omitempty" yaml:"path"`
	Metadata map[string]string `json:"metadata,omitempty" yaml:"metadata"`
	Markdown []byte
	HTML     []byte
}

func (c *content) toHTML() error {
	var b bytes.Buffer
	if err := goldmark.Convert(c.Markdown, &b); err != nil {
		err = fmt.Errorf("Could not render %s as HTML: %v", c.Path, err)
		return err
	}
	c.HTML = b.Bytes()
	return nil
}

type config struct {
	Template string     `yaml:"template" json:"template,omitempty"`
	Content  []*content `yaml:"content" json:"content,omitempty"`
	fs       billy.Filesystem
}

func (c *config) loadContent() error {
	var err error
	for i, content := range c.Content {
		content.Path, err = homedir.Expand(content.Path)
		if err != nil {
			return err
		}
		file, err := c.fs.Open(content.Path)
		if err != nil {
			err = fmt.Errorf("Could not open content file '%s', error occured: %w", content.Path, err)
			return err
		}
		defer file.Close()
		var data bytes.Buffer
		_, err = data.ReadFrom(file)
		if err != nil {
			return fmt.Errorf("could not read content file %s: %s", content.Path, err.Error())
		}
		content.Markdown = data.Bytes()
		err = content.toHTML()
		if err != nil {
			return err
		}
		c.Content[i] = content
	}
	return nil
}

func NewConfig(path string, fs billy.Filesystem) (config, error) {
	c := defaults()
	c.fs = fs
	path, err := homedir.Expand(path)
	if err != nil {
		return c, err
	}
	file, err := fs.Open(path)
	if err != nil {
		err = fmt.Errorf("Could not open file '%s', error occured: %w", path, err)
		return c, err
	}
	defer file.Close()
	var data bytes.Buffer
	_, err = data.ReadFrom(file)
	if err != nil {
		return c, fmt.Errorf("could not read file %s: %s", path, err.Error())
	}
	err = yaml.Unmarshal(data.Bytes(), &c)
	if err != nil {
		return c, fmt.Errorf("could not read data from config file %s: %s", path, err.Error())
	}
	err = c.loadContent()
	if err != nil {
		return c, err
	}
	return c, nil
}

func defaults() config {
	return config{}
}

func RenderTemplate(data []*content, templ string) (bytes.Buffer, error) {
	var w bytes.Buffer
	tpl, err := template.ParseFiles(templ)
	if err != nil {
		return w, fmt.Errorf("error while parsing template: %s", err.Error())
	}
	err = tpl.Execute(&w, data)
	if err != nil {
		return w, fmt.Errorf("error while parsing template: %s", err.Error())
	}
	return w, nil
}
