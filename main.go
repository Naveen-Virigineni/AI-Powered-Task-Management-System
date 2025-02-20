package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	AssignedTo  string `json:"assignedTo"`
	Status      string `json:"status"`
}

var users = make(map[string]string) // In-memory user store
var tasks = make(map[string]Task)  // In-memory task store

func main() {
	r := gin.Default()

	// User Authentication
	r.POST("/signup", signup)
	r.POST("/login", login)

	// Task Management
	r.POST("/tasks", createTask)
	r.GET("/tasks", getTasks)
	r.PUT("/tasks/:id", updateTask)

	// WebSocket for Real-time Updates
	r.GET("/ws", handleWebSocket)

	// Start server
	log.Fatal(r.Run(":8080"))
}

// Signup Handler
func signup(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	users[user.Username] = string(hashedPassword)
	c.JSON(http.StatusOK, gin.H{"message": "User created"})
}

// Login Handler
func login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, exists := users[user.Username]
	if !exists || bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// Create Task Handler
func createTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.ID = time.Now().Format("20060102150405")
	tasks[task.ID] = task
	c.JSON(http.StatusOK, task)
}

// Get Tasks Handler
func getTasks(c *gin.Context) {
	c.JSON(http.StatusOK, tasks)
}

// Update Task Handler
func updateTask(c *gin.Context) {
	id := c.Param("id")
	var updatedTask Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if task, exists := tasks[id]; exists {
		task.Status = updatedTask.Status
		tasks[id] = task
		c.JSON(http.StatusOK, task)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
	}
}

// WebSocket Handler
func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		// Broadcast task updates to all connected clients
		conn.WriteJSON(tasks)
		time.Sleep(2 * time.Second) // Simulate real-time updates
	}
}
