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

func (s *studentRepo) CreateStudent(ctx context.Context, student *biz.Student) (*biz.Student, error) {
	s.data.gormDB.Create(student)
	return student, nil
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

func (s *studentRepo) DeleteStudent(ctx context.Context, student *biz.Student) error {
	s.data.gormDB.Where("id = ?", student.ID).Delete(&student)
	return nil
}

func (s *studentRepo) ListStudent(ctx context.Context,
	pageInfo *biz.PageInfo, student *biz.Student) ([]*biz.Student, *biz.PageInfo, error) {

	skuList := make([]*biz.Student, 0)
	var totalCount int64

	queryParams := biz.Student{
		ID:     student.ID,
		Name:   student.Name,
		Info:   student.Info,
		Status: student.Status,
	}
	s.data.gormDB.Scopes(Paginate(pageInfo.PageNo, pageInfo.PageSize)).Where(queryParams).Find(&skuList)
	s.data.gormDB.Model(&biz.Student{}).Count(&totalCount)

	return skuList, &biz.PageInfo{
		PageNo:     pageInfo.PageNo,
		PageSize:   pageInfo.PageSize,
		TotalCount: totalCount,
	}, nil
}

func NewstudentRepo(data *Data, logger log.Logger) biz.StudentRepo {
	return &studentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
