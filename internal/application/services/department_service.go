package services

import (
	"gin-starter/internal/domain/models/rbac"
	"gin-starter/internal/infra/database"
	"gin-starter/pkg/utils/res"

	"gorm.io/gorm"
)

// DepartmentService 部门服务
type DepartmentService struct {
	db *gorm.DB
}

// NewDepartmentService 创建部门服务实例
func NewDepartmentService() *DepartmentService {
	return &DepartmentService{
		db: database.GetDB(),
	}
}

// CreateDepartment 创建部门
func (s *DepartmentService) CreateDepartment(name, description string, parentID *uint) (*rbac.Department, error) {
	// 检查部门名称是否已存在
	var existingDepartment rbac.Department
	if err := s.db.Where("name = ?", name).First(&existingDepartment).Error; err == nil {
		return nil, res.ErrInvalidParam.WithMessage("部门名称已存在")
	}

	// 如果指定了父部门，检查父部门是否存在
	if parentID != nil {
		var parentDepartment rbac.Department
		if err := s.db.First(&parentDepartment, *parentID).Error; err != nil {
			return nil, res.ErrInvalidParam.WithMessage("父部门不存在")
		}
	}

	// 创建部门对象
	department := &rbac.Department{
		Name:        name,
		Description: description,
		ParentID:    parentID,
	}

	// 保存到数据库
	if err := s.db.Create(department).Error; err != nil {
		return nil, err
	}

	return department, nil
}

// GetAllDepartments 获取所有部门
func (s *DepartmentService) GetAllDepartments() ([]*rbac.Department, error) {
	var departments []*rbac.Department
	if err := s.db.Find(&departments).Error; err != nil {
		return nil, err
	}

	return departments, nil
}

// GetDepartmentByID 根据ID获取部门
func (s *DepartmentService) GetDepartmentByID(id uint) (*rbac.Department, error) {
	var department rbac.Department
	if err := s.db.First(&department, id).Error; err != nil {
		return nil, res.ErrNotFound.WithMessage("部门不存在")
	}

	return &department, nil
}

// UpdateDepartment 更新部门
func (s *DepartmentService) UpdateDepartment(id uint, name, description string, parentID *uint) (*rbac.Department, error) {
	var department rbac.Department
	if err := s.db.First(&department, id).Error; err != nil {
		return nil, res.ErrNotFound.WithMessage("部门不存在")
	}

	// 检查部门名称是否已被其他部门使用
	var existingDepartment rbac.Department
	if err := s.db.Where("name = ? AND id != ?", name, id).First(&existingDepartment).Error; err == nil {
		return nil, res.ErrInvalidParam.WithMessage("部门名称已存在")
	}

	// 如果指定了父部门，检查父部门是否存在
	if parentID != nil {
		var parentDepartment rbac.Department
		if err := s.db.First(&parentDepartment, *parentID).Error; err != nil {
			return nil, res.ErrInvalidParam.WithMessage("父部门不存在")
		}

		// 检查是否试图将部门设置为自己的子部门
		if *parentID == id {
			return nil, res.ErrInvalidParam.WithMessage("不能将部门设置为自己的子部门")
		}
	}

	// 更新部门信息
	department.Name = name
	department.Description = description
	department.ParentID = parentID

	if err := s.db.Save(&department).Error; err != nil {
		return nil, err
	}

	return &department, nil
}

// DeleteDepartment 删除部门
func (s *DepartmentService) DeleteDepartment(id uint) error {
	var department rbac.Department
	if err := s.db.First(&department, id).Error; err != nil {
		return res.ErrNotFound.WithMessage("部门不存在")
	}

	// 检查是否有子部门
	var childDepartments []rbac.Department
	if err := s.db.Where("parent_id = ?", id).Find(&childDepartments).Error; err != nil {
		return err
	}

	if len(childDepartments) > 0 {
		return res.ErrInvalidParam.WithMessage("该部门有子部门，不能删除")
	}

	// 删除部门
	if err := s.db.Delete(&department).Error; err != nil {
		return err
	}

	return nil
}

// GetDepartmentTree 获取部门树形结构
func (s *DepartmentService) GetDepartmentTree() ([]*rbac.Department, error) {
	var departments []rbac.Department
	if err := s.db.Find(&departments).Error; err != nil {
		return nil, err
	}

	// 构建树形结构
	departmentMap := make(map[uint]*rbac.Department)
	for i := range departments {
		departmentMap[departments[i].ID] = &departments[i]
	}

	var rootDepartments []*rbac.Department
	for i := range departments {
		department := &departments[i]
		if department.ParentID == nil {
			rootDepartments = append(rootDepartments, department)
		} else {
			if parent, exists := departmentMap[*department.ParentID]; exists {
				parent.Children = append(parent.Children, *department)
			}
		}
	}

	return rootDepartments, nil
}
