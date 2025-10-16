# Mem0 Agent Memory Example

This example demonstrates how to use the Mem0 memory implementation with LangChainGo to create an AI agent that maintains conversation history and long-term memory.

## Overview

The example shows:
- Setting up Mem0 client and memory
- Creating a simple conversational agent
- Maintaining conversation context across multiple interactions
- Retrieving relevant memories from past conversations

## Prerequisites

1. Get a Mem0 API key from [Mem0](https://mem0.ai)
2. Install and run Ollama with the gemma3:4b model:
   ```bash
   # Install Ollama
   curl -fsSL https://ollama.ai/install.sh | sh
   
   # Pull the model
   ollama pull gemma3:4b
   
   # Start Ollama server
   ollama serve
   ```
3. Set the API key as an environment variable:
   ```bash
   export MEM0_API_KEY="your-api-key-here"
   ```

## Usage

Run the example:
```bash
go run main.go
```

## Key Components

### Memory Setup
```go
// Initialize Mem0 client
client, err := client.NewMemoryClient(client.ClientOptions{
    APIKey: os.Getenv("MEM0_API_KEY"),
})

// Create memory with user ID
memory := mem0.NewMemory(client, "user-123",
    mem0.WithMemoryKey("chat_history"),
    mem0.WithHumanPrefix("User"),
    mem0.WithAIPrefix("Assistant"),
)
```

### Agent Integration
The memory is integrated into LangChainGo chains to maintain context:
- User messages are automatically stored
- AI responses are saved to memory
- Past context is retrieved for each new interaction

### Memory Features
- **Short-term memory**: Recent conversation messages
- **Long-term memory**: Important facts and context extracted by Mem0
- **User-specific**: Each user has their own memory space
- **Persistent**: Memories are stored across sessions

## Example Output

```
=== Mem0 Agent Memory Example ===

User: Hi, I'm John and I love hiking in the mountains
Assistant: Hello John! It's wonderful to meet someone who enjoys mountain hiking. There's something truly special about being up in the mountains, isn't there? The fresh air, the views, and the sense of accomplishment after a good hike. What's your favorite mountain trail you've been on?

User: I recently climbed Mount Whitney
Assistant: That's impressive! Mount Whitney is the highest peak in the contiguous United States at 14,505 feet. That must have been quite a challenging and rewarding experience. How did you prepare for such a high-altitude climb? Did you encounter any difficult conditions during your ascent?

User: What's my name?
Assistant: Your name is John. I remember from when you introduced yourself at the beginning of our conversation.

User: What do I like to do?
Assistant: You mentioned that you love hiking in the mountains, and you recently climbed Mount Whitney, which shows you're quite an accomplished mountaineer!
```

## Files

- `main.go` - Main example application
- `README.md` - This documentation