package mem0

import (
	"context"
	"testing"

	"github.com/bytectlgo/mem0-go/client"
	"github.com/bytectlgo/mem0-go/types"
	"github.com/tmc/langchaingo/llms"
)

// MockMem0Client implements a simple mock for testing
type MockMem0Client struct {
	memories []types.Memory
}

func (m *MockMem0Client) Add(messages interface{}, options types.MemoryOptions) ([]types.Memory, error) {
	mem0Messages, ok := messages.([]types.Message)
	if !ok {
		return nil, nil
	}

	memory := types.Memory{
		ID:       "test-id",
		UserID:   options.UserID,
		Messages: mem0Messages,
		Memory:   "Test memory content",
	}

	m.memories = append(m.memories, memory)
	return []types.Memory{memory}, nil
}

func (m *MockMem0Client) GetAll(options *types.SearchOptions) ([]types.Memory, error) {
	var result []types.Memory
	for _, memory := range m.memories {
		if memory.UserID == options.UserID {
			result = append(result, memory)
		}
	}
	return result, nil
}

func (m *MockMem0Client) DeleteAll(options types.MemoryOptions) error {
	var filtered []types.Memory
	for _, memory := range m.memories {
		if memory.UserID != options.UserID {
			filtered = append(filtered, memory)
		}
	}
	m.memories = filtered
	return nil
}

// Implement other required methods with no-op for testing
func (m *MockMem0Client) Get(memoryID string) (*types.Memory, error) {
	return nil, nil
}

func (m *MockMem0Client) Search(query string, options *types.SearchOptions) ([]types.Memory, error) {
	return m.memories, nil
}

func (m *MockMem0Client) Update(memoryID string, message string) ([]types.Memory, error) {
	return nil, nil
}

func (m *MockMem0Client) Delete(memoryID string) error {
	return nil
}

func (m *MockMem0Client) BatchDelete(memoryIDs []string) error {
	return nil
}

func (m *MockMem0Client) BatchUpdate(memories []types.MemoryUpdateBody) error {
	return nil
}

func (m *MockMem0Client) History(memoryID string) ([]types.MemoryHistory, error) {
	return nil, nil
}

func (m *MockMem0Client) Feedback(payload types.FeedbackPayload) error {
	return nil
}

func (m *MockMem0Client) Users() (*types.AllUsers, error) {
	return nil, nil
}

func (m *MockMem0Client) DeleteUser(entityID string) error {
	return nil
}

func (m *MockMem0Client) DeleteUsers() error {
	return nil
}

func (m *MockMem0Client) GetProject(options types.ProjectOptions) (*types.ProjectResponse, error) {
	return nil, nil
}

func (m *MockMem0Client) UpdateProject(payload types.PromptUpdatePayload) error {
	return nil
}

func (m *MockMem0Client) CreateWebhook(webhook types.WebhookPayload) (*types.Webhook, error) {
	return nil, nil
}

func (m *MockMem0Client) GetWebhooks(projectID string) ([]types.Webhook, error) {
	return nil, nil
}

func (m *MockMem0Client) UpdateWebhook(webhook types.WebhookPayload) error {
	return nil
}

func (m *MockMem0Client) DeleteWebhook(webhookID string) error {
	return nil
}

func createMockMem0Client() *client.MemoryClient {
	// Note: This is a simplified mock. In a real implementation, you would need to
	// properly mock the client interface or use a testing framework that supports it.
	return nil
}

func TestNewMemory(t *testing.T) {
	t.Parallel()

	// Since we can't easily mock the client interface, we'll test the options logic
	m := &Memory{}
	m = applyMem0MemoryOptions(
		WithReturnMessages(false),
		WithMemoryKey("custom_history"),
		WithHumanPrefix("User"),
		WithAIPrefix("Assistant"),
		WithInputKey("input"),
		WithOutputKey("output"),
	)

	if m.ReturnMessages {
		t.Errorf("Expected ReturnMessages to be false")
	}
	if m.MemoryKey != "custom_history" {
		t.Errorf("Expected MemoryKey to be 'custom_history', got %s", m.MemoryKey)
	}
	if m.HumanPrefix != "User" {
		t.Errorf("Expected HumanPrefix to be 'User', got %s", m.HumanPrefix)
	}
	if m.AIPrefix != "Assistant" {
		t.Errorf("Expected AIPrefix to be 'Assistant', got %s", m.AIPrefix)
	}
	if m.InputKey != "input" {
		t.Errorf("Expected InputKey to be 'input', got %s", m.InputKey)
	}
	if m.OutputKey != "output" {
		t.Errorf("Expected OutputKey to be 'output', got %s", m.OutputKey)
	}
}

