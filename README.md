# goprojex

An opinionated project management/import isolation tool. It basically provides a convenient way
to modify and revert your GOPATH, PATH, and shell prompt to avoid version conflicts in imported
go packages.

# Build & Install

    git clone git@github.com:marklap/goprojex.git
    cd goprojex
    go build -o goprojex src/*.go
    cp goprojex $HOME/bin

# Use

View help message

    goprojex -h

Create a project gopath of `./myproj`, an activate script of `./myproj/activate` and a project name
(shown in the prompt) of `[go:myproj]`

    goprojex myproj

Create a project gopath of `./.gopath`, an activate script of `./.gopath/activate` and a project
name of `[go:myproj]`

    goprojex -name myproj .gopath

Activate your project gopath

    . ./myproj/activate

Deactivate your project gopath

    deactivate



# Contribute

Create a template_*.go file for your GOOS and submit a PR
