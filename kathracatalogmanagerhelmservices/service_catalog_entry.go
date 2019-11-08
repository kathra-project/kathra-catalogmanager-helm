package kathracatalogmanagerhelmservices

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
	"golang.org/x/net/html"
	"gopkg.in/yaml.v2"
)

func GetAllCatalogServicesImpl(w http.ResponseWriter, r *http.Request) {
	var entries, err = getAllCatalogEntries()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(entries)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(b))
}

func GetCatalogEntryVersionsImpl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var entries, err = getAllCatalogEntryVersions(vars["providerId"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(entries)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(b))
}

func GetCatalogEntryFromVersionImpl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var entries, err = GetCatalogEntryFromProviderId(vars["providerId"], vars["version"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(entries)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(b))
}

func GetCatalogEntryImpl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var entry, err = GetCatalogEntryFromProviderId(vars["providerId"], "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if entry.ProviderId == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := json.Marshal(entry)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(b))
}

func GetCatalogEntryFromProviderId(providerId string, version string) (CatalogEntry, error) {
	var helmEntries []HelmEntry
	var err error
	if version != "" {
		helmEntries, err = HelmSearchFromVersionInMemory(providerId, version)
	} else {
		helmEntries, err = HelmSearchInMemory(providerId, false)
	}
	if err != nil || len(helmEntries) != 1 {
		return CatalogEntry{}, err
	}
	var catalogEntry = convertHelmEntryToCatalogEntry(helmEntries[0])
	var args, err2 = getArgumentsFromChart(catalogEntry)
	if err2 != nil {
		return CatalogEntry{}, err2
	}
	catalogEntry.Arguments = args

	var readmePath, readmeErr = helmGetFileFromChart(catalogEntry.ProviderId, catalogEntry.Version, "README.md")
	if readmeErr != nil {
		fmt.Printf("error: %v \n", readmeErr)
	}
	if readmePath != "" {
		readmeContent, readmeContentErr := ioutil.ReadFile(readmePath)
		if readmeContentErr != nil {
			fmt.Printf("error: %v \n", readmeContentErr)
		}
		catalogEntry.Documentation = string(readmeContent)
	}
	return catalogEntry, nil
}

func AddCatalogEntryFromTemplateImpl(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t Template
	err := decoder.Decode(&t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var entry, errorCreate = createCatalogEntryFromTemplate(t)
	if errorCreate != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	b, err := json.Marshal(entry)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(b))
}

func createCatalogEntryFromTemplate(template Template) (CatalogEntry, error) {

	var tIsValidErr = template.templateIsValid()
	if tIsValidErr != nil {
		return CatalogEntry{}, tIsValidErr
	}
	var chartDirectory, err = generateFilesFromTemplate(template)
	if err != nil {
		return CatalogEntry{}, err
	}

	var catalogRepository = getKathraCatalogRepository()
	var errPush = pushIntoChartMuseum(catalogRepository, chartDirectory)
	if errPush != nil {
		return CatalogEntry{}, errPush
	}

	var entry = CatalogEntry{
		Name:        template.getValueFromKey("CHART_NAME"),
		Version:     template.getValueFromKey("CHART_VERSION"),
		Description: template.getValueFromKey("CHART_DESCRIPTION"),
		ProviderId:  template.getValueFromKey("CHART_NAME")}

	var repoName, err2 = HelmFindLocalRepository(catalogRepository)
	if err2 != nil {
		return entry, err
	}
	var chartExist, err3 = helmUSearchIfChartExist(repoName, entry.Name, entry.Version)
	if err3 != nil {
		return entry, err3
	}
	if !chartExist {
		return entry, errors.New("Unable to find chart pushed")
	}
	return entry, nil
}

func getAllCatalogEntries() ([]CatalogEntry, error) {
	var helmEntries, err = HelmSearchInMemory("", false)
	catalogEntries := []CatalogEntry{}
	if err != nil {
		return catalogEntries, err
	}
	for i := range helmEntries {
		catalogEntries = append(catalogEntries, convertHelmEntryToCatalogEntry(helmEntries[i]))
	}
	return catalogEntries, nil
}

func getAllCatalogEntryVersions(providerId string) ([]CatalogEntry, error) {
	var helmEntries, err = HelmSearchInMemory(providerId, true)
	catalogEntries := []CatalogEntry{}
	if err != nil {
		return catalogEntries, err
	}
	for i := range helmEntries {
		catalogEntries = append(catalogEntries, convertHelmEntryToCatalogEntry(helmEntries[i]))
	}
	return catalogEntries, nil
}

func convertHelmEntryToCatalogEntry(helmEntry HelmEntry) CatalogEntry {
	return CatalogEntry{Name: helmEntry.Name, Description: helmEntry.Description, Repository: helmEntry.RepositoryURL, ProviderId: helmEntry.LocalName, Version: helmEntry.VersionChart}
}

