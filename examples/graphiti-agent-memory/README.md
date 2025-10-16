# Graphiti Agent Memory Example

This example demonstrates how to use the Graphiti memory implementation with LangChainGo to create an AI agent that maintains conversation history and long-term memory using Zep.

## Overview

The example shows:
- Setting up Zep client and Graphiti memory
- Creating a simple conversational agent
- Maintaining conversation context across multiple interactions
- Retrieving relevant memories and facts from past conversations

## Prerequisites

1. Get a Zep API key from [Zep](https://getzep.com)
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
   export ZEP_API_KEY="your-api-key-here"
   ```

## Usage

Run the example:
```bash
go run main.go
```

## Key Components

### Memory Setup
```go
// Initialize Zep client
client, err := zepClient.NewClient(zepClient.WithAPIKey(os.Getenv("ZEP_API_KEY")))

// Create memory with session ID
memory := graphiti.NewMemory(client, "session-123",
    graphiti.WithMemoryKey("chat_history"),
    graphiti.WithHumanPrefix("User"),
    graphiti.WithAIPrefix("Assistant"),
    graphiti.WithMemoryType(zep.MemoryTypePerpetual),
)
```

### Agent Integration
The memory is integrated into LangChainGo chains to maintain context:
- User messages are automatically stored in Zep
- AI responses are saved to memory
- Past context and relevant facts are retrieved for each new interaction
- Memory persists across sessions using Zep's memory management

### Memory Features
- **Short-term memory**: Recent conversation messages
- **Long-term memory**: Facts and context extracted by Zep's graph-based memory
- **Session-based**: Each conversation session has its own memory space
- **Persistent**: Memories are stored indefinitely in Zep
- **Fact extraction**: Zep automatically extracts and stores important facts

## Example Output

```
=== Graphiti Agent Memory Example ===

User: Hi, I'm Sarah and I work as a software engineer at TechCorp
Assistant: Hello Sarah! It's great to meet you. Working as a software engineer at TechCorp sounds interesting. What kind of projects do you work on there? I imagine you get to tackle some challenging technical problems.

User: I specialize in machine learning and recently deployed a recommendation system
Assistant: That's fascinating, Sarah! Machine learning is such an exciting field, and recommendation systems are particularly interesting because they directly impact user experience. Deploying a recommendation system at TechCorp must have involved working with large datasets and complex algorithms. What kind of recommendations does your system provide?

User: What's my job?
Assistant: You work as a software engineer at TechCorp, where you specialize in machine learning and recently deployed a recommendation system.

User: What do I specialize in?
Assistant: You specialize in machine learning, which you mentioned when you told me about your work at TechCorp and the recommendation system you recently deployed.
```

## Files

- `main.go` - Main example application
- `README.md` - This documentation