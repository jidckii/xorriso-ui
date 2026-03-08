package services

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"xorriso-ui/pkg/models"
)

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

// GetHomeDirectory returns the current user's home directory
func (s *ProjectService) GetHomeDirectory() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "/"
	}
	return home
}

// OpenWithDefault opens a file or directory with the default application via xdg-open
func (s *ProjectService) OpenWithDefault(filePath string) error {
	cmd := exec.Command("xdg-open", filePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return cmd.Start()
}
