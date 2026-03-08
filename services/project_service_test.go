package services

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/jpeg"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"xorriso-ui/pkg/models"
)

// --- Helper function tests ---

func TestDecodeFlash(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0", "No flash"},
		{"1", "Flash fired"},
		{"16", "No flash (off)"},
		{"24", "No flash, auto"},
		{"25", "Flash fired, auto"},
		{"99", "99"}, // unknown value passed through
	}
	for _, tt := range tests {
		got := decodeFlash(tt.input)
		if got != tt.expected {
			t.Errorf("decodeFlash(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestDecodeOrientation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1", "Normal"},
		{"3", "Rotated 180°"},
		{"6", "Rotated 90° CW"},
		{"8", "Rotated 90° CCW"},
		{"42", "42"}, // unknown value passed through
	}
	for _, tt := range tests {
		got := decodeOrientation(tt.input)
		if got != tt.expected {
			t.Errorf("decodeOrientation(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestParseFrameRate(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"30/1", "30 fps"},
		{"25/1", "25 fps"},
		{"30000/1001", "29.97 fps"},
		{"24000/1001", "23.98 fps"},
		{"60/1", "60 fps"},
	}
	for _, tt := range tests {
		got := parseFrameRate(tt.input)
		if got != tt.expected {
			t.Errorf("parseFrameRate(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestFormatBitrate(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"128000", "128 kbps"},
		{"1500000", "1.50 Mbps"},
		{"5000000", "5.00 Mbps"},
		{"320000", "320 kbps"},
	}
	for _, tt := range tests {
		got := formatBitrate(tt.input)
		if got != tt.expected {
			t.Errorf("formatBitrate(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestParseDurationSeconds(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"123.456", 123.456},
		{"0.000000", 0},
		{"3661.5", 3661.5},
	}
	for _, tt := range tests {
		got := parseDurationSeconds(tt.input)
		if got != tt.expected {
			t.Errorf("parseDurationSeconds(%q) = %f, want %f", tt.input, got, tt.expected)
		}
	}
}

// --- EXIF test with synthetic JPEG ---

// createMinimalJPEGWithExif creates a JPEG file with a minimal EXIF APP1 segment
// containing Make and Model tags.
func createMinimalJPEGWithExif(path string, make_, model string) error {
	// Create a small image first
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))

	// Encode to JPEG
	var imgBuf bytes.Buffer
	if err := jpeg.Encode(&imgBuf, img, &jpeg.Options{Quality: 50}); err != nil {
		return err
	}
	imgBytes := imgBuf.Bytes()

	// Build a minimal EXIF APP1 segment
	// Structure: FFD8 (SOI) + FFE1 (APP1) + length + "Exif\0\0" + TIFF header + IFD
	exifPayload := buildMinimalExif(make_, model)

	var out bytes.Buffer
	// SOI marker
	out.Write([]byte{0xFF, 0xD8})
	// APP1 marker
	out.Write([]byte{0xFF, 0xE1})
	// Length (includes the 2 length bytes themselves)
	length := uint16(len(exifPayload) + 2)
	binary.Write(&out, binary.BigEndian, length)
	out.Write(exifPayload)
	// Append rest of JPEG (skip SOI)
	out.Write(imgBytes[2:])

	return os.WriteFile(path, out.Bytes(), 0644)
}

func buildMinimalExif(make_, model string) []byte {
	var buf bytes.Buffer

	// Exif header
	buf.WriteString("Exif\x00\x00")

	// TIFF header (little-endian)
	tiffStart := buf.Len()
	buf.Write([]byte("II"))                             // Little-endian
	binary.Write(&buf, binary.LittleEndian, uint16(42)) // Magic
	binary.Write(&buf, binary.LittleEndian, uint32(8))  // Offset to first IFD (from TIFF start)

	// IFD0 with 2 entries: Make (0x010F) and Model (0x0110)
	numEntries := uint16(2)
	binary.Write(&buf, binary.LittleEndian, numEntries)

	// Calculate data area offset (from TIFF start)
	// IFD: 2 (count) + 2*12 (entries) + 4 (next IFD offset) = 30
	dataOffset := uint32(8 + 2 + 2*12 + 4)

	makeBytes := []byte(make_ + "\x00")
	modelBytes := []byte(model + "\x00")

	// Entry: Make (tag=0x010F, type=ASCII=2)
	binary.Write(&buf, binary.LittleEndian, uint16(0x010F))
	binary.Write(&buf, binary.LittleEndian, uint16(2)) // ASCII
	binary.Write(&buf, binary.LittleEndian, uint32(len(makeBytes)))
	if len(makeBytes) <= 4 {
		var val [4]byte
		copy(val[:], makeBytes)
		buf.Write(val[:])
	} else {
		binary.Write(&buf, binary.LittleEndian, dataOffset)
	}

	makeDataOffset := dataOffset
	if len(makeBytes) > 4 {
		dataOffset += uint32(len(makeBytes))
	}

	// Entry: Model (tag=0x0110, type=ASCII=2)
	binary.Write(&buf, binary.LittleEndian, uint16(0x0110))
	binary.Write(&buf, binary.LittleEndian, uint16(2)) // ASCII
	binary.Write(&buf, binary.LittleEndian, uint32(len(modelBytes)))
	if len(modelBytes) <= 4 {
		var val [4]byte
		copy(val[:], modelBytes)
		buf.Write(val[:])
	} else {
		binary.Write(&buf, binary.LittleEndian, dataOffset)
	}

	// Next IFD offset = 0 (no more IFDs)
	binary.Write(&buf, binary.LittleEndian, uint32(0))

	// Data area
	if len(makeBytes) > 4 {
		buf.Write(makeBytes)
	}
	if len(modelBytes) > 4 {
		buf.Write(modelBytes)
	}

	_ = tiffStart
	_ = makeDataOffset
	return buf.Bytes()
}

func TestFillExifData_WithExif(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test_camera.jpg")

	err := createMinimalJPEGWithExif(path, "TestMake", "TestModel123")
	if err != nil {
		t.Fatalf("Failed to create test JPEG: %v", err)
	}

	props := &FileProperties{}
	fillExifData(path, props)

	if props.CameraMake != "TestMake" {
		t.Errorf("CameraMake = %q, want %q", props.CameraMake, "TestMake")
	}
	if props.CameraModel != "TestModel123" {
		t.Errorf("CameraModel = %q, want %q", props.CameraModel, "TestModel123")
	}
}

func TestFillExifData_NoExif(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test_plain.jpg")

	// Create plain JPEG without EXIF
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	jpeg.Encode(f, img, nil)
	f.Close()

	props := &FileProperties{}
	fillExifData(path, props)

	// Should not crash and fields should be empty
	if props.CameraMake != "" {
		t.Errorf("CameraMake should be empty for plain JPEG, got %q", props.CameraMake)
	}
	if props.CameraModel != "" {
		t.Errorf("CameraModel should be empty for plain JPEG, got %q", props.CameraModel)
	}
}

func TestFillExifData_NonexistentFile(t *testing.T) {
	props := &FileProperties{}
	fillExifData("/nonexistent/file.jpg", props)
	// Should not panic
	if props.CameraMake != "" {
		t.Error("Expected empty CameraMake for nonexistent file")
	}
}

// --- GetFileProperties integration test ---

func TestGetFileProperties_RegularFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("hello world"), 0644)

	svc := NewProjectService()
	props, err := svc.GetFileProperties(path)
	if err != nil {
		t.Fatal(err)
	}

	if props.Name != "test.txt" {
		t.Errorf("Name = %q, want %q", props.Name, "test.txt")
	}
	if props.Size != 11 {
		t.Errorf("Size = %d, want 11", props.Size)
	}
	if props.IsDir {
		t.Error("IsDir should be false")
	}
	if props.Permissions == "" {
		t.Error("Permissions should not be empty")
	}
	if props.ModTime == "" {
		t.Error("ModTime should not be empty")
	}
}

func TestGetFileProperties_Directory(t *testing.T) {
	dir := t.TempDir()
	// Create some files inside
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(dir, "b.txt"), []byte("bb"), 0644)

	svc := NewProjectService()
	props, err := svc.GetFileProperties(dir)
	if err != nil {
		t.Fatal(err)
	}

	if !props.IsDir {
		t.Error("IsDir should be true")
	}
	if props.ItemCount != 2 {
		t.Errorf("ItemCount = %d, want 2", props.ItemCount)
	}
}

func TestGetFileProperties_ImageWithDimensions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.jpg")

	img := image.NewRGBA(image.Rect(0, 0, 100, 50))
	f, _ := os.Create(path)
	jpeg.Encode(f, img, nil)
	f.Close()

	svc := NewProjectService()
	props, err := svc.GetFileProperties(path)
	if err != nil {
		t.Fatal(err)
	}

	if props.ImageWidth != 100 {
		t.Errorf("ImageWidth = %d, want 100", props.ImageWidth)
	}
	if props.ImageHeight != 50 {
		t.Errorf("ImageHeight = %d, want 50", props.ImageHeight)
	}
}

