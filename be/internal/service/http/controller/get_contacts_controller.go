package controller

import (
	"context"

	"github.com/chientranthien/goldenpay/internal/common"
	httpcommon "github.com/chientranthien/goldenpay/internal/common/http"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type (
	GetContactsBody struct {
		Pagination *Pagination `json:"pagination"`
	}

	Contact struct {
		Id            uint64 `json:"id"`
		ContactUserId uint64 `json:"contact_user_id"`
		Name          string `json:"name"`
		Email         string `json:"email"`
	}

	GetContactsData struct {
		Contacts       []Contact   `json:"contacts"`
		NextPagination *Pagination `json:"next_pagination"`
	}

	GetContactsController struct {
		uClient proto.UserServiceClient

		ctx  context.Context
		body *GetContactsBody
		req  httpcommon.Req
	}
)

func NewGetContactsController(uClient proto.UserServiceClient) *GetContactsController {
	return &GetContactsController{uClient: uClient}
}

func (c *GetContactsController) Do() (common.AnyPtr, common.Code) {
	getResp, err := c.uClient.GetContacts(c.ctx, &proto.GetContactsReq{
		Cond: &proto.GetContactsCond{
			User: &proto.GetContactsCond_UserCond{
				Eq: c.req.Metadata.UserId,
			},
		},
		Pagination: c.body.Pagination.ToServicePagination(),
	})

	code := common.GetCodeFromErr(err)
	if !code.Success() {
		common.L().Errorw("getContactsErr", "bode", c.body, "err", err)
		return nil, code
	}

	var contacts []Contact
	for _, contact := range getResp.Contacts {
		getUserResp, err := c.uClient.Get(c.ctx, &proto.GetReq{
			Id: contact.ContactUserId,
		})

		code := common.GetCodeFromErr(err)
		if !code.Success() {
			common.L().Errorw("getContactsErr", "bode", c.body, "err", err)
			return nil, code
		}
		contacts = append(contacts, Contact{
			Id:            contact.Id,
			ContactUserId: contact.ContactUserId,
			Name:          getUserResp.GetUser().GetName(),
			Email:         getUserResp.GetUser().GetEmail(),
		})
	}

	return &GetContactsData{
		Contacts:       contacts,
		NextPagination: FromPagination(getResp.NextPagination),
	}, common.CodeSuccess
}

func (c *GetContactsController) Take(ctx context.Context, req httpcommon.Req) common.Code {
	if req.Metadata.UserId <= 0 {
		return common.CodeInvalidMetadata
	}

	if b, ok := req.Body.(*GetContactsBody); ok {
		c.body = b
		c.ctx = ctx
		c.req = req
	} else {
		return common.CodeBody
	}

	return common.CodeSuccess
}
