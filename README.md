
# A small infocenter service

A Go based REST API service that implements a real-time message broadcasting system using SSE. The service has topics, to which clients can subscribe to receive messages and post new messages to broadcast to all subscribers.

## Features

- **Real-time messaging**: Clients receive messages via SSE
- **Topic-based communication**: Independent topics with different message streams
- **Automatic client timeout**: Clients are automatically disconnected after a configurable timeout period

## Prerequisites

- Go 1.25 or higher
- `github.com/joho/godotenv` package

## Installation
```bash
# Clone the repository
git clone https://github.com/xaxax4/Infocenter_service.git
cd Infocenter_service

# Install dependencies
go mod download

# Run the application
go run ./cmd/infocenter/main.go 
```

## Configuration

Create a `.env` file in the root directory with the following variables:
```env
PORT=selected_server_port
SERVER_SHUTDOWN_TIMEOUT=selected_server_shutdown_timeout
CLIENT_TIMEOUT=selected_client_timeout
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port number | 3000    |
| `SERVER_SHUTDOWN_TIMEOUT` | Graceful shutdown timeout in seconds | 5       |
| `CLIENT_TIMEOUT` | Client connection timeout in seconds | 30      |


## API Endpoints

### Subscribe to Topic (GET)
```http
GET /infocenter/{topic}
```

Establishes a Server-Sent Events connection to receive real-time messages for a specific topic.

**Parameters:**
- `topic` (path parameter) - The specified topic

**Response:**
- Content-Type: `text/event-stream`
- Connection remains open, streaming messages as they are sent
- Client is automatically disconnected after `CLIENT_TIMEOUT` seconds of inactivity


**Example:**
```bash
curl -N http://localhost:3000/infocenter/test
```

**Message Format:**
```
id: <message_id>
data: <message_content>
```

**Timeout Message:**
When a client times out, it receives a final message with ID `-1`, where the content is the connection duration.


### Send Message to Topic (POST)
```http
POST /infocenter/{topic}
```

Broadcasts a message to all clients subscribed to the specified topic.

**Parameters:**
- `topic` (path parameter) - The specified topic

**Request Body:**
- Content-Type: `text/plain`
- Raw message content

**Response:**
- Status: `204 No Content` on success

**Example:**
```bash
curl -X POST http://localhost:3000/infocenter/test \
  -H "Content-Type: text/plain" \
  -d "Hello World!"
```

## Project Structure
```
Infocenter_service/
├── cmd/
|   ├──infocenter/
|		└── main.go 		   # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/          # HTTP request handlers
│   │   ├── middlewares/       # CORS middleware
│   │   └── router/            # Routes
│   ├── models/                # Data models (Topic, Client, Message)
│   └── services/              # MiddleMan service
├── .env                       # Environment configuration
├── go.mod                     # Go module definition
├── go.sum                     # Dependency checksums
└── README.md
```

