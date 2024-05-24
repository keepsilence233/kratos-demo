package biz

//biz 负责业务组装，定义了biz的repo接口

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"kratos-demo/internal/model"
)

// StudentRepo 定义对 struct model.Student 的操作接口：
type StudentRepo interface {
	CreateStudent(context.Context, *model.Student) (*model.Student, error)
	GetStudent(context.Context, *model.Student) (*model.Student, error)
	DeleteStudent(context.Context, *model.Student) error
	UpdateStudent(context.Context, *model.Student) error
	ListStudent(context.Context, *PageInfo, *model.Student) ([]*model.Student, *PageInfo, error)
}

// StudentUsercase 对 model.Student 的操作加上日志
type StudentUsercase struct {
	repo StudentRepo
	log  *log.Helper
}

// NewStudentUsercase 初始化 StudentUsercase
func NewStudentUsercase(repo StudentRepo, logger log.Logger) *StudentUsercase {
	return &StudentUsercase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *StudentUsercase) CreateStudent(ctx context.Context, stu *model.Student) (*model.Student, error) {
	s, err := uc.repo.CreateStudent(ctx, stu)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// GetStudent 业务逻辑编写 类似DDD中的domain层
func (uc *StudentUsercase) GetStudent(ctx context.Context, stu *model.Student) (*model.Student, error) {
	uc.log.WithContext(ctx).Infof("GetStudent: %v", stu.Name)
	s, err := uc.repo.GetStudent(ctx, &model.Student{ID: stu.ID})
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (uc *StudentUsercase) UpdateStudent(ctx context.Context, stu *model.Student) (*model.Student, error) {
	err := uc.repo.UpdateStudent(ctx, stu)
	if err != nil {
		return nil, err
	}
	stu, err = uc.GetStudent(ctx, stu)
	return stu, nil
}

func (uc *StudentUsercase) DeleteStudent(ctx context.Context, stu *model.Student) error {
	return uc.repo.DeleteStudent(ctx, &model.Student{ID: stu.ID})
}

func (uc *StudentUsercase) ListStudent(ctx context.Context, pageInfo *PageInfo, stu *model.Student) ([]*model.Student, *PageInfo, error) {
	s, pageInfo, err := uc.repo.ListStudent(ctx, pageInfo, stu)
	if err != nil {
		return nil, nil, err
	}
	return s, pageInfo, nil
}
