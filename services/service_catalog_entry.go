package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/gomarkdown/markdown"
	apiModel "github.com/kathra-project/kathra-core-model-go/models"
	"golang.org/x/net/html"
	"gopkg.in/yaml.v2"
)

func GetCatalogEntryPackageVersionFromProviderId(providerId string, version string) (*apiModel.CatalogEntryPackageVersion, error) {
	var helmEntries []*HelmEntry
	var err error
	if version != "" {
		helmEntries, err = HelmSearchFromVersionInMemory(providerId, version)
	} else {
		helmEntries, err = HelmSearchInMemory(providerId, false)
	}
	if err != nil || len(helmEntries) != 1 {
		return &apiModel.CatalogEntryPackageVersion{}, err
	}
	var catalogEntryPackageVersion = convertHelmEntryToCatalogEntryPackageVersion(helmEntries[0])
	var args, err2 = getArgumentsFromChart(catalogEntryPackageVersion)
	if err2 != nil {
		return &apiModel.CatalogEntryPackageVersion{}, err2
	}
	catalogEntryPackageVersion.Arguments = args

	var readmePath, readmeErr = GetHelmServiceInstance().helmGetFileFromChart(catalogEntryPackageVersion.CatalogEntryPackage.ProviderID, catalogEntryPackageVersion.Version, "README.md")
	if readmeErr != nil {
		fmt.Printf("error: %v \n", readmeErr)
	}
	if readmePath != "" {
		readmeContent, readmeContentErr := ioutil.ReadFile(readmePath)
		if readmeContentErr != nil {
			fmt.Printf("error: %v \n", readmeContentErr)
		}
		catalogEntryPackageVersion.Documentation = string(readmeContent)
	}
	return catalogEntryPackageVersion, nil
}

/*
func CreateCatalogEntryPackageVersionFromTemplate(template PackageTemplate) (apiModel.CatalogEntryPackageVersion, error) {

	var tIsValidErr = template.templateIsValid()
	if tIsValidErr != nil {
		return apiModel.CatalogEntryPackageVersion{}, tIsValidErr
	}
	var chartDirectory, err = generateFilesFromTemplate(template)
	if err != nil {
		return apiModel.CatalogEntryPackageVersion{}, err
	}

	var catalogRepository = getKathraCatalogRepository()
	var errPush = pushIntoChartMuseum(catalogRepository, chartDirectory)
	if errPush != nil {
		return apiModel.CatalogEntryPackageVersion{}, errPush
	}
	var catalogEntry = apiModel.CatalogEntry{
		Description: template.getValueFromKey("CHART_DESCRIPTION"),
	}

	var catalogEntryPackage = apiModel.CatalogEntryPackage{
		ProviderID:   template.getValueFromKey("CHART_NAME"),
		CatalogEntry: &catalogEntry,
	}
	catalogEntryPackage.Resource.Name = template.getValueFromKey("CHART_NAME")

	var entry = apiModel.CatalogEntryPackageVersion{
		Version:             template.getValueFromKey("CHART_VERSION"),
		CatalogEntryPackage: &catalogEntryPackage}

	var repoName, err2 = HelmFindLocalRepository(catalogRepository)
	if err2 != nil {
		return entry, err
	}
	var chartExist, err3 = helmUSearchIfChartExist(repoName, entry.CatalogEntryPackage.Name, entry.Version)
	if err3 != nil {
		return entry, err3
	}
	if !chartExist {
		return entry, errors.New("Unable to find chart pushed")
	}
	return entry, nil
}
*/
func GetAllCatalogEntryPackage() ([]apiModel.CatalogEntryPackage, error) {
	var helmEntries, err = HelmSearchInMemory("", false)
	catalogEntries := []apiModel.CatalogEntryPackage{}
	if err != nil {
		return catalogEntries, err
	}
	for i := range helmEntries {
		catalogEntries = append(catalogEntries, *convertHelmEntryToCatalogEntryPackage(helmEntries[i]))
	}
	return catalogEntries, nil
}

func GetAllCatalogEntryPackageVersionVersions(providerId string) ([]apiModel.CatalogEntryPackageVersion, error) {
	var helmEntries, err = HelmSearchInMemory(providerId, true)
	catalogEntries := []apiModel.CatalogEntryPackageVersion{}
	if err != nil {
		return catalogEntries, err
	}
	for i := range helmEntries {
		catalogEntries = append(catalogEntries, *convertHelmEntryToCatalogEntryPackageVersion(helmEntries[i]))
	}
	return catalogEntries, nil
}

