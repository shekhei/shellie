package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"shellie/common"
	"shellie/pb"
	"strings"
	"syscall"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"google.golang.org/grpc"
)

var client *openai.Client

func verifyConfig(config common.ServiceConfig) error {
	if config.APIKey == "" {
		return fmt.Errorf("API key is not set")
	}
	if config.ChatCompletionEndpoint == "" {
		return fmt.Errorf("OpenAI endpoint is not set")
	}
	if config.Model == "" {
		return fmt.Errorf("model is not set")
	}
	return nil
}

// endpoint unix://some.sock/

type PromptSuggestionServer struct {
	pb.UnimplementedPromptSuggestionServer
	config *common.Config
}

type Suggestion struct {
	Command     string `json:"command"`
	Explanation string `json:"explanation"`
}

func (s *PromptSuggestionServer) Suggest(ctx context.Context, req *pb.SuggestRequest) (*pb.SuggestResponse, error) {
	startTime := time.Now()
	log.Println("start suggest: ", startTime)
	defer func() {
		log.Println("time taken: ", time.Since(startTime))
	}()
	reqParams := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(`You are a command suggestion assistant. You suggest the next shell command based on context.
			Always respond with exactly 3 command suggestions in the following JSON format, response must be a valid JSON:
			[
				{"command": "complete command here", "explanation": "brief explanation"},
				{"command": "alternative command here", "explanation": "brief explanation"},
				{"command": "another alternative here", "explanation": "brief explanation"}
			]
			When analyzing command history, give more weight to recent commands as they are more relevant to the current context.`),
			openai.UserMessage(fmt.Sprintf("Shell environment: %s", req.GetShell())),
			openai.UserMessage(fmt.Sprintf("Context:\n%s", req.GetContext())),
			openai.UserMessage("Do not wrap the json codes in JSON markers"),
			openai.UserMessage(fmt.Sprintf("Suggest completions that starts with '%s'", req.GetCommand())),
		}),
		Model: openai.F(s.config.Service.Model),
	}
	chatCompletion, err := client.Chat.Completions.New(ctx, reqParams)
	if err != nil {
		return nil, err
	}
	content := chatCompletion.Choices[0].Message.Content
	var suggestions []Suggestion
	log.Println("content: ", content)
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "```json") {
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimSuffix(content, "```")
		log.Println("content: ", content)
	}
	if err := json.Unmarshal([]byte(content), &suggestions); err != nil {
		return nil, err
	}
	return &pb.SuggestResponse{
		Suggestion: suggestions[0].Command,
	}, nil
}
func main() {
	config, err := common.ReadOrCreateConfig()
	if err != nil {
		panic(err.Error())
	}
	if err := verifyConfig(config.Service); err != nil {
		panic(err.Error())
	}
	fmt.Println("Creating client using config: ", config)
	client = openai.NewClient(
		option.WithAPIKey(config.Service.APIKey),
		option.WithBaseURL(config.Service.ChatCompletionEndpoint),
		option.WithProject(config.Service.Project),
		option.WithOrganization(config.Service.Organization),
	)
	grpcServer := grpc.NewServer()
	u, err := url.Parse(config.Service.ListenAddress)
	if err != nil {
		panic(err.Error())
	}
	scheme := u.Scheme
	log.Printf("Listening on %s %s", scheme, u.Host)
	var lis net.Listener
	if scheme != "unix" && scheme != "http" {
		panic(fmt.Errorf("invalid scheme: %s, only supports unix or http", scheme))
	}
	if scheme == "unix" {
		lis, err = net.Listen("unix", u.Path)
	} else {
		lis, err = net.Listen("tcp", u.Host)
	}
	if err != nil {
		panic(err.Error())
	}
	defer lis.Close()
	pb.RegisterPromptSuggestionServer(grpcServer, &PromptSuggestionServer{
		config: config,
	})
	go func() {
		grpcServer.Serve(lis)
	}()
	// Use a buffered channel so we don't miss any signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received.
	s := <-c
	fmt.Println("Got signal:", s)
	grpcServer.Stop()
}
