////////////////////////////////////////////////////////////////////////////////
// The MIT License (MIT)
//
// Copyright (c) 2017 Mark LaPerriere
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

  if [ -n "$VIRTUAL_ENV" ]; then
    unset VIRTUAL_ENV
  fi

  if [ -n "$GOPJX_SRC_PATH" ]; then
    unset GOPJX_SRC_PATH
  fi

  if [ -n "$BASH" -o -n "$ZSH_VERSION" ]; then
      hash -r 2>/dev/null
  fi

  if [ "$1" != "nondestructive" ]; then
    unset -f deactivate
  fi
}

deactivate nondestructive

if [ -z "$VIRTUAL_ENV_DISABLE_PROMPT" ]; then
  export _OLD_GOPJX_PS1="$PS1"
  export PS1="[go:{{.Name}}] $_OLD_GOPJX_PS1"
fi

export _OLD_GOPJX_PATH="$PATH"
export _OLD_GOPJX_GOPATH="$GOPATH"
export GOPATH={{.GoPath}}
export PATH=$GOPATH/bin:$PATH
export VIRTUAL_ENV=$GOPATH
export GOPJX_SRC_PATH={{.SrcPath}}
`