// --- Video metadata test (skipped if ffprobe not available) ---

func TestFillMediaMetadata_WithFFprobe(t *testing.T) {
	// Skip if ffprobe is not installed
	if _, err := exec.LookPath("ffprobe"); err != nil {
		t.Skip("ffprobe not found, skipping video metadata test")
	}

	// Skip if ffmpeg is not installed (needed to create test video)
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not found, skipping video metadata test")
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "test.mp4")

	// Create a minimal test video with ffmpeg
	err := exec.Command("ffmpeg",
		"-f", "lavfi", "-i", "testsrc=duration=2:size=320x240:rate=25",
		"-f", "lavfi", "-i", "sine=frequency=440:duration=2",
		"-c:v", "libx264", "-c:a", "aac",
		"-y", path,
	).Run()
	if err != nil {
		t.Skipf("Failed to create test video: %v", err)
	}

	props := &FileProperties{}
	fillMediaMetadata(path, props)

	if props.VideoCodec == "" {
		t.Error("VideoCodec should not be empty")
	}
	if props.AudioCodec == "" {
		t.Error("AudioCodec should not be empty")
	}
	if props.Duration == "" {
		t.Error("Duration should not be empty")
	}
	if props.ImageWidth != 320 {
		t.Errorf("ImageWidth = %d, want 320", props.ImageWidth)
	}
	if props.ImageHeight != 240 {
		t.Errorf("ImageHeight = %d, want 240", props.ImageHeight)
	}
	if props.FrameRate == "" {
		t.Error("FrameRate should not be empty")
	}

	// Check that video codec is h264
	if !strings.Contains(props.VideoCodec, "264") && props.VideoCodec != "h264" {
		t.Logf("VideoCodec = %q (expected h264)", props.VideoCodec)
	}
}

