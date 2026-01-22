package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"
)

const serverURL = "http://localhost:5000"

// getCurrentUsername gets the current system username
func getCurrentUsername() string {
	if username := os.Getenv("USER"); username != "" {
		return username
	}
	if username := os.Getenv("USERNAME"); username != "" {
		return username
	}
	if currentUser, err := user.Current(); err == nil {
		return currentUser.Username
	}
	if output, err := exec.Command("whoami").Output(); err == nil {
		return strings.TrimSpace(string(output))
	}
	return "unknown"
}

// Model represents a model entity
type Model struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Path         string      `json:"path"`
	SizeBytes    int64       `json:"size_bytes"`
	DownloadedAt string      `json:"downloaded_at"`
	DownloadedBy string      `json:"downloaded_by,omitempty"`
	Status       string      `json:"status"`
	Stats        *ModelStats `json:"stats,omitempty"`
}

// ModelStats represents usage statistics
type ModelStats struct {
	DownloadCount       int    `json:"download_count"`
	AccessCount         int    `json:"access_count"`
	TotalRequests       int    `json:"total_requests"`
	FirstDownloadedBy   string `json:"first_downloaded_by,omitempty"`
	FirstDownloadedAt   string `json:"first_downloaded_at,omitempty"`
	FirstDownloadedFrom string `json:"first_downloaded_from,omitempty"`
}

// SystemStats represents overall system statistics
type SystemStats struct {
	TotalModels    int            `json:"total_models"`
	TotalSizeBytes int64          `json:"total_size_bytes"`
	TotalRequests  int            `json:"total_requests"`
	UniqueUsers    int            `json:"unique_users"`
	RecentActivity []ActivityLog  `json:"recent_activity,omitempty"`
}

