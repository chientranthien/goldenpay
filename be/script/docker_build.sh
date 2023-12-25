#!/usr/bin/env bash
. "./script/common.sh"

usage() {
  cat <<EOF
To build Docker images for all or a particular service. Build target could be any folder under the
'./internal/service' one, such as:

  ./script/docker_build.sh http

Or, if you trigger this script via make:

  make docker/build/$service

You can also build all under './internal/service' by a special target:
'all'.

  make docker/build/all

See "make docker/gen/all" if you are looking for Dockerfile generation
EOF
}

build() {
  service=$1

  tag="chientt1993/goldenpay-be-"$service
  tag=$(echo $tag | sed s/_/-/g)
  docker_file="./internal/service/$service/docker/Dockerfile"
  env=${G_ENV:="dev"}
  arg="G_ENV=$env"

  echo_info "building $docker_file, arg= $arg ,tag= $tag"
  case "$env" in
   dev)
    docker build --progress=plain --build-arg $arg -f $docker_file -t $tag .
    ;;
    prod)
    docker build --progress=plain --build-arg $arg -f $docker_file -t $tag . --no-cache
    ;;
  esac
}

build_all() {
  # Build
  for service in $(ls ./internal/service/ -1) ; do
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
  build "$@"
  exit
  ;;
esac
