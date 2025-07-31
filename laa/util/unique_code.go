package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func GetUniqueCode() (string, error) {
	if runtime.GOOS != "windows" {
		return "", fmt.Errorf("当前操作系统不是Windows，无法获取唯一机器码")
	}

	cpuID, err := getCpuIdentifier()
	if err != nil {
		return "", fmt.Errorf("获取CPU标识符失败: %w", err)
	}
	if cpuID == "" {
		return "", fmt.Errorf("未获取到CPU标识符")
	}

	motherboardSN, err := getMotherboardSerialNumber()
	if err != nil {
		return "", fmt.Errorf("获取主板序列号失败: %w", err)
	}
	if motherboardSN == "" {
		return "", fmt.Errorf("未获取到主板序列号")
	}

	uniqueCode := cpuID + "_" + motherboardSN
	hasher := sha256.New()
	hasher.Write([]byte(uniqueCode))
	hashInBytes := hasher.Sum(nil)
	hashedCode := hex.EncodeToString(hashInBytes)
	return hashedCode, nil
}

func getCpuIdentifier() (string, error) {
	cmdStr := "wmic cpu get ProcessorId"
	output, err := executeCommand(cmdStr)
	if err != nil {
		return "", err
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" && trimmedLine != "ProcessorId" {
			return trimmedLine, nil
		}
	}
	return "", fmt.Errorf("未从wmic output中解析到ProcessorId")
}

func getMotherboardSerialNumber() (string, error) {
	// 对应 Kotlin 中的 `wmic baseboard get serialnumber`
	cmdStr := "wmic baseboard get serialnumber"
	output, err := executeCommand(cmdStr)
	if err != nil {
		return "", err
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" && trimmedLine != "SerialNumber" {
			return trimmedLine, nil
		}
	}
	return "", fmt.Errorf("未从wmic output中解析到SerialNumber")
}

func executeCommand(commandStr string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("cmd", "/C", commandStr)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("执行命令失败: %s, 错误: %w, stderr: %s", commandStr, err, stderr.String())
	}

	return stdout.String(), nil
}
