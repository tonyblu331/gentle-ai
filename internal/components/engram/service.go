package engram

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gentleman-programming/gentle-ai/internal/backup"
	"github.com/gentleman-programming/gentle-ai/internal/storage"
)

// snapshotter is a subset of backup.Snapshotter used for testability.
type snapshotter interface {
	Create(snapshotDir string, paths []string) (backup.Manifest, error)
}

// DataDirService provides safe, snapshot-guarded operations on the Engram
// data directory. Every mutating operation creates a backup before changing
// anything; the snapshot is visible in the Backups TUI screen.
type DataDirService struct {
	homeDir     string
	backupRoot  string
	snapshotter snapshotter
}

// NewDataDirService creates a DataDirService with the default backup.Snapshotter.
func NewDataDirService(homeDir string) DataDirService {
	return DataDirService{
		homeDir:     homeDir,
		backupRoot:  filepath.Join(homeDir, ".gentle-ai", "backups"),
		snapshotter: backup.NewSnapshotter(),
	}
}

// CopyTo copies all SQLite artifacts from currentDir to dst without removing the source.
// A snapshot of existing source files is created first. Returns the snapshot manifest.
func (s DataDirService) CopyTo(currentDir, dst string) (backup.Manifest, error) {
	paths := SQLiteArtifactPaths(currentDir)
	snap, err := s.snapshot(paths)
	if err != nil {
		return backup.Manifest{}, fmt.Errorf("snapshot before copy: %w", err)
	}
	if err := CopySQLiteArtifacts(currentDir, dst); err != nil {
		return snap, fmt.Errorf("copy: %w", err)
	}
	return snap, nil
}

// MoveTo copies all SQLite artifacts from currentDir to dst, then removes every
// artifact from the source directory. The source is only removed after the copy succeeds.
func (s DataDirService) MoveTo(currentDir, dst string) (backup.Manifest, error) {
	snap, err := s.CopyTo(currentDir, dst)
	if err != nil {
		return snap, err
	}
	if err := RemoveSQLiteArtifacts(currentDir); err != nil {
		return snap, fmt.Errorf("remove source after move: %w", err)
	}
	return snap, nil
}

// Delete creates a snapshot of existing SQLite files, then removes every artifact.
// Returns the snapshot manifest so the caller can show the backup ID.
func (s DataDirService) Delete(dataDir string) (backup.Manifest, error) {
	paths := SQLiteArtifactPaths(dataDir)
	snap, err := s.snapshot(paths)
	if err != nil {
		return backup.Manifest{}, fmt.Errorf("snapshot before delete: %w", err)
	}
	if err := RemoveSQLiteArtifacts(dataDir); err != nil {
		return snap, fmt.Errorf("delete: %w", err)
	}
	return snap, nil
}

// DiskSpaceOK reports whether the volume probed by dstProbePath has enough free
// space to hold a copy of all existing SQLite artifacts under srcDataDir.
// Returns (ok, needed bytes, available bytes, error).
func (s DataDirService) DiskSpaceOK(srcDataDir, dstProbePath string) (bool, int64, int64, error) {
	paths, err := ExistingSQLiteArtifacts(srcDataDir)
	if err != nil {
		return false, 0, 0, err
	}
	if len(paths) == 0 {
		return false, 0, 0, fmt.Errorf("no sqlite artifacts under %q", srcDataDir)
	}
	var needed int64
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			return false, 0, 0, fmt.Errorf("stat source artifact %q: %w", p, err)
		}
		needed += info.Size()
	}
	avail, err := storage.AvailableBytes(dstProbePath)
	if err != nil {
		return false, 0, 0, fmt.Errorf("check available space at %q: %w", dstProbePath, err)
	}
	return avail > needed, needed, avail, nil
}

// snapshot creates a timestamped backup of the given paths under the backup root.
// Missing paths are recorded in the manifest but omitted from the archive (see backup.Snapshotter).
func (s DataDirService) snapshot(paths []string) (backup.Manifest, error) {
	if err := os.MkdirAll(s.backupRoot, 0o755); err != nil {
		return backup.Manifest{}, fmt.Errorf("create backup root %q: %w", s.backupRoot, err)
	}
	snapshotDir := filepath.Join(s.backupRoot, time.Now().UTC().Format("20060102150405.000000000"))
	return s.snapshotter.Create(snapshotDir, paths)
}
