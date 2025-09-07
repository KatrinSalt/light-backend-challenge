package workflow

import (
	"github.com/KatrinSalt/backend-challenge-go/api"
)

// approvalChannel represents the channel for sending approval requests.
type approvalChannel string

// approver represents the approver and the channel for sending approval requests.
type approver struct {
	approver        api.Approver
	approvalChannel approvalChannel
}
