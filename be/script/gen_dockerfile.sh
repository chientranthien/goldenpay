#!/usr/bin/env bash
. "./script/common.sh"

gen() {
  service=$1
  echo_info "Generating Dockerfile for $service"

  target_dir="./internal/service/$service/docker/"
  mkdir -p $target_dir
  target=$target_dir"Dockerfile"
  sed s/_service/$service/g ./script/dockerfile_template  > $target
}

gen_all() {
  # Build
  for service in $(ls ./internal/service/ -1) ; do
    # build the service
    gen "$service"
  done
}

gen_all