package kathracatalogmanagerhelmservices

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sort"
)

type TemplateService interface {
	getTemplates() []Template
	templateIsValid() (bool, error)
	generateFilesFromTemplate(t Template) (string, error)
}

func (t Template) getValueFromKey(key string) string {
	sort.Slice(t.Arguments, func(i, j int) bool {
		return t.Arguments[i].Key <= t.Arguments[j].Key
	})
	iKey := sort.Search(len(t.Arguments), func(i int) bool {
		return string(t.Arguments[i].Key) >= key
	})
	return t.Arguments[iKey].Value
}
func (t Template) getConstraintFromKey(key string) string {
	sort.Slice(t.Arguments, func(i, j int) bool {
		return t.Arguments[i].Key <= t.Arguments[j].Key
	})
	iKey := sort.Search(len(t.Arguments), func(i int) bool {
		return string(t.Arguments[i].Key) >= key
	})
	return t.Arguments[iKey].Contrainst
}

// Return collection of @Template
func getTemplates() []Template {
	var chartName = TemplateArgument{Key: "CHART_NAME", Value: "", Contrainst: "[A-Za-z0-9]"}
	var chartVersion = TemplateArgument{Key: "CHART_VERSION", Value: "", Contrainst: "[0-9]+\\.[0-9]+\\.[0-9]+"}
	var chartDescription = TemplateArgument{Key: "CHART_DESCRIPTION", Value: "", Contrainst: "[A-Za-z0-9]"}
	var appVersion = TemplateArgument{Key: "APP_VERSION", Value: "", Contrainst: "[0-9]+\\.[0-9]+\\.[0-9]+"}
	var imageName = TemplateArgument{Key: "IMAGE_NAME", Value: ".+"}
	var imageTag = TemplateArgument{Key: "IMAGE_TAG", Value: ".+"}
	var registryHost = TemplateArgument{Key: "REGISTRY_HOST", Value: ".+"}
	var restServiceTemplate = Template{Name: "RestApiService", Arguments: []TemplateArgument{chartName, chartVersion, chartDescription, appVersion, imageName, imageTag, registryHost}}
	return []Template{restServiceTemplate}
}

func GetTemplatesImpl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	b, err := json.Marshal(getTemplates())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(b))
}

func (templateToCheck Template) templateIsValid() error {
	var templates = getTemplates()
	for _, template := range templates {
		// find template
		if template.Name == templateToCheck.Name {
			for _, arg := range template.Arguments {
				if arg.Contrainst == "" {
					continue
				}
				regexContrainst, _ := regexp.Compile(arg.Contrainst)
				var valueSetted string = templateToCheck.getValueFromKey(arg.Key)
				if !regexContrainst.MatchString(valueSetted) {
					return errors.New("Template '" + template.Name + "' is not valid: argument '" + arg.Key + "' doesn't respect contrainst '" + arg.Contrainst + "'. Value defined : " + valueSetted)
				}
			}
			return nil
		}
	}
	return errors.New("Template '" + templateToCheck.Name + "' not found ")
}

func generateFilesFromTemplate(t Template) (string, error) {
	dir, err := ioutil.TempDir(os.TempDir(), "kathra-catalogmanager-")
	if err != nil {
		log.Println(err)
		return "", err
	}
	var dirWithSrc = dir + "/" + t.getValueFromKey("CHART_NAME")
	Dir("./templates/"+t.Name, dirWithSrc)
	for _, arg := range t.Arguments {
		cmd := exec.Command("/bin/bash", "-c", "find "+dirWithSrc+" -type f -exec sed -i -e 's/${"+arg.Key+"}/"+arg.Value+"/g' {} \\;")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Println(err)
			return "", err
		}
	}
	println("Template generated in " + dirWithSrc)
	return dirWithSrc, nil
}

// File copies a single file from src to dst
func File(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

// Dir copies a whole directory recursively
func Dir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = Dir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = File(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
