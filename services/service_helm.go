package services

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

var repositories = []HelmCatalogRepository{}

func HelmInitKathraRepository() {
	var err error
	repositories, err = getAllRepositoriesFromSettings()
	if err != nil {
		log.Panic("Err: Get all repository from config file")
	}
	for repo := range repositories {
		var repoLocalName, err = HelmFindLocalRepository(repositories[repo])
		if err != nil {
			log.Panic("Err: Unable to configure local chart repository " + repositories[repo].Name + " " + repositories[repo].Url)
		} else {
			log.Println("Info: Repository " + repositories[repo].Url + " configured on Helm, local-name : " + repoLocalName)
		}
	}
}

func getAllRepositoriesFromSettings() ([]HelmCatalogRepository, error) {
	var configFile = os.Getenv("REPOSITORIES_CONFIG")
	if configFile == "" {
		configFile = "repositories.yaml"
	}
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("read kubeconfig: %v", err)
	}
	var config []HelmCatalogRepository
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("unmarshal kubeconfig: %v", err)
	}
	return config, nil
}

func getKathraCatalogRepository() HelmCatalogRepository {
	for repo := range repositories {
		if repositories[repo].Name == "kathra" {
			return repositories[repo]
		}
	}
	log.Panic("Err: Unable to configure local chart repository with name 'kathra'")
	return HelmCatalogRepository{}
}

func pushIntoChartMuseum(catalogRepository HelmCatalogRepository, chartDirectory string) error {

	var localRepositoryName, err = HelmFindLocalRepository(catalogRepository)
	if err != nil {
		return err
	}
	err = helmPackage(chartDirectory)
	if err != nil {
		return err
	}
	err = helmPush(chartDirectory, localRepositoryName)
	if err != nil {
		return err
	}
	return nil
}

func helmPackage(chartDirectory string) error {
	cmd := exec.Command("/bin/bash", "-c", "cd "+chartDirectory+" && "+helmBinary+" package . ")
	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr
	err := cmd.Run()
	if err != nil {
		log.Println(stdErr.String())
		log.Println(err)
		return err
	}
	return nil
}

func HelmUpdate() error {
	cmd := exec.Command(helmBinary, "repo", "update")
	var stdErr, stdOut bytes.Buffer
	cmd.Stderr = &stdErr
	cmd.Stdout = &stdOut
	err := cmd.Run()
	log.Println(stdOut.String())
	if err != nil {
		log.Println(stdErr.String())
		log.Println(err)
		return err
	}
	return nil
}

type HelmEntry struct {
	Name          string
	LocalName     string
	VersionChart  string
	VersionApp    string
	Description   string
	RepositoryURL string
}
type HelmCatalogRepository struct {
	Name     string `json:"name,omitempty"`
	Url      string `json:"url,omitempty"`
	Username string `json:"usernmae,omitempty"`
	Password string `json:"password,omitempty"`
}

var entriesCached = []HelmEntry{}
var entriesAllVersionsCached = []HelmEntry{}
var helmBinary = "helm3"

var chartDownloadCacheDirectory = os.TempDir() + "/kathra-catalogmanager-helm/cacheChart"

func helmSearch(searchOpt string) ([]HelmEntry, error) {

	entries := []HelmEntry{}
	cmd := exec.Command("/bin/bash", "-c", helmBinary+" search repo "+searchOpt+" | tail -n +2")
	stdout, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(stdout)

	cmd.Start()
	for scanner.Scan() {
		ucl := scanner.Text()
		var lineSplitted = strings.Split(ucl, "\t")
		if strings.TrimSpace(lineSplitted[0]) == "" {
			continue
		}
		var nameSplited = strings.Split(strings.TrimSpace(lineSplitted[0]), "/")
		if len(nameSplited) == 1 {
			log.Println("warn: unable to parse line : " + ucl)
		}
		var entry = HelmEntry{Name: nameSplited[1], LocalName: strings.TrimSpace(lineSplitted[0]), VersionApp: strings.TrimSpace(lineSplitted[2]), VersionChart: strings.TrimSpace(lineSplitted[1]), Description: strings.TrimSpace(lineSplitted[3])}
		entry.RepositoryURL = helmFindHelmRepositoryFromChartName(entry.LocalName).Url
		entries = append(entries, entry)
	}
	cmd.Wait()

	return entries, nil
}

func helmFindHelmRepositoryFromChartName(chartName string) HelmCatalogRepository {
	var lineSplitted = strings.Split(chartName, "/")
	var localRepoIdentifier = lineSplitted[0]

	for i := range repositories {
		if localRepoIdentifier == repositories[i].Name {
			return repositories[i]
		}
	}
	return HelmCatalogRepository{}
}

func HelmRepoList() ([]HelmCatalogRepository, error) {

	repositories := []HelmCatalogRepository{}
	cmd := exec.Command("/bin/bash", "-c", helmBinary+" repo list | tail -n +2")
	stdout, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(stdout)

	cmd.Start()
	for scanner.Scan() {
		ucl := scanner.Text()
		var lineSplitted = strings.Split(ucl, "\t")
		var repo = HelmCatalogRepository{Name: strings.TrimSpace(lineSplitted[0]), Url: strings.TrimSpace(lineSplitted[1])}
		repositories = append(repositories, repo)
	}
	cmd.Wait()

	return repositories, nil
}

