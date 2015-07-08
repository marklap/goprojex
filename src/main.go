////////////////////////////////////////////////////////////////////////////////
// The MIT License (MIT)
//
// Copyright (c) 2015 Mark LaPerriere
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
////////////////////////////////////////////////////////////////////////////////

package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "text/template"
)

const Version = "0.0.1"

var (
    version      bool
    help         bool
    project_name string
)

type activateTmplVars struct {
    Name   string
    GoPath string
}

func init() {
    flag.BoolVar(&version, "version", false, "Print version and exit")
    flag.StringVar(&project_name, "name", "", "A custom project name to display in the prompt (otherwise the base of the GOPATH is used)")
}

func filepathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    } else if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

func deriveName(gopath, name string) string {
    if name != "" {
        return name
    } else {
        return filepath.Base(gopath)
    }
}

func createGoPathDir(gopath string) error {
    exists, err := filepathExists(gopath)
    if err != nil {
        return err
    }
    if !exists {
        if err := os.MkdirAll(gopath, os.FileMode(0775)); err != nil {
            return err
        }
    }
    return nil
}

func createActivateScript(scriptPath string, tmpl *template.Template, vars *activateTmplVars) error {
    script, err := os.Create(scriptPath)
    if err != nil {
        fmt.Println("Failed to create activate script at path" + scriptPath)
        return err
    }
    defer script.Close()

    err = script.Chmod(0755)
    if err != nil {
        return err
    }

    err = tmpl.Execute(script, vars)
    if err != nil {
        return err
    }

    return nil
}

func createProj(gopath, name string) error {
    project_name := deriveName(gopath, name)

    absGopath, err := filepath.Abs(gopath)
    if err != nil {
        fmt.Println("Failed to derive absolute path from " + gopath)
        return err
    }
    activateScriptPath := filepath.Join(absGopath, "activate")

    tmplVars := &activateTmplVars{project_name, absGopath}
    activateTmpl := template.Must(template.New("activate").Parse(activateTmplSrc))

    err = createGoPathDir(absGopath)
    if err != nil {
        fmt.Println("Faile to create GOPATH dir")
        return err
    }

    err = createActivateScript(activateScriptPath, activateTmpl, tmplVars)
    if err != nil {
        fmt.Println("Failed to create activate script")
        return err
    }

    return nil
}

func main() {
    flag.Parse()
    if version {
        fmt.Println(Version)
        return
    }

    if flag.NArg() == 0 {
        fmt.Println("\nERROR: Must specify your desired GOROOT path\n")
        flag.Usage()
        return
    }

    err := createProj(flag.Arg(0), project_name)
    if err != nil {
        fmt.Println("Failed to create Go project")
        return
    }

}
