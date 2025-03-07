package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"shellie/common"
	"shellie/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// process command shell type from arguments
var shell = flag.String("shell", "", "shell type")
var command = flag.String("command", "", "command")
var pwd = flag.String("pwd", "", "current directory")

func verifyConfig(config common.ClientConfig) error {
	if config.ServerAddress == "" {
		return fmt.Errorf("server address is not set")
	}
	return nil
}

func main() {
	flag.Parse()
	// gets a list of command from stdin, newline denotes a new command
	// the last command is the command that the user is currently entering
	config, err := common.ReadOrCreateConfig()
	if err != nil {
		panic(err.Error())
	}
	if err := verifyConfig(config.Client); err != nil {
		panic(err.Error())
	}
	f, err := os.OpenFile("/tmp/shellie-client.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
	u, err := url.Parse(config.Client.ServerAddress)
	if err != nil {
		panic(err.Error())
	}
	scheme := u.Scheme
	if scheme != "unix" && scheme != "http" {
		panic(fmt.Errorf("invalid scheme: %s, only supports unix or http", scheme))
	}
	conn, err := grpc.NewClient(config.Client.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
	client := pb.NewPromptSuggestionClient(conn)
	log.Printf("shell (%s), command (%s), pwd (%s)", *shell, *command, *pwd)
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
