#!/bin/sh
set -e

COMMAND="${1:-server}"

if [ $COMMAND == "server" ]; then
  echo "Starting server..."
  {{{ .CmdServer }}}
{{{ range $k, $v := .OneOffs }}}
elif [ $COMMAND == "{{{$k}}}" ]; then
  echo "Running {{{$k}}}..."
  {{{$v}}}
{{{ end }}}
else
  echo "Usage: entrypoint.sh [CMD]"
  exit 1
fi
