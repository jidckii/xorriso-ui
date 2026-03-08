package services

import (
	"encoding/json"
	"os"
	"path/filepath"
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
		Version:  1,
		Name:     name,
		VolumeID: volumeID,
		Entries:  []models.FileEntry{},
		ISOOptions: models.ISOOptions{
			UDF:       true,
			ISOLevel:  3,
			RockRidge: true,
			Joliet:    true,
			MD5:       true,
		},
		BurnOptions: models.BurnOptions{
			Speed:      "auto",
			Verify:     true,
			Eject:      true,
			BurnMode:   "auto",
			Padding:    300,
			CleanupISO: true,
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

// AddFiles adds files/directories to the project.
// Directories are added recursively — each file inside gets its own entry.
func (s *ProjectService) AddFiles(project *models.Project, sourcePaths []string, destDir string) (*models.Project, error) {
	for _, src := range sourcePaths {
		info, err := os.Stat(src)
		if err != nil {
			continue
		}

		if info.IsDir() {
			// Add directory entry itself
			size, _ := dirSize(src)
			dirEntry := models.FileEntry{
				SourcePath: src,
				DestPath:   filepath.Join(destDir, filepath.Base(src)),
				Name:       filepath.Base(src),
				IsDir:      true,
				Size:       size,
				ModTime:    info.ModTime().UnixMilli(),
			}
			project.Entries = append(project.Entries, dirEntry)

			// Recursively add all files inside
			baseDest := filepath.Join(destDir, filepath.Base(src))
			_ = filepath.Walk(src, func(path string, fi os.FileInfo, err error) error {
				if err != nil || path == src {
					return nil
				}
				rel, _ := filepath.Rel(src, path)
				project.Entries = append(project.Entries, models.FileEntry{
					SourcePath: path,
					DestPath:   filepath.Join(baseDest, rel),
					Name:       fi.Name(),
					IsDir:      fi.IsDir(),
					Size:       fi.Size(),
					ModTime:    fi.ModTime().UnixMilli(),
				})
				return nil
			})
		} else {
			project.Entries = append(project.Entries, models.FileEntry{
				SourcePath: src,
				DestPath:   filepath.Join(destDir, filepath.Base(src)),
				Name:       filepath.Base(src),
				IsDir:      false,
				Size:       info.Size(),
				ModTime:    info.ModTime().UnixMilli(),
			})
		}
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

// RemoveEntries удаляет несколько записей из проекта по списку destPaths
func (s *ProjectService) RemoveEntries(project *models.Project, destPaths []string) (*models.Project, error) {
	toRemove := make(map[string]struct{}, len(destPaths))
	for _, p := range destPaths {
		toRemove[p] = struct{}{}
	}

	filtered := make([]models.FileEntry, 0, len(project.Entries))
	for _, e := range project.Entries {
		if _, ok := toRemove[e.DestPath]; !ok {
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
