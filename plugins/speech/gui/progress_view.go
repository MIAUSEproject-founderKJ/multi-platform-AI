//MIAUSEproject-founderKJ/multi-platform-AI/plugins/speech/gui/progress_view.go

package gui

// Conceptual layout for a "Reflective" boot screen
func (v *ProgressView) Update(update hmi.ProgressUpdate) {
    // 1. Update Progress Bar
    v.ProgressBar.SetValue(update.percentage)
    
    // 2. Reflect specific hardware status
    v.StatusText.SetText(update.message)
    
    // 3. Visual "Pulse"
    if update.critical {
        v.Container.SetBackgroundColor(colors.WarningRed)
    } else {
        v.Container.SetBackgroundColor(colors.DeepSpaceBlue)
    }
}