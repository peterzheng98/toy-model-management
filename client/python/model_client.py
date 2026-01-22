#!/usr/bin/env python3
"""
Model Management System - Python CLI Client

A command-line client for interacting with the Model Management Server.
"""

import argparse
import getpass
import json
import sys
import requests
from typing import Optional, Dict, List, Any


DEFAULT_SERVER_URL = "http://localhost:5000"


class ModelClient:
    """Client for the Model Management API"""
    
    def __init__(self, server_url: str = DEFAULT_SERVER_URL):
        self.server_url = server_url
        self.api_base = f"{server_url}/api"
        self.session = requests.Session()
        self.session.timeout = 300  # 5 minutes for downloads
    
    def _handle_response(self, response: requests.Response) -> Dict[str, Any]:
        """Handle API response and return JSON data"""
        try:
            data = response.json()
            return data
        except json.JSONDecodeError:
            return {
                'success': False,
                'error': f'Invalid JSON response: {response.text}'
            }
    
    def list_models(self) -> List[Dict[str, Any]]:
        """List all models"""
        try:
            response = self.session.get(f"{self.api_base}/models")
            data = self._handle_response(response)
            
            if data.get('success'):
                return data.get('models', [])
            else:
                raise Exception(data.get('error', 'Unknown error'))
        except requests.RequestException as e:
            raise Exception(f"Failed to list models: {e}")
    
    def get_model(self, model_id: str) -> Dict[str, Any]:
        """Get a specific model by ID"""
        try:
            response = self.session.get(f"{self.api_base}/models/{model_id}")
            data = self._handle_response(response)
            
            if data.get('success'):
                return data.get('model', {})
            else:
                raise Exception(data.get('error', 'Unknown error'))
        except requests.RequestException as e:
            raise Exception(f"Failed to get model: {e}")
    
    def download_model(self, model_name: str, username: Optional[str] = None) -> Dict[str, Any]:
        """Download a model from Hugging Face"""
        if username is None:
            username = getpass.getuser()
        
        try:
            response = self.session.post(
                f"{self.api_base}/models/download",
                json={
                    'model_name': model_name,
                    'username': username
                }
            )
            data = self._handle_response(response)
            
            if data.get('success'):
                return data
            else:
                raise Exception(data.get('error', 'Unknown error'))
        except requests.RequestException as e:
            raise Exception(f"Failed to download model: {e}")
    
    def delete_model(self, model_id: str) -> None:
        """Delete a model"""
        try:
            response = self.session.delete(f"{self.api_base}/models/{model_id}")
            data = self._handle_response(response)
            
            if not data.get('success'):
                raise Exception(data.get('error', 'Unknown error'))
        except requests.RequestException as e:
            raise Exception(f"Failed to delete model: {e}")
    
    def get_stats(self) -> Dict[str, Any]:
        """Get system statistics"""
        try:
            response = self.session.get(f"{self.api_base}/stats")
            data = self._handle_response(response)
            
            if data.get('success'):
                return data.get('stats', {})
            else:
                raise Exception(data.get('error', 'Unknown error'))
        except requests.RequestException as e:
            raise Exception(f"Failed to get statistics: {e}")
    
    def health_check(self) -> bool:
        """Check server health"""
        try:
            response = self.session.get(f"{self.api_base}/health")
            data = self._handle_response(response)
            return data.get('success', False)
        except requests.RequestException:
            return False


def format_bytes(bytes_size: int) -> str:
    """Format bytes to human-readable format"""
    if bytes_size == 0:
        return "0 B"
    
    units = ['B', 'KB', 'MB', 'GB', 'TB']
    unit_index = 0
    size = float(bytes_size)
    
    while size >= 1024 and unit_index < len(units) - 1:
        size /= 1024
        unit_index += 1
    
    return f"{size:.2f} {units[unit_index]}"


def cmd_list(args):
    """List all models"""
    client = ModelClient(args.server)
    
    try:
        models = client.list_models()
        
        if not models:
            print("No models found")
            return 0
        
        # Print header
        print(f"{'NAME':<40} {'STATUS':<10} {'SIZE':<12} {'DOWNLOADS':<12} {'FIRST BY':<20}")
        print("-" * 94)
        
        # Print models
        for model in models:
            stats = model.get('stats', {})
            downloads = stats.get('download_count', 0)
            first_by = stats.get('first_downloaded_by', 'N/A')
            
            print(
                f"{model['name']:<40} "
                f"{model['status']:<10} "
                f"{format_bytes(model['size_bytes']):<12} "
                f"{downloads:<12} "
                f"{first_by:<20}"
            )
        
        return 0
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        return 1


