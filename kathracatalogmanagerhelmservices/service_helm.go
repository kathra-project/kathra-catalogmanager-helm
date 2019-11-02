package kathracatalogmanagerhelmservices

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
	"os/exec"
	"strings"
)

func getHelmCatalogRepository() HelmCatalogRepository {
	var repoName = os.Getenv("KATHRA_REPO_NAME")
	if repoName == "" {
		repoName = "kathra-local"
	}
	return HelmCatalogRepository{
		Name:     os.Getenv("KATHRA_REPO_NAME"),
		Url:      os.Getenv("KATHRA_REPO_URL"),
		Username: os.Getenv("KATHRA_REPO_CREDENTIAL_ID"),
		Password: os.Getenv("KATHRA_REPO_SECRET")}
}

func pushIntoChartMuseum(catalogRepository HelmCatalogRepository, chartDirectory string) error {

	var localRepositoryName, err = helmFindLocalRepository(catalogRepository)
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
	cmd := exec.Command("/bin/bash", "-c", "cd "+chartDirectory+" && helm package . ")
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
	cmd := exec.Command("helm", "update")
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
	Name         string
	VersionChart string
	VersionApp   string
	Description  string
}

var chartDownloadCacheDirectory = os.TempDir() + "/kathra-catalogmanager-helm/cacheChart"

func HelmSearch(chartName string) ([]HelmEntry, error) {

	entries := []HelmEntry{}
	cmd := exec.Command("/bin/bash", "-c", "helm search "+chartName+" | tail -n +2")
	stdout, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(stdout)

	cmd.Start()
	for scanner.Scan() {
		ucl := scanner.Text()
		var lineSplitted = strings.Split(ucl, "\t")
		var entry = HelmEntry{Name: strings.TrimSpace(lineSplitted[0]), VersionApp: strings.TrimSpace(lineSplitted[2]), VersionChart: strings.TrimSpace(lineSplitted[1]), Description: strings.TrimSpace(lineSplitted[3])}
		entries = append(entries, entry)
	}
	cmd.Wait()

	return entries, nil
}

func helmUSearchIfChartExist(repositoryName string, chartName string, chartVersion string) (bool, error) {
	cmd := exec.Command("/bin/bash", "-c", "helm search "+repositoryName+" -l  | awk '{if (($1 == \""+repositoryName+"/"+chartName+"\") && ($2 == \""+chartVersion+"\")) {print $1;}}'")
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
	cmd := exec.Command("/bin/bash", "-c", "cd "+chartDirectory+" && helm push . "+repositoryName)
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

func helmFindLocalRepository(catalogRepository HelmCatalogRepository) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", "helm repo list  | awk '{if ($2 == \""+catalogRepository.Url+"\") {print $1;}}'")
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
		cmdAddRepo := exec.Command("/bin/bash", "-c", "helm repo add "+catalogRepository.Name+" --username="+catalogRepository.Username+" --password="+catalogRepository.Password+" "+catalogRepository.Url+"")
		errAddRepo := cmdAddRepo.Run()
		cmdAddRepo.Stderr = &stdErr
		if errAddRepo != nil {
			println("Err: " + stdErr.String())
			log.Println(errAddRepo)
			return "", err
		}
		repoName = catalogRepository.Name
	}
	return repoName, nil
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
		cmd := exec.Command("/bin/bash", "-c", "cd "+chartDownloadCacheDirectory+" && helm fetch "+chartName+" "+versionOpt+" --untar --untardir "+directoryChart+" ")
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
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
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
