package app

var (
	audio = []string{".mp3", ".wav", ".flac"}
	video = []string{".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv", ".webm"}
)

var formats = append(audio, video...)
