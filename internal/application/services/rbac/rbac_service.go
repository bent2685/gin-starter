package rbac

import (
	"gin-starter/internal/infra/database"
	"strconv"

	casbin2 "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

type RBACService struct {
	enforcer *casbin2.Enforcer
}

var rbacService *RBACService

func InitRBAC() error {
	text := `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _
g2 = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == "super_admin" || g(r.sub, p.sub) || g2(r.sub, p.sub) || r.sub == p.sub && r.obj == p.obj && r.act == p.act
`
	m, err := model.NewModelFromString(text)
	if err != nil {
		return err
	}
	db := database.GetDB()
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return err
	}
	e, err := casbin2.NewEnforcer(m, a)
	if err != nil {
		return err
	}
	if err := e.LoadPolicy(); err != nil {
		return err
	}
	rbacService = &RBACService{enforcer: e}
	return nil
}

func GetRBAC() *RBACService {
	return rbacService
}

func AddPolicy(sub, obj, act string) (bool, error) {
	return rbacService.enforcer.AddPolicy(sub, obj, act)
}

func AddRoleForUser(user, role string) (bool, error) {
	return rbacService.enforcer.AddGroupingPolicy(user, role)
}

func DeleteRoleForUser(user, role string) (bool, error) {
	return rbacService.enforcer.RemoveGroupingPolicy(user, role)
}

func GetRolesForUser(user string) ([]string, error) {
	return rbacService.enforcer.GetRolesForUser(user)
}

func AddDepartmentForUser(user, department string) (bool, error) {
	return rbacService.enforcer.AddNamedGroupingPolicy("g2", user, department)
}

func DeleteDepartmentForUser(user, department string) (bool, error) {
	return rbacService.enforcer.RemoveNamedGroupingPolicy("g2", user, department)
}

func GetDepartmentsForUser(user string) ([]string, error) {
	policies, err := rbacService.enforcer.GetGroupingPolicy()
	if err != nil {
		return nil, err
	}
	var departments []string
	for _, policy := range policies {
		if len(policy) >= 2 && policy[0] == user {
			departments = append(departments, policy[1])
		}
	}
	return departments, nil
}

func AddPermissionForRole(role, resource, action string) (bool, error) {
	return rbacService.enforcer.AddPolicy(role, resource, action)
}

func AddPermissionForDepartment(department, resource, action string) (bool, error) {
	return rbacService.enforcer.AddPolicy(department, resource, action)
}

func Enforce(sub, obj, act string) (bool, error) {
	return rbacService.enforcer.Enforce(sub, obj, act)
}

func LoadPolicy() error {
	return rbacService.enforcer.LoadPolicy()
}

func SavePolicy() error {
	return rbacService.enforcer.SavePolicy()
}

func GetUserID(userID uint) string {
	return strconv.FormatUint(uint64(userID), 10)
}
