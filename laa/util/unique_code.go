package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

func GetUniqueCode() (string, error) {
	if runtime.GOOS != "windows" {
		return "", fmt.Errorf("当前操作系统不是Windows，无法获取唯一机器码")
	}
	cpuID, _ := getCpuIdentifier()
	motherboardSN, _ := getMotherboardSerialNumber()
	addressMAC, _ := getFirstMacAddress()
	var uniqueCode string
	switch {
	case cpuID != "" && motherboardSN != "":
		uniqueCode = cpuID + "_" + motherboardSN
	case cpuID != "" && addressMAC != "":
		uniqueCode = cpuID + "_" + addressMAC
	case motherboardSN != "" && addressMAC != "":
		uniqueCode = motherboardSN + "_" + addressMAC
	case cpuID != "":
		uniqueCode = cpuID
	case motherboardSN != "":
		uniqueCode = motherboardSN
	case addressMAC != "":
		uniqueCode = addressMAC
	}
	if uniqueCode == "" {
		return "", fmt.Errorf("未能获取到任何可用的唯一机器标识符")
	}
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
	return "", fmt.Errorf("拿不到CPU序列号")
}

func getMotherboardSerialNumber() (string, error) {
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
	return "", fmt.Errorf("拿不到主板序列号")
}

func getFirstMacAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if len(iface.HardwareAddr) > 0 {
			return iface.HardwareAddr.String(), nil
		}
	}

	return "", fmt.Errorf("未找到有效的 MAC 地址")
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
