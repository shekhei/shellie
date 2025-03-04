package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"shellie/pb"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"google.golang.org/grpc"
)

type Config struct {
	OpenAIEndpoint string `toml:"openai_endpoint"`
	APIKey         string `toml:"api_key"`
	Organization   string `toml:"organization"`
	Project        string `toml:"project"`
	Model          string `toml:"model"`
}

var client *openai.Client
var config Config

func verifyConfig(config Config) error {
	if config.APIKey == "" {
		return fmt.Errorf("API key is not set")
	}
	if config.OpenAIEndpoint == "" {
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
			openai.UserMessage(fmt.Sprintf("Suggest completions for this partial command: %s", req.GetCommand())),
		}),
		Model: openai.F(config.Model),
	}
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), reqParams)
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
	// Read config from ~/.config/promptsuggestion.json
	HOME := os.Getenv("HOME")
	// create default config
	// check if ~/.config/promptsuggestion.json exists
	if _, err := os.Stat(fmt.Sprintf("%s/.config/promptsuggestion.toml", HOME)); os.IsNotExist(err) {
		config := Config{
			OpenAIEndpoint: "https://api.openai.com/v1",
		}
		content, err := toml.Marshal(config)
		if err != nil {
			panic(err.Error())
		}
		os.WriteFile(fmt.Sprintf("%s/.config/promptsuggestion.toml", HOME), content, 0644)
	}
	configBytes, err := os.ReadFile(fmt.Sprintf("%s/.config/promptsuggestion.toml", HOME))
	if err != nil {
		panic(err.Error())
	}
	if err := toml.Unmarshal(configBytes, &config); err != nil {
		panic(err.Error())
	}
	if err := verifyConfig(config); err != nil {
		panic(err.Error())
	}
	fmt.Println("Creating client using config: ", config)
	client = openai.NewClient(
		option.WithAPIKey(config.APIKey),
		option.WithBaseURL(config.OpenAIEndpoint),
		option.WithProject(config.Project),
		option.WithOrganization(config.Organization),
	)
	grpcServer := grpc.NewServer()
	lis, err := net.Listen("unix", "/tmp/promptsuggestion.sock")
	defer func() {
		os.Remove("/tmp/promptsuggestion.sock")
	}()
	if err != nil {
		panic(err.Error())
	}
	pb.RegisterPromptSuggestionServer(grpcServer, &PromptSuggestionServer{})
	grpcServer.Serve(lis)
}
