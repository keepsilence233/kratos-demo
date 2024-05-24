package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"kratos-demo/internal/biz"
	"kratos-demo/internal/model"
)

var (
	WHERE_LIKE = "%"
)

type studentRepo struct {
	data *Data
	log  *log.Helper
}

func (s *studentRepo) CreateStudent(ctx context.Context, student *model.Student) (*model.Student, error) {
	s.data.gormDB.Create(student)
	return student, nil
}

func (s *studentRepo) GetStudent(ctx context.Context, student *model.Student) (*model.Student, error) {
	//return &model.Student{Name: "lisi"}, nil
	var stu model.Student
	s.data.gormDB.Where("id = ?", student.ID).First(&stu) // 这里使用了 gorm
	return &model.Student{
		ID:        stu.ID,
		Name:      stu.Name,
		Status:    stu.Status,
		Info:      stu.Info,
		UpdatedAt: stu.UpdatedAt,
		CreatedAt: stu.CreatedAt,
	}, nil
}

func (s *studentRepo) DeleteStudent(ctx context.Context, student *model.Student) error {
	s.data.gormDB.Where("id = ?", student.ID).Delete(&student)
	return nil
}

func (s *studentRepo) UpdateStudent(ctx context.Context, student *model.Student) error {
	s.data.gormDB.Model(&student).Updates(model.Student{Name: student.Name, Info: student.Info, Status: student.Status})
	return nil
}

func (s *studentRepo) ListStudent(ctx context.Context,
	pageInfo *biz.PageInfo, queryParams *model.Student) ([]*model.Student, *biz.PageInfo, error) {

	skuList := make([]*model.Student, 0)
	var totalCount int64

	s.data.gormDB.Scopes(Paginate(pageInfo.PageNo, pageInfo.PageSize)).
		Where(queryParams).Find(&skuList)
	//Where("id = ? and name like ? and info = ? and status = ?",
	//	queryParams.ID, queryParams.Name+WHERE_LIKE, queryParams.Info, queryParams.Status).Find(&skuList)
	s.data.gormDB.Model(&model.Student{}).Count(&totalCount)

	return skuList, &biz.PageInfo{
		PageNo:     pageInfo.PageNo,
		PageSize:   pageInfo.PageSize,
		TotalCount: totalCount,
	}, nil
}

//func (s *studentRepo) ListStudent(ctx context.Context,
//	pageInfo *biz.PageInfo, queryParams *model.Student) ([]*model.Student, *biz.PageInfo, error) {
//
//	query, u := gplus.NewQuery[model.Student]()
//	page := gplus.NewPage[model.Student](int(pageInfo.PageNo), int(pageInfo.PageSize))
//
//	if queryParams.Name != "" {
//		query.LikeRight(&u.Name, queryParams.Name)
//	}
//	if queryParams.Status > 0 {
//		query.Eq(&u.Status, queryParams.Status)
//	}
//
//	page, _ = gplus.SelectPage(page, query)
//
//	return page.Records, &biz.PageInfo{
//		PageNo:     int32(page.Current),
//		PageSize:   int32(page.Size),
//		TotalCount: page.Total,
//	}, nil
//}

func NewstudentRepo(data *Data, logger log.Logger) biz.StudentRepo {
	return &studentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
