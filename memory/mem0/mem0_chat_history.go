package mem0

import (
	"context"
	"fmt"
	"log"

	"github.com/bytectlgo/mem0-go/client"
	"github.com/bytectlgo/mem0-go/types"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

// ChatMessageHistory is a struct that stores chat messages using Mem0.
type ChatMessageHistory struct {
	Mem0Client  *client.MemoryClient
	UserID      string
	HumanPrefix string
	AIPrefix    string
}

// Statically assert that Mem0ChatMessageHistory implement the chat message history interface.
var _ schema.ChatMessageHistory = &ChatMessageHistory{}

// NewMem0ChatMessageHistory creates a new Mem0ChatMessageHistory using chat message options.
func NewMem0ChatMessageHistory(mem0Client *client.MemoryClient, userID string, options ...ChatMessageHistoryOption) *ChatMessageHistory {
	messageHistory := applyMem0ChatHistoryOptions(options...)
	messageHistory.Mem0Client = mem0Client
	messageHistory.UserID = userID
	return messageHistory
}

func (h *ChatMessageHistory) messagesFromMem0Messages(mem0Messages []types.Memory) []llms.ChatMessage {
	var chatMessages []llms.ChatMessage
	for _, mem0Memory := range mem0Messages {
		for _, message := range mem0Memory.Messages {
			switch message.Role {
			case "user":
				chatMessages = append(chatMessages, llms.HumanChatMessage{Content: message.Content})
			case "assistant":
				chatMessages = append(chatMessages, llms.AIChatMessage{Content: message.Content})
			case "tool", "function":
				chatMessages = append(chatMessages, llms.ToolChatMessage{Content: message.Content})
			default:
				log.Print(fmt.Errorf("unknown role: %s", message.Role))
				continue
			}
		}
	}
	return chatMessages
}

func (h *ChatMessageHistory) messagesToMem0Messages(messages []llms.ChatMessage) []types.Message {
	var mem0Messages []types.Message
	for _, m := range messages {
		mem0Message := types.Message{
			Content: m.GetContent(),
		}
		switch m.GetType() {
		case llms.ChatMessageTypeHuman:
			mem0Message.Role = "user"
		case llms.ChatMessageTypeAI:
			mem0Message.Role = "assistant"
		case llms.ChatMessageTypeFunction:
			mem0Message.Role = "function"
		case llms.ChatMessageTypeTool:
			mem0Message.Role = "tool"
		default:
			log.Print(fmt.Errorf("unknown message type: %s", m.GetType()))
			continue
		}
		mem0Messages = append(mem0Messages, mem0Message)
	}
	return mem0Messages
}

// Messages returns all messages stored.
func (h *ChatMessageHistory) Messages(ctx context.Context) ([]llms.ChatMessage, error) {
	searchOptions := &types.SearchOptions{
		MemoryOptions: types.MemoryOptions{
			UserID: h.UserID,
		},
	}
	mem0Memories, err := h.Mem0Client.GetAll(searchOptions)
	if err != nil {
		return nil, err
	}

	messages := h.messagesFromMem0Messages(mem0Memories)

	// Add system context if available
	var systemPromptContent string
	for _, memory := range mem0Memories {
		if memory.Memory != "" {
			systemPromptContent += fmt.Sprintf("%s\n", memory.Memory)
		}
	}

	if systemPromptContent != "" {
		// Add system prompt to the beginning of the messages.
		messages = append(
			[]llms.ChatMessage{
				llms.SystemChatMessage{
					Content: systemPromptContent,
				},
			},
			messages...,
		)
	}

	return messages, nil
}

// AddAIMessage adds an AIMessage to the chat message history.
func (h *ChatMessageHistory) AddAIMessage(ctx context.Context, text string) error {
	mem0Messages := h.messagesToMem0Messages(
		[]llms.ChatMessage{
			llms.AIChatMessage{Content: text},
		},
	)

	memoryOptions := types.MemoryOptions{
		UserID: h.UserID,
	}

	_, err := h.Mem0Client.Add(mem0Messages, memoryOptions)
	if err != nil {
		return err
	}
	return nil
}

// AddUserMessage adds a user message to the chat message history.
func (h *ChatMessageHistory) AddUserMessage(ctx context.Context, text string) error {
	mem0Messages := h.messagesToMem0Messages(
		[]llms.ChatMessage{
			llms.HumanChatMessage{Content: text},
		},
	)

	memoryOptions := types.MemoryOptions{
		UserID: h.UserID,
	}

	_, err := h.Mem0Client.Add(mem0Messages, memoryOptions)
	if err != nil {
		return err
	}
	return nil
}

func (h *ChatMessageHistory) Clear(ctx context.Context) error {
	memoryOptions := types.MemoryOptions{
		UserID: h.UserID,
	}

	err := h.Mem0Client.DeleteAll(memoryOptions)
	if err != nil {
		return err
	}
	return nil
}

func (h *ChatMessageHistory) AddMessage(ctx context.Context, message llms.ChatMessage) error {
	mem0Messages := h.messagesToMem0Messages([]llms.ChatMessage{message})

	memoryOptions := types.MemoryOptions{
		UserID: h.UserID,
	}

	_, err := h.Mem0Client.Add(mem0Messages, memoryOptions)
	if err != nil {
		return err
	}
	return nil
}

func (*ChatMessageHistory) SetMessages(_ context.Context, _ []llms.ChatMessage) error {
	return nil
}
