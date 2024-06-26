#!/usr/bin/env bash
. "./script/common.sh"

usage() {
  cat <<EOF
To apply on k8s for all or a particular service. Target could be any folder under the
'./internal/service' one, such as:

  ./script/k8s_apply.sh http

Or, if you trigger this script via make:

  make k8s/apply/$service

You can also apply all under './internal/service' by a special target:
'all'.

  make k8s/apply/all
EOF
}

apply() {
  service=$1

  k8s_file="./internal/service/$service/k8s/dep.yaml"
  echo_info "Applying $k8s_file"
  kubectl apply -f $k8s_file
}

apply_all() {
  apply_middleware

  for service in $(ls ./internal/service/ -1) ; do
    apply "$service"
  done
}

apply_middleware() {
  for m in $(ls ./k8s/ -1) ; do
    for f in $(ls ./k8s/$m -1) ; do
      echo_info "Applying $m/$f"
      kubectl apply -f k8s/"$m"/"$f"
    done
  done

#  echo_info "Applying MySQL PV"
#
#  echo_info "Applying MySQL"
#  kubectl apply -f k8s/mysql_dep.yaml
#
#  echo_info "applying Kafka"
#  kubectl apply -f k8s/kafka_dep.yaml
#
#  echo_info "applying Loki"
#  kubectl apply -f k8s/loki.yaml
#
#  echo_info "Applying Grafana PV"
#  kubectl apply -f k8s/grafana_pv.yaml
#
#  echo_info "applying Grafana"
#  kubectl apply -f k8s/grafana_dep.yaml
#
#  echo_info "applying Promtail"
#  kubectl apply -f k8s/promtail_ds.yaml
}

apply_service() {
  for service in $(ls ./internal/service/ -1) ; do
    apply "$service"
  done
}

case "$1" in
all)
  apply_all
  exit ;;

middleware)
  apply_middleware
  exit ;;
service)
  apply_service
  exit ;;

-h | --help)
  usage
  exit ;;

*)
  apply "$@"
  exit ;;
esac
