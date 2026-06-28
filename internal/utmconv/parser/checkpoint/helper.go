package checkpoint

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func b2s(b bool) string {
	return strconv.FormatBool(b)
}

func itoa(i int) string {
	return strconv.Itoa(i)
}

func join(ss []string) string {
	return strings.Join(ss, ";")
}

func joinNames(objs []CPObject) string {
	var sb strings.Builder
	for i, o := range objs {
		sb.WriteString(o.Name)
		if i != 0 {
			sb.WriteRune(';')
		}
	}
	return sb.String()
}

func extractToTemp(tarPath, tempDir string) error {
	file, err := os.Open(tarPath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	zr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	tr := tar.NewReader(zr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading tar archive: %w", err)
		}

		targetPath := filepath.Join(tempDir, hdr.Name)

		if !strings.HasPrefix(targetPath,
			filepath.Clean(tempDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", hdr.Name)
		}

		if hdr.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return err
			}
			outFile, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			_, err = io.Copy(outFile, tr)
			closeErr := outFile.Close()
			if err != nil {
				return fmt.Errorf("failed to write file %s: %w", hdr.Name, err)
			}
			if closeErr != nil {
				return fmt.Errorf("failed to close file: %s: %w",
					hdr.Name, closeErr)
			}
		}
	}
	return nil
}

func processTarGzip(tarFile string, fn func(tempDir string) error) error {
	tempDir, err := os.MkdirTemp("", "utmconv-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer func() {
		slog.Info("Cleaning up temporary directory", slog.String("dir", tempDir))
		os.RemoveAll(tempDir)
	}()

	err = extractToTemp(tarFile, tempDir)
	if err != nil {
		return err
	}

	slog.Info("successfully extracted to directory", slog.String("dir", tempDir))

	return fn(tempDir)
}
