package biz

//biz 负责业务组装，定义了biz的repo接口

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// Student is a Student model.
type Student struct {
	ID        int32
	Name      string
	Info      string
	Status    int32
	UpdatedAt time.Time
	CreatedAt time.Time
}

// StudentRepo 定义对 struct student 的操作接口：
type StudentRepo interface {
	CreateStudent(context.Context, *Student) (*Student, error)
	GetStudent(context.Context, *Student) (*Student, error)
	DeleteStudent(context.Context, *Student) error
	UpdateStudent(context.Context, *Student) error
	ListStudent(context.Context, *PageInfo, *Student) ([]*Student, *PageInfo, error)
}

// StudentUsercase 对 student 的操作加上日志
type StudentUsercase struct {
	repo StudentRepo
	log  *log.Helper
}

// NewStudentUsercase 初始化 StudentUsercase
func NewStudentUsercase(repo StudentRepo, logger log.Logger) *StudentUsercase {
	return &StudentUsercase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *StudentUsercase) CreateStudent(ctx context.Context, stu *Student) (*Student, error) {
	s, err := uc.repo.CreateStudent(ctx, stu)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// GetStudent 业务逻辑编写 类似DDD中的domain层
func (uc *StudentUsercase) GetStudent(ctx context.Context, stu *Student) (*Student, error) {
	uc.log.WithContext(ctx).Infof("GetStudent: %v", stu.Name)
	s, err := uc.repo.GetStudent(ctx, &Student{ID: stu.ID})
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (uc *StudentUsercase) UpdateStudent(ctx context.Context, stu *Student) (*Student, error) {
	err := uc.repo.UpdateStudent(ctx, stu)
	if err != nil {
		return nil, err
	}
	stu, err = uc.GetStudent(ctx, stu)
	return stu, nil
}

func (uc *StudentUsercase) DeleteStudent(ctx context.Context, stu *Student) error {
	return uc.repo.DeleteStudent(ctx, &Student{ID: stu.ID})
}

func (uc *StudentUsercase) ListStudent(ctx context.Context, pageInfo *PageInfo, stu *Student) ([]*Student, *PageInfo, error) {
	s, pageInfo, err := uc.repo.ListStudent(ctx, pageInfo, stu)
	if err != nil {
		return nil, nil, err
	}
	return s, pageInfo, nil
}
