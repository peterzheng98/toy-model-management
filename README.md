# Toy Model Management System

A comprehensive model management system with a Python Flask server and Go client for downloading, storing, and managing machine learning models from Hugging Face.

## Features

- **Flask Server** with RESTful API
- **Modern Web UI** for users and administrators
- **Hugging Face Integration** for model downloads
- **Model Existence Checking** before downloads
- **CRUD Operations** for model management
- **Go CLI Client** for programmatic access
- **Configurable Storage** via mount point
- **Usage Tracking** - Records username and IP for all requests
- **Statistics Dashboard** - Total folder size, download counts, unique users
- **User Attribution** - Track who requested each model download first

## Architecture

### Server (Python/Flask)
- **Backend API**: RESTful endpoints for model operations
- **Frontend**: Beautiful, responsive UI with user and admin views
- **Storage**: Configurable mount point for model storage
- **Database**: JSON-based simple database for model metadata

### Clients

**Go Client** (`client/go/`)
- **CLI Tool**: Command-line interface for all operations
- **HTTP Client**: Communicates with Flask server API
- **CRUD Support**: List, get, download, update, and delete models

**Python Client** (`client/python/`)
- **CLI Tool**: Python-based command-line interface
- **HTTP Client**: Uses requests library
- **Same Features**: All CRUD operations matching Go client functionality

## Installation

### Prerequisites
- Python 3.8+
- Go 1.21+
- Hugging Face account (optional, for private models)

### Server Setup

1. Install Python dependencies:
```bash
pip install -r requirements.txt
```

2. Set the models mount point (optional):
```bash
export MODELS_MOUNT_POINT=/path/to/your/models
```

3. Run the server:
```bash
cd server
python app.py
```

The server will start on `http://localhost:5000`

### Client Setup

#### Go Client

1. Build the Go client:
```bash
cd client/go
go build -o model-client main.go
```

2. The client binary `model-client` is now ready to use!

#### Python Client

1. Install Python client dependencies:
```bash
cd client/python
pip install -r requirements.txt
```

2. The client script `model_client.py` is now ready to use!

## Usage

### Web Interface

#### User View
Visit `http://localhost:5000/` to:
- Browse all available models
- View model details (size, download date, status)
- See model metadata in a clean interface

#### Admin Panel
Visit `http://localhost:5000/admin` to:
- Download new models from Hugging Face
- View all models in a table format
- Delete existing models
- Monitor download status

### CLI Clients

Both Go and Python clients support the same commands:

#### List all models

**Go:**
```bash
./client/go/model-client list
```

**Python:**
```bash
python client/python/model_client.py list
```

#### Get a specific model

**Go:**
```bash
./client/go/model-client get -id bert-base-uncased
```

**Python:**
```bash
python client/python/model_client.py get -i bert-base-uncased
```

#### Download a model from Hugging Face

**Go:**
```bash
./client/go/model-client download -name google/flan-t5-small
```

**Python:**
```bash
python client/python/model_client.py download -n google/flan-t5-small
```

Examples of model names:
- `bert-base-uncased`
- `google/flan-t5-small`
- `facebook/opt-125m`
- `gpt2`

#### Delete a model

**Go:**
```bash
./client/go/model-client delete -id google_flan-t5-small
```

**Python:**
```bash
python client/python/model_client.py delete -i google_flan-t5-small
```

#### Get system statistics

**Go:**
```bash
./client/go/model-client stats
```

**Python:**
```bash
python client/python/model_client.py stats
```

#### Check server health

**Go:**
```bash
./client/go/model-client health
```

**Python:**
```bash
python client/python/model_client.py health
```

#### Use custom server URL

**Go:**
```bash
./client/go/model-client list -server http://custom-server:5000
```

**Python:**
```bash
python client/python/model_client.py -s http://custom-server:5000 list
```

## API Endpoints

### GET `/api/models`
List all models with usage statistics

**Response:**
```json
{
  "success": true,
  "models": [
    {
      "id": "bert-base-uncased",
      "name": "bert-base-uncased",
      "size_bytes": 440473133,
      "status": "ready",
      "stats": {
        "download_count": 5,
        "access_count": 12,
        "total_requests": 17,
        "first_downloaded_by": "john",
        "first_downloaded_from": "192.168.1.100"
      }
    }
  ]
}
```

### GET `/api/models/<model_id>`
Get a specific model

**Response:**
```json
{
  "success": true,
  "model": {...}
}
```

### POST `/api/models/download`
Download a model from Hugging Face

**Request Body:**
```json
{
  "model_name": "bert-base-uncased",
  "username": "john"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Model downloaded successfully",
  "model": {
    "id": "bert-base-uncased",
    "downloaded_by": "john",
    "stats": {
      "download_count": 1,
      "first_downloaded_by": "john",
      "first_downloaded_from": "192.168.1.100"
    }
  }
}
```

