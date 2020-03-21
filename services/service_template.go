package services

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sort"

	apiModel "github.com/kathra-project/kathra-core-model-go/models"
)

type TemplateService interface {
	getTemplates() []apiModel.PackageTemplate
	templateIsValid(t apiModel.PackageTemplate) (bool, error)
	generateFilesFromTemplate(t apiModel.PackageTemplate) (string, error)
	getValueFromKey(t apiModel.PackageTemplate, key string) string
	GetConstraintFromKey(t apiModel.PackageTemplate, key string) string
}

func getValueFromKey(t apiModel.PackageTemplate, key string) string {
	sort.Slice(t.Arguments, func(i, j int) bool {
		return t.Arguments[i].Key <= t.Arguments[j].Key
	})
	iKey := sort.Search(len(t.Arguments), func(i int) bool {
		return string(t.Arguments[i].Key) >= key
	})
	return t.Arguments[iKey].Value
}
func GetConstraintFromKey(t apiModel.PackageTemplate, key string) string {
	sort.Slice(t.Arguments, func(i, j int) bool {
		return t.Arguments[i].Key <= t.Arguments[j].Key
	})
	iKey := sort.Search(len(t.Arguments), func(i int) bool {
		return string(t.Arguments[i].Key) >= key
	})
	return t.Arguments[iKey].Contrainst
}

// Return collection of @Template
func getTemplates() []apiModel.PackageTemplate {
	var chartName = apiModel.PackageTemplateArgument{Key: "CHART_NAME", Value: "", Contrainst: "[A-Za-z0-9]"}
	var chartVersion = apiModel.PackageTemplateArgument{Key: "CHART_VERSION", Value: "", Contrainst: "[0-9]+\\.[0-9]+\\.[0-9]+"}
	var chartDescription = apiModel.PackageTemplateArgument{Key: "CHART_DESCRIPTION", Value: "", Contrainst: "[A-Za-z0-9]"}
	var appVersion = apiModel.PackageTemplateArgument{Key: "APP_VERSION", Value: "", Contrainst: "[0-9]+\\.[0-9]+\\.[0-9]+"}
	var imageName = apiModel.PackageTemplateArgument{Key: "IMAGE_NAME", Value: ".+"}
	var imageTag = apiModel.PackageTemplateArgument{Key: "IMAGE_TAG", Value: ".+"}
	var registryHost = apiModel.PackageTemplateArgument{Key: "REGISTRY_HOST", Value: ".+"}
	var source = apiModel.PackageTemplateArgument{Key: "SOURCE_URL", Value: ".+"}
	var website = apiModel.PackageTemplateArgument{Key: "HOME_URL", Value: ".+"}
	var icon = apiModel.PackageTemplateArgument{Key: "ICON_URL", Value: ".+"}
	var arguments = []*apiModel.PackageTemplateArgument{&chartName, &chartVersion, &chartDescription, &appVersion, &imageName, &imageTag, &registryHost, &source, &icon, &website}
	var restServiceTemplate = apiModel.PackageTemplate{Name: "RestApiService", Arguments: arguments}
	return []apiModel.PackageTemplate{restServiceTemplate}
}

func templateIsValid(templateToCheck apiModel.PackageTemplate) error {
	var templates = getTemplates()
	for _, template := range templates {
		// find template
		if template.Name == templateToCheck.Name {
			for _, arg := range template.Arguments {
				if arg.Contrainst == "" {
					continue
				}
				regexContrainst, _ := regexp.Compile(arg.Contrainst)
				var valueSetted string = getValueFromKey(templateToCheck, arg.Key)
				if !regexContrainst.MatchString(valueSetted) {
					return errors.New("Template '" + template.Name + "' is not valid: argument '" + arg.Key + "' doesn't respect contrainst '" + arg.Contrainst + "'. Value defined : " + valueSetted)
				}
			}
			return nil
		}
	}
	return errors.New("Template '" + templateToCheck.Name + "' not found ")
}

func generateFilesFromTemplate(t apiModel.PackageTemplate) (string, error) {
	dir, err := ioutil.TempDir(os.TempDir(), "kathra-catalogmanager-")
	if err != nil {
		log.Println(err)
		return "", err
	}
	var dirWithSrc = dir + "/" + getValueFromKey(t, "CHART_NAME")
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
