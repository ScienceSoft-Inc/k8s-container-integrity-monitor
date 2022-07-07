package api

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/k8s-container-integrity-monitor/internal/core/consts"

	"github.com/sirupsen/logrus"
)

// SearchFilePath searches for all files in the given directory
func SearchFilePath(ctx context.Context, commonPath string, jobs chan<- string, logger *logrus.Logger) {
	_, cancel := context.WithTimeout(ctx, consts.TimeOut*time.Second)
	defer cancel()
	err := filepath.Walk(commonPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			jobs <- path
		}
		if err != nil {
			logger.Error("err while going to path files", err)
			return err
		}

		return nil
	})
	close(jobs)

	if err != nil {
		logger.Error("not exist directory path", err)
		return
	}
}

// Result launching an infinite loop of receiving and outputting to Stdout the result and signal control
func Result(ctx context.Context, results chan HashData, c chan os.Signal) []HashData {
	var allHashData []HashData
	for {
		select {
		case hashData, ok := <-results:
			if !ok {
				return allHashData
			}
			allHashData = append(allHashData, hashData)
		case <-c:
			fmt.Println("exit program")
			return []HashData{}
		case <-ctx.Done():
		}
	}
}
