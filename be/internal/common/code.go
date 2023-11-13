package common

import (
	"google.golang.org/grpc/codes"

	"github.com/chientranthien/goldenpay/internal/proto"
)

var (
	CodeSuccess = &proto.Code{
		Id:  int32(codes.OK),
		Msg: "success",
	}
	CodeExisted = &proto.Code{
		Id:  int32(codes.AlreadyExists),
		Msg: "existed",
	}
	CodeUnknown = &proto.Code{
		Id:  int32(codes.Unknown),
		Msg: "unknown",
	}

	AllCodes = map[int32]*proto.Code{
		CodeExisted.Id: CodeExisted,
	}
)

func GetCode(c int32) *proto.Code {
	codeObj, ok := AllCodes[int32(c)]
	if !ok {
		return CodeUnknown
	}

	return codeObj
}
