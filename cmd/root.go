/*
Copyright © 2025 NAME HERE darshbothra007@gmail.com
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// first we'll define the task i.e. what will it have
type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type TaskManager struct {
	Tasks    []Task `json:"tasks"`
	NextID   int    `json:"next_id"`
	FilePath string `json:"-"`
}

var (
	taskFile string
	showAll  bool
	taskman  *TaskManager
	m        = make(map[int]string)
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "todo-cli.com",
	Short: "A personal task manager",
	Long: `taskman is a CLI task manager that helps you organize your work.

Store tasks locally with priorities, mark them complete, and keep
track of your productivity over time.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun:  loadTodo,
	PersistentPostRun: saveTodo,
}

var addCmd = &cobra.Command{
	Use:   "add [description]",
	Short: "Add the todos",
	Long:  "Add your todos with detailed description",
	Args:  cobra.MinimumNArgs(1),
	Run:   addTodo,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the todos",
	Long:  "List your todos",
	Run:   listTodo,
}

var listByIdCmd = &cobra.Command{
	Use:   "list [id]",
	Short: "List the todos by id",
	Long:  "List your todos by id",
	Args:  cobra.MinimumNArgs(1),
	Run:   listTodoById,
}

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete the todos by id",
	Long:  "Delete your todos by id",
	Args:  cobra.MinimumNArgs(1),
	Run:   deleteTodoById,
}

var completedCmd = &cobra.Command{
	Use:   "complete [id]",
	Short: "Complete the todos by id",
	Long:  "Complete your todos by id",
	Args:  cobra.MinimumNArgs(1),
	Run:   completeTodoById,
}

func addTodo(cmd *cobra.Command, args []string) {
	des := strings.Join(args, " ")

	task := Task{
		ID:          taskman.NextID,
		Description: des,
		Completed:   false,
		CreatedAt:   time.Now(),
	}

	m[taskman.NextID] = des
	taskman.Tasks = append(taskman.Tasks, task)
	taskman.NextID++

	fmt.Printf("Added task #%d: %s\n", task.ID, task.Description)
}

func listTodo(cmd *cobra.Command, args []string) {
	if len(taskman.Tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	fmt.Printf("%-4s %-10s %-50s %s\n", "ID", "STATUS", "DESCRIPTION", "CREATED")
	fmt.Println(strings.Repeat("-", 90))
	for _, x := range taskman.Tasks {
		status := "Pending"
		if x.Completed {
			status = "Done"
		}
		fmt.Printf("%-4d %-10s %-50s %-10s\n", 
			x.ID, 
			status, 
			x.Description,
			x.CreatedAt.Format("2006-01-02"))
	}
}

func listTodoById(cmd *cobra.Command, args []string) {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Invalid task id %d", id)
		os.Exit(1)
	}
	val, exists := m[id]

	if exists {
		fmt.Printf("%s", val)
	} else {
		fmt.Println("No task for this id exists")
	}
}

func deleteTodoById(cmd *cobra.Command, args []string) {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Invalid task id %d", id)
		os.Exit(1)
	}

	for i, x := range taskman.Tasks {
		if x.ID == id {
			taskman.Tasks = append(taskman.Tasks[:i], taskman.Tasks[i+1:]...)
			delete(m, x.ID)
			fmt.Printf("Deleted task #%d: %s\n", id, x.Description)
			return
		}
	}

	fmt.Printf("Task #%d not found\n", id)
	os.Exit(1)
}

func completeTodoById(cmd *cobra.Command, args []string) {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Invalid task id %d", id)
		os.Exit(1)
	}

	for _, x := range taskman.Tasks {
		if x.ID == id {
			if x.Completed {
				fmt.Printf("Task #%d is already completed\n", id)
				return
			}
			now := time.Now()
			x.Completed = true
			x.CompletedAt = &now
			fmt.Printf("Completed task #%d: %s\n", id, x.Description)
			return
		}
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.todo-cli.com.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().StringVar(&taskFile, "file", "", "task file (default is $HOME/.taskman.json)")
	listCmd.Flags().BoolVarP(&showAll, "all", "a", false, "show completed tasks too")
	rootCmd.AddCommand(addCmd, listCmd, listByIdCmd, completedCmd, deleteCmd)

	setupConfig()
}

func setupConfig() {
	viper.SetConfigName("taskman")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.SetDefault("file", filepath.Join(os.Getenv("HOME"), ".taskman.json"))
	viper.ReadInConfig()
}

func getTodoFile() string {
	if taskFile != "" {
		return taskFile
	}

	return viper.GetString("file")
}

func loadTodo(cmd *cobra.Command, args []string) {
	file := getTodoFile()
	taskman = &TaskManager{
		Tasks:    make([]Task, 0),
		NextID:   1,
		FilePath: file,
	}

	if data, err := os.ReadFile(file); err == nil {
		json.Unmarshal(data, taskman)
	}
}

func saveTodo(cmd *cobra.Command, args []string) {
	data, err := json.MarshalIndent(taskman, "", " ")

	if err != nil {
		fmt.Printf("Their is an error in saving file : %v\n", err)
	}
	os.WriteFile(taskman.FilePath, data, 0644)
}
