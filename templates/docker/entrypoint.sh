#!/bin/sh
set -e

COMMAND="${1:-server}"

if [ $COMMAND == "server" ]; then
  echo "Starting server..."
  echo "TODO"
{{ range .ECSOneOffs }}
elif [ $COMMAND == "{{.}}" ]; then
  echo "Running {{.}}..."
  echo "TODO"
{{ end }}
else
  echo "Usage: entrypoint.sh [CMD]"
  exit 1
fi
