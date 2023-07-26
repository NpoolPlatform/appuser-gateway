package history

import (
	"context"
	"encoding/json"
	"fmt"

	historymwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user/login/history"
)

func Prepare(body string) (interface{}, error) {
	req := []*historymwpb.HistoryReq{}
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		return nil, err
	}
	return req, nil
}

func Apply(ctx context.Context, req interface{}) error {
	history, ok := req.([]*historymwpb.HistoryReq)
	if !ok {
		return fmt.Errorf("invalid request in apply")
	}
	fmt.Println("history", history)

	// generate notif
	return nil
}
