package services

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/image/draw"
	"golang.org/x/image/webp"

	"xorriso-ui/pkg/models"
)

type ProjectService struct{}

func NewProjectService() *ProjectService {
	return &ProjectService{}
}

// NewProject creates a new empty project
func (s *ProjectService) NewProject(name string, volumeID string) *models.Project {
	return &models.Project{
		Name:     name,
		VolumeID: volumeID,
		Entries:  []models.FileEntry{},
		ISOOptions: models.ISOOptions{
			RockRidge: true,
			Joliet:    true,
			MD5:       true,
		},
		BurnOptions: models.BurnOptions{
			Speed:    "auto",
			Verify:   true,
			Eject:    true,
			BurnMode: "auto",
			Padding:  300,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// SaveProject saves the project to its file path
func (s *ProjectService) SaveProject(project *models.Project) error {
	if project.FilePath == "" {
		return nil
	}
	project.UpdatedAt = time.Now()
	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(project.FilePath, data, 0644)
}

// SaveProjectAs saves the project to a new file path
func (s *ProjectService) SaveProjectAs(project *models.Project, filePath string) error {
	project.FilePath = filePath
	return s.SaveProject(project)
}

// OpenProject loads a project from file
func (s *ProjectService) OpenProject(filePath string) (*models.Project, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var project models.Project
	if err := json.Unmarshal(data, &project); err != nil {
		return nil, err
	}
	project.FilePath = filePath
	return &project, nil
}

// GetHomeDirectory returns the current user's home directory
func (s *ProjectService) GetHomeDirectory() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "/"
	}
	return home
}

// BrowseDirectory returns contents of a local directory
func (s *ProjectService) BrowseDirectory(path string) ([]models.FileEntry, error) {
	if path == "" {
		path = "/"
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var result []models.FileEntry
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		fe := models.FileEntry{
			SourcePath: filepath.Join(path, entry.Name()),
			Name:       entry.Name(),
			IsDir:      entry.IsDir(),
			Size:       info.Size(),
		}
		result = append(result, fe)
	}

	// Sort: directories first, then by name
	sort.Slice(result, func(i, j int) bool {
		if result[i].IsDir != result[j].IsDir {
			return result[i].IsDir
		}
		return result[i].Name < result[j].Name
	})

	return result, nil
}

// AddFiles adds files/directories to the project
func (s *ProjectService) AddFiles(project *models.Project, sourcePaths []string, destDir string) (*models.Project, error) {
	for _, src := range sourcePaths {
		info, err := os.Stat(src)
		if err != nil {
			continue
		}

		destPath := filepath.Join(destDir, filepath.Base(src))
		entry := models.FileEntry{
			SourcePath: src,
			DestPath:   destPath,
			Name:       filepath.Base(src),
			IsDir:      info.IsDir(),
			Size:       info.Size(),
		}

		if info.IsDir() {
			size, _ := dirSize(src)
			entry.Size = size
		}

		project.Entries = append(project.Entries, entry)
	}
	project.UpdatedAt = time.Now()
	return project, nil
}

// RemoveEntry removes a file entry from the project by dest path
func (s *ProjectService) RemoveEntry(project *models.Project, destPath string) (*models.Project, error) {
	filtered := make([]models.FileEntry, 0, len(project.Entries))
	for _, e := range project.Entries {
		if e.DestPath != destPath {
			filtered = append(filtered, e)
		}
	}
	project.Entries = filtered
	project.UpdatedAt = time.Now()
	return project, nil
}

// CalculateSize returns total size of all entries in bytes
func (s *ProjectService) CalculateSize(project *models.Project) (int64, error) {
	var total int64
	for _, e := range project.Entries {
		total += e.Size
	}
	return total, nil
}

// MountPoint represents a mounted filesystem location for quick navigation
type MountPoint struct {
	Label string `json:"label"`
	Path  string `json:"path"`
	Icon  string `json:"icon"`
}

// ListMountPoints returns Home + mounted external drives/USB devices
func (s *ProjectService) ListMountPoints() ([]MountPoint, error) {
	var points []MountPoint

	// Always include Home first
	home, err := os.UserHomeDir()
	if err != nil {
		home = "/"
	}
	points = append(points, MountPoint{Label: "Home", Path: home, Icon: "home"})

	// Parse /proc/mounts for external drives
	f, err := os.Open("/proc/mounts")
	if err != nil {
		return points, nil
	}
	defer f.Close()

	// Filesystem types to skip (virtual/pseudo filesystems)
	skipFS := map[string]bool{
		"sysfs": true, "proc": true, "devtmpfs": true, "devpts": true,
		"tmpfs": true, "securityfs": true, "cgroup": true, "cgroup2": true,
		"pstore": true, "efivarfs": true, "bpf": true, "autofs": true,
		"mqueue": true, "hugetlbfs": true, "debugfs": true, "tracefs": true,
		"fusectl": true, "configfs": true, "ramfs": true, "rpc_pipefs": true,
		"nfsd": true, "overlay": true, "nsfs": true, "fuse.portal": true,
		"fuse.gvfsd-fuse": true,
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 3 {
			continue
		}
		device := unescapeMountPath(fields[0])
		mountPath := unescapeMountPath(fields[1])
		fsType := fields[2]

		// Skip virtual filesystems
		if skipFS[fsType] {
			continue
		}

		// Only include mounts under /media/, /mnt/, /run/media/
		if !strings.HasPrefix(mountPath, "/media/") &&
			!strings.HasPrefix(mountPath, "/mnt/") &&
			!strings.HasPrefix(mountPath, "/run/media/") {
			continue
		}

		// Skip if device doesn't look like a real block device
		if !strings.HasPrefix(device, "/dev/") {
			continue
		}

		label := filepath.Base(mountPath)
		points = append(points, MountPoint{Label: label, Path: mountPath, Icon: "usb"})
	}

	return points, nil
}

// unescapeMountPath decodes octal escape sequences used in /proc/mounts
// (e.g. \040 = space, \011 = tab, \134 = backslash, \012 = newline)
func unescapeMountPath(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+3 < len(s) &&
			s[i+1] >= '0' && s[i+1] <= '3' &&
			s[i+2] >= '0' && s[i+2] <= '7' &&
			s[i+3] >= '0' && s[i+3] <= '7' {
			val := (s[i+1]-'0')*64 + (s[i+2]-'0')*8 + (s[i+3] - '0')
			b.WriteByte(val)
			i += 3
		} else {
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

// imageExtensions maps lowercase extensions to MIME types for preview support
var imageExtensions = map[string]string{
	".jpg": "image/jpeg", ".jpeg": "image/jpeg",
	".png": "image/png", ".gif": "image/gif",
	".webp": "image/webp", ".bmp": "image/bmp",
}

// GetImagePreview returns a base64 data URL thumbnail for the given image file.
// maxSize is the maximum width/height in pixels (aspect ratio preserved).
// Returns empty string for non-image files or files larger than 20MB.
func (s *ProjectService) GetImagePreview(filePath string, maxSize int) (string, error) {
	if maxSize <= 0 {
		maxSize = 200
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	_, ok := imageExtensions[ext]
	if !ok {
		return "", nil
	}

	// Skip files larger than 20MB
	info, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}
	if info.Size() > 20*1024*1024 {
		return "", nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var src image.Image
	switch ext {
	case ".jpg", ".jpeg":
		src, err = jpeg.Decode(f)
	case ".png":
		src, err = png.Decode(f)
	case ".gif":
		src, err = gif.Decode(f)
	case ".webp":
		src, err = webp.Decode(f)
	case ".bmp":
		src, _, err = image.Decode(f)
	default:
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("decode %s: %w", ext, err)
	}

	// Calculate thumbnail dimensions preserving aspect ratio
	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w <= maxSize && h <= maxSize {
		// Image is already small enough, encode directly
	} else if w > h {
		h = h * maxSize / w
		w = maxSize
	} else {
		w = w * maxSize / h
		h = maxSize
	}

	// Resize using bilinear interpolation
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)

	// Encode as JPEG for smaller size
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 80}); err != nil {
		return "", err
	}

	dataURL := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
	return dataURL, nil
}

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}
