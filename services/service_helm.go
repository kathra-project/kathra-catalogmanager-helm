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
	"sync"

	"gopkg.in/yaml.v2"
)

type HelmService struct {
	Repositories []*HelmCatalogRepository
}

var helmServiceInstance *HelmService
var onceHelmServiceInstance sync.Once

func GetHelmServiceInstance() *HelmService {
	onceHelmServiceInstance.Do(func() {
		helmServiceInstance = NewHelmService()
	})
	return helmServiceInstance
}

func NewHelmService() *HelmService {
	var svc = HelmService{Repositories: []*HelmCatalogRepository{}}
	var repositoriesFromFiles, errFromFiles = getAllRepositoriesFromSettings()
	if errFromFiles != nil {
		log.Panic("Err: Get all repository from config file")
	}
	for i := range repositoriesFromFiles {
		svc.initRepository(&repositoriesFromFiles[i])
		svc.Repositories = append(svc.Repositories, &repositoriesFromFiles[i])
	}
	return &svc
}

func (svc *HelmService) initRepository(repository *HelmCatalogRepository) {
	log.Println("initRepository:" + repository.Name + " -> " + repository.Url)
	var repoLocalName, err = svc.HelmFindLocalRepository(repository)
	if err != nil {
		log.Panic("Err: Unable to configure local chart repository " + repository.Name + " " + repository.Url)
	} else {
		log.Println("Info: Repository " + repository.Url + " configured on Helm, local-name : " + repoLocalName)
	}
}

func (svc *HelmService) UpdateFromResourceManager() {

	var repositoriesFromResourceManager, errResourceManager = getHelmRepositoryFromResourceManager()
	if errResourceManager != nil {
		log.Panic("Err: Get all repository from resourcemanager")
	}
	for i := range repositoriesFromResourceManager {
		var existing = false
		for z := range svc.Repositories {
			if svc.Repositories[z].Name == repositoriesFromResourceManager[i].Name {
				existing = true
				break
			}
		}
		if !existing {
			svc.initRepository(&repositoriesFromResourceManager[i])
			svc.Repositories = append(svc.Repositories, &repositoriesFromResourceManager[i])
		}
	}
}

