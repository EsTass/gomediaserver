package main

import (
    "archive/zip"
    //"fmt"
    "io"
    //"log"
    "os"
    "path/filepath"
    "strings"
	"os/exec"
    
    //go get "github.com/ricochet2200/go-disk-usage/du"
    "github.com/ricochet2200/go-disk-usage/du"
)

func Unzip(src string, dest string) bool {
    result := true
    var filenames []string

    r, err := zip.OpenReader(src)
    if err != nil {
        return false
    }
    defer r.Close()

    for _, f := range r.File {

        fpath := filepath.Join(dest, f.Name)

        if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
            return false
        }

        filenames = append(filenames, fpath)

        if f.FileInfo().IsDir() {
            os.MkdirAll(fpath, os.ModePerm)
            continue
        }

        if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
            return false
        }

        outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
        if err != nil {
            return false
        }

        rc, err := f.Open()
        if err != nil {
            return false
        }

        _, err = io.Copy(outFile, rc)

        outFile.Close()
        rc.Close()

        if err != nil {
            return false
        }
    }
    return result
}

func Unrar(src string, dest string) bool {
    result := false
    command := "unrar"
    args := []string{
        "t",
        src,
    }
    _, err := exec.Command( command, args... ).Output()
    if err == nil {
        args := []string{
            "x",
            src,
            dest,
        }
        _, err := exec.Command( command, args... ).Output()
        if err == nil {
            result = true
        }
    }
    
    return result
}

func Un7z(src string, dest string) bool {
    result := false
    command := "7z"
    args := []string{
        "t",
        src,
    }
    _, err := exec.Command( command, args... ).Output()
    if err == nil {
        args := []string{
            //"-o'" + dest + "'",
            "-o" + dest + "",
            "x",
            src,
        }
        _, err := exec.Command( command, args... ).Output()
        if err == nil {
            result = true
        }
    }
    
    return result
}

//FREE SPACE

func getFreeSpace( path string ) int64 {
    var result int64 = 0
    usage := du.NewDiskUsage(path)
    result = int64(usage.Available())
    return result
}