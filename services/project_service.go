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
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"golang.org/x/image/draw"
	"golang.org/x/image/webp"

	"github.com/rwcarlsen/goexif/exif"

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

// OpenWithDefault opens a file or directory with the default application via xdg-open
func (s *ProjectService) OpenWithDefault(filePath string) error {
	cmd := exec.Command("xdg-open", filePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return cmd.Start()
}

// FileProperties contains detailed file metadata
type FileProperties struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Size        int64  `json:"size"`
	IsDir       bool   `json:"isDir"`
	Permissions string `json:"permissions"`
	Owner       string `json:"owner"`
	Group       string `json:"group"`
	ModTime     string `json:"modTime"`
	AccessTime  string `json:"accessTime"`
	MimeType    string `json:"mimeType"`
	// Image-specific
	ImageWidth  int `json:"imageWidth,omitempty"`
	ImageHeight int `json:"imageHeight,omitempty"`
	// Directory-specific
	ItemCount int `json:"itemCount,omitempty"`
	// EXIF metadata (photos)
	CameraModel  string `json:"cameraModel,omitempty"`
	CameraMake   string `json:"cameraMake,omitempty"`
	FNumber      string `json:"fNumber,omitempty"`
	ExposureTime string `json:"exposureTime,omitempty"`
	ISOSpeed     string `json:"isoSpeed,omitempty"`
	FocalLength  string `json:"focalLength,omitempty"`
	FocalLength35 string `json:"focalLength35,omitempty"`
	Flash        string `json:"flash,omitempty"`
	DateTaken    string `json:"dateTaken,omitempty"`
	Orientation  string `json:"orientation,omitempty"`
	Software     string `json:"software,omitempty"`
	// Video metadata (via ffprobe)
	VideoCodec    string `json:"videoCodec,omitempty"`
	AudioCodec    string `json:"audioCodec,omitempty"`
	Duration      string `json:"duration,omitempty"`
	VideoBitrate  string `json:"videoBitrate,omitempty"`
	AudioBitrate  string `json:"audioBitrate,omitempty"`
	FrameRate     string `json:"frameRate,omitempty"`
	SampleRate    string `json:"sampleRate,omitempty"`
	AudioChannels string `json:"audioChannels,omitempty"`
}

// GetFileProperties returns detailed metadata for a file or directory
func (s *ProjectService) GetFileProperties(filePath string) (*FileProperties, error) {
	info, err := os.Lstat(filePath)
	if err != nil {
		return nil, err
	}

	props := &FileProperties{
		Name:        info.Name(),
		Path:        filePath,
		Size:        info.Size(),
		IsDir:       info.IsDir(),
		Permissions: info.Mode().Perm().String(),
		ModTime:     info.ModTime().Format(time.RFC3339),
	}

	// Get owner/group from syscall stat
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		props.Owner = fmt.Sprintf("%d", stat.Uid)
		props.Group = fmt.Sprintf("%d", stat.Gid)
		props.AccessTime = time.Unix(stat.Atim.Sec, stat.Atim.Nsec).Format(time.RFC3339)
	}

	if info.IsDir() {
		// Count items in directory
		entries, err := os.ReadDir(filePath)
		if err == nil {
			props.ItemCount = len(entries)
		}
		// Calculate total directory size
		size, _ := dirSize(filePath)
		props.Size = size
	} else {
		// Detect MIME type via `file --mime-type`
		props.MimeType = detectMimeType(filePath)

		// Get image dimensions if applicable
		ext := strings.ToLower(filepath.Ext(filePath))
		if _, ok := imageExtensions[ext]; ok {
			w, h := getImageDimensions(filePath, ext)
			props.ImageWidth = w
			props.ImageHeight = h
		}

		// Extract EXIF metadata for photos
		if ext == ".jpg" || ext == ".jpeg" || ext == ".tiff" || ext == ".tif" {
			fillExifData(filePath, props)
		}

		// Extract video/audio metadata via ffprobe
		if strings.HasPrefix(props.MimeType, "video/") || strings.HasPrefix(props.MimeType, "audio/") {
			fillMediaMetadata(filePath, props)
		}
	}

	return props, nil
}