func getHelmRepositoryFromResourceManager() ([]HelmCatalogRepository, error) {
	var helmCatalogRepositories = []HelmCatalogRepository{}

	var resourceManagerService = NewResourceManagerService()
	var binaryRepositoriesResourceManager = resourceManagerService.getBinaryRepositories()
	for i := range binaryRepositoriesResourceManager {
		if binaryRepositoriesResourceManager[i].Type != "HELM" {
			continue
		}
		if binaryRepositoriesResourceManager[i].URL == "" {
			log.Println("BinaryRepository:" + binaryRepositoriesResourceManager[i].Resource.Name + " doesn't have URL")
			continue
		}
		if binaryRepositoriesResourceManager[i].Group == nil {
			log.Println("BinaryRepository:" + binaryRepositoriesResourceManager[i].Resource.Name + " doesn't have group")
			continue
		}
		var group = resourceManagerService.getGroupById(binaryRepositoriesResourceManager[i].Group.Resource.ID)
		if group.TechnicalUser == nil {
			log.Println("Group:" + group.Resource.Name + " doesn't have technicalUser")
			continue
		}

		var user = resourceManagerService.getUserById(group.TechnicalUser.Resource.ID)
		if user.Metadata["HARBOR_SECRET_CLI"] == nil {
			log.Println("User:" + user.Resource.Name + " doesn't have property HARBOR_SECRET_CLI")
			continue
		}
		var password = fmt.Sprintf("%v", user.Metadata["HARBOR_SECRET_CLI"])
		var url = binaryRepositoriesResourceManager[i].URL
		helmCatalogRepositories = append(helmCatalogRepositories, HelmCatalogRepository{Name: binaryRepositoriesResourceManager[i].Name, Url: url, Username: user.Resource.Name, Password: password})
	}

	return helmCatalogRepositories, nil
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

func (svc *HelmService) getKathraCatalogRepository() *HelmCatalogRepository {
	for repo := range svc.Repositories {
		if svc.Repositories[repo].Name == "kathra" {
			return svc.Repositories[repo]
		}
	}
	log.Panic("Err: Unable to configure local chart repository with name 'kathra'")
	return nil
}

func (svc *HelmService) pushIntoChartMuseum(catalogRepository *HelmCatalogRepository, chartDirectory string) error {

	var localRepositoryName, err = svc.HelmFindLocalRepository(catalogRepository)
	if err != nil {
		return err
	}
	err = svc.helmPackage(chartDirectory)
	if err != nil {
		return err
	}
	err = svc.helmPush(chartDirectory, localRepositoryName)
	if err != nil {
		return err
	}
	return nil
}

func (svc *HelmService) helmPackage(chartDirectory string) error {
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

func (svc *HelmService) HelmUpdate() error {
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
	Icon          string
}
type HelmCatalogRepository struct {
	Name     string `json:"name,omitempty"`
	Url      string `json:"url,omitempty"`
	Username string `json:"usernmae,omitempty"`
	Password string `json:"password,omitempty"`
}

var entriesCached = []*HelmEntry{}
var entriesAllVersionsCached = []*HelmEntry{}
var helmBinary = "helm"

var chartDownloadCacheDirectory = os.TempDir() + "/kathra-catalogmanager-helm/cacheChart"

func (svc *HelmService) helmSearch(searchOpt string) ([]*HelmEntry, error) {

	entries := []*HelmEntry{}
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
		entry.RepositoryURL = svc.helmFindHelmRepositoryFromChartName(entry.LocalName).Url
		entries = append(entries, &entry)
	}
	cmd.Wait()

	return entries, nil
}

func (svc *HelmService) helmFindHelmRepositoryFromChartName(chartName string) *HelmCatalogRepository {
	var lineSplitted = strings.Split(chartName, "/")
	var localRepoIdentifier = lineSplitted[0]

	for i := range svc.Repositories {
		if localRepoIdentifier == svc.Repositories[i].Name {
			return svc.Repositories[i]
		}
	}
	return nil
}

func HelmRepoList() ([]*HelmCatalogRepository, error) {

	repositories := []*HelmCatalogRepository{}
	cmd := exec.Command("/bin/bash", "-c", helmBinary+" repo list | tail -n +2")
	stdout, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(stdout)

	cmd.Start()
	for scanner.Scan() {
		ucl := scanner.Text()
		var lineSplitted = strings.Split(ucl, "\t")
		var repo = HelmCatalogRepository{Name: strings.TrimSpace(lineSplitted[0]), Url: strings.TrimSpace(lineSplitted[1])}
		repositories = append(repositories, &repo)
	}
	cmd.Wait()

	return repositories, nil
}

func HelmSearchInMemory(localName string, allversion bool) ([]*HelmEntry, error) {
	var allEntries []*HelmEntry
	if allversion {
		allEntries = append(entriesAllVersionsCached)
	} else {
		allEntries = append(entriesCached)
	}
	log.Println("localName: " + localName)
	if localName != "" {
		var entriesFiltered = []*HelmEntry{}
		for i := range allEntries {
			if allEntries[i].LocalName == localName {
				entriesFiltered = append(entriesFiltered, allEntries[i])
			}
		}
		return entriesFiltered, nil
	}
	return allEntries, nil
}

func HelmSearchFromVersionInMemory(localName string, version string) ([]*HelmEntry, error) {
	var entriesFiltered = []*HelmEntry{}
	var allEntries = append(entriesAllVersionsCached)
	log.Println("localName: " + localName)
	for i := range allEntries {
		if allEntries[i].LocalName == localName && allEntries[i].VersionChart == version {
			entriesFiltered = append(entriesFiltered, allEntries[i])
		}
	}
	return entriesFiltered, nil
}

func (svc *HelmService) HelmLoadAllInMemory() {

	var repoList, err3 = HelmRepoList()
	if err3 != nil {
		log.Println(err3)
	} else {
		svc.Repositories = repoList
	}

	var found, err = svc.helmSearch("")
	if err != nil {
		log.Println(err)
	} else {
		entriesCached = found
		go func() {
			for index := range entriesCached {
				entriesCached[index].Icon, _ = getIconFromChart(entriesCached[index].LocalName, entriesCached[index].VersionChart)
			}
		}()
	}
	var foundAllVersion, err2 = svc.helmSearch("-l")
	if err2 != nil {
		log.Println(err2)
	} else {
		entriesAllVersionsCached = foundAllVersion
	}

}

func (svc *HelmService) helmUSearchIfChartExist(repositoryName string, chartName string, chartVersion string) (bool, error) {
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
func (svc *HelmService) helmPush(chartDirectory string, repositoryName string) error {
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

func (svc *HelmService) HelmFindLocalRepository(catalogRepository *HelmCatalogRepository) (string, error) {
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
		var errAddRepo = svc.HelmAddRepository(catalogRepository)
		if errAddRepo != nil {
			log.Println(errAddRepo)
			return "", errAddRepo
		}
		repoName = catalogRepository.Name
	}
	return repoName, nil
}

func (svc *HelmService) HelmAddRepository(catalogRepository *HelmCatalogRepository) error {
	cmd := helmBinary + " repo add " + catalogRepository.Name + " " + catalogRepository.Url + ""
	if catalogRepository.Username != "" {
		cmd = cmd + " --username=" + catalogRepository.Username
	}
	if catalogRepository.Password != "" {
		cmd = cmd + " --password=" + catalogRepository.Password
	}

	log.Println(cmd)
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

func (svc *HelmService) helmDownloadChart(chartName string, chartVersion string) (string, error) {
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
	if err == nil || os.IsNotExist(err) == false {
		return true
	} else {
		return false
	}
}

func (svc *HelmService) helmGetFileFromChart(chartName string, chartVersion string, filePath string) (string, error) {

	var directoryChart, err = svc.helmDownloadChart(chartName, chartVersion)
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
