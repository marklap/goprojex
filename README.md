# goprojex

An opinionated project management/import isolation tool. It basically provides a convenient way
to modify and revert your GOPATH, PATH, and shell prompt to avoid version conflicts in imported
go packages.

# Build & Install

    git clone git@github.com:marklap/goprojex.git
    cd goprojex
    go build -o goprojex *.go
    cp goprojex $HOME/bin

# Use

View help message

    goprojex -h

Create a default workspace in the current directory:

    goprojex

Create a workspace in `./here` with a source tree of 'my/code':

    goprojex -ws ./here -src my/code

Activate your workspace

    . ./.gopjx/activate

Deactivate your workspace

    deactivate



# Contribute

Create a template_*.go file for your GOOS and submit a PR
