package internal

import (
	"github.com/lvyahui8/goenum"
)

// Role 参考 https://docs.gitlab.com/ee/user/permissions.html
type Role struct {
	goenum.Enum
	perms []Permission
}

func (r *Role) UnmarshalJSON(data []byte) error {
	role, err := goenum.Unmarshal[Role](data)
	if err != nil {
		return err
	}
	*r = role
	return nil
}

func (r *Role) HasPerm(p Permission) bool {
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
	Issues        = goenum.NewEnum[Module]("Issues", Module{perms: []Permission{AddLabels, AddTopic}, basePath: "/issues/"})
	MergeRequests = goenum.NewEnum[Module]("MergeRequests", Module{perms: []Permission{ViewMergeRequest, ApproveMergeRequest, DeleteMergeRequest}, basePath: "/merge/"})
)

// 定义角色
var (
	Reporter  = goenum.NewEnum[Role]("Reporter", Role{perms: []Permission{ViewMergeRequest}})
	Developer = goenum.NewEnum[Role]("Developer", Role{perms: []Permission{AddLabels, AddTopic, ViewMergeRequest}})
	Owner     = goenum.NewEnum[Role]("Owner", Role{perms: []Permission{AddLabels, AddTopic, ViewMergeRequest, ApproveMergeRequest, DeleteMergeRequest}}) // 可以考虑给Owner单独定义一个All的权限
)
