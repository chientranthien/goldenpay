gen_user_proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./internal/proto/user_service.proto

gen_wallet_proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./internal/proto/wallet_service.proto

gen_proto: gen_user_proto gen_wallet_proto

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
