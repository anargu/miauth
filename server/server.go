package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/dgrijalva/jwt-go"

	"github.com/anargu/miauth"
	pb "github.com/anargu/miauth/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var r *gin.Engine

func InitServer() {
	isDebug := os.Getenv("DEBUG") == "true"
	if !isDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	r = gin.Default()
	r.Static("/public", "./public")
	r.LoadHTMLGlob("./public/*.html")
	BindAPIRoutes()

	fmt.Printf("miauth server running at port: %v\n", miauth.Config.Port)
	err := r.Run(fmt.Sprintf(":%s", miauth.Config.Port))
	if err != nil {
		log.Fatal(err)
	}
}

func BindAPIRoutes() {
	authAPI := r.Group("/auth")
	{
		authAPI.POST("/login", LoginEndpoint)
		authAPI.POST("/signup", SignupEndpoint)
		authAPI.GET("/verify", verifyEndpoint)
		authAPI.POST("/token/refresh", tokenRefreshEndpoint)
	}
	forgotAPI := r.Group("/forgot")
	{
		forgotAPI.POST("/request", ForgotRequestEndpoint)
		forgotAPI.GET("/reset", forgotResetEndpoint)
		forgotAPI.POST("/reset", forgotResetEndpoint)
	}
	adminAPI := r.Group("/admin")
	{
		adminAPI.PUT("/update/templates", updateTemplatesEndpoint)
		adminAPI.POST("/revoke_all", RevokeAllEndpoint)
	}
}

type ErrorResponsePayload struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	UserMessage      string `json:"user_message"`
}

func SendError(c *gin.Context, code int, err error) {
	c.JSON(code, err)
	return
}

func ErrorResponse(c *gin.Context, code int, err error, description string, userMessage string) {
	var _err = err
	if _err == nil {
		_err = errors.New(description)
	}
	c.JSON(code, miauth.ErrorMessage{
		Name:             _err.Error(),
		ErrorDescription: description,
		UserMessage:      userMessage,
	})
	return
}

// server is used to implement helloworld.GreeterServer.
type grpcServer struct {
	pb.UnimplementedMiAuthServer
}

func (s *grpcServer) VerifyToken(ctx context.Context, in *pb.ValidationInput) (*pb.ValidationResult, error) {
	log.Printf("Received: %v", in.AccessToken)

	if in != nil {
		var tk *jwt.Token
		var err error
		if tk, err = miauth.VerifyAccessToken(in.AccessToken); err != nil || !tk.Valid {
			if err == nil {
				err = errors.New("invalid token")
				return nil, err
			} else {
				return nil, err
			}
		} else {
			claims := tk.Claims.(jwt.MapClaims)
			return &pb.ValidationResult{IsOk: true, UserEmail: claims["user_email"].(string), UserMiauthID: claims["userId"].(string)}, nil
		}
	}
	return nil, errors.New("not input")
}

func InitGrpcServer() {
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		panic("Not GRPC_PORT setted")
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMiAuthServer(s, &grpcServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
