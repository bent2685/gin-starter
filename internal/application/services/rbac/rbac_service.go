package rbac

import (
	"gin-starter/internal/infra/database"

	casbin2 "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

// RBACService RBAC服务
type RBACService struct {
	enforcer *casbin2.Enforcer
}

// NewRBACService 创建RBAC服务实例
func NewRBACService() (*RBACService, error) {
	// 定义RBAC模型
	text := `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`

	m, err := model.NewModelFromString(text)
	if err != nil {
		return nil, err
	}

	// 获取数据库连接
	db := database.GetDB()

	// 使用Casbin官方的GORM适配器
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	// 创建执行器
	e, err := casbin2.NewEnforcer(m, a)
	if err != nil {
		return nil, err
	}

	// 加载策略
	err = e.LoadPolicy()
	if err != nil {
		return nil, err
	}

	return &RBACService{
		enforcer: e,
	}, nil
}

// AddPolicy 添加策略
func (s *RBACService) AddPolicy(sub, obj, act string) (bool, error) {
	return s.enforcer.AddPolicy(sub, obj, act)
}

// AddRoleForUser 为用户添加角色
func (s *RBACService) AddRoleForUser(user, role string) (bool, error) {
	return s.enforcer.AddRoleForUser(user, role)
}

// DeleteRoleForUser 删除用户的角色
func (s *RBACService) DeleteRoleForUser(user, role string) (bool, error) {
	return s.enforcer.DeleteRoleForUser(user, role)
}

// GetRolesForUser 获取用户的角色
func (s *RBACService) GetRolesForUser(user string) ([]string, error) {
	return s.enforcer.GetRolesForUser(user)
}

// Enforce 验证权限
func (s *RBACService) Enforce(sub, obj, act string) (bool, error) {
	return s.enforcer.Enforce(sub, obj, act)
}

// LoadPolicy 加载策略
func (s *RBACService) LoadPolicy() error {
	return s.enforcer.LoadPolicy()
}

// SavePolicy 保存策略
func (s *RBACService) SavePolicy() error {
	return s.enforcer.SavePolicy()
}
