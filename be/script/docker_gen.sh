#!/usr/bin/env bash
. "./script/common.sh"

usage() {
  cat <<EOF
To generate Dockerfile for all or a particular service. Target could be any folder under the
'./internal/service' one, such as:

  ./script/docker_gen.sh http

Or, if you trigger this script via make:

  make docker/gen/$service

You can also gen all under './internal/service' by a special target:
'all'.

  make docker/gen/all
EOF
}

gen() {
  service=$1

  target_dir="./internal/service/$service/docker/"
  mkdir -p $target_dir
  target=$target_dir"Dockerfile"

  echo_info "Generating Dockerfile for $service , target= $target"
  sed s/__service/$service/g ./script/dockerfile_template.Dockerfile  > $target
}

gen_all() {
  # Build
  for service in $(ls ./internal/service/ -1) ; do
    # build the service
    gen "$service"
  done
}



case "$1" in
all)
  gen_all
  exit ;;

-h | --help)
  usage
  exit
  ;;
*)
  gen "$@"
  exit
  ;;
esac