// ActivityLog represents an activity log entry
type ActivityLog struct {
	Timestamp string `json:"timestamp"`
	Action    string `json:"action"`
	ModelID   string `json:"model_id,omitempty"`
	Username  string `json:"username"`
	IPAddress string `json:"ip_address"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Success       bool         `json:"success"`
	Message       string       `json:"message,omitempty"`
	Error         string       `json:"error,omitempty"`
	Models        []Model      `json:"models,omitempty"`
	Model         *Model       `json:"model,omitempty"`
	Stats         *SystemStats `json:"stats,omitempty"`
	AlreadyExists bool         `json:"already_exists,omitempty"`
}

func healthCheck() bool {
	fmt.Println("Checking server health...")

	resp, err := http.Get(serverURL + "/api/health")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return false
	}

	if apiResp.Success {
		fmt.Println("Server is healthy!")
		return true
	}

	fmt.Println("Server is not healthy")
	return false
}

func getStats() {
	fmt.Println("\nGetting system statistics...")

	resp, err := http.Get(serverURL + "/api/stats")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}

	if apiResp.Success && apiResp.Stats != nil {
		stats := apiResp.Stats
		sizeMB := float64(stats.TotalSizeBytes) / (1024 * 1024)
		fmt.Println("Statistics:")
		fmt.Printf("  Total Models: %d\n", stats.TotalModels)
		fmt.Printf("  Total Size: %.2f MB\n", sizeMB)
		fmt.Printf("  Total Requests: %d\n", stats.TotalRequests)
		fmt.Printf("  Unique Users: %d\n", stats.UniqueUsers)
		
		if len(stats.RecentActivity) > 0 {
			fmt.Printf("\n  Recent Activity (last %d):\n", len(stats.RecentActivity))
			for _, activity := range stats.RecentActivity {
				fmt.Printf("    - %s by %s from %s\n", 
					activity.Action, activity.Username, activity.IPAddress)
			}
		}
	} else {
		fmt.Printf("Error: %s\n", apiResp.Error)
	}
}

func listModels() {
	fmt.Println("\nListing all models...")

	resp, err := http.Get(serverURL + "/api/models")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}

	if apiResp.Success {
		if len(apiResp.Models) > 0 {
			fmt.Printf("Found %d model(s):\n", len(apiResp.Models))
			for _, model := range apiResp.Models {
				sizeMB := float64(model.SizeBytes) / (1024 * 1024)
				fmt.Printf("  - %s\n", model.Name)
				fmt.Printf("    ID: %s\n", model.ID)
				fmt.Printf("    Size: %.2f MB\n", sizeMB)
				fmt.Printf("    Status: %s\n", model.Status)
				
				if model.Stats != nil {
					fmt.Printf("    Downloads: %d\n", model.Stats.DownloadCount)
					if model.Stats.FirstDownloadedBy != "" {
						fmt.Printf("    First By: %s\n", model.Stats.FirstDownloadedBy)
					}
				}
				fmt.Println()
			}
		} else {
			fmt.Println("  No models found")
		}
	} else {
		fmt.Printf("Error: %s\n", apiResp.Error)
	}
}

func downloadModel(modelName string) {
	username := getCurrentUsername()
	
	fmt.Printf("\nDownloading model: %s\n", modelName)
	fmt.Printf("   Requester: %s\n", username)
	fmt.Println("   This may take a while...")

	payload := map[string]string{
		"model_name": modelName,
		"username":   username,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	resp, err := http.Post(
		serverURL+"/api/models/download",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}

	if apiResp.Success {
		if apiResp.AlreadyExists {
			fmt.Println("Model already exists!")
		} else {
			fmt.Println("Model downloaded successfully!")
		}

		if apiResp.Model != nil {
			sizeMB := float64(apiResp.Model.SizeBytes) / (1024 * 1024)
			fmt.Printf("  ID: %s\n", apiResp.Model.ID)
			fmt.Printf("  Size: %.2f MB\n", sizeMB)
			fmt.Printf("  Path: %s\n", apiResp.Model.Path)
			
			if apiResp.Model.Stats != nil && apiResp.Model.Stats.FirstDownloadedBy != "" {
				fmt.Printf("  First Downloaded By: %s\n", apiResp.Model.Stats.FirstDownloadedBy)
			}
		}
	} else {
		fmt.Printf("Download failed: %s\n", apiResp.Error)
	}
}

func getModel(modelID string) {
	fmt.Printf("\nGetting model details: %s\n", modelID)

	resp, err := http.Get(serverURL + "/api/models/" + modelID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}

	if apiResp.Success && apiResp.Model != nil {
		model := apiResp.Model
		sizeMB := float64(model.SizeBytes) / (1024 * 1024)
		fmt.Println("Model found!")
		fmt.Printf("  Name:         %s\n", model.Name)
		fmt.Printf("  ID:           %s\n", model.ID)
		fmt.Printf("  Size:         %.2f MB\n", sizeMB)
		fmt.Printf("  Status:       %s\n", model.Status)
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
	} else {
		fmt.Printf("Error: %s\n", apiResp.Error)
	}
}

func deleteModel(modelID string) {
	fmt.Printf("\nDeleting model: %s\n", modelID)

	req, err := http.NewRequest("DELETE", serverURL+"/api/models/"+modelID, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}

	if apiResp.Success {
		fmt.Println("Model deleted successfully!")
	} else {
		fmt.Printf("Delete failed: %s\n", apiResp.Error)
	}
}

func main() {
	username := getCurrentUsername()
	
	fmt.Println("============================================================")
	fmt.Println("Model Management System - Go Example")
	fmt.Printf("Running as user: %s\n", username)
	fmt.Println("============================================================")

	// Check server health
	if !healthCheck() {
		fmt.Println("\nError: Server is not running!")
		fmt.Println("Start the server with: ./start_server.sh")
		return
	}

	// Get system statistics
	getStats()

	// List existing models
	listModels()

	// Example: Download a small model (uncomment to test)
	// Note: This will actually download the model, which may take time
	// downloadModel("gpt2")

	// Example: Get a specific model (uncomment after downloading)
	// getModel("gpt2")

	// Example: Delete a model (uncomment to test)
	// deleteModel("gpt2")

	fmt.Println("\n============================================================")
	fmt.Println("Example complete!")
	fmt.Println("============================================================")
}
