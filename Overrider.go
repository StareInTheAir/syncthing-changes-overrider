package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"crypto/tls"
	"log"
	"os"
)

type OverriderConfig struct {
	SyncthingAddress            string
	ApiKey                      string
	IgnoreInvalidSslCertificate bool
	OverrideAllFolders          bool
	OverrideFolderIds           []string
}

type RestSystemConfig struct {
	Folders []Folder
}
type Folder struct {
	Id string
}

type RestDbFolder struct {
	NeedBytes       int
	NeedDeletes     int
	NeedDirectories int
	NeedFiles       int
	NeedSymlinks    int
}

var logOut = log.New(os.Stdout, "", log.LstdFlags)
var logErr = log.New(os.Stderr, "", log.LstdFlags)

func main() {
	var config OverriderConfig

	jsonBytes, err := ioutil.ReadFile("OverriderConfig.json")

	dieOnError(err)

	json.Unmarshal(jsonBytes, &config)
	logOut.Println(config)

	var client *http.Client
	if config.IgnoreInvalidSslCertificate {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	} else {
		client = http.DefaultClient
	}

	if config.OverrideAllFolders {
		req := createSyncthingHttpRequest(config, "GET", "/rest/system/config")
		resp, err := client.Do(req)
		dieOnError(err)
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		dieOnError(err)

		var jsonSystemConfig RestSystemConfig
		json.Unmarshal(body, &jsonSystemConfig)

		var folders []string
		for _, folder := range jsonSystemConfig.Folders {
			folders = append(folders, folder.Id)
		}
		overrideDirtySyncthingFolders(config, client, folders)
	} else {
		// only overwrite folder listed in config
		overrideDirtySyncthingFolders(config, client, config.OverrideFolderIds)
	}

}
func dieOnError(err error) {
	if err != nil {
		logErr.Fatalln(err)
	}
}

func createSyncthingHttpRequest(config OverriderConfig, method string, httpCommand string) *http.Request {
	req, err := http.NewRequest(method, config.SyncthingAddress + httpCommand, nil)
	dieOnError(err)
	req.Header.Add("X-API-Key", config.ApiKey)
	return req
}

func overrideDirtySyncthingFolders(config OverriderConfig, client *http.Client, folders []string) {
	overwroteChanges := false
	for _, folder := range folders {
		req := createSyncthingHttpRequest(config, "GET", "/rest/db/status?folder=" + folder)
		resp, err := client.Do(req)
		dieOnError(err)
		defer resp.Body.Close()
		jsonDbFolder, err := ioutil.ReadAll(resp.Body)
		var dbFolderInfo RestDbFolder
		json.Unmarshal(jsonDbFolder, &dbFolderInfo)
		if dbFolderInfo.NeedBytes > 0 || dbFolderInfo.NeedDeletes > 0 ||
			dbFolderInfo.NeedDirectories > 0 || dbFolderInfo.NeedFiles > 0 ||
			dbFolderInfo.NeedSymlinks > 0 {
			req := createSyncthingHttpRequest(config, "POST", "/rest/db/override?folder=" + folder)
			resp, err := client.Do(req)
			dieOnError(err)
			defer resp.Body.Close()
			logOut.Println("Overwrote changes of folder " + folder)
			overwroteChanges = true
		}
	}
	if !overwroteChanges {
		logOut.Println("No changes to override")
	}
}