func TestFillMediaMetadata_NoFFprobe(t *testing.T) {
	// Test with nonexistent file - should not panic
	props := &FileProperties{}
	fillMediaMetadata("/nonexistent/video.mp4", props)

	if props.VideoCodec != "" {
		t.Error("VideoCodec should be empty for nonexistent file")
	}
}

// --- detectMimeType test ---

func TestDetectMimeType(t *testing.T) {
	dir := t.TempDir()

	// Create a plain text file
	txtPath := filepath.Join(dir, "test.txt")
	os.WriteFile(txtPath, []byte("hello"), 0644)
	mime := detectMimeType(txtPath)
	if !strings.Contains(mime, "text/plain") {
		t.Errorf("detectMimeType(txt) = %q, want text/plain", mime)
	}

	// Create a JPEG file
	jpgPath := filepath.Join(dir, "test.jpg")
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	f, _ := os.Create(jpgPath)
	jpeg.Encode(f, img, nil)
	f.Close()
	mime = detectMimeType(jpgPath)
	if !strings.Contains(mime, "image/jpeg") {
		t.Errorf("detectMimeType(jpg) = %q, want image/jpeg", mime)
	}
}

// --- BrowseDirectory test ---

func TestBrowseDirectory(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0644)
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644)
	os.Mkdir(filepath.Join(dir, "zdir"), 0755)

	svc := NewProjectService()
	entries, err := svc.BrowseDirectory(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 3 {
		t.Fatalf("Expected 3 entries, got %d", len(entries))
	}

	// Directories should come first
	if !entries[0].IsDir {
		t.Error("First entry should be a directory")
	}
	if entries[0].Name != "zdir" {
		t.Errorf("First entry name = %q, want %q", entries[0].Name, "zdir")
	}

	// Files sorted alphabetically
	if entries[1].Name != "a.txt" {
		t.Errorf("Second entry name = %q, want %q", entries[1].Name, "a.txt")
	}
	if entries[2].Name != "b.txt" {
		t.Errorf("Third entry name = %q, want %q", entries[2].Name, "b.txt")
	}
}

