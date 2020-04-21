package main

import (
	"encoding/json"
	"log"
	"gRPC-PostgreSQL-REST/proto"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"google.golang.org/grpc"
)

type user struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PhoneNo      string `json:"phoneNo"`
	Password     string `json:"password"`
	Organization string `json:"organisation"`
}

var conn, _ = grpc.Dial("localhost:4040", grpc.WithInsecure())

func main() {

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
	}))
	e.POST("/create", create)
	e.POST("/login", login)
	e.GET("/getUser", getUser)
	e.PUT("/updateUser", updateUser)
	e.DELETE("/deleteUser", deleteUser)

	e.Start(":8080")
}

func updateUser(ctx echo.Context) error {
	var user user
	client := proto.NewCrudServiceClient(conn)
	reqErr := json.NewDecoder(ctx.Request().Body).Decode(&user)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	defer ctx.Request().Body.Close()
	req := ctx.Request().Header
	token := req.Get("token")
	if reqErr != nil {
		response := `{code:404, message:"Json Error"}`
		return json.NewEncoder(ctx.Response()).Encode(response)
	}
	rpcReq := &proto.UpdateRequest{PhoneNo: user.PhoneNo, Organization: user.Organization, Token: token}
	if response, err := client.UpdateUser(ctx.Request().Context(), rpcReq); err == nil {
		log.Println(response)
		return json.NewEncoder(ctx.Response()).Encode(response)
	} else {
		response := `{code:404, message:"Server Error"}`
		return json.NewEncoder(ctx.Response()).Encode(response)
	}
}

func deleteUser(ctx echo.Context) error {
	client := proto.NewCrudServiceClient(conn)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	defer ctx.Request().Body.Close()
	req := ctx.Request().Header
	token := req.Get("token")
	rpcReq := &proto.TokenRequest{Token: token}
	if response, err := client.DeleteUser(ctx.Request().Context(), rpcReq); err == nil {
		log.Println(response)
		return json.NewEncoder(ctx.Response()).Encode(response)
	}
	response := `{code:404, message:"Json Error"}`
	return json.NewEncoder(ctx.Response()).Encode(response)
}

func getUser(ctx echo.Context) error {
	client := proto.NewCrudServiceClient(conn)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	defer ctx.Request().Body.Close()
	req := ctx.Request().Header
	token := req.Get("token")
	rpcReq := &proto.TokenRequest{Token: token}
	if response, err := client.GetUser(ctx.Request().Context(), rpcReq); err == nil {
		log.Println(response)
		return json.NewEncoder(ctx.Response()).Encode(response)
	}
	response := `{code:404, message:"Json Error"}`
	return json.NewEncoder(ctx.Response()).Encode(response)
}

func login(ctx echo.Context) error {
	var user user
	client := proto.NewCrudServiceClient(conn)
	reqErr := json.NewDecoder(ctx.Request().Body).Decode(&user)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	defer ctx.Request().Body.Close()
	if reqErr != nil {
		response := `{code:404, message:"Json Error"}`
		return json.NewEncoder(ctx.Response()).Encode(response)
	}
	req := &proto.LoginRequest{Email: user.Email, Password: user.Password}
	if response, err := client.Login(ctx.Request().Context(), req); err == nil {
		log.Println(response)
		return json.NewEncoder(ctx.Response()).Encode(response)
	}
	response := `{code:404, message:"Json Error"}`
	return json.NewEncoder(ctx.Response()).Encode(response)
}

func create(ctx echo.Context) error {
	var user user
	client := proto.NewCrudServiceClient(conn)
	reqErr := json.NewDecoder(ctx.Request().Body).Decode(&user)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	defer ctx.Request().Body.Close()
	if reqErr != nil {
		response := `{code:404, message:"Json Error"}`
		return json.NewEncoder(ctx.Response()).Encode(response)
	}
	req := &proto.CreateRequest{Name: user.Name, Email: user.Email, PhoneNo: user.PhoneNo, Password: user.Password, Organization: user.Organization}
	if response, err := client.Create(ctx.Request().Context(), req); err == nil {
		log.Println(response)
		return json.NewEncoder(ctx.Response()).Encode(response)
	} else {
		response := `{code:404, message:"Server Error"}`
		return json.NewEncoder(ctx.Response()).Encode(response)
	}
}
