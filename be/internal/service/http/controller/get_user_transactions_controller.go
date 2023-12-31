package controller

import (
	"context"

	"github.com/chientranthien/goldenpay/internal/common"
	httpcommon "github.com/chientranthien/goldenpay/internal/common/http"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type (
	Pagination struct {
		Val     int64  `json:"val"`
		Limit   uint32 `json:"limit"`
		HasMore bool   `json:"has_more"`
	}

	GetUserTransactionsBody struct {
		Pagination *Pagination `json:"pagination"`
	}

	User struct {
		Id    uint64 `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	Transaction struct {
		Id     uint64 `json:"id"`
		From   User   `json:"from"`
		To     User   `json:"to"`
		Amount int64  `json:"amount"`
		Status uint64 `json:"status"`
		Ctime  uint64 `json:"ctime"`
	}

	GetUserTransactionsData struct {
		Transactions   []*Transaction `json:"transactions"`
		NextPagination *Pagination    `json:"next_pagination"`
	}

	GetUserTransactionsController struct {
		uClient proto.UserServiceClient
		wClient proto.WalletServiceClient

		ctx  context.Context
		body *GetUserTransactionsBody
		req  httpcommon.Req
	}
)

func NewGetUserTransactionsController(
	uClient proto.UserServiceClient,
	wClient proto.WalletServiceClient,
) *GetUserTransactionsController {
	return &GetUserTransactionsController{uClient: uClient, wClient: wClient}
}

func FromPagination(p *proto.Pagination) *Pagination {
	if p == nil {
		return nil
	}

	return &Pagination{Val: p.Val, Limit: p.Limit, HasMore: p.HasMore}
}

func (c GetUserTransactionsController) Do() (common.AnyPtr, common.Code) {
	getResp, err := c.wClient.GetUserTransactions(c.ctx, &proto.GetUserTransactionsReq{
		Cond: &proto.GetUserTransactionsCond{
			User: &proto.GetUserTransactionsCond_UserCond{Eq: c.req.Metadata.UserId},
		},
		Pagination: c.body.Pagination.ToServicePagination(),
	})

	code := common.GetCodeFromErr(err)
	if !code.Success() {
		common.L().Errorw("getUserTransactionsErr", "bode", c.body, "err", err)
		return nil, code
	}

	var transactions []*Transaction
	if len(getResp.GetTransactions()) > 0 {
		var ids []uint64
		for _, t := range getResp.Transactions {
			ids = append(ids, t.FromUser, t.ToUser)
			transactions = append(transactions, &Transaction{
				Id:     t.Id,
				Amount: t.Amount,
				Status: t.Status,
				Ctime:  t.Ctime,
			})
		}
		batchResp, err := c.uClient.GetBatch(c.ctx, &proto.GetBatchReq{Ids: ids})
		code = common.GetCodeFromErr(err)
		if !code.Success() {
			common.L().Errorw("userUserBatchErr", "bode", c.body, "err", err)
			return nil, code
		}

		for i, user := range batchResp.Users {
			if i%2 == 0 {
				transactions[i/2].From = User{
					Id:    user.Id,
					Name:  user.Name,
					Email: user.Email,
				}
			} else {
				transactions[i/2].To = User{
					Id:    user.Id,
					Name:  user.Name,
					Email: user.Email,
				}
			}
		}

	}

	data := &GetUserTransactionsData{
		Transactions:   transactions,
		NextPagination: FromPagination(getResp.NextPagination),
	}

	return data, common.CodeSuccess
}

func (c *GetUserTransactionsController) Take(ctx context.Context, req httpcommon.Req) common.Code {
	if req.Metadata.UserId <= 0 {
		return common.CodeInvalidMetadata
	}

	if b, ok := req.Body.(*GetUserTransactionsBody); ok {
		c.body = b
		c.ctx = ctx
		c.req = req
	} else {
		return common.CodeBody
	}

	return common.CodeSuccess
}

func (p Pagination) ToServicePagination() *proto.Pagination {
	return &proto.Pagination{
		Val:     p.Val,
		Limit:   p.Limit,
		HasMore: p.HasMore,
	}
}
