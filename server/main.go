package main

import (
	"context"
	"errors"
	"gRPC-PostgreSQL-REST/proto"
	"gRPC-PostgreSQL-REST/server/helpers"
	"gRPC-PostgreSQL-REST/server/interfaces"
	"gRPC-PostgreSQL-REST/server/lib/database"
	"gRPC-PostgreSQL-REST/server/model"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func main() {

	log.Println("Starting User Service...")
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}

	err = initGormClient()
	if err != nil {
		log.Fatalf("DB failure", err)
	}

	srv := grpc.NewServer()
	proto.RegisterCrudServiceServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}

}

func (s *server) Create(ctx context.Context, request *proto.CreateRequest) (*proto.Response, error) {
	var user model.User
	user.Name, user.Email, user.PhoneNo, user.Organisation, user.Password = request.GetName(), request.GetEmail(), request.GetPhoneNo(), request.GetOrganization(), request.GetPassword()
	user.Password = helpers.Encrypt(user.Password)

	userExists := interfaces.DBEngine.CheckUser(user.Email)
	if userExists != nil {
		log.Println(userExists)
		return &proto.Response{Code: 400, Message: "User Already exists"}, errors.New("User already exists")
	}
	err := interfaces.DBEngine.CreateUser(user)
	if err != nil {
		return &proto.Response{Code: 400, Message: "DB Error"}, errors.New("DB Error")
	}
	return &proto.Response{Code: 200, Message: "OK"}, nil
}

func (s *server) Login(ctx context.Context, request *proto.LoginRequest) (*proto.TokenResponse, error) {
	var user model.User
	email, password := request.GetEmail(), request.GetPassword()
	password = helpers.Encrypt(password)
	userExists, userID := interfaces.DBEngine.Authenticate(email, password)
	if userExists != nil {
		log.Println(userExists)
		return &proto.TokenResponse{Code: 400, Token: ""}, nil
	}
	user.UserID = userID
	user.Email = email
	user.Password = password
	token, tokenErr := helpers.CreateToken(user)

	if tokenErr != nil {
		return &proto.TokenResponse{Code: 400, Token: ""}, nil
	}
	return &proto.TokenResponse{Code: 200, Token: token}, nil
}

func (s *server) GetUser(ctx context.Context, request *proto.TokenRequest) (*proto.UserResponse, error) {
	var user model.User
	mapClaims, tokenErr := helpers.ValidateToken(request.GetToken())
	if tokenErr != nil {
		log.Println(tokenErr)
		return &proto.UserResponse{}, errors.New("Invalid token")
	}

	userID := int(mapClaims["userID"].(float64))

	user, err := interfaces.DBEngine.GetUser(userID)
	if err != nil {
		return &proto.UserResponse{}, errors.New("Invalid token")
	}

	return &proto.UserResponse{Name: user.Name, Email: user.Email, PhoneNo: user.PhoneNo, Organization: user.Organisation}, nil
}

func (s *server) UpdateUser(ctx context.Context, request *proto.UpdateRequest) (*proto.Response, error) {
	phoneNo, organisation := request.GetPhoneNo(), request.GetOrganization()
	mapClaims, tokenErr := helpers.ValidateToken(request.GetToken())
	if tokenErr != nil {
		log.Println(tokenErr)
		return &proto.Response{Code: 400, Message: "Invalid Token"}, errors.New("Invalid token")
	}

	userID := int(mapClaims["userID"].(float64))

	err := interfaces.DBEngine.UpdateUser(phoneNo, organisation, userID)
	if err != nil {
		return &proto.Response{Code: 400, Message: "DB error"}, errors.New("DB error")
	}

	return &proto.Response{Code: 200, Message: "OK"}, nil
}

func (s *server) DeleteUser(ctx context.Context, request *proto.TokenRequest) (*proto.Response, error) {
	mapClaims, tokenErr := helpers.ValidateToken(request.GetToken())
	if tokenErr != nil {
		log.Println(tokenErr)
		return &proto.Response{Code: 400, Message: "Invalid Token"}, errors.New("Invalid token")
	}

	userID := int(mapClaims["userID"].(float64))

	err := interfaces.DBEngine.DeleteUser(userID)
	if err != nil {
		return &proto.Response{Code: 400, Message: "DB error"}, errors.New("DB error")
	}

	return &proto.Response{Code: 200, Message: "OK"}, nil
}

func initGormClient() error {
	log.Println("Initiating DB conn")
	var config model.DBConfig
	err := godotenv.Load()
	if err != nil {
		return err
	}

	config.User = os.Getenv("DBUSER")
	log.Println(config.User)
	config.DBName = os.Getenv("DB")
	config.Password = os.Getenv("PASSWORD")
	config.Host = os.Getenv("HOST")
	config.Port = os.Getenv("PORT")
	interfaces.DBEngine = new(database.DBRepo)
	err = interfaces.DBEngine.DBConnect(config)
	return err
}