### PUT `/api/models/<model_id>`
Update model metadata

**Request Body:**
```json
{
  "status": "ready"
}
```

### DELETE `/api/models/<model_id>`
Delete a model

**Response:**
```json
{
  "success": true,
  "message": "Model deleted successfully"
}
```

### GET `/api/stats`
Get overall system statistics

**Response:**
```json
{
  "success": true,
  "stats": {
    "total_models": 5,
    "total_size_bytes": 2147483648,
    "total_requests": 150,
    "unique_users": 12,
    "recent_activity": [
      {
        "timestamp": "2026-01-22T10:30:00Z",
        "action": "download",
        "model_id": "bert-base-uncased",
        "username": "john",
        "ip_address": "192.168.1.100"
      }
    ]
  }
}
```

### GET `/api/health`
Health check endpoint

**Response:**
```json
{
  "success": true,
  "status": "healthy",
  "mount_point": "/path/to/models"
}
```

## Configuration

### Environment Variables

- `MODELS_MOUNT_POINT`: Directory where models will be stored (default: `./models`)

Example:
```bash
export MODELS_MOUNT_POINT=/mnt/shared/models
python server/app.py
```

### Hugging Face Authentication

For private models, set your Hugging Face token:
```bash
export HUGGINGFACE_TOKEN=your_token_here
```

Or use `huggingface-cli login` before running the server.

## Project Structure

```
toy-model-management/
├── server/
│   ├── app.py                 # Flask application
│   └── templates/
│       ├── index.html         # User view UI
│       └── admin.html         # Admin panel UI
├── client/
│   ├── go/
│   │   ├── main.go            # Go client implementation
│   │   └── go.mod             # Go module file
│   └── python/
│       ├── model_client.py    # Python client implementation
│       ├── requirements.txt   # Python client dependencies
│       └── README.md          # Python client documentation
├── requirements.txt           # Server dependencies
├── .gitignore                 # Git ignore rules
├── README.md                  # This file
├── QUICKSTART.md              # Quick start guide
├── CHANGELOG.md               # Version history
├── start_server.sh            # Server startup script
└── build_client.sh            # Go client build script
```

## Model Storage

Models are stored in the configured mount point with the following structure:

```
models/
├── models_db.json          # Model metadata database
├── bert-base-uncased/      # Model files
├── google_flan-t5-small/   # Model files (/ replaced with _)
└── ...
```

## Features in Detail

### Usage Tracking
Every request to the server is logged with:
- **Username** - Automatically detected from the system (whoami) or provided by web UI
- **IP Address** - Client IP address
- **Timestamp** - When the request occurred
- **Action** - What operation was performed (list, get, download, delete)

### Statistics Dashboard
The web UI displays real-time statistics:
- **Total Models** - Number of models stored
- **Total Folder Size** - Combined size of all models
- **Total Requests** - All API requests made to the system
- **Unique Users** - Number of different users who have interacted with the system
- **Download Count** - Per-model download statistics
- **First Downloader** - Who requested each model first, including their IP address

### Model Existence Checking
Before downloading, the server checks if the model already exists in the mount point to avoid redundant downloads.

### Automatic Model ID Generation
Model IDs are automatically generated from model names by replacing `/` with `_` (e.g., `google/flan-t5-small` → `google_flan-t5-small`)

### Automatic Username Detection
The Go client automatically detects the current system username using:
1. Environment variables (`USER`, `USERNAME`)
2. Go's `user.Current()` function
3. `whoami` command as fallback

### Download Progress
The admin panel shows real-time status updates during model downloads.

### Responsive UI
Both user and admin interfaces are fully responsive and work on mobile devices.

## Troubleshooting

### Server won't start
- Check if port 5000 is available
- Ensure Python dependencies are installed
- Verify mount point directory permissions

### Model download fails
- Check internet connection
- Verify Hugging Face model name is correct
- For private models, ensure authentication is set up
- Check available disk space

### Client connection error
- Ensure server is running on the specified URL
- Check firewall settings
- Verify network connectivity

## Development

### Adding New Features
The codebase is structured for easy extension:

- **Server endpoints**: Add new routes in `server/app.py`
- **Client commands**: Add new subcommands in `client/main.go`
- **UI features**: Modify templates in `server/templates/`

### Testing
Test the API endpoints:
```bash
# Health check
curl http://localhost:5000/api/health

# List models
curl http://localhost:5000/api/models

# Download model
curl -X POST http://localhost:5000/api/models/download \
  -H "Content-Type: application/json" \
  -d '{"model_name": "bert-base-uncased"}'
```

## License

This is a toy project for demonstration purposes.

## Contributing

This is a toy project, but feel free to fork and modify as needed!
