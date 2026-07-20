package exporter

import (
	"encoding/json"
	"os"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/git"
)

func (e *Exporter) writeJSON(path string, analysis *git.Analysis) error {
	data, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
