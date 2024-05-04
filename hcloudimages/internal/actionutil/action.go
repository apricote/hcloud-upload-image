package actionutil

import (
	"context"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func Settle(ctx context.Context, client hcloud.IActionClient, actions ...*hcloud.Action) (successActions []*hcloud.Action, errorActions []*hcloud.Action, err error) {
	err = client.WaitForFunc(ctx, func(update *hcloud.Action) error {
		switch update.Status {
		case hcloud.ActionStatusSuccess:
			successActions = append(successActions, update)
		case hcloud.ActionStatusError:
			errorActions = append(errorActions, update)
		}

		return nil
	}, actions...)
	if err != nil {
		return nil, nil, err
	}

	return successActions, errorActions, nil
}
