package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"shellie/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// process command shell type from arguments
var shell = flag.String("shell", "", "shell type")
var command = flag.String("command", "", "command")
var pwd = flag.String("pwd", "", "current directory")

func main() {
	// gets a list of command from stdin, newline denotes a new command
	// the last command is the command that the user is currently entering
	f, err := os.OpenFile("/tmp/promptsuggestion.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	ctxInput := string(stdin)
	conn, err := grpc.NewClient("unix:///tmp/promptsuggestion.sock", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
	client := pb.NewPromptSuggestionClient(conn)
	response, err := client.Suggest(context.Background(), &pb.SuggestRequest{
		Shell:   *shell,
		Context: ctxInput,
		Command: *command,
		Pwd:     *pwd,
	})
	if err != nil {
		panic(err.Error())
	}
	topSuggestion := response.Suggestion
	fmt.Println(topSuggestion)
}
