gen_user_proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./internal/proto/user_service.proto

gen_wallet_proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./internal/proto/wallet_service.proto

gen_proto: gen_user_proto gen_wallet_proto
