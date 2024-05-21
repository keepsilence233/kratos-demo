package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v12 "kratos-demo/api/student/v1"
)

// 定义 struct Student
type Student struct {
	ID      string
	Name    string
	Sayname string
}

// 定义对 struct student 的操作接口：
type StudentRepo interface {
	GetStudent(context.Context, *Student) (*Student, error)
	ListStudent(context.Context, *Student) (*Student, error)
}

// 对 student 的操作加上日志
type StudentUsercase struct {
	repo StudentRepo
	log  *log.Helper
}

// 初始化 StudentUsercase
func NewStudentUsercase(repo StudentRepo, logger log.Logger) *StudentUsercase {
	return &StudentUsercase{repo: repo, log: log.NewHelper(logger)}
}

// 编写 GetStudent 方法，也就是一些业务逻辑编写 类似DDD中的domain层
func (uc *StudentUsercase) GetStudent(ctx context.Context, stu *v12.GetStudentRequest) (*v12.GetStudentReply, error) {
	uc.log.WithContext(ctx).Infof("GetStudent: %v", stu.Name)
	_, err := uc.repo.GetStudent(ctx, &Student{Name: stu.Name})
	if err != nil {
		return nil, err
	}
	return &v12.GetStudentReply{Name: stu.Name}, nil
}

// 编写 ListStudent 方法，也就是一些业务逻辑编写
func (uc *StudentUsercase) ListStudent(ctx context.Context, stu *Student) (*Student, error) {
	uc.log.WithContext(ctx).Infof("ListStudent: %v", stu.ID)
	return uc.repo.ListStudent(ctx, stu)
}
