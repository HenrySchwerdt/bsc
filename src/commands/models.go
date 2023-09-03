package commands

type ProjectSettings struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Entry        string            `json:"entry"`
	BuildDir     string            `json:"buildDir"`
	ArtifactName string            `json:"artifactName"`
	Commands     map[string]string `json:"commands"`
	Deps         []string          `json:"deps"`
	License      string            `json:"license"`
}
