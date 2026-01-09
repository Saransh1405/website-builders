# Architecture & Flow Documentation

This document explains the complete architecture, flow, and logic of the Website Builder Backend service.

## Table of Contents

1. [Project Overview](#project-overview)
2. [Architecture Overview](#architecture-overview)
3. [Complete Request Flow](#complete-request-flow)
4. [Frontend Integration](#frontend-integration)
5. [Service Logic Breakdown](#service-logic-breakdown)
6. [Data Flow Diagrams](#data-flow-diagrams)
7. [API Contracts](#api-contracts)

---

## Project Overview

The Website Builder Backend is an AI-powered code generation service that:
- Accepts natural language prompts from users
- Determines the intent (what type of code to generate)
- Loads relevant templates from the filesystem
- Calls Claude AI API to generate code
- Post-processes the output (adds config files, package.json, etc.)
- Streams the response back to the frontend in real-time

---

## Architecture Overview

### Clean Architecture Layers

```
┌─────────────────────────────────────────┐
│         HTTP Layer (Gin)                │
│  ┌──────────────┐  ┌──────────────┐    │
│  │   Handler    │  │  Middleware  │    │
│  └──────────────┘  └──────────────┘    │
└─────────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│         Service Layer                   │
│  ┌────────────────────────────────┐    │
│  │  CodeGenerationService         │    │
│  │  (Orchestrator)                │    │
│  └────────────────────────────────┘    │
│  ┌──────────┐ ┌──────────┐            │
│  │  AI      │ │ Intent   │            │
│  │ Service  │ │ Parser   │            │
│  └──────────┘ └──────────┘            │
│  ┌──────────────────────────┐         │
│  │   PostProcessor          │         │
│  └──────────────────────────┘         │
└─────────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│         Repository Layer                │
│  ┌──────────────────────────────┐      │
│  │   TemplateRepository         │      │
│  │   (File System Access)       │      │
│  └──────────────────────────────┘      │
└─────────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│         Domain Layer                    │
│  ┌──────────┐ ┌──────────┐            │
│  │ Entities │ │Interfaces│            │
│  └──────────┘ └──────────┘            │
└─────────────────────────────────────────┘
```

### Key Components

1. **Handler Layer**: HTTP request handling, validation, SSE streaming
2. **Service Layer**: Business logic orchestration
3. **Repository Layer**: Data access (templates from filesystem)
4. **Domain Layer**: Core entities and interfaces

---

## Complete Request Flow

### Flow Diagram (Step-by-Step)

```
Frontend Request
    │
    ├─→ POST /api/v1/generate/stream
    │
    ↓
┌─────────────────────────────────────┐
│ 1. HTTP Handler                     │
│    - Validate request body          │
│    - Set SSE headers                │
│    - Create request context         │
└─────────────────────────────────────┘
    │
    ↓
┌─────────────────────────────────────┐
│ 2. CodeGenerationService            │
│    GenerateCodeStream()             │
└─────────────────────────────────────┘
    │
    ├─→ Step 2.1: Parse Intent
    │   ┌──────────────────────────┐
    │   │ IntentParserService      │
    │   │ - Analyze prompt         │
    │   │ - Detect component type  │
    │   │ - Extract component name │
    │   │ - Identify features      │
    │   └──────────────────────────┘
    │
    ├─→ Step 2.2: Load Template
    │   ┌──────────────────────────┐
    │   │ TemplateRepository       │
    │   │ - Select template path   │
    │   │ - Read from filesystem   │
    │   │ - Return template content│
    │   └──────────────────────────┘
    │
    ├─→ Step 2.3: Build System Prompt
    │   ┌──────────────────────────┐
    │   │ CodeGenerationService    │
    │   │ - Combine intent         │
    │   │ - Add requirements       │
    │   │ - Include options        │
    │   └──────────────────────────┘
    │
    ├─→ Step 2.4: Generate Code
    │   ┌──────────────────────────┐
    │   │ ClaudeAIService          │
    │   │ - Call Claude API        │
    │   │ - Stream response chunks │
    │   │ - Handle errors          │
    │   └──────────────────────────┘
    │   │
    │   ┌──────────────────────────┐
    │   │ Claude API (External)    │
    │   │ - Process prompt         │
    │   │ - Generate code          │
    │   │ - Stream back            │
    │   └──────────────────────────┘
    │
    └─→ Step 2.5: Post-Process
        ┌──────────────────────────┐
        │ PostProcessorService     │
        │ - Extract code blocks    │
        │ - Generate package.json  │
        │ - Create config files    │
        │ - Organize file structure│
        └──────────────────────────┘
    │
    ↓
┌─────────────────────────────────────┐
│ 3. HTTP Handler (Streaming)         │
│    - Format SSE events              │
│    - Send chunks to client          │
│    - Handle completion              │
└─────────────────────────────────────┘
    │
    ↓
Frontend receives streaming response
```

### Detailed Flow Steps

#### Step 1: Request Reception (HTTP Handler)

**File**: `internal/handler/http_handler.go`

1. Frontend sends POST request to `/api/v1/generate/stream`
2. Request body contains:
   ```json
   {
     "prompt": "Create a React button component with Tailwind CSS",
     "options": {
       "useTypeScript": true,
       "styleLibrary": "tailwind"
     }
   }
   ```
3. Handler validates request using Gin's binding
4. Sets SSE headers:
   - `Content-Type: text/event-stream`
   - `Cache-Control: no-cache`
   - `Connection: keep-alive`
5. Creates context with timeout
6. Calls `CodeGenerationService.GenerateCodeStream()`

#### Step 2: Intent Parsing

**File**: `internal/service/intent_parser_service.go`

**Logic**:
- Analyzes the prompt text for keywords
- Detects component type:
  - "react component" → `react-component`
  - "vue component" → `vue-component`
  - "hook" → `react-hook`
- Extracts component name using regex patterns:
  - "Create a Button component" → `Button`
  - "Build LoginForm" → `LoginForm`
- Identifies features:
  - "useState", "state" → `state` feature
  - "tailwind" → `tailwind` feature
  - "typescript" → `typescript` feature
- Returns `Intent` struct with confidence score

**Example Output**:
```go
Intent{
    Type: "react-component",
    ComponentName: "Button",
    Features: ["tailwind", "typescript"],
    Confidence: 0.9,
}
```

#### Step 3: Template Loading

**File**: `internal/repository/template_repository.go`

**Logic**:
- Based on parsed intent, selects template path:
  - `react-component` → `templates/react/component.tsx.template`
  - `react-hook` → `templates/react/hook.ts.template`
  - `vue-component` → `templates/vue/component.vue.template`
- Validates path (prevents directory traversal attacks)
- Reads template file from filesystem
- Returns template content as string

**Template Example**:
```tsx
import React from 'react';

interface Props {
  // Define your component props here
}

export const Component: React.FC<Props> = ({}) => {
  return (
    <div className="component">
      {/* Your component implementation */}
    </div>
  );
};
```

#### Step 4: System Prompt Building

**File**: `internal/service/code_generation_service.go`

**Logic**:
- Combines intent information
- Builds detailed system prompt:
  - Instructions for AI (use functional components, hooks, etc.)
  - TypeScript requirements if requested
  - Styling library instructions (Tailwind, styled-components, CSS)
  - Component name if extracted
  - Required features list
  - Best practices reminders

**Example System Prompt**:
```
You are an expert frontend developer specializing in React.
Generate production-ready, clean, and maintainable code.

Requirements:
- Use functional components with hooks
- Use TypeScript with proper type definitions
- Use Tailwind CSS for styling
- Component name should be: Button
- Follow React best practices
- Required features: tailwind, typescript
```

#### Step 5: AI Code Generation

**File**: `internal/service/ai_service.go`

**Logic**:
- Prepares Claude API request:
  - Model: `claude-3-5-sonnet-20241022` (configurable)
  - System prompt: Built in step 4
  - User message: Original prompt + template reference
  - Max tokens: 4096
  - Stream: true
- Makes HTTP POST request to Claude API
- Processes streaming response:
  - Each chunk contains text delta
  - Sends chunks to channel as they arrive
  - Handles errors and completion
- Returns two channels:
  - `codeChan`: Stream of code chunks
  - `errorChan`: Stream of errors

**Claude API Request**:
```json
{
  "model": "claude-3-5-sonnet-20241022",
  "max_tokens": 4096,
  "system": "You are an expert...",
  "messages": [
    {
      "role": "user",
      "content": "Template reference:\n```\n...\n```\n\nUser request: Create a React button..."
    }
  ],
  "stream": true
}
```

#### Step 6: Post-Processing

**File**: `internal/service/post_processor_service.go`

**Logic**:
1. **Extract Code**:
   - Removes markdown code blocks if present
   - Cleans up formatting
   - Extracts pure code

2. **Generate package.json**:
   - Adds React dependencies
   - Adds TypeScript if requested
   - Adds Tailwind if requested
   - Includes build scripts

3. **Generate Config Files**:
   - `tailwind.config.js` if using Tailwind
   - `tsconfig.json` if using TypeScript

4. **Organize Files**:
   - Determines component file path
   - Creates file structure
   - Returns `GenerationResponse` with all files

**Output Structure**:
```go
GenerationResponse{
    ComponentCode: "import React...",
    PackageJSON: "{...}",
    ConfigFiles: {
        "tailwind.config.js": "...",
        "tsconfig.json": "..."
    },
    Files: [
        {
            Path: "src/components/button.tsx",
            Content: "...",
            Type: "component"
        },
        // ... more files
    ]
}
```

#### Step 7: Streaming Response

**File**: `internal/handler/http_handler.go`

**Logic**:
- Receives chunks from `CodeGenerationService`
- Formats as SSE events:
  ```
  event: message
  data: {"componentCode":"import React..."}
  
  event: message
  data: {"componentCode":"import React...\nexport const Button..."}
  
  event: complete
  data: "true"
  ```
- Flushes after each event
- Handles errors by sending error events
- Closes connection on completion

---

## Frontend Integration

### Connection Setup

#### 1. Basic Configuration

Frontend needs to configure API base URL:
```javascript
const API_BASE_URL = 'http://localhost:8080/api/v1';
```

#### 2. SSE Streaming Connection

**Using EventSource (for GET requests)**:
```javascript
// Note: EventSource only supports GET, so you might need fetch with streaming
```

**Using Fetch API with Streaming** (Recommended):
```javascript
async function generateCode(prompt, options) {
  const response = await fetch(`${API_BASE_URL}/generate/stream`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      prompt: prompt,
      options: options
    })
  });

  const reader = response.body.getReader();
  const decoder = new TextDecoder();

  while (true) {
    const { done, value } = await reader.read();
    if (done) break;

    const chunk = decoder.decode(value);
    const lines = chunk.split('\n\n');

    for (const line of lines) {
      if (line.startsWith('data: ')) {
        const data = JSON.parse(line.slice(6));
        handleCodeChunk(data);
      } else if (line.startsWith('event: ')) {
        const eventType = line.slice(7);
        handleEvent(eventType);
      }
    }
  }
}
```

#### 3. React Hook Example

```typescript
import { useState, useCallback } from 'react';

export function useCodeGeneration() {
  const [code, setCode] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const generateCode = useCallback(async (prompt: string, options: any) => {
    setLoading(true);
    setError(null);
    setCode('');

    try {
      const response = await fetch('http://localhost:8080/api/v1/generate/stream', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ prompt, options }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const reader = response.body?.getReader();
      const decoder = new TextDecoder();

      if (!reader) {
        throw new Error('No response body');
      }

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        const chunk = decoder.decode(value);
        const lines = chunk.split('\n\n');

        for (const line of lines) {
          if (line.startsWith('event: message')) {
            // Next line will be data
            continue;
          }
          
          if (line.startsWith('data: ')) {
            const jsonStr = line.slice(6);
            if (jsonStr === 'true' || jsonStr === 'false') {
              // Completion event
              setLoading(false);
              break;
            }
            
            try {
              const data = JSON.parse(jsonStr);
              if (data.componentCode) {
                setCode(data.componentCode);
              }
            } catch (e) {
              console.error('Failed to parse SSE data:', e);
            }
          }
          
          if (line.startsWith('event: error')) {
            // Next line will be error data
            continue;
          }
          
          if (line.startsWith('event: complete')) {
            setLoading(false);
            break;
          }
        }
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
      setLoading(false);
    }
  }, []);

  return { code, loading, error, generateCode };
}
```

#### 4. Usage in Component

```typescript
function CodeGenerator() {
  const { code, loading, error, generateCode } = useCodeGeneration();

  const handleGenerate = () => {
    generateCode('Create a React button component with Tailwind CSS', {
      useTypeScript: true,
      styleLibrary: 'tailwind',
    });
  };

  return (
    <div>
      <button onClick={handleGenerate} disabled={loading}>
        {loading ? 'Generating...' : 'Generate Code'}
      </button>
      {error && <div className="error">{error}</div>}
      {code && <pre><code>{code}</code></pre>}
    </div>
  );
}
```

### CORS Configuration

Backend CORS is configured to allow frontend origins:

**Environment Variables**:
```env
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOW_HEADERS=Content-Type,Authorization
```

This allows:
- Frontend running on `localhost:5173` (Vite default)
- Frontend running on `localhost:3000` (React/CRA default)

### Non-Streaming Endpoint

For simpler use cases, you can use the non-streaming endpoint:

```typescript
async function generateCodeSimple(prompt: string, options: any) {
  const response = await fetch(`${API_BASE_URL}/generate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ prompt, options }),
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  const data = await response.json();
  return data; // Complete GenerationResponse
}
```

---

## Service Logic Breakdown

### 1. CodeGenerationService

**Purpose**: Main orchestrator that coordinates the entire code generation flow.

**Key Methods**:

#### `GenerateCodeStream(ctx, req)`
- Orchestrates streaming generation
- Returns two channels: `responseChan` and `errorChan`
- Flow:
  1. Parse intent
  2. Load template
  3. Build system prompt
  4. Stream AI generation (forward chunks as they arrive)
  5. Post-process final code
  6. Send complete response

**Concurrency**: Runs in goroutine to allow streaming

#### `GenerateCode(ctx, req)`
- Non-streaming version
- Same flow but waits for complete response
- Returns final `GenerationResponse`

#### `selectTemplate(intent)`
- Maps intent type to template path
- Returns template file path

#### `buildSystemPrompt(intent, options)`
- Constructs detailed system prompt for Claude
- Includes all requirements and constraints

---

### 2. IntentParserService

**Purpose**: Analyze user prompts to understand what code to generate.

**Key Methods**:

#### `ParseIntent(ctx, prompt)`
- Analyzes prompt text
- Returns structured `Intent` object

**Detection Logic**:

1. **Component Type Detection**:
   - Keywords: "react component", "vue component", "hook"
   - Patterns: `.jsx`, `.tsx`, `.vue` file extensions
   - Default: `react-component`

2. **Component Name Extraction**:
   - Pattern: "Create a [Name] component"
   - Pattern: "[Name] component"
   - Regex matching for capitalized words
   - Fallback: "Component"

3. **Feature Detection**:
   - State: "useState", "state"
   - Hooks: "hooks", "useEffect"
   - Styling: "tailwind", "styled-components"
   - TypeScript: "typescript", "ts"
   - API: "api", "fetch", "async"

**Confidence Scoring**:
- High (0.9): Explicit keywords found
- Medium (0.7): Patterns matched
- Low (0.5): Default fallback

---

### 3. ClaudeAIService

**Purpose**: Interface with Claude AI API for code generation.

**Key Methods**:

#### `GenerateCodeStream(ctx, prompt, systemPrompt, templateContent)`
- Makes streaming HTTP request to Claude API
- Processes Server-Sent Events from Claude
- Sends code chunks to channel
- Handles errors and completion

**Implementation Details**:
- Uses standard HTTP client (no SDK dependency)
- Headers:
  - `x-api-key`: Claude API key
  - `anthropic-version`: API version
  - `Content-Type`: application/json
- Parses SSE chunks from Claude
- Extracts text deltas from `content_block_delta` events

#### `GenerateCode(ctx, prompt, systemPrompt, templateContent)`
- Non-streaming version
- Waits for complete response
- Returns full generated code string

**Error Handling**:
- API errors (non-200 status)
- Network errors
- Parsing errors
- Context cancellation

---

### 4. PostProcessorService

**Purpose**: Process and enhance AI-generated code.

**Key Methods**:

#### `ProcessCode(ctx, code, intent, options)`
- Main post-processing function
- Returns complete `GenerationResponse`

**Processing Steps**:

1. **Code Extraction**:
   - Removes markdown code blocks
   - Extracts pure code
   - Cleans formatting

2. **Component Name Extraction**:
   - Parses code for component name
   - Uses regex to find `export const/function ComponentName`
   - Fallback to intent component name

3. **File Path Generation**:
   - Determines appropriate file path
   - Example: `src/components/button.tsx`
   - Based on intent type and component name

4. **package.json Generation**:
   - Base React dependencies
   - Adds TypeScript if requested
   - Adds Tailwind if requested
   - Adds styled-components if requested
   - Includes build scripts (Vite, etc.)

5. **Config Files Generation**:
   - `tailwind.config.js`: Tailwind configuration
   - `tsconfig.json`: TypeScript configuration

6. **Response Assembly**:
   - Creates file structure
   - Organizes all files
   - Returns complete response

**Helper Methods**:
- `extractCodeFromMarkdown()`: Remove code blocks
- `extractComponentNameFromCode()`: Parse component name
- `getComponentPath()`: Generate file path
- `generatePackageJSON()`: Create package.json
- `generateTailwindConfig()`: Create Tailwind config
- `generateTSConfig()`: Create TypeScript config

---

### 5. TemplateRepository

**Purpose**: Load template files from filesystem.

**Key Methods**:

#### `LoadTemplate(ctx, templatePath)`
- Loads template file from disk
- Validates path (security)
- Returns file contents

**Security**:
- Prevents directory traversal attacks
- Validates paths are within base directory
- Uses `filepath.IsLocal()` and path prefix checking

#### `ListTemplates(ctx)`
- Lists all available templates
- Recursively walks templates directory
- Returns list of template paths

#### `TemplateExists(ctx, templatePath)`
- Checks if template exists
- Performs same security checks
- Returns boolean

---

## Data Flow Diagrams

### Request Flow (Text)

```
Frontend → HTTP Handler
         ↓
    Validate Request
         ↓
    CodeGenerationService.GenerateCodeStream()
         ↓
    ┌─────────────────────────────────────┐
    │ Intent Parser                       │
    │ Prompt → Intent                     │
    └─────────────────────────────────────┘
         ↓
    ┌─────────────────────────────────────┐
    │ Template Repository                 │
    │ Intent → Template Path → File       │
    └─────────────────────────────────────┘
         ↓
    ┌─────────────────────────────────────┐
    │ Build System Prompt                 │
    │ Intent + Options → Prompt String    │
    └─────────────────────────────────────┘
         ↓
    ┌─────────────────────────────────────┐
    │ Claude AI Service                   │
    │ Prompt → HTTP Request               │
    │         → Claude API                │
    │         → Stream Chunks             │
    └─────────────────────────────────────┘
         ↓
    ┌─────────────────────────────────────┐
    │ Collect Stream Chunks               │
    │ Chunks → Full Code String           │
    └─────────────────────────────────────┘
         ↓
    ┌─────────────────────────────────────┐
    │ Post Processor                      │
    │ Code → Extract                      │
    │      → Generate package.json        │
    │      → Generate configs             │
    │      → Create file structure        │
    └─────────────────────────────────────┘
         ↓
    ┌─────────────────────────────────────┐
    │ Format SSE Events                   │
    │ Response → SSE Format               │
    └─────────────────────────────────────┘
         ↓
    Frontend (Streaming)
```

### Data Structures Flow

```
GenerationRequest (from frontend)
    ↓
Intent (parsed)
    ↓
Template Content (loaded)
    ↓
System Prompt (built)
    ↓
Claude API Request
    ↓
Code Chunks (streamed)
    ↓
Complete Code String
    ↓
GenerationResponse (with files, configs)
    ↓
SSE Events (to frontend)
```

---

## API Contracts

### Request Format

**Endpoint**: `POST /api/v1/generate/stream`

**Headers**:
```
Content-Type: application/json
```

**Body**:
```json
{
  "prompt": "Create a React button component with Tailwind CSS",
  "componentType": "react",  // Optional
  "options": {
    "useTypeScript": true,
    "includeTests": false,
    "styleLibrary": "tailwind"  // tailwind, css, styled-components
  }
}
```

### Response Format (SSE)

**Headers**:
```
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive
```

**Events**:

1. **Message Event** (streaming chunks):
```
event: message
data: {"componentCode":"import React..."}

event: message
data: {"componentCode":"import React...\nexport const Button..."}
```

2. **Error Event** (if error occurs):
```
event: error
data: {"error":{"code":"GENERATION_ERROR","message":"..."}}
```

3. **Complete Event** (when done):
```
event: complete
data: "true"
```

### Final Response Structure

The last message event contains complete response:

```json
{
  "componentCode": "import React from 'react';\n...",
  "packageJson": "{\"name\":\"generated-project\",...}",
  "configFiles": {
    "tailwind.config.js": "...",
    "tsconfig.json": "..."
  },
  "files": [
    {
      "path": "src/components/button.tsx",
      "content": "...",
      "type": "component"
    },
    {
      "path": "package.json",
      "content": "...",
      "type": "config"
    }
  ]
}
```

### Error Response Format

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": [
      {
        "field": "fieldName",
        "issue": "Specific issue",
        "value": "Invalid value"
      }
    ],
    "traceId": "request-trace-id",
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

### Health Check

**Endpoint**: `GET /api/v1/health`

**Response**:
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "version": "1.0.0",
  "services": {
    "api": "operational",
    "claude": "operational",
    "templates": "operational"
  }
}
```

---

## Summary

This architecture follows clean architecture principles with clear separation of concerns:

1. **HTTP Layer**: Handles requests, validation, and response formatting
2. **Service Layer**: Contains all business logic
3. **Repository Layer**: Abstracts data access
4. **Domain Layer**: Defines core entities and contracts

The flow is:
1. Request → Intent Parsing
2. Intent → Template Loading
3. Intent + Template → AI Generation
4. AI Output → Post-Processing
5. Final Response → Streaming to Frontend

Each service has a single responsibility and communicates through well-defined interfaces, making the system testable and maintainable.

---

## Simple Example: "make a todo app"

Let's walk through a complete example with the prompt **"make a todo app"** to see exactly how the request flows through the system.

### Step 1: Frontend Sends Request

**Frontend Code**:
```typescript
fetch('http://localhost:8080/api/v1/generate/stream', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    prompt: "make a todo app",
    options: {
      useTypeScript: true,
      styleLibrary: "tailwind"
    }
  })
})
```

**Request Body**:
```json
{
  "prompt": "make a todo app",
  "options": {
    "useTypeScript": true,
    "styleLibrary": "tailwind"
  }
}
```

---

### Step 2: HTTP Handler Receives Request

**Location**: `internal/handler/http_handler.go`

**What Happens**:
1. Validates JSON structure ✓
2. Sets SSE headers
3. Calls `codeGenService.GenerateCodeStream()`

**Data Passed**:
```go
req = GenerationRequest{
    Prompt: "make a todo app",
    Options: RequestOptions{
        UseTypeScript: true,
        StyleLibrary: "tailwind",
    },
}
```

---

### Step 3: Intent Parser Analyzes Prompt

**Location**: `internal/service/intent_parser_service.go`

**What Happens**:
- Analyzes: "make a todo app"
- Detects: No explicit "React component" but "app" keyword
- Extracts: Component name might be "Todo" or "TodoApp"
- Detects features: None explicitly mentioned, uses options

**Output**:
```go
Intent{
    Type: "react-component",  // Default fallback
    ComponentName: "TodoApp", // Extracted from "todo app"
    Features: ["tailwind", "typescript"], // From options
    Confidence: 0.6, // Medium confidence (not explicit)
}
```

---

### Step 4: Template Repository Loads Template

**Location**: `internal/repository/template_repository.go`

**What Happens**:
- Intent type: `react-component`
- Template path selected: `templates/react/component.tsx.template`
- Reads file from filesystem

**Template Loaded**:
```tsx
import React from 'react';

interface Props {
  // Define your component props here
}

export const Component: React.FC<Props> = ({}) => {
  return (
    <div className="component">
      {/* Your component implementation */}
    </div>
  );
};

export default Component;
```

---

### Step 5: System Prompt Built

**Location**: `internal/service/code_generation_service.go`

**What Happens**:
- Combines intent + options + template
- Creates detailed system prompt

**System Prompt Created**:
```
You are an expert frontend developer specializing in React.
Generate production-ready, clean, and maintainable code.

Requirements:
- Use functional components with hooks
- Use TypeScript with proper type definitions
- Use Tailwind CSS for styling
- Component name should be: TodoApp
- Follow React best practices and modern patterns
- Include proper error handling
- Make components reusable and well-documented
- Required features: tailwind, typescript
```

**User Message to Claude**:
```
Template reference:
```
tsx
import React from 'react';

interface Props {
  // Define your component props here
}

export const Component: React.FC<Props> = ({}) => {
  return (
    <div className="component">
      {/* Your component implementation */}
    </div>
  );
};

export default Component;
```

User request: make a todo app
```

---

### Step 6: Claude AI Generates Code

**Location**: `internal/service/ai_service.go`

**What Happens**:
1. Sends HTTP POST to Claude API:
```json
{
  "model": "claude-3-5-sonnet-20241022",
  "max_tokens": 4096,
  "system": "You are an expert frontend developer...",
  "messages": [
    {
      "role": "user",
      "content": "Template reference: ... User request: make a todo app"
    }
  ],
  "stream": true
}
```

2. Receives streaming chunks:
```
Chunk 1: "import React, { useState } from 'react';\n\n"
Chunk 2: "interface Todo {\n  id: number;\n  text: string;\n  completed: boolean;\n}\n\n"
Chunk 3: "export const TodoApp: React.FC = () => {\n  const [todos, setTodos] = useState<Todo[]>([]);\n"
Chunk 4: "  const [input, setInput] = useState('');\n\n  const addTodo = () => {..."
...
```

3. Each chunk is sent immediately to frontend via SSE

**Accumulated Code** (after all chunks):
```tsx
import React, { useState } from 'react';

interface Todo {
  id: number;
  text: string;
  completed: boolean;
}

export const TodoApp: React.FC = () => {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [input, setInput] = useState('');

  const addTodo = () => {
    if (input.trim()) {
      setTodos([...todos, {
        id: Date.now(),
        text: input,
        completed: false
      }]);
      setInput('');
    }
  };

  const toggleTodo = (id: number) => {
    setTodos(todos.map(todo =>
      todo.id === id ? { ...todo, completed: !todo.completed } : todo
    ));
  };

  return (
    <div className="max-w-md mx-auto mt-8 p-4">
      <h1 className="text-2xl font-bold mb-4">Todo App</h1>
      <div className="flex gap-2 mb-4">
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyPress={(e) => e.key === 'Enter' && addTodo()}
          className="flex-1 px-4 py-2 border rounded"
          placeholder="Add a todo..."
        />
        <button
          onClick={addTodo}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          Add
        </button>
      </div>
      <ul className="space-y-2">
        {todos.map(todo => (
          <li
            key={todo.id}
            className={`flex items-center gap-2 p-2 border rounded ${
              todo.completed ? 'line-through text-gray-500' : ''
            }`}
          >
            <input
              type="checkbox"
              checked={todo.completed}
              onChange={() => toggleTodo(todo.id)}
            />
            <span>{todo.text}</span>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default TodoApp;
```

---

### Step 7: Post-Processor Enhances Output

**Location**: `internal/service/post_processor_service.go`

**What Happens**:

1. **Extracts Code**: Removes any markdown, keeps pure TSX
2. **Generates package.json**:
```json
{
  "name": "generated-project",
  "version": "0.1.0",
  "private": true,
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "typescript": "^5.0.0",
    "@types/react": "^18.2.0",
    "@types/react-dom": "^18.2.0",
    "tailwindcss": "^3.3.0",
    "autoprefixer": "^10.4.14",
    "postcss": "^8.4.24"
  },
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "preview": "vite preview"
  }
}
```

3. **Generates tailwind.config.js**:
```js
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
```

4. **Generates tsconfig.json**:
```json
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "jsx": "react-jsx",
    "strict": true
  },
  "include": ["src"]
}
```

5. **Organizes Files**:
```go
Files: [
    {
        Path: "src/components/todoapp.tsx",
        Content: "import React... (full code)",
        Type: "component"
    },
    {
        Path: "package.json",
        Content: "{...}",
        Type: "config"
    },
    {
        Path: "tailwind.config.js",
        Content: "...",
        Type: "config"
    },
    {
        Path: "tsconfig.json",
        Content: "...",
        Type: "config"
    }
]
```

---

### Step 8: Response Streamed to Frontend

**Location**: `internal/handler/http_handler.go`

**SSE Events Sent to Frontend**:

**Event 1** (first chunk):
```
event: message
data: {"componentCode":"import React, { useState } from 'react';\n\n"}
```

**Event 2** (more code):
```
event: message
data: {"componentCode":"import React, { useState } from 'react';\n\ninterface Todo {\n  id: number;\n  text: string;\n  completed: boolean;\n}\n\n"}
```

**... (more events as code streams)**

**Final Event** (complete response):
```
event: message
data: {"componentCode":"...full code...","packageJson":"{...}","configFiles":{"tailwind.config.js":"...","tsconfig.json":"..."},"files":[{"path":"src/components/todoapp.tsx","content":"...","type":"component"},...]}
```

**Completion Event**:
```
event: complete
data: "true"
```

---

### Step 9: Frontend Receives and Displays

**Frontend React Component**:
```typescript
function CodeViewer() {
  const [code, setCode] = useState('');
  
  useEffect(() => {
    // Stream handling code from earlier
    // Each SSE event updates `code` state
  }, []);

  return (
    <div>
      <pre><code>{code}</code></pre>
      {/* Code appears character by character as it streams */}
    </div>
  );
}
```

**User Sees**:
1. Code starts appearing: `import React...`
2. More code streams in: `interface Todo...`
3. Component code appears: `export const TodoApp...`
4. Finally, complete file structure with package.json, configs, etc.

---

### Summary Timeline

```
0ms    → Frontend sends: "make a todo app"
5ms    → Handler validates request
10ms   → Intent Parser: "react-component, TodoApp"
15ms   → Template loaded: component.tsx.template
20ms   → System prompt built
25ms   → Request sent to Claude API
100ms  → First chunk arrives: "import React..."
150ms  → Chunk streamed to frontend (SSE event)
200ms  → Second chunk: "interface Todo..."
250ms  → Chunk streamed to frontend
...    → More chunks streaming
2000ms → All chunks received
2050ms → Post-processing: Generate package.json, configs
2100ms → Final response sent
2150ms → Completion event sent
2200ms → Frontend displays complete code
```

---

### Complete Data Transformation

```
Input:  "make a todo app" + { useTypeScript: true, styleLibrary: "tailwind" }
  ↓
Intent: { Type: "react-component", ComponentName: "TodoApp", Features: ["tailwind", "typescript"] }
  ↓
Template: Basic React component template
  ↓
AI Output: Full TodoApp component with useState, styling, etc.
  ↓
Post-Processed: Code + package.json + tailwind.config.js + tsconfig.json
  ↓
Output: Complete project structure ready to use
```

This example shows how a simple 3-word prompt transforms into a complete, production-ready React component with all necessary configuration files!
