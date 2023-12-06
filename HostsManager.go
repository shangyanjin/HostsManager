/*
HostManager provides a lightweight command-line tool for managing the hosts file on Windows.
Author: Shang Yanjin
Email: shangyanjin@msn.com
*/
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// 获取系统的 hosts 文件路径
	hostsPath := getHostsFilePath()

	// 检查是否没有传递参数
	if len(os.Args) == 1 {
		displayHostsFileContent(hostsPath)
		fmt.Println("Usage: dns.exe [action] [domain] [ip]")
		fmt.Println("Example: dns.exe add example.com 127.0.0.1")
		return
	}

	// 检查参数数量是否正确
	if len(os.Args) < 3 {
		fmt.Println("Error: Insufficient number of arguments. Use 'dns.exe [action] [domain] [ip]'.")
		return
	}

	// 解析命令行参数
	action := os.Args[1]
	domain := os.Args[2]
	var ip string
	if len(os.Args) == 4 {
		ip = os.Args[3]
	} else {
		ip = "127.0.0.1"
	}

	// 检查 IP 地址的合法性
	if !isValidIP(ip) {
		fmt.Println("Error: Invalid IP address.")
		return
	}

	// 根据动作执行相应的操作
	switch action {
	case "add":
		err := addEntry(domain, ip, hostsPath)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Entry added successfully.")
			displayHostsFileContent(hostsPath)
		}
	case "list":
		listEntries(hostsPath)
	case "edit":
		err := editEntry(domain, ip, hostsPath)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Entry edited successfully.")
			displayHostsFileContent(hostsPath)
		}
	case "del":
		err := deleteEntry(domain, hostsPath)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Entry deleted successfully.")
			displayHostsFileContent(hostsPath)
		}
	default:
		fmt.Println("Invalid action. Use 'add', 'list', 'edit', or 'del'.")
	}
}

// 获取系统 hosts 文件路径
func getHostsFilePath() string {
	hostsPath := os.Getenv("SystemRoot") + "\\System32\\drivers\\etc\\hosts"
	return hostsPath
}

// 显示 hosts 文件内容
func displayHostsFileContent(hostsPath string) {
	file, err := os.Open(hostsPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	fmt.Println("Current hosts file content:")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading hosts file:", err)
	}
}

// 添加条目到 hosts 文件
func addEntry(domain, ip, hostsPath string) error {
	file, err := os.OpenFile(hostsPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	entry := fmt.Sprintf("%s\t%s\n", ip, domain)
	_, err = file.WriteString(entry)
	return err
}

// 列出 hosts 文件中的所有条目
func listEntries(hostsPath string) {
	file, err := os.Open(hostsPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	fmt.Println("Current hosts file content:")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading hosts file:", err)
	}
}

// 编辑已存在的 hosts 文件条目
func editEntry(domain, ip, hostsPath string) error {
	file, err := os.Open(hostsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, domain) {
			// 替换已存在的 IP 地址
			line = fmt.Sprintf("%s\t%s", ip, domain)
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	file, err = os.Create(hostsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// 从 hosts 文件中删除指定的条目
func deleteEntry(domain, hostsPath string) error {
	file, err := os.Open(hostsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, domain) {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	file, err = os.Create(hostsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// 检查 IP 地址的合法性
func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}