// detectMimeType uses the `file` command to get MIME type
func detectMimeType(path string) string {
	out, err := exec.Command("file", "--mime-type", "-b", path).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// getImageDimensions returns width and height without fully decoding the image
func getImageDimensions(path string, ext string) (int, int) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0
	}
	defer f.Close()

	var cfg image.Config
	switch ext {
	case ".jpg", ".jpeg":
		cfg, err = jpeg.DecodeConfig(f)
	case ".png":
		cfg, err = png.DecodeConfig(f)
	case ".gif":
		cfg, err = gif.DecodeConfig(f)
	case ".webp":
		cfg, err = webp.DecodeConfig(f)
	default:
		cfg, _, err = image.DecodeConfig(f)
	}
	if err != nil {
		return 0, 0
	}
	return cfg.Width, cfg.Height
}

// fillExifData reads EXIF metadata from a JPEG/TIFF file and populates FileProperties
func fillExifData(filePath string, props *FileProperties) {
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		return
	}

	if tag, err := x.Get(exif.Model); err == nil {
		props.CameraModel = strings.Trim(tag.String(), "\"")
	}
	if tag, err := x.Get(exif.Make); err == nil {
		props.CameraMake = strings.Trim(tag.String(), "\"")
	}
	if tag, err := x.Get(exif.FNumber); err == nil {
		if num, den, err := tag.Rat2(0); err == nil && den != 0 {
			val := float64(num) / float64(den)
			props.FNumber = fmt.Sprintf("f/%.1f", val)
		}
	}
	if tag, err := x.Get(exif.ExposureTime); err == nil {
		if num, den, err := tag.Rat2(0); err == nil && den != 0 {
			val := float64(num) / float64(den)
			if val < 1 {
				// Show as fraction: 1/xxx s
				props.ExposureTime = fmt.Sprintf("1/%d s", int(0.5+1.0/val))
			} else {
				props.ExposureTime = fmt.Sprintf("%.4g s", val)
			}
		}
	}
	if tag, err := x.Get(exif.ISOSpeedRatings); err == nil {
		props.ISOSpeed = tag.String()
	}
	if tag, err := x.Get(exif.FocalLength); err == nil {
		if num, den, err := tag.Rat2(0); err == nil && den != 0 {
			val := float64(num) / float64(den)
			props.FocalLength = fmt.Sprintf("%.4g mm", val)
		}
	}
	if tag, err := x.Get(exif.FocalLengthIn35mmFilm); err == nil {
		props.FocalLength35 = tag.String() + " mm"
	}
	if tag, err := x.Get(exif.Flash); err == nil {
		props.Flash = decodeFlash(tag.String())
	}
	if tag, err := x.Get(exif.Software); err == nil {
		props.Software = strings.Trim(tag.String(), "\"")
	}
	if tag, err := x.Get(exif.Orientation); err == nil {
		props.Orientation = decodeOrientation(tag.String())
	}
	if t, err := x.DateTime(); err == nil {
		props.DateTaken = t.Format(time.RFC3339)
	}
}

// decodeFlash converts EXIF flash value to human-readable string
func decodeFlash(val string) string {
	val = strings.TrimSpace(val)
	switch val {
	case "0":
		return "No flash"
	case "1":
		return "Flash fired"
	case "5":
		return "Flash fired, no strobe return"
	case "7":
		return "Flash fired, strobe return"
	case "8":
		return "No flash, compulsory"
	case "9":
		return "Flash fired, compulsory"
	case "16":
		return "No flash (off)"
	case "24":
		return "No flash, auto"
	case "25":
		return "Flash fired, auto"
	default:
		return val
	}
}

