package dto

// CreateDepartmentRequest 创建部门请求
type CreateDepartmentRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"max=255"`
	ParentID    *uint  `json:"parent_id"`
}

// UpdateDepartmentRequest 更新部门请求
type UpdateDepartmentRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"max=255"`
	ParentID    *uint  `json:"parent_id"`
}

// DepartmentResponse 部门响应
type DepartmentResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	ParentID    *uint                `json:"parent_id"`
	Parent      *DepartmentResponse  `json:"parent,omitempty"`
	Children    []DepartmentResponse `json:"children,omitempty"`
	CreatedAt   int64                `json:"created_at"`
	UpdatedAt   int64                `json:"updated_at"`
}
