# Quick Start Guide

Get started with the Model Management System in 5 minutes!

## Prerequisites

- **Python 3.8+** installed
- **Go 1.21+** installed (for the client)
- Internet connection (for downloading models from Hugging Face)

## Step 1: Install Python Dependencies

```bash
pip install -r requirements.txt
```

## Step 2: Start the Server

### Option A: Using the start script (recommended)
```bash
./start_server.sh
```

### Option B: Manual start
```bash
cd server
python3 app.py
```

The server will start at `http://localhost:5000`

## Step 3: Access the Web Interface

Open your browser and visit:

- **User View**: http://localhost:5000/
  - Browse available models
  - View model details
  - See system statistics (total models, folder size, request counts)
  - View who downloaded each model first
  
- **Admin Panel**: http://localhost:5000/admin
  - Download new models (requires username)
  - Manage existing models
  - Delete models
  - View usage statistics and download counts

## Step 4: Build the Clients (Optional)

### Go Client

**Option A: Using the build script (recommended)**
```bash
./build_client.sh
```

**Option B: Manual build**
```bash
cd client/go
go build -o model-client main.go
```

### Python Client

**Install dependencies:**
```bash
cd client/python
pip install -r requirements.txt
```

## Step 5: Use the Clients

You can use either the Go client or the Python client. Both have the same functionality.

### Go Client

**List all models:**
```bash
./client/go/model-client list
```

**Download a model:**
```bash
./client/go/model-client download -name gpt2
```

**Get model details:**
```bash
./client/go/model-client get -id gpt2
```

**Delete a model:**
```bash
./client/go/model-client delete -id gpt2
```

**Check server health:**
```bash
./client/go/model-client health
```

### Python Client

**List all models:**
```bash
python client/python/model_client.py list
```

**Download a model:**
```bash
python client/python/model_client.py download -n gpt2
```

**Get model details:**
```bash
python client/python/model_client.py get -i gpt2
```

**Get system statistics:**
```bash
python client/python/model_client.py stats
```

**Delete a model:**
```bash
python client/python/model_client.py delete -i gpt2
```

**Check server health:**
```bash
python client/python/model_client.py health
```

### Popular Models to Try
- `gpt2` (548 MB)
- `distilbert-base-uncased` (268 MB)
- `google/flan-t5-small` (308 MB)

## Configuration

### Set Custom Models Directory

By default, models are stored in `./models`. To use a different location:

```bash
export MODELS_MOUNT_POINT=/path/to/your/models
./start_server.sh
```

### For Private Hugging Face Models

If you need to download private models:

```bash
# Option 1: Set token environment variable
export HUGGINGFACE_TOKEN=your_token_here

# Option 2: Login with Hugging Face CLI
pip install huggingface-hub[cli]
huggingface-cli login
```

## Testing the System

### Using the Web UI

1. Visit http://localhost:5000/admin
2. Enter your username (will be saved for future use)
3. Enter a model name (e.g., `gpt2`)
4. Click "Download Model"
5. Wait for the download to complete
6. View the model in the user interface at http://localhost:5000/
7. See statistics including who downloaded it first and usage count

### Using the CLI Clients

Both clients automatically detect your system username.

**Go Client:**
```bash
# Download a small model (username auto-detected)
./client/go/model-client download -name gpt2

# List all models (shows download counts and first downloader)
./client/go/model-client list

# Get model details (includes usage statistics)
./client/go/model-client get -id gpt2
```

**Python Client:**
```bash
# Download a small model (username auto-detected)
python client/python/model_client.py download -n gpt2

# List all models (shows download counts and first downloader)
python client/python/model_client.py list

# Get model details (includes usage statistics)
python client/python/model_client.py get -i gpt2
```

### Using Python Example

```bash
cd examples
python3 python_example.py
```

### Using Go Example

```bash
cd examples
go run go_example.go
```

## Common Issues

### Port 5000 already in use
```bash
# Find process using port 5000
lsof -ti:5000

# Kill the process (Mac/Linux)
kill -9 $(lsof -ti:5000)

# Or use a different port
export FLASK_RUN_PORT=8000
./start_server.sh
```

### Python dependencies not installed
```bash
pip install -r requirements.txt
```

### Go build fails
```bash
# Make sure you're in the correct directory
cd client/go
go mod tidy
go build -o model-client main.go
```

### Python client import errors
```bash
# Install dependencies
cd client/python
pip install -r requirements.txt
```

## Next Steps

- Read the full [README.md](README.md) for detailed documentation
- Check out the [examples/](examples/) directory for code samples
- Explore the API endpoints documentation in README.md

## Need Help?

- Check the server logs for errors
- Ensure you have enough disk space for model downloads
- Verify internet connection for Hugging Face access
- Make sure Python and Go are properly installed

## Example Workflow

Here's a complete workflow to test the system:

```bash
# 1. Start the server
./start_server.sh

# 2. In a new terminal, build the client
./build_client.sh

# 3. Check server health
./client/model-client health

# 4. Download a small model
./client/model-client download -name gpt2

# 5. List models
./client/model-client list

# 6. View in browser
# Open http://localhost:5000/ to see the model in the UI

# 7. Delete the model
./client/model-client delete -id gpt2
```

Enjoy using the Model Management System!
