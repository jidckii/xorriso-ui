package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"

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
