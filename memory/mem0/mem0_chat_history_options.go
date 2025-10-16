package mem0

// ChatMessageHistoryOption is a function for creating new chat message history
// with other than the default values.
type ChatMessageHistoryOption func(m *ChatMessageHistory)

// WithChatHistoryHumanPrefix is an option for specifying the human prefix. Will be passed as role for the message to mem0.
func WithChatHistoryHumanPrefix(humanPrefix string) ChatMessageHistoryOption {
	return func(b *ChatMessageHistory) {
		b.HumanPrefix = humanPrefix
	}
}

// WithChatHistoryAIPrefix is an option for specifying the AI prefix. Will be passed as role for the message to mem0.
func WithChatHistoryAIPrefix(aiPrefix string) ChatMessageHistoryOption {
	return func(b *ChatMessageHistory) {
		b.AIPrefix = aiPrefix
	}
}

func applyMem0ChatHistoryOptions(options ...ChatMessageHistoryOption) *ChatMessageHistory {
	h := &ChatMessageHistory{
		HumanPrefix: "Human",
		AIPrefix:    "AI",
	}

	for _, option := range options {
		option(h)
	}

	return h
}