type ChartYaml struct {
	Icon string `yaml:"icon"`
}

func getIconFromChart(providerId string, version string) (string, error) {
	var yamlFilePath, err = GetHelmServiceInstance().helmGetFileFromChart(providerId, version, "Chart.yaml")
	if err != nil {
		return "", err
	}
	yamlFile, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		log.Printf("Unable to read Chart.yaml for %v : %v", providerId, err)
		return "", nil
	}
	var chartYaml ChartYaml
	err = yaml.Unmarshal(yamlFile, &chartYaml)
	if err != nil {
		log.Printf("Unmarshal Chart.yaml for %v : %v", providerId, err)
		return "", nil
	}
	return chartYaml.Icon, nil
}

func getArgumentsFromChart(catalogEntryPackageVersion *apiModel.CatalogEntryPackageVersion) ([]*apiModel.CatalogEntryArgument, error) {
	var questionFile, err = GetHelmServiceInstance().helmGetFileFromChart(catalogEntryPackageVersion.CatalogEntryPackage.ProviderID, catalogEntryPackageVersion.Version, "questions.yml")
	if err != nil {
		return nil, err
	}
	if questionFile != "" {
		return convertQuestionFileYamlToCatalogEntryPackageVersionArgumentArray(questionFile)
	}
	var readmeFile, err2 = GetHelmServiceInstance().helmGetFileFromChart(catalogEntryPackageVersion.CatalogEntryPackage.ProviderID, catalogEntryPackageVersion.Version, "README.md")
	if err2 != nil {
		return nil, err2
	}
	if readmeFile != "" {
		return convertReadmeToCatalogEntryPackageVersionArgumentArray(readmeFile)
	}
	return []*apiModel.CatalogEntryArgument{}, nil
}

func convertHelmEntryToCatalogEntryPackage(helmEntry *HelmEntry) *apiModel.CatalogEntryPackage {
	var catalogEntry = apiModel.CatalogEntry{}
	catalogEntry.Resource.Name = helmEntry.Name
	catalogEntry.Description = helmEntry.Description
	catalogEntry.Icon = helmEntry.Icon
	versions := []*apiModel.CatalogEntryPackageVersion{}
	versions = append(versions, &apiModel.CatalogEntryPackageVersion{Version: helmEntry.VersionChart})
	return &apiModel.CatalogEntryPackage{CatalogEntry: &catalogEntry, URL: helmEntry.RepositoryURL, ProviderID: helmEntry.LocalName, Versions: versions, PackageType: apiModel.PACKAGETYPEHELM}
}

func convertHelmEntryToCatalogEntryPackageVersion(helmEntry *HelmEntry) *apiModel.CatalogEntryPackageVersion {
	var catalogEntry = apiModel.CatalogEntry{}
	catalogEntry.Resource.Name = helmEntry.Name
	catalogEntry.Description = helmEntry.Description
	catalogEntry.Icon = helmEntry.Icon
	var catalogEntryPackage = apiModel.CatalogEntryPackage{CatalogEntry: &catalogEntry, URL: helmEntry.RepositoryURL, ProviderID: helmEntry.LocalName, PackageType: apiModel.PACKAGETYPEHELM}
	return &apiModel.CatalogEntryPackageVersion{Version: helmEntry.VersionChart, CatalogEntryPackage: &catalogEntryPackage}
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

func convertQuestionFileYamlToCatalogEntryPackageVersionArgumentArray(questionFile string) ([]*apiModel.CatalogEntryArgument, error) {
	catalogEntries := []*apiModel.CatalogEntryArgument{}
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
		catalogEntries = append(catalogEntries, &apiModel.CatalogEntryArgument{Label: question.Label, Description: question.Description, Key: question.Variable, Value: question.DefaultV, Contrainst: constraint})
	}

	return catalogEntries, nil
}

func convertReadmeToCatalogEntryPackageVersionArgumentArray(readmeFile string) ([]*apiModel.CatalogEntryArgument, error) {
	catalogEntries := []*apiModel.CatalogEntryArgument{}
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
			catalogEntries = append(catalogEntries, &apiModel.CatalogEntryArgument{Label: name, Description: description, Key: name, Value: defaultValue, Contrainst: constraint})
		}
	}
	return catalogEntries, nil
}
