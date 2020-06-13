package metadata

type GeneralData struct {
	Size uint64 `json:"size"`
	CreationTime uint64 `json:"creation_time"`
	ModificationTime uint64 `json:"modification_time"`
	Format string `json:"format"`
	Name string `json:"name"`
}

type fileMetadata struct {
	Url string `json:"url"`
	General GeneralData `json:"general"`
}
