package render

import (
	"bytes"
	"io/ioutil"
	"path"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"

	assets "github.com/openshift/ibm-roks-toolkit/pkg/assets"
)

type renderContext struct {
	outputDir     string
	params        interface{}
	funcs         template.FuncMap
	manifestFiles []string
	manifests     map[string]string
}

func newRenderContext(params interface{}, outputDir string) *renderContext {
	renderContext := &renderContext{
		params:    params,
		outputDir: outputDir,
		manifests: make(map[string]string),
	}
	return renderContext
}

func (c *renderContext) setFuncs(f template.FuncMap) {
	c.funcs = f
}

func (c *renderContext) renderManifests() error {
	for _, f := range c.manifestFiles {
		outputFile := filepath.Join(c.outputDir, path.Base(f))
		content, err := c.substituteParams(c.params, f)
		if err != nil {
			return errors.Wrapf(err, "cannot render %s", f)
		}
		ioutil.WriteFile(outputFile, []byte(content), 0600)
	}

	for name, content := range c.manifests {
		outputFile := filepath.Join(c.outputDir, name)
		ioutil.WriteFile(outputFile, []byte(content), 0600)
	}

	return nil
}

func (c *renderContext) addManifestFiles(name ...string) {
	c.manifestFiles = append(c.manifestFiles, name...)
}

func (c *renderContext) addManifest(name, content string) {
	c.manifests[name] = content
}

func (c *renderContext) substituteParams(data interface{}, fileName string) (string, error) {
	out := &bytes.Buffer{}
	asset := assets.MustAsset(fileName)
	t := template.Must(template.New("template").Funcs(c.funcs).Parse(string(asset)))
	err := t.Execute(out, data)
	if err != nil {
		panic(err.Error())
	}
	return out.String(), nil
}
