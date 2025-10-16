package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/0xDezzy/langchaingo-memory/memory/graphiti"
	"github.com/getzep/zep-go"
	zepClient "github.com/getzep/zep-go/client"
	"github.com/getzep/zep-go/option"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
)

// SimpleAgent represents a basic conversational agent with memory
type SimpleAgent struct {
	llm    llms.Model
	memory schema.Memory
}

// NewSimpleAgent creates a new agent with the specified LLM and memory
func NewSimpleAgent(llm llms.Model, memory schema.Memory) *SimpleAgent {
	return &SimpleAgent{
		llm:    llm,
		memory: memory,
	}
}

// Chat handles a conversation turn with the agent
func (a *SimpleAgent) Chat(ctx context.Context, userMessage string) (string, error) {
	// Save user message to memory
	err := a.memory.SaveContext(ctx, map[string]any{
		"input": userMessage,
	}, map[string]any{})
	if err != nil {
		return "", fmt.Errorf("failed to save user message: %w", err)
	}

	// Get conversation history from memory
	memoryVars, err := a.memory.LoadMemoryVariables(ctx, map[string]any{})
	if err != nil {
		return "", fmt.Errorf("failed to load memory: %w", err)
	}

	// Build prompt with memory context
	history, ok := memoryVars["chat_history"].([]llms.ChatMessage)
	if !ok {
		history = []llms.ChatMessage{}
	}

	// Create messages for LLM
	messageContents := make([]llms.MessageContent, 0, len(history)+1)

	// Convert history messages to MessageContent
	for _, msg := range history {
		switch msg.GetType() {
		case llms.ChatMessageTypeHuman:
			messageContents = append(messageContents, llms.TextParts(llms.ChatMessageTypeHuman, msg.GetContent()))
		case llms.ChatMessageTypeAI:
			messageContents = append(messageContents, llms.TextParts(llms.ChatMessageTypeAI, msg.GetContent()))
		case llms.ChatMessageTypeSystem:
			messageContents = append(messageContents, llms.TextParts(llms.ChatMessageTypeSystem, msg.GetContent()))
		}
	}

	// Add current user message
	messageContents = append(messageContents, llms.TextParts(llms.ChatMessageTypeHuman, userMessage))

	// Generate response
	var response strings.Builder
	_, err = a.llm.GenerateContent(ctx, messageContents, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		response.Write(chunk)
		return nil
	}))
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	assistantResponse := response.String()

	// Save assistant response to memory
	err = a.memory.SaveContext(ctx, map[string]any{}, map[string]any{
		"output": assistantResponse,
	})
	if err != nil {
		return "", fmt.Errorf("failed to save assistant message: %w", err)
	}

	return assistantResponse, nil
}

// ClearMemory clears the agent's memory
func (a *SimpleAgent) ClearMemory(ctx context.Context) error {
	return a.memory.Clear(ctx)
}

func main() {
	fmt.Println("=== Graphiti Agent Memory Example ===")
	fmt.Println()

	// Check for API key
	apiKey := os.Getenv("ZEP_API_KEY")
	if apiKey == "" {
		log.Fatal("ZEP_API_KEY environment variable is not set")
	}

	// Initialize Zep client
	zepClient := zepClient.NewClient(option.WithAPIKey(apiKey))

	// Create memory with session ID
	memory := graphiti.NewMemory(zepClient, "session-123",
		graphiti.WithMemoryKey("chat_history"),
		graphiti.WithHumanPrefix("User"),
		graphiti.WithAIPrefix("Assistant"),
		graphiti.WithInputKey("input"),
		graphiti.WithOutputKey("output"),
		graphiti.WithMemoryType(zep.MemoryTypePerpetual),
	)

	// Initialize Ollama LLM
	llm, err := ollama.New(ollama.WithModel("gemma3:4b"))
	if err != nil {
		log.Fatalf("Failed to create LLM: %v", err)
	}

	// Create agent
	agent := NewSimpleAgent(llm, memory)

	ctx := context.Background()

	// Example conversation
	conversations := []string{
		"Hi, I'm Sarah and I work as a software engineer at TechCorp",
		"I specialize in machine learning and recently deployed a recommendation system",
		"What's my job?",
		"What do I specialize in?",
	}

	for i, message := range conversations {
		fmt.Printf("User: %s\n", message)

		response, err := agent.Chat(ctx, message)
		if err != nil {
			log.Printf("Error in conversation %d: %v", i+1, err)
			continue
		}

		fmt.Printf("Assistant: %s\n", response)
		fmt.Println()
	}

	// Demonstrate memory persistence
	fmt.Println("=== Memory Persistence Demo ===")
	fmt.Println("Starting new conversation session...")
	fmt.Println()

	// Create new agent instance with same session ID to demonstrate persistence
	newAgent := NewSimpleAgent(llm, memory)

	persistenceTest := []string{
		"Do you remember my name?",
		"Where do I work?",
	}

	for i, message := range persistenceTest {
		fmt.Printf("User: %s\n", message)

		response, err := newAgent.Chat(ctx, message)
		if err != nil {
			log.Printf("Error in persistence test %d: %v", i+1, err)
			continue
		}

		fmt.Printf("Assistant: %s\n", response)
		fmt.Println()
	}

	fmt.Println("=== Example Complete ===")
}
