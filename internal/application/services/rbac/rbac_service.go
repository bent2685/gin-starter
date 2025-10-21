package rbac

import (
	"gin-starter/internal/infra/database"
	"strconv"

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
	// 定义RBAC模型，支持角色和部门
	text := `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _  # 用户-角色关系
g2 = _, _ # 用户-部门关系

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == "super_admin" || g(r.sub, p.sub) || g2(r.sub, p.sub) || r.sub == p.sub && r.obj == p.obj && r.act == p.act
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

// AddPolicy 添加策略 (角色/部门的权限策略)
func (s *RBACService) AddPolicy(sub, obj, act string) (bool, error) {
	return s.enforcer.AddPolicy(sub, obj, act)
}

// AddRoleForUser 为用户添加角色
func (s *RBACService) AddRoleForUser(user, role string) (bool, error) {
	return s.enforcer.AddGroupingPolicy(user, role)
}

// DeleteRoleForUser 删除用户的角色
func (s *RBACService) DeleteRoleForUser(user, role string) (bool, error) {
	return s.enforcer.RemoveGroupingPolicy(user, role)
}

// GetRolesForUser 获取用户的角色
func (s *RBACService) GetRolesForUser(user string) ([]string, error) {
	return s.enforcer.GetRolesForUser(user)
}

// AddDepartmentForUser 为用户添加部门
func (s *RBACService) AddDepartmentForUser(user, department string) (bool, error) {
	// 使用g2来表示用户-部门关系
	return s.enforcer.AddNamedGroupingPolicy("g2", user, department)
}

// DeleteDepartmentForUser 删除用户的部门
func (s *RBACService) DeleteDepartmentForUser(user, department string) (bool, error) {
	return s.enforcer.RemoveNamedGroupingPolicy("g2", user, department)
}

// GetDepartmentsForUser 获取用户的部门
func (s *RBACService) GetDepartmentsForUser(user string) ([]string, error) {
	// 获取所有分组策略
	policies, err := s.enforcer.GetGroupingPolicy()
	if err != nil {
		return nil, err
	}

	// 筛选出用户所属的部门
	var departments []string
	for _, policy := range policies {
		// g分组策略格式: [user, role/department]
		// 我们需要区分角色和部门，这里简单地返回所有部门
		if len(policy) >= 2 && policy[0] == user {
			departments = append(departments, policy[1])
		}
	}

	return departments, nil
}

// AddPermissionForRole 为角色添加权限
func (s *RBACService) AddPermissionForRole(role, resource, action string) (bool, error) {
	return s.enforcer.AddPolicy(role, resource, action)
}

// AddPermissionForDepartment 为部门添加权限
func (s *RBACService) AddPermissionForDepartment(department, resource, action string) (bool, error) {
	return s.enforcer.AddPolicy(department, resource, action)
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

// GetUserID 获取用户ID字符串
func GetUserID(userID uint) string {
	return strconv.FormatUint(uint64(userID), 10)
}
