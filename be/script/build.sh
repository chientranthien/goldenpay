#!/usr/bin/env bash

. "$(git rev-parse --show-toplevel || echo ".")/be/script/common.sh"

build_user_service() {
 echo_info "Building user_service"
 go build -o bin/user_service/user_service ./internal/service/user/main/...
 echo_finish
}

build_wallet_service() {
 echo_info "Building wallet_service"
 go build -o bin/wallet_service/wallet_service ./internal/service/wallet/main/...
 echo_finish
}

build_http_gateway() {
 echo_info "Building http_gateway"
 go build -o bin/http_gateway/http_gateway ./internal/service/http/main/...
 echo_finish
}

build_event_handler() {
 echo_info "Building event_handler"
 go build -o bin/event_handler/event_handler ./internal/service/event_handler/main/...
 echo_finish
}

build() {
  build_user_service
  build_wallet_service
  build_http_gateway
  build_event_handler
}

build