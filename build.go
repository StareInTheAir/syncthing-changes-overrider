// Part of this file where taken from Syncthing's build.go (mostly packZip)
// https://github.com/syncthing/syncthing/blob/987718baf89f78718b44f345429c9412920645d8/build.go

// +build ignore

package main

import (
	"os/exec"
	"fmt"
	"log"
	"archive/zip"
	"os"
	"path/filepath"
	"io"
	"compress/flate"
)

const (
	OutputFolder = "./bin"
	ProjectPath = "github.com/StareInTheAir/syncthing-changes-overrider/Overrider"
	BinaryName = "syncthing-changes-overrider"
)

type archiveFile struct {
	src string
	dst string
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: go run build.go <version>")
	}
	version := os.Args[1]

	os.RemoveAll(OutputFolder)
	err := os.Mkdir(OutputFolder, 0775)
	if err != nil {
		log.Fatalln(err)
	}

	// MAC
	//buildAndPackage("darwin", "386", version, false)
	buildAndPackage("darwin", "amd64", version, false)

	// LINUX
	buildAndPackage("linux", "386", version, false)
	buildAndPackage("linux", "amd64", version, false)
	buildAndPackage("linux", "arm", version, false)
	buildAndPackage("linux", "arm64", version, false)

	// WINDOWS
	buildAndPackage("windows", "386", version, false)
	buildAndPackage("windows", "amd64", version, false)
	buildAndPackage("windows", "386", version, true)
	buildAndPackage("windows", "amd64", version, true)
}

func buildAndPackage(goOs string, goArch string, version string, windowsFaceless bool) {
	if windowsFaceless {
		fmt.Printf("%s_%s_faceless: building", goOs, goArch)
	} else {
		fmt.Printf("%s_%s: building", goOs, goArch)
	}
	command := exec.Command("go", "build", "-o", getOutputName(goOs, goArch, windowsFaceless), "-ldflags", getLdFlags(windowsFaceless), ProjectPath)
	os.Setenv("GOOS", goOs)
	os.Setenv("GOARCH", goArch)
	output, err := command.CombinedOutput()
	if err != nil {
		log.Fatalf("Error while building: %s\n%s\n", err, string(output))
	}
	fmt.Print(", zipping")
	var zipName string
	if windowsFaceless {
		zipName = fmt.Sprintf("%s/%s-%s_%s_faceless-%s.zip", OutputFolder, BinaryName, goOs, goArch, version)
	} else {
		zipName = fmt.Sprintf("%s/%s-%s_%s-%s.zip", OutputFolder, BinaryName, goOs, goArch, version)
	}
	executableName := BinaryName
	if goOs == "windows" {
		executableName += ".exe"
	}
	zipFiles := []archiveFile{{getOutputName(goOs, goArch, windowsFaceless), executableName},
		{"./OverriderConfig-default.json", "OverriderConfig-default.json"}}
	packZip(zipName, zipFiles)
	fmt.Println(", done")
}

func packZip(out string, files []archiveFile) {
	fd, err := os.Create(out)
	if err != nil {
		log.Fatalln(err)
	}
	defer fd.Close()

	zw := zip.NewWriter(fd)
	zw.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})
	defer zw.Close()

	for _, f := range files {
		sf, err := os.Open(f.src)
		if err != nil {
			log.Fatalln(err)
		}

		info, err := sf.Stat()
		if err != nil {
			log.Fatalln(err)
		}

		fh, err := zip.FileInfoHeader(info)
		if err != nil {
			log.Fatalln(err)
		}
		fh.Name = filepath.ToSlash(f.dst)
		fh.Method = zip.Deflate

		of, err := zw.CreateHeader(fh)
		if err != nil {
			log.Fatalln(err)
		}
		_, err = io.Copy(of, sf)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func getOutputName(goOs string, goArch string, windowsFaceless bool) (output string) {
	if windowsFaceless {
		output = fmt.Sprintf("%s/%s_%s_faceless/%s", OutputFolder, goOs, goArch, BinaryName)
	} else {
		output = fmt.Sprintf("%s/%s_%s/%s", OutputFolder, goOs, goArch, BinaryName)
	}

	if goOs == "windows" {
		output += ".exe"
	}
	return output
}

func getLdFlags(windowsFaceless bool) (flags string) {
	flags = "-w" // no debug symbols
	if windowsFaceless {
		flags += " -H windowsgui"
	}
	return flags
}