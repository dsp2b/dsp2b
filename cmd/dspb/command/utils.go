package command

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/codfrm/cago/pkg/logger"
	"io"
	"os"
)

func ReadRepository() (*Repository, error) {
	// 读取仓库信息
	repo := &Repository{}
	data, err := os.ReadFile("dspb.json")
	if err != nil {
		logger.Default().Error("请先初始化仓库信息: dspb init")
		return nil, err
	}
	if err := json.Unmarshal(data, repo); err != nil {
		return nil, err
	}
	return repo, nil
}

func SaveRepository(repo *Repository) error {
	data, err := json.MarshalIndent(repo, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("dspb.json", data, os.ModePerm)
}

func Hash(data []byte) (string, error) {
	hash := sha1.New()
	if _, err := hash.Write(data); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func HashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func printBlueColor(str string) {
	fmt.Printf("\033[34m%s\033[0m", str)
}

func printYellowColor(str string) {
	fmt.Printf("\033[33m%s\033[0m", str)
}

func printRedColor(str string) {
	fmt.Printf("\033[31m%s\033[0m", str)
}
