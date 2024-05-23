package service

import (
	"context"
	pb "kratos-demo/api/student/v1"
	"kratos-demo/internal/biz"
)

// 通过proto文件 直接生成对应的service代码 使用-t指定生成目录
// kratos proto server api/helloworld/v1/student.proto -t internal/service
type StudentService struct {
	pb.UnimplementedStudentServer

	uc *biz.StudentUsercase
}

func NewStudentService(uc *biz.StudentUsercase) *StudentService {
	return &StudentService{uc: uc}
}

func (s *StudentService) CreateStudent(ctx context.Context, req *pb.CreateStudentRequest) (*pb.CreateStudentReply, error) {
	stu, err := s.uc.CreateStudent(ctx, &biz.Student{
		Name:   req.Name,
		Info:   req.Info,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateStudentReply{
		ID:     stu.ID,
		Name:   stu.Name,
		Status: stu.Status,
		Info:   stu.Info,
	}, nil
}
func (s *StudentService) UpdateStudent(ctx context.Context, req *pb.UpdateStudentRequest) (*pb.UpdateStudentReply, error) {
	return &pb.UpdateStudentReply{}, nil
}
func (s *StudentService) DeleteStudent(ctx context.Context, req *pb.DeleteStudentRequest) (*pb.DeleteStudentReply, error) {
	err := s.uc.DeleteStudent(ctx, &biz.Student{
		ID: req.ID,
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteStudentReply{}, nil
}
func (s *StudentService) GetStudent(ctx context.Context, req *pb.GetStudentRequest) (*pb.GetStudentReply, error) {
	stu, err := s.uc.GetStudent(ctx, &biz.Student{ID: req.ID, Name: req.Name})
	if err != nil {
		return nil, err
	}
	return &pb.GetStudentReply{
		ID:     stu.ID,
		Name:   stu.Name,
		Status: stu.Status,
		Info:   stu.Info,
		//UpdatedAt: strconv.FormatInt(stu.CreatedAt.Unix(), 10),
		UpdatedAt: stu.UpdatedAt.String(),
		CreatedAt: stu.CreatedAt.String(),
	}, nil
}
func (s *StudentService) ListStudent(ctx context.Context, req *pb.ListStudentRequest) (*pb.ListStudentReply, error) {
	skuList, pageInfo, err := s.uc.ListStudent(ctx,
		&biz.PageInfo{PageNo: req.PageNo, PageSize: req.PageSize},
		&biz.Student{ID: req.Id, Name: req.Name, Info: req.Info, Status: req.Status})
	if err != nil {
		return nil, err
	}

	var studentReplyList []*pb.StudentReply
	for _, student := range skuList {
		studentReplyList = append(studentReplyList, &pb.StudentReply{
			ID:     student.ID,
			Name:   student.Name,
			Info:   student.Info,
			Status: student.Status,
		})
	}

	return &pb.ListStudentReply{
		PageNo:     pageInfo.PageNo,
		PageSize:   pageInfo.PageSize,
		TotalCount: pageInfo.TotalCount,
		Data:       studentReplyList,
	}, nil
}
