# Model Management Python Client

A command-line client for interacting with the Model Management Server.

## Installation

1. Install dependencies:
```bash
pip install -r requirements.txt
```

2. Make the client executable (optional):
```bash
chmod +x model_client.py
```

## Usage

### List all models
```bash
python model_client.py list
```

### Get a specific model
```bash
python model_client.py get -i bert-base-uncased
```

### Download a model from Hugging Face
```bash
python model_client.py download -n google/flan-t5-small
```

Popular models to try:
- `gpt2` (548 MB)
- `distilbert-base-uncased` (268 MB)
- `google/flan-t5-small` (308 MB)
- `bert-base-uncased` (440 MB)

### Delete a model
```bash
python model_client.py delete -i bert-base-uncased
```

### Get system statistics
```bash
python model_client.py stats
```

### Check server health
```bash
python model_client.py health
```

### Use custom server URL
```bash
python model_client.py -s http://custom-server:5000 list
```

## Features

- Automatic username detection (uses `getpass.getuser()`)
- Formatted output with aligned columns
- Human-readable file sizes
- Usage statistics display
- Error handling with helpful messages
- Support for custom server URLs
- Long timeout for model downloads (5 minutes)

## Command Reference

### Global Options
- `-s, --server` - Server URL (default: http://localhost:5000)

### Commands

#### `list`
List all models with their status, size, download count, and first downloader.

#### `get -i MODEL_ID`
Get detailed information about a specific model, including usage statistics.

**Options:**
- `-i, --id` - Model ID (required)

#### `download -n MODEL_NAME`
Download a model from Hugging Face. The username is automatically detected.

**Options:**
- `-n, --name` - Model name from Hugging Face (required)

#### `delete -i MODEL_ID`
Delete a model from the server.

**Options:**
- `-i, --id` - Model ID to delete (required)

#### `stats`
Display system-wide statistics including total models, storage size, request counts, and recent activity.

#### `health`
Check if the server is healthy and responding.

## Examples

```bash
# List all models
python model_client.py list

# Download a model
python model_client.py download -n bert-base-uncased

# Get model details
python model_client.py get -i bert-base-uncased

# View system statistics
python model_client.py stats

# Delete a model
python model_client.py delete -i bert-base-uncased

# Check server health
python model_client.py health

# Use with a different server
python model_client.py -s http://192.168.1.100:5000 list
```

## Output Format

### List Command
```
NAME                                     STATUS      SIZE         DOWNLOADS    FIRST BY            
----------------------------------------------------------------------------------------------
bert-base-uncased                        ready       440.47 MB    3            john                
google/flan-t5-small                     ready       308.24 MB    1            jane                
```

### Get Command
```
Model Details:
  ID:           bert-base-uncased
  Name:         bert-base-uncased
  Status:       ready
  Size:         440.47 MB
  Path:         ./models/bert-base-uncased
  Downloaded:   2026-01-22T10:30:00.000000

Usage Statistics:
  Downloads:    3
  Accesses:     12
  Total Reqs:   15
  First By:     john
  First From:   192.168.1.100
  First At:     2026-01-22T10:30:00.000000
```

## Error Handling

The client provides clear error messages for common issues:
- Connection errors (server not running)
- Invalid model names
- Model not found
- Network timeouts
- Invalid JSON responses

Exit codes:
- `0` - Success
- `1` - Error occurred

## Requirements

- Python 3.6 or higher
- `requests` library