func HelmSearchInMemory(localName string, allversion bool) ([]HelmEntry, error) {
	var allEntries []HelmEntry
	if allversion {
		allEntries = append(entriesAllVersionsCached)
	} else {
		allEntries = append(entriesCached)
	}
	log.Println("localName: " + localName)
	if localName != "" {
		var entriesFiltered = []HelmEntry{}
		for i := range allEntries {
			if allEntries[i].LocalName == localName {
				entriesFiltered = append(entriesFiltered, allEntries[i])
			}
		}
		return entriesFiltered, nil
	}
	return allEntries, nil
}

func HelmSearchFromVersionInMemory(localName string, version string) ([]HelmEntry, error) {
	var entriesFiltered = []HelmEntry{}
	var allEntries = append(entriesAllVersionsCached)
	log.Println("localName: " + localName)
	for i := range allEntries {
		if allEntries[i].LocalName == localName && allEntries[i].VersionChart == version {
			entriesFiltered = append(entriesFiltered, allEntries[i])
		}
	}
	return entriesFiltered, nil
}

func HelmLoadAllInMemory() {

	var repoList, err3 = HelmRepoList()
	if err3 != nil {
		log.Println(err3)
	} else {
		repositories = repoList
	}

	var found, err = helmSearch("")
	if err != nil {
		log.Println(err)
	} else {
		entriesCached = found
	}
	var foundAllVersion, err2 = helmSearch("-l")
	if err2 != nil {
		log.Println(err2)
	} else {
		entriesAllVersionsCached = foundAllVersion
	}

}

func helmUSearchIfChartExist(repositoryName string, chartName string, chartVersion string) (bool, error) {
	cmd := exec.Command("/bin/bash", "-c", helmBinary+" search repo "+repositoryName+" -l  | awk '{if (($1 == \""+repositoryName+"/"+chartName+"\") && ($2 == \""+chartVersion+"\")) {print $1;}}'")
	var out bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stdErr
	err := cmd.Run()
	if err != nil {
		println("Err: " + stdErr.String())
		log.Println(err)
		return false, err
	}
	var chartFound = out.String()
	return chartFound == repositoryName+"/"+chartName, nil
}
func helmPush(chartDirectory string, repositoryName string) error {
	cmd := exec.Command("/bin/bash", "-c", "cd "+chartDirectory+" && "+helmBinary+" push . "+repositoryName)
	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr
	err := cmd.Run()
	if err != nil {
		log.Println(stdErr.String())
		log.Println(err)
		return err
	}
	return nil
}

func HelmFindLocalRepository(catalogRepository HelmCatalogRepository) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", helmBinary+" repo list  | awk '{if ($2 == \""+catalogRepository.Url+"\") {print $1;}}'")
	var out bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stdErr
	err := cmd.Run()
	if err != nil {
		println("Err: " + stdErr.String())
		log.Println(err)
		return "", err
	}
	var repoName = out.String()
	if repoName == "" {
		println("Unable to find repository with url : " + catalogRepository.Url + ", add new repository ")
		var errAddRepo = HelmAddRepository(catalogRepository)
		if errAddRepo != nil {
			log.Println(errAddRepo)
			return "", errAddRepo
		}
		repoName = catalogRepository.Name
	}
	return repoName, nil
}

func HelmAddRepository(catalogRepository HelmCatalogRepository) error {
	cmd := helmBinary + " repo add " + catalogRepository.Name + " " + catalogRepository.Url + ""
	if catalogRepository.Username != "" {
		cmd = cmd + " --username=" + catalogRepository.Username
	}
	if catalogRepository.Password != "" {
		cmd = cmd + " --password=" + catalogRepository.Password
	}

	cmdAddRepo := exec.Command("/bin/bash", "-c", cmd)
	errAddRepo := cmdAddRepo.Run()
	var stdErr bytes.Buffer
	cmdAddRepo.Stderr = &stdErr
	if errAddRepo != nil {
		println("Err: " + stdErr.String())
		log.Println(errAddRepo)
		return errAddRepo
	}
	return nil
}

func helmDownloadChart(chartName string, chartVersion string) (string, error) {
	if !exists(chartDownloadCacheDirectory) {
		os.MkdirAll(chartDownloadCacheDirectory, os.ModePerm)
	}

	hasher := md5.New()
	hasher.Write([]byte(chartName + chartVersion))
	var directoryChart = chartDownloadCacheDirectory + "/" + hex.EncodeToString(hasher.Sum(nil))

	if !exists(directoryChart) {
		var versionOpt = ""
		if chartVersion != "" {
			versionOpt = "--version=\"" + chartVersion + "\""
		}
		println("Download chart " + chartName + "@" + chartVersion + " into directory " + directoryChart)
		cmd := exec.Command("/bin/bash", "-c", "cd "+chartDownloadCacheDirectory+" && "+helmBinary+" fetch "+chartName+" "+versionOpt+" --untar --untardir "+directoryChart+" ")
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Println(err)
			return "", err
		}
	} else {
		println("Chart " + chartName + "@" + chartVersion + " already exist in cache " + directoryChart)
	}
	return directoryChart, nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func helmGetFileFromChart(chartName string, chartVersion string, filePath string) (string, error) {

	var directoryChart, err = helmDownloadChart(chartName, chartVersion)
	if err != nil {
		log.Println(err)
		return "", err
	}

	cmdFindFile := exec.Command("/bin/bash", "-c", "cd "+directoryChart+" && find . -name \""+filePath+"\" | head -n 1 ")
	var out bytes.Buffer
	cmdFindFile.Stdout = &out
	if err := cmdFindFile.Run(); err != nil {
		log.Println(err)
		return "", err
	}
	if out.String() == "" {
		return "", nil
	}
	return strings.TrimSpace(directoryChart + "/" + out.String()), nil
}
