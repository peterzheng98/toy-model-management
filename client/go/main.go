package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"text/tabwriter"
	"time"
)

const (
	defaultServerURL = "http://localhost:5000"
)

// Model represents a model entity
type Model struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Path         string         `json:"path"`
	SizeBytes    int64          `json:"size_bytes"`
	DownloadedAt string         `json:"downloaded_at"`
	DownloadedBy string         `json:"downloaded_by,omitempty"`
	Status       string         `json:"status"`
	UpdatedAt    string         `json:"updated_at,omitempty"`
	Stats        *ModelStats    `json:"stats,omitempty"`
}

// ModelStats represents usage statistics for a model
type ModelStats struct {
	DownloadCount      int    `json:"download_count"`
	AccessCount        int    `json:"access_count"`
	TotalRequests      int    `json:"total_requests"`
	FirstDownloadedBy  string `json:"first_downloaded_by,omitempty"`
	FirstDownloadedAt  string `json:"first_downloaded_at,omitempty"`
	FirstDownloadedFrom string `json:"first_downloaded_from,omitempty"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Success       bool    `json:"success"`
	Message       string  `json:"message,omitempty"`
	Error         string  `json:"error,omitempty"`
	Models        []Model `json:"models,omitempty"`
	Model         *Model  `json:"model,omitempty"`
	AlreadyExists bool    `json:"already_exists,omitempty"`
}

// Client is the HTTP client for the model management server
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new client
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = defaultServerURL
	}
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute, // Long timeout for downloads
		},
	}
}

// ListModels retrieves all models from the server
func (c *Client) ListModels() ([]Model, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/api/models")
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Error)
	}

	return apiResp.Models, nil
}

// GetModel retrieves a specific model by ID
func (c *Client) GetModel(modelID string) (*Model, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/api/models/" + modelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get model: %w", err)
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Error)
	}

	return apiResp.Model, nil
}

// DownloadModel requests the server to download a model from Hugging Face
func (c *Client) DownloadModel(modelName, username string) (*Model, error) {
	payload := map[string]string{
		"model_name": modelName,
		"username":   username,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.HTTPClient.Post(
		c.BaseURL+"/api/models/download",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to download model: %w", err)
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Error)
	}

	return apiResp.Model, nil
}

// DeleteModel deletes a model from the server
func (c *Client) DeleteModel(modelID string) error {
	req, err := http.NewRequest("DELETE", c.BaseURL+"/api/models/"+modelID, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete model: %w", err)
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return fmt.Errorf("API error: %s", apiResp.Error)
	}

	return nil
}

// UpdateModel updates model metadata
func (c *Client) UpdateModel(modelID string, updates map[string]interface{}) (*Model, error) {
	jsonData, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(
		"PUT",
		c.BaseURL+"/api/models/"+modelID,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update model: %w", err)
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Error)
	}

	return apiResp.Model, nil
}

// HealthCheck checks if the server is healthy
func (c *Client) HealthCheck() error {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/api/health")
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return fmt.Errorf("server unhealthy")
	}

	return nil
}

// getCurrentUsername gets the current system username
func getCurrentUsername() string {
	// Try to get username from environment
	if username := os.Getenv("USER"); username != "" {
		return username
	}
	if username := os.Getenv("USERNAME"); username != "" {
		return username
	}
	
	// Try to get from user package
	if currentUser, err := user.Current(); err == nil {
		return currentUser.Username
	}
	
	// Try whoami command as fallback
	if output, err := exec.Command("whoami").Output(); err == nil {
		return strings.TrimSpace(string(output))
	}
	
	return "unknown"
}

// formatBytes formats bytes to human-readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func main() {
	// Define subcommands
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	healthCmd := flag.NewFlagSet("health", flag.ExitOnError)

	// Server URL flag for all commands
	serverURL := ""
	listCmd.StringVar(&serverURL, "server", defaultServerURL, "Server URL")
	getCmd.StringVar(&serverURL, "server", defaultServerURL, "Server URL")
	downloadCmd.StringVar(&serverURL, "server", defaultServerURL, "Server URL")
	deleteCmd.StringVar(&serverURL, "server", defaultServerURL, "Server URL")
	healthCmd.StringVar(&serverURL, "server", defaultServerURL, "Server URL")

	// Command-specific flags
	getModelID := getCmd.String("id", "", "Model ID")
	downloadModelName := downloadCmd.String("name", "", "Model name from Hugging Face")
	deleteModelID := deleteCmd.String("id", "", "Model ID to delete")

	if len(os.Args) < 2 {
		fmt.Println("Model Management Client")
		fmt.Println("\nUsage:")
		fmt.Println("  client <command> [options]")
		fmt.Println("\nCommands:")
		fmt.Println("  list          List all models")
		fmt.Println("  get           Get a specific model")
		fmt.Println("  download      Download a model from Hugging Face")
		fmt.Println("  delete        Delete a model")
		fmt.Println("  health        Check server health")
		fmt.Println("\nExamples:")
		fmt.Println("  client list")
		fmt.Println("  client get -id bert-base-uncased")
		fmt.Println("  client download -name google/flan-t5-small")
		fmt.Println("  client delete -id bert-base-uncased")
		fmt.Println("  client health")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "list":
		listCmd.Parse(os.Args[2:])
		client := NewClient(serverURL)

		models, err := client.ListModels()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(models) == 0 {
			fmt.Println("No models found")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tSTATUS\tSIZE\tDOWNLOADS\tFIRST BY")
		for _, model := range models {
			downloads := 0
			firstBy := "N/A"
			if model.Stats != nil {
				downloads = model.Stats.DownloadCount
				if model.Stats.FirstDownloadedBy != "" {
					firstBy = model.Stats.FirstDownloadedBy
				}
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
				model.Name,
				model.Status,
				formatBytes(model.SizeBytes),
				downloads,
				firstBy,
			)
		}
		w.Flush()

	case "get":
		getCmd.Parse(os.Args[2:])
		if *getModelID == "" {
			fmt.Fprintln(os.Stderr, "Error: -id is required")
			getCmd.PrintDefaults()
			os.Exit(1)
		}

		client := NewClient(serverURL)
		model, err := client.GetModel(*getModelID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Model Details:\n")
		fmt.Printf("  ID:           %s\n", model.ID)
		fmt.Printf("  Name:         %s\n", model.Name)
		fmt.Printf("  Status:       %s\n", model.Status)
		fmt.Printf("  Size:         %s\n", formatBytes(model.SizeBytes))
		fmt.Printf("  Path:         %s\n", model.Path)
		fmt.Printf("  Downloaded:   %s\n", model.DownloadedAt)
		
		if model.Stats != nil {
			fmt.Printf("\nUsage Statistics:\n")
			fmt.Printf("  Downloads:    %d\n", model.Stats.DownloadCount)
			fmt.Printf("  Accesses:     %d\n", model.Stats.AccessCount)
			fmt.Printf("  Total Reqs:   %d\n", model.Stats.TotalRequests)
			if model.Stats.FirstDownloadedBy != "" {
				fmt.Printf("  First By:     %s\n", model.Stats.FirstDownloadedBy)
				fmt.Printf("  First From:   %s\n", model.Stats.FirstDownloadedFrom)
				fmt.Printf("  First At:     %s\n", model.Stats.FirstDownloadedAt)
			}
		}

	case "download":
		downloadCmd.Parse(os.Args[2:])
		if *downloadModelName == "" {
			fmt.Fprintln(os.Stderr, "Error: -name is required")
			downloadCmd.PrintDefaults()
			os.Exit(1)
		}

		client := NewClient(serverURL)
		username := getCurrentUsername()
		
		fmt.Printf("Requesting download of model: %s\n", *downloadModelName)
		fmt.Printf("Requester: %s\n", username)
		fmt.Println("This may take a while...")

		model, err := client.DownloadModel(*downloadModelName, username)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nModel downloaded successfully!\n")
		fmt.Printf("  ID:          %s\n", model.ID)
		fmt.Printf("  Name:        %s\n", model.Name)
		fmt.Printf("  Size:        %s\n", formatBytes(model.SizeBytes))
		fmt.Printf("  Path:        %s\n", model.Path)
		if model.Stats != nil && model.Stats.FirstDownloadedBy != "" {
			fmt.Printf("  First By:    %s\n", model.Stats.FirstDownloadedBy)
		}

	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if *deleteModelID == "" {
			fmt.Fprintln(os.Stderr, "Error: -id is required")
			deleteCmd.PrintDefaults()
			os.Exit(1)
		}

		client := NewClient(serverURL)
		fmt.Printf("Deleting model: %s\n", *deleteModelID)

		err := client.DeleteModel(*deleteModelID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Model deleted successfully!")

	case "health":
		healthCmd.Parse(os.Args[2:])
		client := NewClient(serverURL)

		err := client.HealthCheck()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Server is healthy!")

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}
