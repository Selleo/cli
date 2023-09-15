package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ADR struct {
	Title string
	Date  string
	ID    string
}

func (adr *ADR) Dir() string {
	return filepath.Join("docs", "adr")
}

func (adr *ADR) CreateDir() error {
	err := os.MkdirAll(adr.Dir(), 0755)
	return err
}

func (adr *ADR) FindAll() ([]int, error) {
	entries, err := os.ReadDir(adr.Dir())
	if err != nil {
		return nil, err
	}

	var ids []int
	pattern := regexp.MustCompile(`^(?P<id>\d{4})-(?P<title>.*?)\.(?:(md|MD))$`)
	for _, entry := range entries {
		fmt.Println(entry.Name())
		if pattern.MatchString(entry.Name()) {
			matches := pattern.FindStringSubmatch(entry.Name())
			num := matches[1]
			id, err := strconv.Atoi(num)
			if err != nil {
				return nil, err
			}
			ids = append(ids, id)
		}
	}
	sort.Ints(ids)

	return ids, nil
}

func (adr *ADR) CreateNew(filename string) (*os.File, error) {
	f, err := os.Create(fmt.Sprintf("%s/%s", adr.Dir(), filename))
	if err != nil {
		return nil, err
	}
	_ = f.Chmod(0755)
	return f, nil
}

func (adr *ADR) NextID() (int, error) {
	ids, err := adr.FindAll()
	if err != nil {
		return -1, err
	}

	if len(ids) == 0 {
		return 1, nil
	}

	return ids[len(ids)-1] + 1, nil
}

func (adr *ADR) Render(t *TemplateRenderer) error {
	err := adr.CreateDir()
	if err != nil {
		return err
	}
	id, err := adr.NextID()
	if err != nil {
		return err
	}

	// needed to render the template
	adr.Date = time.Now().Format("2006-01-02")
	adr.ID = fmt.Sprintf("%04d", id)

	title := strings.ToLower(adr.Title)
	title = strings.ReplaceAll(title, " ", "-")
	filename := fmt.Sprintf("%s-%s.md", adr.ID, adr.Title)

	f, err := adr.CreateNew(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.Render(f, "templates/adr/new.md", adr)

	return err
}
