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

push() {
  service=$1

  tag="chientt1993/goldenpay-be-"$service
  echo_info "pushing docker,tag= $tag"
  docker push $tag
}

push_all() {
  # Build
  for service in $(ls ./internal/service/ -1) ; do
    # build the service
    push "$service"
  done
}

case "$1" in
all)
  push_all
  exit ;;

-h | --help)
  usage
  exit
  ;;
*)
  push "$@"
  exit
  ;;
esac
