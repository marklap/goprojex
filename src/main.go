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

// Version is the current semver of this project
const Version = "0.1.0"
const goprojexDir = ".gopjx"

var workspaceSkel = []string{
	goprojexDir,
	"scripts",
	"src",
}

var srcSkel = []string{
	"build",
	"cmd",
	"configs",
	"docs",
	"examples",
	"scripts",
	"test",
	"vendor",
}

func workspaceExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func projectName(name, path string) string {
	if name == "" {
		return filepath.Base(path)
	}
	return name
}

func workspacePath(path string) (string, error) {
	if path == "" {
		cwd, err := os.Getwd()
		if err != nil {
			cwd = "."
		}

		path = cwd
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return path, nil
}

func sourcePath(wsPath, srcPath string) string {
	if srcPath == "" {
		parent := filepath.Dir(wsPath)
		grandParent := filepath.Dir(parent)
		srcPath = filepath.Join(parent, grandParent)
	}
	return filepath.Clean(filepath.Join(wsPath, "src", srcPath))
}

func createWorkspace(path, name string) error {
	wsExists, err := workspaceExists(path)
	if err != nil {
		return err
	}
	if wsExists {
		return fmt.Errorf("ERROR: Workspace already exists")
	}

	for _, p := range workspaceSkel {
		if err := os.Mkdir(filepath.Join(path, p), os.FileMode(0775)); err != nil {
			return err
		}
	}

	err = createActivateScript(
		template.Must(template.New("activate").Parse(activateTmplSrc)),
		path,
		name,
	)
	if err != nil {
		fmt.Println("Failed to create activate script")
		return err
	}

	return nil
}

func createActivateScript(tmpl *template.Template, path, name string) error {
	script, err := os.Create(filepath.Join(path, goprojexDir))
	if err != nil {
		return err
	}
	defer script.Close()

	err = script.Chmod(0755)
	if err != nil {
		return err
	}

	vars := struct {
		Name   string
		GoPath string
	}{
		name,
		path,
	}

	err = tmpl.Execute(script, vars)
	if err != nil {
		return err
	}

	return nil
}

func createSrcTree(path string) error {
	for _, p := range srcSkel {
		if err := os.Mkdir(filepath.Join(path, p), os.FileMode(0775)); err != nil {
			return err
		}
	}
	return nil
}

func goprojex(wsPath, srcPath, pName string) (err error) {
	wsPath, err = workspacePath(wsPath)
	if err != nil {
		fmt.Println("ERROR: Failed to determine workspace path")
		return err
	}

	srcPath = sourcePath(wsPath, srcPath)

	pName = projectName(pName, wsPath)

	err = createWorkspace(wsPath, pName)
	if err != nil {
		fmt.Println("Failed to create workspace")
		return err
	}

	err = createSrcTree(srcPath)
	if err != nil {
		fmt.Println("Failed to create src tree")
		return err
	}

	return nil
}

func main() {

	workspacePath := flag.String("path", "", "Path to where the workspace will be created (othwerwise CWD)")
	sourcePath := flag.String("src", "", "A customer path to use for the src tree (otherwise determined from workspace parent directories")
	projectName := flag.String("name", "", "A custom project name to display in the prompt (otherwise base of the workspace path)")
	showVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(Version)
		return
	}

	err := goprojex(*workspacePath, *sourcePath, *projectName)
	if err != nil {
		fmt.Println("Failed to create Go project")
		fmt.Println(err)
		return
	}

}
