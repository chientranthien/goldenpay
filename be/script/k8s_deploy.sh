#!/usr/bin/env bash
. "./script/common.sh"

usage() {
  cat <<EOF
To deploy on k8s for all or a particular service. Target could be any folder under the
'./internal/service' one, such as:

  ./script/k8s_deploy.sh http

Or, if you trigger this script via make:

  make k8s/deploy/$service

You can also deploy all under './internal/service' by a special target:
'all'.

  make k8s/deploy/all
EOF
}

deploy() {
  service=$1

  k8s_file="./internal/service/$service/k8s/dep.yaml"
  echo_info "deploying $k8s_file"
  kubectl apply -f $k8s_file
}

deploy_all() {
  deploy_middleware

  for service in $(ls ./internal/service/ -1) ; do
    deploy "$service"
  done
}

deploy_middleware() {
  echo_info "Deploying MySQL PV"
  kubectl apply -f k8s/mysql_pv.yaml

  echo_info "Deploying MySQL"
  kubectl apply -f k8s/mysql_dep.yaml

  echo_info "Deploying Kafka"
  kubectl apply -f k8s/kafka_dep.yaml
}

deploy_service() {
  for service in $(ls ./internal/service/ -1) ; do
    deploy "$service"
  done
}

case "$1" in
all)
  deploy_all
  exit ;;

middleware)
  deploy_middleware
  exit ;;
service)
  deploy_service
  exit ;;

-h | --help)
  usage
  exit ;;

*)
  deploy "$@"
  exit ;;
esac