// --- UnescapeMountPath test ---

// --- Project CRUD тесты ---

func TestNewProject(t *testing.T) {
	svc := NewProjectService()
	project := svc.NewProject("Test Project", "TEST_VOL")

	if project.Name != "Test Project" {
		t.Errorf("Name = %q, want %q", project.Name, "Test Project")
	}
	if project.VolumeID != "TEST_VOL" {
		t.Errorf("VolumeID = %q, want %q", project.VolumeID, "TEST_VOL")
	}
	if len(project.Entries) != 0 {
		t.Errorf("Entries should be empty, got %d", len(project.Entries))
	}
	// Проверяем defaults ISOOptions
	if !project.ISOOptions.UDF {
		t.Error("ISOOptions.UDF should be true by default")
	}
	if project.ISOOptions.ISOLevel != 3 {
		t.Errorf("ISOLevel = %d, want 3", project.ISOOptions.ISOLevel)
	}
	if !project.ISOOptions.RockRidge {
		t.Error("RockRidge should be true by default")
	}
	if !project.ISOOptions.Joliet {
		t.Error("Joliet should be true by default")
	}
	if !project.ISOOptions.MD5 {
		t.Error("MD5 should be true by default")
	}
	// Проверяем defaults BurnOptions
	if project.BurnOptions.Speed != "auto" {
		t.Errorf("Speed = %q, want auto", project.BurnOptions.Speed)
	}
	if !project.BurnOptions.Verify {
		t.Error("Verify should be true by default")
	}
	if !project.BurnOptions.Eject {
		t.Error("Eject should be true by default")
	}
	if project.BurnOptions.Padding != 300 {
		t.Errorf("Padding = %d, want 300", project.BurnOptions.Padding)
	}
	if project.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestSaveProject_EmptyPath(t *testing.T) {
	svc := NewProjectService()
	project := svc.NewProject("Test", "VOL")
	// FilePath пустой -- SaveProject не должен падать
	err := svc.SaveProject(project)
	if err != nil {
		t.Fatalf("SaveProject with empty path should return nil, got: %v", err)
	}
}

func TestSaveAndOpenProject(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.xorriso-project")

	svc := NewProjectService()
	project := svc.NewProject("Roundtrip Test", "RT_VOL")
	project.Entries = []models.FileEntry{
		{SourcePath: "/tmp/file.txt", DestPath: "/file.txt", Name: "file.txt", Size: 1024},
	}
	project.ISOOptions.HFSPlus = true

	// Сохраняем
	err := svc.SaveProjectAs(project, filePath)
	if err != nil {
		t.Fatalf("SaveProjectAs: %v", err)
	}

	// Загружаем обратно
	loaded, err := svc.OpenProject(filePath)
	if err != nil {
		t.Fatalf("OpenProject: %v", err)
	}

	if loaded.Name != "Roundtrip Test" {
		t.Errorf("Name = %q, want %q", loaded.Name, "Roundtrip Test")
	}
	if loaded.VolumeID != "RT_VOL" {
		t.Errorf("VolumeID = %q, want %q", loaded.VolumeID, "RT_VOL")
	}
	if loaded.FilePath != filePath {
		t.Errorf("FilePath = %q, want %q", loaded.FilePath, filePath)
	}
	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].SourcePath != "/tmp/file.txt" {
		t.Errorf("Entry SourcePath = %q, want /tmp/file.txt", loaded.Entries[0].SourcePath)
	}
	if !loaded.ISOOptions.HFSPlus {
		t.Error("ISOOptions.HFSPlus should be true")
	}
}

