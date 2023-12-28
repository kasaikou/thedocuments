package core

type DirectoryArtifact struct {
	Path             string
	DirectoryContent ArtifactObject
	Content          []ArtifactObject
}
