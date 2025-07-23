gen_common_proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
		./internal/proto/common.proto

gen_user_proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./internal/proto/user_service.proto

gen_wallet_proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./internal/proto/wallet_service.proto

gen_chat_proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./internal/proto/chat_service.proto

gen_proto: gen_common_proto gen_user_proto gen_wallet_proto gen_chat_proto

build/%:
	@script/build.sh $*

docker/gen/%:
	@script/docker_gen.sh $*

docker/build/%: build/%
	@script/docker_build.sh $*

docker/push/%:
	@script/docker_push.sh $*

k8s/gen/%:
	@script/k8s/gen.sh $*

k8s/apply/%:
	@script/k8s/apply.sh $*

k8s/delete/%:
	@script/k8s/delete.sh $*

# Local development targets
run/local:
	@script/run_local.sh start

stop/local:
	@script/run_local.sh stop

restart/local:
	@script/run_local.sh restart

status/local:
	@script/run_local.sh status

logs/local:
	@script/run_local.sh logs

clean/local:
	@script/run_local.sh clean
