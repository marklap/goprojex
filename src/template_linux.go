package main

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
