package cctx

import "context"

const (
	TenantID = "TenantId" // 租户id
	UserID   = "UserId"   // 用户id
)

func GetTenantIdFromContext(ctx context.Context) string {
	tenantID, ok := ctx.Value(TenantID).(string)
	if !ok {
		return ""
	}
	return tenantID
}
