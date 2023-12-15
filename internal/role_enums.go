package internal

import "github.com/lvyahui8/goenum"

func castList[T any](items ...any) (res []T) {
	for _, item := range items {
		if v, ok := item.(T); ok {
			res = append(res, v)
		}
	}
	return
}

// Role 参考 https://docs.gitlab.com/ee/user/permissions.html
type Role struct {
	goenum.Enum
	perms []Permission
}

func (r Role) Init(args ...any) any {
	r.perms = castList[Permission](args...)
	return r
}

func (r Role) HasPerm(p Permission) bool {
	for _, perm := range r.perms {
		if p.Equals(perm) {
			return true
		}
	}
	return false
}

type Module struct {
	goenum.Enum
	perms    []Permission
	basePath string
}

func (m Module) Init(args ...any) any {
	m.perms = args[0].([]Permission)
	m.basePath = args[1].(string)
	return m
}

func (m Module) GetPerms() []Permission {
	return m.perms
}

func (m Module) BasePath() string {
	return m.basePath
}

type Permission struct {
	goenum.Enum
}

// 定义权限
var (
	AddLabels           = goenum.NewEnum[Permission]("AddLabels")
	AddTopic            = goenum.NewEnum[Permission]("AddTopic")
	ViewMergeRequest    = goenum.NewEnum[Permission]("ViewMergeRequest")
	ApproveMergeRequest = goenum.NewEnum[Permission]("ApproveMergeRequest")
	DeleteMergeRequest  = goenum.NewEnum[Permission]("DeleteMergeRequest")
)

// 定义模块
var (
	Issues        = goenum.NewEnum[Module]("Issues", []Permission{AddLabels, AddTopic}, "/issues/")
	MergeRequests = goenum.NewEnum[Module]("MergeRequests", []Permission{ViewMergeRequest, ApproveMergeRequest, DeleteMergeRequest}, "/merge/")
)

// 定义角色
var (
	Reporter  = goenum.NewEnum[Role]("Reporter", ViewMergeRequest)
	Developer = goenum.NewEnum[Role]("Developer", AddLabels, AddTopic, ViewMergeRequest)
	Owner     = goenum.NewEnum[Role]("Owner", AddLabels, AddTopic, ViewMergeRequest, ApproveMergeRequest, DeleteMergeRequest) // 可以考虑给Owner单独定义一个All的权限
)
