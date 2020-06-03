package main

import (
	"RPC/blog/blogpb"
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error While Connecting to Client:%v \n", err)
	}
	defer cc.Close()
	c := blogpb.NewBlogServiceClient(cc)
	blog := &blogpb.Blog{

		AuthorId: "ansh",
		Title:    "New Blog",
		Content:  "New Content",
	}
	creatBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Error While Connecting to Client:%v \n", err)
	}
	fmt.Printf("Blog has been created : %v", creatBlogRes)
	blogID := creatBlogRes.GetBlog().GetId()
	fmt.Println("Reading Blogs Now", blogID)

	// _, readErr := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "scdcscdsc"})

	// if readErr != nil {
	// 	fmt.Printf("Error While Reading Blogs : %v", readErr)
	// }

	resp, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: blogID})
	if err2 != nil {
		fmt.Printf("Error ===================: %v", err2)
	}
	fmt.Printf("Response after Reading Blogs : %v", resp)
	newblog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Changed Autor",
		Title:    "New Edited Blog (My first blog)",
		Content:  "Content of first blog",
	}
	_, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: newblog,
	})
	if updateErr != nil {
		fmt.Printf("Error happend while updating: %v", updateErr)
	} else {

		// deleteResp, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
		// 	BlogId: blogID,
		// })
		// if deleteErr != nil {
		// 	fmt.Printf("Error happend while updating: %v", deleteErr)
		// } else {
		// 	fmt.Printf("Blog Deleted: %v", deleteResp)
		// }
		ListBlogs, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
		if err != nil {
			log.Fatalln("Error while fetching all blogs")
		}
		for {
			res, err := ListBlogs.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatalf("Error while fetching all blogs:%v", err)
			} else {
				fmt.Println(res.GetBlog())
			}
		}
	}

}
