package mem0

// MemoryOption Mem0MemoryOption is a function for creating new buffer
// with other than the default values.
type MemoryOption func(b *Memory)

// WithReturnMessages is an option for specifying should it return messages.
func WithReturnMessages(returnMessages bool) MemoryOption {
	return func(b *Memory) {
		b.ReturnMessages = returnMessages
	}
}

// WithInputKey is an option for specifying the input key.
func WithInputKey(inputKey string) MemoryOption {
	return func(b *Memory) {
		b.InputKey = inputKey
	}
}

// WithOutputKey is an option for specifying the output key.
func WithOutputKey(outputKey string) MemoryOption {
	return func(b *Memory) {
		b.OutputKey = outputKey
	}
}

// WithHumanPrefix is an option for specifying the human prefix. Will be passed as role for the message to mem0.
func WithHumanPrefix(humanPrefix string) MemoryOption {
	return func(b *Memory) {
		b.HumanPrefix = humanPrefix
	}
}

// WithAIPrefix is an option for specifying the AI prefix. Will be passed as role for the message to mem0.
func WithAIPrefix(aiPrefix string) MemoryOption {
	return func(b *Memory) {
		b.AIPrefix = aiPrefix
	}
}

// WithMemoryKey is an option for specifying the memory key.
func WithMemoryKey(memoryKey string) MemoryOption {
	return func(b *Memory) {
		b.MemoryKey = memoryKey
	}
}

func applyMem0MemoryOptions(opts ...MemoryOption) *Memory {
	m := &Memory{
		ReturnMessages: true,
		InputKey:       "",
		OutputKey:      "",
		HumanPrefix:    "Human",
		AIPrefix:       "AI",
		MemoryKey:      "history",
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}
