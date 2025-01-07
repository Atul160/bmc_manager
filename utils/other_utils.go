package utils

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"time"

	"github.com/go-ping/ping"
	"github.com/sfreiberg/simplessh"
)

func loadEnvVars() map[string]string {
	envVars := make(map[string]string)
	for _, env := range os.Environ() {
		parts := strings.Split(env, "=")
		key := parts[0]
		value := parts[1]
		envVars[key] = value
	}
	return envVars
}

// Helper function to get the current time in a formatted string
func GetCurrentTime() string {
	now := time.Now()
	return now.Format("2006-01-02 15:04:05")
}

func TCPPing(ip string, port string) bool {
	timeout := time.Second * 2
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
	if err != nil {
		fmt.Printf("TCP ping failed: %v\n", err)
		log.Printf("TCP ping failed: %v\n", err)
		return false
	}
	conn.Close()
	return true
}

// Fallback method using shell command
func shellPing(ip string) bool {
	cmd := exec.Command("ping", "-c", "2", "-t", "2", ip)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Shell ping failed: %v\n", err)
		log.Printf("Shell ping failed: %v\n", err)
		return false
	}
	return true
}

// Ping function to check the availability of the IP
func PingIP(ip string) bool {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		fmt.Printf("Error creating pinger: %v. Falling back to shell ping.\n", err)
		log.Printf("Error creating pinger: %v. Falling back to shell ping.\n", err)
		return shellPing(ip) // Fallback method
	}

	pinger.SetPrivileged(true)
	pinger.Count = 2 // Send 2 pings
	pinger.Timeout = time.Second * 2

	err = pinger.Run()
	if err != nil {
		fmt.Printf("Ping failed: %v\n", err)
		log.Printf("Ping failed: %v\n", err)
		return false
	}

	stats := pinger.Statistics()

	return stats.PacketsRecv > 0
}

// ResolveDNS takes an IP address as input and returns the FQDN
func ResolveDNS(ipAddress string) (string, error) {
	// Perform reverse DNS lookup using net.LookupAddr
	names, err := net.LookupAddr(ipAddress)
	if err != nil {
		return "", fmt.Errorf("failed to resolve DNS for %s: %v", ipAddress, err)
	}

	// Check if any domain names were found
	if len(names) == 0 {
		return "", fmt.Errorf("no FQDN found for IP address %s", ipAddress)
	}

	// Return the first resolved domain name (FQDN)
	return names[0], nil
}

// Function to create SSH connection
func CreateSSHConnection(ip string, username string, password string) (*simplessh.Client, error) {
	// Create a new SSH client
	var client *simplessh.Client
	var err error

	client, err = simplessh.ConnectWithPassword(fmt.Sprintf("%s:22", ip), username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH connection: %v", err)
	}
	return client, nil
}

// Function to execute SSH command
func ExecuteSSHCommand(client *simplessh.Client, command string) (string, error) {
	defer client.Close()
	// Execute the command
	output, err := client.Exec(command)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %v", err)
	}

	return string(output), nil
}

func ReadNestedMap(nestedMap map[string]interface{}, depth int) {
	if depth <= 0 {
		return
	}

	for key, value := range nestedMap {
		switch reflect.TypeOf(value).Kind() {
		case reflect.Map:
			fmt.Printf("Key: %s, Value: (nested map)\n", key)
			ReadNestedMap(value.(map[string]interface{}), depth-1)
		case reflect.Slice:
			fmt.Printf("Key: %s, Value: (array)\n", key)
			// Handle array elements if needed
		default:
			fmt.Printf("Key: %s, Value: %v\n", key, value)
		}
	}
}
