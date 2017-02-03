# syncthing-changes-overrider

When using folder type "Master" (also called "Send Only Folders" or `type="readonly"` in `config.xml`) on a [Syncthing](https://github.com/syncthing/syncthing) folder, sometimes there are phantom changes, triggering the "Override Changes" button to appear in the Web UI. This button usually appears if a remote device has changed the content of a "Master" folder and it triggers overriding of these changes on the remote device. Where these phantom changes originate from, I cannot say, maybe differences in file metadata (timestamps, permissions).

This program automates the overriding of these changes by using the [REST API of syncthing](https://docs.syncthing.net/dev/rest.html). It is meant to be triggered by a cronjob. Pending changes to a folder are indicated by syncthing with `needBytes`, `needDeletes`, `need*` JSON fields returned from https://docs.syncthing.net/rest/db-status-get.html. Other used API calls are

* https://docs.syncthing.net/rest/system-config-get.html
* https://docs.syncthing.net/rest/db-override-post.html

## Configuration
All settings are saved in `OverriderConfig.json` in the same directory as the binary file. A template JSON is given with [OverriderConfig-default.json](OverriderConfig-default.json). You will need to change the [API key](https://docs.syncthing.net/dev/rest.html#api-key) so that this application can use the REST API. If `overrideAllFolders` is set to `true`, changes of all folders are overridden, if set to `false` only the folders listed in `overrideFolderIds` are overridden. Other config parameters should be self-explanatory.

## Download
Managed via [GitHub releases](../../releases).

A special faceless version is available for Windows, that doesn't open a console window when `syncthing-changes-overrider.exe` is executed. Log messages are still written to `log.txt` when `logToFile` is set to `true`.