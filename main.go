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
	"html/template"
	"os"
	"path/filepath"
)

// Version is the current semver of this project
const Version = "0.1.1"

// GoProjexDir is the directory used for the activate script
const GoProjexDir = ".gopjx"

// DefaultFileMod is the default set of permissions
const DefaultFileMod = 0755

var workspaceSkel = []string{
	GoProjexDir,
	"scripts",
	"src",
}

var sourceSkel = []string{
	"build",
	"cmd",
	"configs",
	"docs",
	"examples",
	"scripts",
	"test",
	"vendor",
}

var errSkeletonInitNotSafe = fmt.Errorf("skeleton not safe to initialize")

// Project describes the project
type Project struct {
	Name string
}

// NewProject creates a new project value
func NewProject(name string, ws *Workspace) *Project {
	if name == "" {
		name = filepath.Base(ws.Path)
	}
	return &Project{name}
}

// Skel describes a directory with other directories in it
type Skel struct {
	Path string
	Dirs []string
}

// IsSafe determines if it's save to create the direcotry tree described by Skel
func (s Skel) IsSafe() bool {
	for _, d := range s.Dirs {
		_, err := os.Stat(filepath.Join(s.Path, d))
		if err == nil || !os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Create creates all the skeleton directories
func (s *Skel) Create() error {
	for _, d := range s.Dirs {
		path := filepath.Join(s.Path, d)
		if err := os.MkdirAll(path, os.FileMode(DefaultFileMod)); err != nil {
			return fmt.Errorf("failed to create workspace skeleton dir: %s", d)
		}
		fmt.Println("created:", path)
	}
	return nil
}

// Workspace describes the workspace Skel
type Workspace struct {
	Skel
}

// NewWorkspace sets the path depending on the supplied argument
func NewWorkspace(path string, dirs []string) (*Workspace, error) {
	if path == "" {
		cwd, err := os.Getwd()
		if err != nil {
			cwd = "."
		}

		path = cwd
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("could not determine absolute path of workspace")
	}

	return &Workspace{Skel{path, dirs}}, nil
}

// CreateActivateScript creates the activate script
func (w *Workspace) CreateActivateScript(projectName, srcPath string) error {
	tmpl, err := template.New("activate").Parse(activateTmplSrc)
	if err != nil {
		return fmt.Errorf("failed to compile activate template")
	}

	path := filepath.Join(w.Path, GoProjexDir, "activate")
	script, err := os.Create(path)
	if err != nil {
		return err
	}
	defer script.Close()

	err = script.Chmod(DefaultFileMod)
	if err != nil {
		return err
	}

	vars := struct {
		Name    string
		GoPath  string
		SrcPath string
	}{
		projectName,
		w.Path,
		srcPath,
	}

	err = tmpl.Execute(script, vars)
	if err != nil {
		return err
	}

	fmt.Println("created activate script:", path)

	return nil

}

// Source describes the source directory Skel
type Source struct {
	Skel
}

// NewSource creates a new Source skel
func NewSource(ws *Workspace, path string, dirs []string) *Source {
	if path == "" {
		workspace := filepath.Base(ws.Path)
		parent := filepath.Base(filepath.Dir(ws.Path))
		path = filepath.Join(parent, workspace)
	}
	return &Source{
		Skel{
			filepath.Join(ws.Path, "src", path),
			dirs,
		},
	}
}

// GoProjex creates the goprojex directory structure
func GoProjex(wsPath, srcPath, name string) error {
	ws, err := NewWorkspace(wsPath, workspaceSkel)
	if err != nil {
		return err
	}

	if !ws.IsSafe() {
		return fmt.Errorf("one or more workspace directories exist")
	}

	src := NewSource(ws, srcPath, sourceSkel)
	if !src.IsSafe() {
		return fmt.Errorf("one ore more workspace source directories empty")
	}

	err = ws.Create()
	if err != nil {
		return err
	}

	err = src.Create()
	if err != nil {
		return err
	}

	proj := NewProject(name, ws)
	err = ws.CreateActivateScript(proj.Name, src.Path)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	ws := flag.String("ws", "", "path to create the workspace (othwerwise CWD)")
	src := flag.String("src", "", "source path to create within the workspace (otherwise derived from workspace path")
	name := flag.String("name", "", "name to display in shell prompt (otherwise base of the workspace path)")
	version := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *version {
		fmt.Println(Version)
		return
	}

	err := GoProjex(*ws, *src, *name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}
