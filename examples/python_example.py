#!/usr/bin/env python3
"""
Example script demonstrating how to interact with the Model Management API using Python
"""

import requests
import json
import time
import getpass
import socket

# Server configuration
SERVER_URL = "http://localhost:5000"
API_BASE = f"{SERVER_URL}/api"

# Get current username
CURRENT_USER = getpass.getuser()


def health_check():
    """Check if server is healthy"""
    print("Checking server health...")
    response = requests.get(f"{API_BASE}/health")
    data = response.json()
    
    if data['success']:
        print(f"Server is healthy!")
        print(f"  Mount point: {data.get('mount_point', 'N/A')}")
        return True
    else:
        print("Server is not healthy")
        return False


def list_models():
    """List all models"""
    print("\nListing all models...")
    response = requests.get(f"{API_BASE}/models")
    data = response.json()
    
    if data['success']:
        models = data['models']
        if models:
            print(f"Found {len(models)} model(s):")
            for model in models:
                size_mb = model['size_bytes'] / (1024 * 1024)
                print(f"  - {model['name']}")
                print(f"    ID: {model['id']}")
                print(f"    Size: {size_mb:.2f} MB")
                print(f"    Status: {model['status']}")
                print()
        else:
            print("  No models found")
        return models
    else:
        print(f"Error: {data.get('error', 'Unknown error')}")
        return []


def download_model(model_name, username=None):
    """Download a model from Hugging Face"""
    if username is None:
        username = CURRENT_USER
    
    print(f"\nDownloading model: {model_name}")
    print(f"   Requester: {username}")
    print("   This may take a while...")
    
    response = requests.post(
        f"{API_BASE}/models/download",
        json={
            "model_name": model_name,
            "username": username
        }
    )
    data = response.json()
    
    if data['success']:
        if data.get('already_exists'):
            print(f"Model already exists!")
        else:
            print(f"Model downloaded successfully!")
        
        model = data['model']
        size_mb = model['size_bytes'] / (1024 * 1024)
        print(f"  ID: {model['id']}")
        print(f"  Size: {size_mb:.2f} MB")
        print(f"  Path: {model['path']}")
        
        # Print statistics if available
        if 'stats' in model:
            stats = model['stats']
            print(f"\n  Usage Statistics:")
            print(f"    Downloads: {stats.get('download_count', 0)}")
            print(f"    Total Requests: {stats.get('total_requests', 0)}")
            if stats.get('first_downloaded_by'):
                print(f"    First Downloaded By: {stats['first_downloaded_by']}")
                print(f"    First Downloaded From: {stats.get('first_downloaded_from', 'N/A')}")
        
        return model
    else:
        print(f"Download failed: {data.get('error', 'Unknown error')}")
        return None


def get_model(model_id):
    """Get details of a specific model"""
    print(f"\nGetting model details: {model_id}")
    
    response = requests.get(f"{API_BASE}/models/{model_id}")
    data = response.json()
    
    if data['success']:
        model = data['model']
        size_mb = model['size_bytes'] / (1024 * 1024)
        print(f"Model found!")
        print(f"  Name: {model['name']}")
        print(f"  ID: {model['id']}")
        print(f"  Size: {size_mb:.2f} MB")
        print(f"  Status: {model['status']}")
        print(f"  Path: {model['path']}")
        print(f"  Downloaded: {model['downloaded_at']}")
        
        # Print statistics if available
        if 'stats' in model:
            stats = model['stats']
            print(f"\n  Usage Statistics:")
            print(f"    Downloads: {stats.get('download_count', 0)}")
            print(f"    Accesses: {stats.get('access_count', 0)}")
            print(f"    Total Requests: {stats.get('total_requests', 0)}")
            if stats.get('first_downloaded_by'):
                print(f"    First Downloaded By: {stats['first_downloaded_by']}")
                print(f"    First Downloaded At: {stats.get('first_downloaded_at', 'N/A')}")
                print(f"    First Downloaded From: {stats.get('first_downloaded_from', 'N/A')}")
        
        return model
    else:
        print(f"Error: {data.get('error', 'Unknown error')}")
        return None


def get_stats():
    """Get overall system statistics"""
    print(f"\nGetting system statistics...")
    
    response = requests.get(f"{API_BASE}/stats")
    data = response.json()
    
    if data['success']:
        stats = data['stats']
        total_size_mb = stats['total_size_bytes'] / (1024 * 1024)
        print(f"Statistics:")
        print(f"  Total Models: {stats['total_models']}")
        print(f"  Total Size: {total_size_mb:.2f} MB")
        print(f"  Total Requests: {stats['total_requests']}")
        print(f"  Unique Users: {stats['unique_users']}")
        
        if stats.get('recent_activity'):
            print(f"\n  Recent Activity (last {len(stats['recent_activity'])}):")
            for activity in stats['recent_activity']:
                print(f"    - {activity['action']} by {activity['username']} from {activity['ip_address']}")
        
        return stats
    else:
        print(f"Error: {data.get('error', 'Unknown error')}")
        return None


def delete_model(model_id):
    """Delete a model"""
    print(f"\nDeleting model: {model_id}")
    
    response = requests.delete(f"{API_BASE}/models/{model_id}")
    data = response.json()
    
    if data['success']:
        print(f"Model deleted successfully!")
        return True
    else:
        print(f"Delete failed: {data.get('error', 'Unknown error')}")
        return False


def main():
    """Run example workflow"""
    print("=" * 60)
    print("Model Management System - Python Example")
    print(f"Running as user: {CURRENT_USER}")
    print("=" * 60)
    
    # Check server health
    if not health_check():
        print("\nError: Server is not running!")
        print("Start the server with: ./start_server.sh")
        return
    
    # Get system statistics
    get_stats()
    
    # List existing models
    list_models()
    
    # Example: Download a small model (uncomment to test)
    # Note: This will actually download the model, which may take time
    # download_model("gpt2")
    
    # Example: Get a specific model (uncomment after downloading)
    # get_model("gpt2")
    
    # Example: Delete a model (uncomment to test)
    # delete_model("gpt2")
    
    print("\n" + "=" * 60)
    print("Example complete!")
    print("=" * 60)


if __name__ == "__main__":
    main()
