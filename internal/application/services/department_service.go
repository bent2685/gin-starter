package services

import (
	"gin-starter/internal/domain/models/rbac"
	"gin-starter/internal/infra/database"
	"gin-starter/pkg/utils/res"
)

type DepartmentService struct{}

var Department = &DepartmentService{}

func (s *DepartmentService) CreateDepartment(name, description string, parentID *uint) (*rbac.Department, error) {
	db := database.GetDB()
	var existingDepartment rbac.Department
	if err := db.Where("name = ?", name).First(&existingDepartment).Error; err == nil {
		return nil, res.ErrInvalidParam.WithMessage("部门名称已存在")
	}
	if parentID != nil {
		var parentDepartment rbac.Department
		if err := db.First(&parentDepartment, *parentID).Error; err != nil {
			return nil, res.ErrInvalidParam.WithMessage("父部门不存在")
		}
	}
	department := &rbac.Department{
		Name:        name,
		Description: description,
		ParentID:    parentID,
	}
	if err := db.Create(department).Error; err != nil {
		return nil, err
	}
	return department, nil
}

func (s *DepartmentService) GetAllDepartments() ([]*rbac.Department, error) {
	db := database.GetDB()
	var departments []*rbac.Department
	if err := db.Find(&departments).Error; err != nil {
		return nil, err
	}
	return departments, nil
}

func (s *DepartmentService) GetDepartmentByID(id uint) (*rbac.Department, error) {
	db := database.GetDB()
	var department rbac.Department
	if err := db.First(&department, id).Error; err != nil {
		return nil, res.ErrNotFound.WithMessage("部门不存在")
	}
	return &department, nil
}

func (s *DepartmentService) UpdateDepartment(id uint, name, description string, parentID *uint) (*rbac.Department, error) {
	db := database.GetDB()
	var department rbac.Department
	if err := db.First(&department, id).Error; err != nil {
		return nil, res.ErrNotFound.WithMessage("部门不存在")
	}
	var existingDepartment rbac.Department
	if err := db.Where("name = ? AND id != ?", name, id).First(&existingDepartment).Error; err == nil {
		return nil, res.ErrInvalidParam.WithMessage("部门名称已存在")
	}
	if parentID != nil {
		var parentDepartment rbac.Department
		if err := db.First(&parentDepartment, *parentID).Error; err != nil {
			return nil, res.ErrInvalidParam.WithMessage("父部门不存在")
		}
		if *parentID == id {
			return nil, res.ErrInvalidParam.WithMessage("不能将部门设置为自己的子部门")
		}
	}
	department.Name = name
	department.Description = description
	department.ParentID = parentID
	if err := db.Save(&department).Error; err != nil {
		return nil, err
	}
	return &department, nil
}

func (s *DepartmentService) DeleteDepartment(id uint) error {
	db := database.GetDB()
	var department rbac.Department
	if err := db.First(&department, id).Error; err != nil {
		return res.ErrNotFound.WithMessage("部门不存在")
	}
	var childDepartments []rbac.Department
	if err := db.Where("parent_id = ?", id).Find(&childDepartments).Error; err != nil {
		return err
	}
	if len(childDepartments) > 0 {
		return res.ErrInvalidParam.WithMessage("该部门有子部门，不能删除")
	}
	return db.Delete(&department).Error
}

func (s *DepartmentService) GetDepartmentTree() ([]*rbac.Department, error) {
	db := database.GetDB()
	var departments []rbac.Department
	if err := db.Find(&departments).Error; err != nil {
		return nil, err
	}
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