func TestAddFiles_SingleFile(t *testing.T) {
	dir := t.TempDir()
	testFile := filepath.Join(dir, "hello.txt")
	os.WriteFile(testFile, []byte("hello world"), 0644)

	svc := NewProjectService()
	project := svc.NewProject("Test", "VOL")

	project, err := svc.AddFiles(project, []string{testFile}, "/")
	if err != nil {
		t.Fatalf("AddFiles: %v", err)
	}

	if len(project.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(project.Entries))
	}

	entry := project.Entries[0]
	if entry.Name != "hello.txt" {
		t.Errorf("Name = %q, want hello.txt", entry.Name)
	}
	if entry.SourcePath != testFile {
		t.Errorf("SourcePath = %q, want %q", entry.SourcePath, testFile)
	}
	if entry.DestPath != "/hello.txt" {
		t.Errorf("DestPath = %q, want /hello.txt", entry.DestPath)
	}
	if entry.IsDir {
		t.Error("IsDir should be false for a file")
	}
	if entry.Size != 11 {
		t.Errorf("Size = %d, want 11", entry.Size)
	}
}

func TestAddFiles_Directory(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "mydir")
	os.MkdirAll(subDir, 0755)
	os.WriteFile(filepath.Join(subDir, "a.txt"), []byte("aaa"), 0644)
	os.WriteFile(filepath.Join(subDir, "b.txt"), []byte("bbb"), 0644)

	svc := NewProjectService()
	project := svc.NewProject("Test", "VOL")

	project, err := svc.AddFiles(project, []string{subDir}, "/")
	if err != nil {
		t.Fatalf("AddFiles: %v", err)
	}

	// Ожидаем: директория + 2 файла внутри = 3 записи
	if len(project.Entries) != 3 {
		t.Fatalf("expected 3 entries (dir + 2 files), got %d", len(project.Entries))
	}

	// Первая запись -- сама директория
	if !project.Entries[0].IsDir {
		t.Error("first entry should be a directory")
	}
	if project.Entries[0].Name != "mydir" {
		t.Errorf("first entry Name = %q, want mydir", project.Entries[0].Name)
	}
}

func TestRemoveEntry(t *testing.T) {
	svc := NewProjectService()
	project := svc.NewProject("Test", "VOL")
	project.Entries = []models.FileEntry{
		{SourcePath: "/a", DestPath: "/a.txt", Name: "a.txt", Size: 10},
		{SourcePath: "/b", DestPath: "/b.txt", Name: "b.txt", Size: 20},
		{SourcePath: "/c", DestPath: "/c.txt", Name: "c.txt", Size: 30},
	}

	project, err := svc.RemoveEntry(project, "/b.txt")
	if err != nil {
		t.Fatalf("RemoveEntry: %v", err)
	}

	if len(project.Entries) != 2 {
		t.Fatalf("expected 2 entries after removal, got %d", len(project.Entries))
	}

	for _, e := range project.Entries {
		if e.DestPath == "/b.txt" {
			t.Error("entry /b.txt should have been removed")
		}
	}
}

func TestCalculateSize(t *testing.T) {
	svc := NewProjectService()
	project := svc.NewProject("Test", "VOL")
	project.Entries = []models.FileEntry{
		{Size: 1000},
		{Size: 2000},
		{Size: 3000},
	}

	total, err := svc.CalculateSize(project)
	if err != nil {
		t.Fatalf("CalculateSize: %v", err)
	}
	if total != 6000 {
		t.Errorf("total = %d, want 6000", total)
	}
}

func TestCalculateSize_Empty(t *testing.T) {
	svc := NewProjectService()
	project := svc.NewProject("Test", "VOL")

	total, err := svc.CalculateSize(project)
	if err != nil {
		t.Fatalf("CalculateSize: %v", err)
	}
	if total != 0 {
		t.Errorf("total = %d, want 0 for empty project", total)
	}
}

func TestUnescapeMountPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/mnt/simple", "/mnt/simple"},
		{"/run/media/user/My\\040Drive", "/run/media/user/My Drive"},
		{"/mnt/a\\134b", "/mnt/a\\b"},           // backslash
		{"/mnt/tab\\011here", "/mnt/tab\there"}, // tab
		{"/mnt/multi\\040word\\040path", "/mnt/multi word path"},
	}
	for _, tt := range tests {
		got := unescapeMountPath(tt.input)
		if got != tt.expected {
			t.Errorf("unescapeMountPath(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
