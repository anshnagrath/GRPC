package main

import (
	"RPC/blog/blogpb"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	objectid "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
}
type blogItem struct {
	ID       objectid.ObjectID `bson:"_id:,omitempty"`
	AuthorID string            `bson:"author_id"`

	Content string `bson:"content"`
	Title   string `bson:"title"`
}

var collection *mongo.Collection

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := req.GetBlog()
	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
	}
	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal Server Error:%v", err),
		)
	}

	oid, ok := res.InsertedID.(objectid.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Can't Convert to oid:%v", err),
		)

	}
	resp := &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Content:  blog.GetContent(),
			Title:    blog.GetTitle(),
		},
	}
	return resp, nil

}
func (*server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	blogId := req.GetBlogId()
	_id, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		log.Fatalf("Error while converting the blog id : %v", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Error while coverting string to hex"))
	}
	data := &blogItem{}
	filter := bson.M{"_id": _id}
	res := collection.FindOne(context.Background(), filter)
	if decodeerr := res.Decode(data); decodeerr != nil {
		log.Fatalf("No Result found: %v", decodeerr)
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("No Result found"))
	}
	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorID,
			Title:    data.Title,
			Content:  data.Content,
		},
	}, nil

}
func (*server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	blog := req.GetBlog()

	blogId, err := primitive.ObjectIDFromHex(blog.GetId())

	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Unable to Convert the String to Hex"))
	}

	filter := bson.M{"_id": blogId}
	_, replaceError := collection.ReplaceOne(context.Background(), filter, req.GetBlog())
	if replaceError != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Unable to Convert the String to Hex"))
	}

	fetch := collection.FindOne(context.Background(), filter)
	bgItem := &blogItem{}
	decodeErr := fetch.Decode(bgItem)
	if decodeErr != nil {
		fmt.Println("DecodeErr", decodeErr)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Unable to decode the updated response"))
	}

	resp := &blogpb.UpdateBlogResponse{Blog: &blogpb.Blog{
		Id:       bgItem.ID.Hex(),
		AuthorId: bgItem.AuthorID,
		Title:    bgItem.Title,
		Content:  bgItem.Content,
	}}

	return resp, nil

}
func (*server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	blogId := req.GetBlogId()
	Hex, hexerr := primitive.ObjectIDFromHex(blogId)
	if hexerr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Unable to Convert the String to Hex"))
	}
	filter := bson.M{"_id": Hex}
	res, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Error While delete the  document"))
	}
	if res.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Not found in the DataBase"))
	}
	resp := &blogpb.DeleteBlogResponse{
		BlogId: blogId,
	}
	return resp, nil
}

func (*server) ListBlog(req *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {
	cursor, err := collection.Find(context.Background(), bson.D{})
	fmt.Println(cursor)
	if err != nil {
		fmt.Println(err)
		return status.Errorf(codes.Internal, fmt.Sprintf("Unable to fetch"))
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		data := &blogItem{}
		err := cursor.Decode(data)
		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("Decode Error"))
		}
		stream.Send(&blogpb.ListBlogResponse{Blog: &blogpb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorID,
			Title:    data.Title,
			Content:  data.Content,
		}})

		
	}
	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Curesor error:%v", err))
	}
	return nil

}
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("Error while Connecting to mongodb:%v", err)
	}
	contErr := client.Connect(context.TODO())
	if contErr != nil {
		log.Fatalf("Error while Connecting to mongodb:%v", contErr)
	}
	collection = client.Database("mydb").Collection("blog")
	//To Disconnect Mongo
	defer client.Disconnect(context.TODO())

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Error while making a tcp connection:%v", err)
	}
	s := grpc.NewServer()
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Stating the server")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Error while making a tcp connection:%v", err)
		}
	}()
	// Blocking Channel
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	// Blocking the code
	<-ch
	fmt.Println("Stopping the sever")
	s.Stop()
	fmt.Println("Clossing the listner")
	lis.Close()

}
