package live_reload

type Events interface {
	// Called when a livereload is triggered
	ReloadTriggered(targets []string)
}
