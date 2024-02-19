#!/usr/bin/env bash
. "./script/common.sh"

usage() {
  cat <<EOF
To delete on k8s for all or a particular service. Target could be any folder under the
'./internal/service' one, such as:

  ./script/k8s_delete.sh http

Or, if you trigger this script via make:

  make k8s/delete/$service

You can also deploy all under './internal/service' by a special target:
'all'.

  make k8s/delete/all
EOF
}

delete() {
  service=$1

  k8s_file="./internal/service/$service/k8s/dep.yaml"
  echo_info "deleting $k8s_file"
  kubectl delete -f $k8s_file
}

delete_all() {
  delete_middleware

  for service in $(ls ./internal/service/ -1) ; do
    delete "$service"
  done
}

alias k8s_delete_cmd='kubectl delete -f'
delete_middleware() {
  echo_info "deleting MySQL"
  kubectl delete -f k8s/mysql_dep.yaml

  echo_info "deleting Kafka"
  kubectl delete -f k8s/kafka_dep.yaml

  echo_info "deleting Promtail"
  kubectl delete -f k8s/promtail_ds.yaml

  echo_info "deleting Grafana"
  kubectl delete -f k8s/grafana_dep.yaml

  echo_info "deleting Loki"
  kubectl delete -f k8s/loki.yaml
}

delete_service() {
  for service in $(ls ./internal/service/ -1) ; do
    delete "$service"
  done
}

case "$1" in
all)
  delete_all
  exit ;;

middleware)
  delete_middleware
  exit ;;
service)
  delete_service
  exit ;;

-h | --help)
  usage
  exit ;;

*)
  delete "$@"
  exit ;;
esac