def cmd_get(args):
    """Get a specific model"""
    client = ModelClient(args.server)
    
    try:
        model = client.get_model(args.id)
        
        print("Model Details:")
        print(f"  ID:           {model['id']}")
        print(f"  Name:         {model['name']}")
        print(f"  Status:       {model['status']}")
        print(f"  Size:         {format_bytes(model['size_bytes'])}")
        print(f"  Path:         {model['path']}")
        print(f"  Downloaded:   {model['downloaded_at']}")
        
        if 'stats' in model:
            stats = model['stats']
            print("\nUsage Statistics:")
            print(f"  Downloads:    {stats.get('download_count', 0)}")
            print(f"  Accesses:     {stats.get('access_count', 0)}")
            print(f"  Total Reqs:   {stats.get('total_requests', 0)}")
            
            if stats.get('first_downloaded_by'):
                print(f"  First By:     {stats['first_downloaded_by']}")
                print(f"  First From:   {stats.get('first_downloaded_from', 'N/A')}")
                print(f"  First At:     {stats.get('first_downloaded_at', 'N/A')}")
        
        return 0
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        return 1


def cmd_download(args):
    """Download a model"""
    client = ModelClient(args.server)
    username = getpass.getuser()
    
    print(f"Requesting download of model: {args.name}")
    print(f"Requester: {username}")
    print("This may take a while...")
    
    try:
        result = client.download_model(args.name, username)
        model = result.get('model', {})
        
        if result.get('already_exists'):
            print("\nModel already exists!")
        else:
            print("\nModel downloaded successfully!")
        
        print(f"  ID:          {model.get('id', 'N/A')}")
        print(f"  Name:        {model.get('name', 'N/A')}")
        print(f"  Size:        {format_bytes(model.get('size_bytes', 0))}")
        print(f"  Path:        {model.get('path', 'N/A')}")
        
        stats = model.get('stats', {})
        if stats.get('first_downloaded_by'):
            print(f"  First By:    {stats['first_downloaded_by']}")
        
        return 0
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        return 1


def cmd_delete(args):
    """Delete a model"""
    client = ModelClient(args.server)
    
    print(f"Deleting model: {args.id}")
    
    try:
        client.delete_model(args.id)
        print("Model deleted successfully!")
        return 0
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        return 1


def cmd_stats(args):
    """Get system statistics"""
    client = ModelClient(args.server)
    
    try:
        stats = client.get_stats()
        
        print("System Statistics:")
        print(f"  Total Models:    {stats.get('total_models', 0)}")
        print(f"  Total Size:      {format_bytes(stats.get('total_size_bytes', 0))}")
        print(f"  Total Requests:  {stats.get('total_requests', 0)}")
        print(f"  Unique Users:    {stats.get('unique_users', 0)}")
        
        recent = stats.get('recent_activity', [])
        if recent:
            print(f"\nRecent Activity (last {len(recent)}):")
            for activity in recent:
                print(f"  - {activity['action']} by {activity['username']} from {activity['ip_address']}")
        
        return 0
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        return 1


def cmd_health(args):
    """Check server health"""
    client = ModelClient(args.server)
    
    if client.health_check():
        print("Server is healthy!")
        return 0
    else:
        print("Server is not healthy", file=sys.stderr)
        return 1


def main():
    """Main entry point"""
    parser = argparse.ArgumentParser(
        description='Model Management System - Python CLI Client',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  %(prog)s list
  %(prog)s get -i bert-base-uncased
  %(prog)s download -n google/flan-t5-small
  %(prog)s delete -i bert-base-uncased
  %(prog)s stats
  %(prog)s health
        """
    )
    
    parser.add_argument(
        '-s', '--server',
        default=DEFAULT_SERVER_URL,
        help=f'Server URL (default: {DEFAULT_SERVER_URL})'
    )
    
    subparsers = parser.add_subparsers(dest='command', help='Available commands')
    
    # List command
    parser_list = subparsers.add_parser('list', help='List all models')
    parser_list.set_defaults(func=cmd_list)
    
    # Get command
    parser_get = subparsers.add_parser('get', help='Get a specific model')
    parser_get.add_argument('-i', '--id', required=True, help='Model ID')
    parser_get.set_defaults(func=cmd_get)
    
    # Download command
    parser_download = subparsers.add_parser('download', help='Download a model from Hugging Face')
    parser_download.add_argument('-n', '--name', required=True, help='Model name from Hugging Face')
    parser_download.set_defaults(func=cmd_download)
    
    # Delete command
    parser_delete = subparsers.add_parser('delete', help='Delete a model')
    parser_delete.add_argument('-i', '--id', required=True, help='Model ID to delete')
    parser_delete.set_defaults(func=cmd_delete)
    
    # Stats command
    parser_stats = subparsers.add_parser('stats', help='Get system statistics')
    parser_stats.set_defaults(func=cmd_stats)
    
    # Health command
    parser_health = subparsers.add_parser('health', help='Check server health')
    parser_health.set_defaults(func=cmd_health)
    
    args = parser.parse_args()
    
    if not args.command:
        parser.print_help()
        return 1
    
    return args.func(args)


if __name__ == '__main__':
    sys.exit(main())