func TestNewMem0ChatMessageHistory(t *testing.T) {
	t.Parallel()

	h := NewMem0ChatMessageHistory(nil, "test-user")

	if h.UserID != "test-user" {
		t.Errorf("Expected UserID to be 'test-user', got %s", h.UserID)
	}
	if h.HumanPrefix != "Human" {
		t.Errorf("Expected default HumanPrefix to be 'Human', got %s", h.HumanPrefix)
	}
	if h.AIPrefix != "AI" {
		t.Errorf("Expected default AIPrefix to be 'AI', got %s", h.AIPrefix)
	}
}

func TestNewMem0ChatMessageHistoryWithOptions(t *testing.T) {
	t.Parallel()

	h := NewMem0ChatMessageHistory(
		nil,
		"test-user",
		WithChatHistoryHumanPrefix("User"),
		WithChatHistoryAIPrefix("Bot"),
	)

	if h.HumanPrefix != "User" {
		t.Errorf("Expected HumanPrefix to be 'User', got %s", h.HumanPrefix)
	}
	if h.AIPrefix != "Bot" {
		t.Errorf("Expected AIPrefix to be 'Bot', got %s", h.AIPrefix)
	}
}

func TestMessagesFromMem0Messages(t *testing.T) {
	t.Parallel()

	h := &ChatMessageHistory{}

	mem0Memories := []types.Memory{
		{
			Messages: []types.Message{
				{Role: "user", Content: "Hello"},
				{Role: "assistant", Content: "Hi there"},
				{Role: "function", Content: "Function result"},
			},
		},
	}

	chatMessages := h.messagesFromMem0Messages(mem0Memories)

	if len(chatMessages) != 3 {
		t.Errorf("Expected 3 messages, got %d", len(chatMessages))
	}

	if human, ok := chatMessages[0].(llms.HumanChatMessage); !ok || human.Content != "Hello" {
		t.Errorf("Expected first message to be HumanChatMessage with content 'Hello'")
	}

	if ai, ok := chatMessages[1].(llms.AIChatMessage); !ok || ai.Content != "Hi there" {
		t.Errorf("Expected second message to be AIChatMessage with content 'Hi there'")
	}

	if tool, ok := chatMessages[2].(llms.ToolChatMessage); !ok || tool.Content != "Function result" {
		t.Errorf("Expected third message to be ToolChatMessage with content 'Function result'")
	}
}

func TestMessagesToMem0Messages(t *testing.T) {
	t.Parallel()

	h := &ChatMessageHistory{
		HumanPrefix: "User",
		AIPrefix:    "Bot",
	}

	chatMessages := []llms.ChatMessage{
		llms.HumanChatMessage{Content: "Hello"},
		llms.AIChatMessage{Content: "Hi there"},
		llms.FunctionChatMessage{Content: "Function result"},
	}

	mem0Messages := h.messagesToMem0Messages(chatMessages)

	if len(mem0Messages) != 3 {
		t.Errorf("Expected 3 messages, got %d", len(mem0Messages))
	}

	if mem0Messages[0].Content != "Hello" || mem0Messages[0].Role != "user" {
		t.Errorf("Expected first message to be user role with content 'Hello'")
	}

	if mem0Messages[1].Content != "Hi there" || mem0Messages[1].Role != "assistant" {
		t.Errorf("Expected second message to be assistant role with content 'Hi there'")
	}

	if mem0Messages[2].Content != "Function result" || mem0Messages[2].Role != "function" {
		t.Errorf("Expected third message to be function role with content 'Function result'")
	}
}

func TestMemoryVariables(t *testing.T) {
	t.Parallel()

	m := &Memory{}
	m = applyMem0MemoryOptions(WithMemoryKey("custom_key"))

	ctx := context.Background()
	variables := m.MemoryVariables(ctx)

	expected := []string{"custom_key"}
	if len(variables) != 1 || variables[0] != "custom_key" {
		t.Errorf("Expected variables %v, got %v", expected, variables)
	}
}

func TestGetMemoryKey(t *testing.T) {
	t.Parallel()

	m := &Memory{}
	m = applyMem0MemoryOptions(WithMemoryKey("test_key"))

	ctx := context.Background()
	key := m.GetMemoryKey(ctx)

	if key != "test_key" {
		t.Errorf("Expected memory key 'test_key', got %s", key)
	}
}

// TestChatMessageHistoryMethods tests the ChatMessageHistory implementation methods
func TestChatMessageHistoryMethods(t *testing.T) {
	t.Parallel()

	t.Run("SetMessages", func(t *testing.T) {
		h := &ChatMessageHistory{}
		err := h.SetMessages(context.Background(), []llms.ChatMessage{})
		if err != nil {
			t.Errorf("SetMessages should return nil, got %v", err)
		}
	})
}
