package media

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ImageManager interface {
	Clean(filenames []string) error
	Download(ctx context.Context, url string) (string, error)
}

type ImageManagerService struct {
	logger *zap.Logger
	path   string
}

func NewImageManagerService(dest string, logger *zap.Logger) *ImageManagerService {
	return &ImageManagerService{logger, dest}
}

func (d *ImageManagerService) Clean(filenames []string) error {
	for _, filename := range filenames {
		err := os.Remove(path.Join(d.path, filename))
		d.logger.Error("error from media manager", zap.Error(err))
	}
	return nil
}

func (d *ImageManagerService) Download(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}
	if res.StatusCode > 299 {
		return "", fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
	}
	fileID := uuid.NewString()
	err = os.WriteFile(path.Join(d.path, fileID), body, 0666)
	if err != nil {
		return "", err
	}
	return fileID, nil
}
