#!/usr/bin/env bash

. "./script/common.sh"

usage() {
  cat <<EOF
Build binary artifacts this service. Build target could be any folder under the
'./internal/service' one, such as:

  ./script/build.sh service

Or, if you trigger this script via make:

  make build/$service

You can also build all packages under './internal/service' into binaries by a special target:
'all'.

  make build/all

To do cross compilation, such as build binary for Linux while working on Mac

  GOOS=linux make build/all
EOF
}


[[ -s "/home/chien_tran/.gvm/scripts/gvm" ]] && source "/home/chien_tran/.gvm/scripts/gvm"
gvm use go1.22
build() {
  service=$1
  echo_info "Building $service"
  out="./bin/$service/exc"
  target="./internal/service/$service/main/..."
  go build -ldflags="$GO_LDFLAGS"  -v -o $out $target
  ls -lah -d $out
}

check_target() {
  service="$1"
  if [ -z $1 ]; then
    usage
    exit
  fi

  target="./internal/service/$service/main"
  if [ -d "$target" ]; then
    build "$service"
  else
    echo_warn "$service is not exists! Here's how to use build script"
    usage
  fi
}

build_all() {
  # Build
  for service in $(ls ./internal/service/ -1) ; do
    # skip the folder if there's no go file in it
    ls "./internal/service/$service/main/"*.go >/dev/null 2>&1 || (echo_warn "skipped $service" && continue)
    # build the service
    build "$service"
  done
}

case "$1" in
all)
  build_all
  exit ;;

-h | --help)
  usage
  exit
  ;;
*)
  check_target "$@"
  exit
  ;;
esac