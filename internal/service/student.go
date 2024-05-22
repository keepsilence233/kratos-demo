package service

import (
	"context"
	"kratos-demo/internal/biz"

	pb "kratos-demo/api/student/v1"
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
	return &pb.CreateStudentReply{}, nil
}
func (s *StudentService) UpdateStudent(ctx context.Context, req *pb.UpdateStudentRequest) (*pb.UpdateStudentReply, error) {
	return &pb.UpdateStudentReply{}, nil
}
func (s *StudentService) DeleteStudent(ctx context.Context, req *pb.DeleteStudentRequest) (*pb.DeleteStudentReply, error) {
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
	}, nil
}
func (s *StudentService) ListStudent(ctx context.Context, req *pb.ListStudentRequest) (*pb.ListStudentReply, error) {
	return &pb.ListStudentReply{}, nil
}
