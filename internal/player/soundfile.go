package player

import (
	"os"
	"sort"
	"time"
)

// ----------------------------------------------

// SoundFileList is an arary of *SoundFile
// provoding sort functions.
type SoundFileList []*SoundFile

// SortByDate sorts the lists sound files
// by their last modified date descending.
func (sfl SoundFileList) SortByDate() {
	by(func(sf1, sf2 *SoundFile) bool {
		return sf1.LastModified.After(sf2.LastModified)
	}).Sort(sfl)
}

// SortByName sorts the lists sound files
// by their ident names alphabetically descending.
func (sfl SoundFileList) SortByName() {
	by(func(sf1, sf2 *SoundFile) bool {
		return sf1.Name < sf2.Name
	}).Sort(sfl)
}

// ----------------------------------------------

// SoundFile contains the name (without file extension),
// path and last modified date of a local sound file.
type SoundFile struct {
	Name         string    `json:"name"`
	Path         string    `json:"-"`
	LastModified time.Time `json:"last_modified"`
}

// NewSoundFile creates an instance of SoundFile by
// passed name and path and gets the last modified
// time of the file.
func NewSoundFile(name, path string) (*SoundFile, error) {
	s, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	file := &SoundFile{
		Name:         name,
		Path:         path,
		LastModified: s.ModTime(),
	}

	return file, err
}

// ----------------------------------------------

type soundFileListSorter struct {
	sfl SoundFileList
	by  by
}

func (s *soundFileListSorter) Len() int {
	return len(s.sfl)
}

func (s *soundFileListSorter) Swap(i, j int) {
	s.sfl[i], s.sfl[j] = s.sfl[j], s.sfl[i]
}

func (s *soundFileListSorter) Less(i, j int) bool {
	return s.by(s.sfl[i], s.sfl[j])
}

// ----------------------------------------------

type by func(sf1, sf2 *SoundFile) bool

func (b by) Sort(sfl SoundFileList) {
	ps := &soundFileListSorter{
		sfl: sfl,
		by:  b,
	}
	sort.Sort(ps)
}
