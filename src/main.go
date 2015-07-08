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

const activateTmplSrc = `#!/bin/sh
function deactivate() {
  if [ -n "$_OLD_GOPJX_PATH" ]; then
    export PATH="$_OLD_GOPJX_PATH"
    unset _OLD_GOPJX_PATH
  fi

  if [ -n "$_OLD_GOPJX_GOPATH" ]; then
    export GOPATH="$_OLD_GOPJX_GOPATH"
    unset _OLD_GOPJX_GOPATH
  fi

  if [ -n "$_OLD_GOPJX_PS1" ]; then
    export PS1="$_OLD_GOPJX_PS1"
    unset _OLD_GOPJX_PS1
  fi

  if [ "$1" != "nondestructive" ]; then
    unset -f deactivate
  fi
}

deactivate nondestructive

export _OLD_GOPJX_PATH="$PATH"
export _OLD_GOPJX_GOPATH="$GOPATH"
export _OLD_GOPJX_PS1="$PS1"
export PS1="[go:{{.Name}}] $_OLD_GOPJX_PS1"
export GOPATH={{.GoPath}}
export PATH=$GOPATH/bin:$PATH
`

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