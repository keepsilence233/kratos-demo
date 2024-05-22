package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"kratos-demo/internal/biz"
)

type studentRepo struct {
	data *Data
	log  *log.Helper
}

func (s *studentRepo) GetStudent(ctx context.Context, student *biz.Student) (*biz.Student, error) {
	//return &biz.Student{Name: "lisi"}, nil
	var stu biz.Student
	s.data.gormDB.Where("id = ?", student.ID).First(&stu) // 这里使用了 gorm
	return &biz.Student{
		ID:        stu.ID,
		Name:      stu.Name,
		Status:    stu.Status,
		Info:      stu.Info,
		UpdatedAt: stu.UpdatedAt,
		CreatedAt: stu.CreatedAt,
	}, nil
}

func (s *studentRepo) ListStudent(ctx context.Context, student *biz.Student) (*biz.Student, error) {
	//TODO implement me
	return nil, nil
}

func NewstudentRepo(data *Data, logger log.Logger) biz.StudentRepo {
	return &studentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
