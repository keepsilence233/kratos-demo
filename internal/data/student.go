package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"kratos-demo/internal/biz"
)

var (
	WHERE_LIKE = "%"
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

func (s *studentRepo) UpdateStudent(ctx context.Context, student *biz.Student) error {
	s.data.gormDB.Model(&student).Updates(biz.Student{Name: student.Name, Info: student.Info, Status: student.Status})
	return nil
}

func (s *studentRepo) ListStudent(ctx context.Context,
	pageInfo *biz.PageInfo, queryParams *biz.Student) ([]*biz.Student, *biz.PageInfo, error) {

	skuList := make([]*biz.Student, 0)
	var totalCount int64

	s.data.gormDB.Scopes(Paginate(pageInfo.PageNo, pageInfo.PageSize)).
		Where(queryParams).Find(&skuList)
	//Where("id = ? and name like ? and info = ? and status = ?",
	//	queryParams.ID, queryParams.Name+WHERE_LIKE, queryParams.Info, queryParams.Status).Find(&skuList)
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
