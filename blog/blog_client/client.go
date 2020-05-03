package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/johnnyaustor/go-grpc/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog client")

	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v\n", err)
	}

	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// CREATE BLOG
	fmt.Println("Createing the blog")
	blog := &blogpb.Blog{
		AuthorId: "Agus",
		Title:    "Helo Go Mongo",
		Content:  "Content helo go mongo",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})

	if err != nil {
		log.Fatalf("Unexpected error: %v\n", err)
	}
	fmt.Printf("Blog created: %v\n", createBlogRes)
	blogID := createBlogRes.GetBlog().GetId()

	// READ BLOG
	fmt.Println("Reading the blog")
	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "5eadc96c15ffa100cab55cb",
	})
	if err2 != nil {
		fmt.Printf("Error happened while reading: %v\n", err)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		fmt.Printf("Error happened while reading: %v\n", err)
	}

	fmt.Printf("Blog was read: %v\n", readBlogRes)

	// UPDATE BLOG
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Agus Triant",
		Title:    "Helo Go Mongo (edit)",
		Content:  "Content helo go mongo (edit)",
	}
	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if updateErr != nil {
		fmt.Printf("Error Happened while updating: %v\n", updateErr)
	}
	fmt.Printf("Blog Was update %v\n", updateRes)

	// DELETE BLOG
	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})
	if deleteErr != nil {
		fmt.Printf("Error Happened while Deleting: %v\n", deleteErr)
	}
	fmt.Printf("Blog Was deleting %v\n", deleteRes)

	// LIST BLOG
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling LISTBlog rpc: %v\n", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v\n", err)
		}
		fmt.Println(res.GetBlog())
	}
}
