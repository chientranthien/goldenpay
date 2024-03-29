package common

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Code struct {
	Id  int32  `json:"id"`
	Msg string `json:"msg"`
}

func NewCode(id int32, msg string) Code {
	return Code{Id: id, Msg: msg}
}

func (c Code) Success() bool {
	return c.Id == CodeSuccess.Id
}

var (
	CodeSuccess = Code{
		Id:  int32(codes.OK),
		Msg: "success",
	}
	CodeExisted = Code{
		Id:  int32(codes.AlreadyExists),
		Msg: "existed",
	}
	CodeUnknown = Code{
		Id:  int32(codes.Unknown),
		Msg: "unknown",
	}
	CodeNotFound = Code{
		Id:  int32(codes.NotFound),
		Msg: "not found",
	}
	CodeInvalidArgument = Code{
		Id:  int32(codes.InvalidArgument),
		Msg: "invalid argument",
	}
	CodeUnauthenticated = Code{
		Id:  int32(codes.Unauthenticated),
		Msg: "unauthenticated",
	}
	CodeProceeded = Code{
		Id:  18,
		Msg: "proceeded",
	}
	CodeAuthenticated = Code{
		Id:  19,
		Msg: "authenticated",
	}
	CodeBody = Code{
		Id:  20,
		Msg: "body",
	}
	CodeInvalidMetadata = Code{
		Id:  21,
		Msg: "invalid metadata",
	}

	AllCodes = map[int32]Code{
		CodeExisted.Id:         CodeExisted,
		CodeSuccess.Id:         CodeSuccess,
		CodeUnknown.Id:         CodeUnknown,
		CodeNotFound.Id:        CodeNotFound,
		CodeInvalidArgument.Id: CodeInvalidArgument,
		CodeUnauthenticated.Id: CodeUnauthenticated,
		CodeProceeded.Id:       CodeProceeded,
	}
)

func GetCode(c int32) Code {
	codeObj, ok := AllCodes[c]
	if !ok {
		return CodeUnknown
	}

	return codeObj
}
func GetCodeFromErr(e error) Code {
	return GetCode(int32(status.Code(e)))
}