func getArgumentsFromChart(catalogEntry CatalogEntry) ([]CatalogEntryArgument, error) {
	var questionFile, err = helmGetFileFromChart(catalogEntry.ProviderId, catalogEntry.Version, "questions.yml")
	if err != nil {
		return nil, err
	}
	if questionFile != "" {
		return convertQuestionFileYamlToCatalogEntryArgumentArray(questionFile)
	}
	var readmeFile, err2 = helmGetFileFromChart(catalogEntry.ProviderId, catalogEntry.Version, "README.md")
	if err2 != nil {
		return nil, err2
	}
	if readmeFile != "" {
		return convertReadmeToCatalogEntryArgumentArray(readmeFile)
	}
	return []CatalogEntryArgument{}, nil
}

type rancherQuestionFile struct {
	Categories []string   `yaml:"categories,omitempty"`
	Questions  []question `yaml:"questions,omitempty"`
}

type question struct {
	Variable          string     `yaml:"variable,omitempty"`
	DefaultV          string     `yaml:"default,omitempty"`
	TypeV             string     `yaml:"type,omitempty"`
	Min               int        `yaml:"min,omitempty"`
	Max               int        `yaml:"max,omitempty"`
	Label             string     `yaml:"label,omitempty"`
	Group             string     `yaml:"group,omitempty"`
	Description       string     `yaml:"description,omitempty"`
	ShowSubquestionIf string     `yaml:"show_subquestion_if,omitempty"`
	Subquestions      []question `yaml:"subquestions,omitempty"`
	Options           []string   `yaml:"options,omitempty"`
	Required          bool       `yaml:"required,omitempty"`
}

func convertQuestionFileYamlToCatalogEntryArgumentArray(questionFile string) ([]CatalogEntryArgument, error) {
	catalogEntries := []CatalogEntryArgument{}
	var rancherQuestion rancherQuestionFile
	reader, err := os.Open(questionFile)
	if err != nil {
		log.Printf("error: %v", err)
	}
	buf, _ := ioutil.ReadAll(reader)
	errUnmarshal := yaml.Unmarshal([]byte(buf), &rancherQuestion)
	if errUnmarshal != nil {
		log.Printf("error: %v", errUnmarshal)
		return nil, errUnmarshal
	}

	for i := range rancherQuestion.Questions {
		var constraint = ""
		var question = rancherQuestion.Questions[i]
		switch question.TypeV {
		case "int":
			constraint = "^[0-9]+$"
			break
		case "string":
			constraint = "^.+$"
			break
		case "hostname":
			constraint = "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$"
			break
		case "boolean":
			constraint = "^(true|false)$"
			break
		case "enum":
			constraint = "^(" + strings.Join(question.Options, "|") + ")$"
			break
		default:
			constraint = "^.+$"
			break
		}
		catalogEntries = append(catalogEntries, CatalogEntryArgument{Label: question.Label, Description: question.Description, Key: question.Variable, Value: question.DefaultV, Contrainst: constraint})
	}

	return catalogEntries, nil
}

func convertReadmeToCatalogEntryArgumentArray(readmeFile string) ([]CatalogEntryArgument, error) {
	catalogEntries := []CatalogEntryArgument{}
	reader, err := os.Open(readmeFile)
	if err != nil {
		log.Printf("error: %v", err)
	}
	buf, _ := ioutil.ReadAll(reader)
	outputHTML := markdown.ToHTML([]byte(buf), nil, nil)
	var doc, errParse = htmlquery.Parse(strings.NewReader(string(outputHTML)))
	if errParse != nil {
		log.Printf("error: %v", errParse)
	}
	var listTh = htmlquery.Find(doc, "/html/body/table/thead/tr/th")
	var thParameter *html.Node
	for iTh := range listTh {
		if htmlquery.InnerText(listTh[iTh]) == "Parameter" {
			thParameter = listTh[iTh]
			break
		}
	}
	if thParameter != nil {
		booleanRgx, _ := regexp.Compile("^(true|false)$")
		integerRgx, _ := regexp.Compile("^[0-9]+$")
		linesParameters := htmlquery.Find(thParameter, "../../../tbody/tr")
		for iTr := range linesParameters {
			parametersTd := htmlquery.Find(linesParameters[iTr], "td")
			if len(parametersTd) < 3 {
				continue
			}
			var name = htmlquery.InnerText(parametersTd[0])
			var description = htmlquery.InnerText(parametersTd[1])
			var defaultValue = htmlquery.InnerText(parametersTd[2])
			var constraint = ""
			if defaultValue == "nil" {
				defaultValue = ""
			}
			if booleanRgx.MatchString(defaultValue) {
				constraint = "^(true|false)$"
			} else if integerRgx.MatchString(defaultValue) {
				constraint = "^[0-9]+$"
			} else {
				constraint = "^.*$"
			}
			catalogEntries = append(catalogEntries, CatalogEntryArgument{Label: name, Description: description, Key: name, Value: defaultValue, Contrainst: constraint})
		}
	}
	return catalogEntries, nil
}