// decodeOrientation converts EXIF orientation value to human-readable string
func decodeOrientation(val string) string {
	val = strings.TrimSpace(val)
	switch val {
	case "1":
		return "Normal"
	case "2":
		return "Flipped horizontal"
	case "3":
		return "Rotated 180°"
	case "4":
		return "Flipped vertical"
	case "5":
		return "Transposed"
	case "6":
		return "Rotated 90° CW"
	case "7":
		return "Transverse"
	case "8":
		return "Rotated 90° CCW"
	default:
		return val
	}
}

// ffprobeStream represents a stream from ffprobe JSON output
type ffprobeStream struct {
	CodecName    string `json:"codec_name"`
	CodecType    string `json:"codec_type"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	RFrameRate   string `json:"r_frame_rate"`
	BitRate      string `json:"bit_rate"`
	SampleRate   string `json:"sample_rate"`
	Channels     int    `json:"channels"`
}

type ffprobeFormat struct {
	Duration string `json:"duration"`
	BitRate  string `json:"bit_rate"`
}

type ffprobeOutput struct {
	Streams []ffprobeStream `json:"streams"`
	Format  ffprobeFormat   `json:"format"`
}

// fillMediaMetadata runs ffprobe to extract video/audio metadata
func fillMediaMetadata(filePath string, props *FileProperties) {
	out, err := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_streams",
		"-show_format",
		filePath,
	).Output()
	if err != nil {
		return
	}

	var probe ffprobeOutput
	if err := json.Unmarshal(out, &probe); err != nil {
		return
	}

	// Parse duration
	if probe.Format.Duration != "" {
		if secs := parseDurationSeconds(probe.Format.Duration); secs > 0 {
			h := int(secs) / 3600
			m := (int(secs) % 3600) / 60
			s := int(secs) % 60
			if h > 0 {
				props.Duration = fmt.Sprintf("%d:%02d:%02d", h, m, s)
			} else {
				props.Duration = fmt.Sprintf("%d:%02d", m, s)
			}
		}
	}

	for _, stream := range probe.Streams {
		switch stream.CodecType {
		case "video":
			if props.VideoCodec == "" {
				props.VideoCodec = stream.CodecName
				if stream.Width > 0 && stream.Height > 0 {
					props.ImageWidth = stream.Width
					props.ImageHeight = stream.Height
				}
				if stream.RFrameRate != "" {
					props.FrameRate = parseFrameRate(stream.RFrameRate)
				}
				if stream.BitRate != "" {
					props.VideoBitrate = formatBitrate(stream.BitRate)
				}
			}
		case "audio":
			if props.AudioCodec == "" {
				props.AudioCodec = stream.CodecName
				if stream.SampleRate != "" {
					props.SampleRate = stream.SampleRate + " Hz"
				}
				if stream.Channels > 0 {
					props.AudioChannels = fmt.Sprintf("%d", stream.Channels)
				}
				if stream.BitRate != "" {
					props.AudioBitrate = formatBitrate(stream.BitRate)
				}
			}
		}
	}
}

// parseDurationSeconds parses a string like "123.456" to float64 seconds
func parseDurationSeconds(s string) float64 {
	var v float64
	fmt.Sscanf(s, "%f", &v)
	return v
}

// parseFrameRate converts "30000/1001" to "29.97 fps" or "25/1" to "25 fps"
func parseFrameRate(s string) string {
	parts := strings.Split(s, "/")
	if len(parts) == 2 {
		var num, den float64
		fmt.Sscanf(parts[0], "%f", &num)
		fmt.Sscanf(parts[1], "%f", &den)
		if den > 0 {
			fps := num / den
			if fps == float64(int(fps)) {
				return fmt.Sprintf("%d fps", int(fps))
			}
			return fmt.Sprintf("%.2f fps", fps)
		}
	}
	return s + " fps"
}

// formatBitrate converts "1234567" (bps) to "1.23 Mbps" or "128000" to "128 kbps"
func formatBitrate(s string) string {
	var bps float64
	fmt.Sscanf(s, "%f", &bps)
	if bps <= 0 {
		return s
	}
	if bps >= 1_000_000 {
		return fmt.Sprintf("%.2f Mbps", bps/1_000_000)
	}
	return fmt.Sprintf("%.0f kbps", bps/1000)
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
