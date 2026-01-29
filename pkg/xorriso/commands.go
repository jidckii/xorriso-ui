package xorriso

// CommandBuilder constructs safe xorriso command-line arguments
type CommandBuilder struct {
	args []string
}

func NewCommand() *CommandBuilder {
	return &CommandBuilder{}
}

func (b *CommandBuilder) add(args ...string) *CommandBuilder {
	b.args = append(b.args, args...)
	return b
}

// Basic settings
func (b *CommandBuilder) PktOutput() *CommandBuilder       { return b.add("-pkt_output", "on") }
func (b *CommandBuilder) Device(dev string) *CommandBuilder { return b.add("-dev", dev) }
func (b *CommandBuilder) InDevice(dev string) *CommandBuilder { return b.add("-indev", dev) }
func (b *CommandBuilder) OutDevice(dev string) *CommandBuilder { return b.add("-outdev", dev) }

// Information queries
func (b *CommandBuilder) Devices() *CommandBuilder         { return b.add("-device_links") }
func (b *CommandBuilder) TOC() *CommandBuilder              { return b.add("-toc") }
func (b *CommandBuilder) ListFormats() *CommandBuilder      { return b.add("-list_formats") }
func (b *CommandBuilder) ListSpeeds() *CommandBuilder       { return b.add("-list_speeds") }
func (b *CommandBuilder) ListProfiles(which string) *CommandBuilder {
	return b.add("-list_profiles", which)
}
func (b *CommandBuilder) TellMediaSpace() *CommandBuilder   { return b.add("-tell_media_space") }
func (b *CommandBuilder) CheckDrive() *CommandBuilder       { return b.add("-checkdrive") }
func (b *CommandBuilder) PrintSize() *CommandBuilder        { return b.add("-print_size") }
func (b *CommandBuilder) PVDInfo() *CommandBuilder          { return b.add("-pvd_info") }

// ISO options
func (b *CommandBuilder) VolumeID(id string) *CommandBuilder { return b.add("-volid", id) }
func (b *CommandBuilder) RockRidge(on bool) *CommandBuilder {
	if on {
		return b.add("-rockridge", "on")
	}
	return b.add("-rockridge", "off")
}
func (b *CommandBuilder) Joliet(on bool) *CommandBuilder {
	if on {
		return b.add("-joliet", "on")
	}
	return b.add("-joliet", "off")
}
func (b *CommandBuilder) MD5(mode string) *CommandBuilder  { return b.add("-md5", mode) }
func (b *CommandBuilder) ForBackup() *CommandBuilder        { return b.add("-for_backup") }

// File operations
func (b *CommandBuilder) Map(source, dest string) *CommandBuilder {
	return b.add("-map", source, dest)
}
func (b *CommandBuilder) Add(paths ...string) *CommandBuilder {
	args := []string{"-add"}
	args = append(args, paths...)
	args = append(args, "--")
	return b.add(args...)
}

// Write operations
func (b *CommandBuilder) WriteSpeed(speed string) *CommandBuilder {
	return b.add("-speed", speed)
}
func (b *CommandBuilder) Dummy(on bool) *CommandBuilder {
	if on {
		return b.add("-dummy", "on")
	}
	return b.add("-dummy", "off")
}
func (b *CommandBuilder) Close(on bool) *CommandBuilder {
	if on {
		return b.add("-close", "on")
	}
	return b.add("-close", "off")
}
func (b *CommandBuilder) StreamRecording(on bool) *CommandBuilder {
	if on {
		return b.add("-stream_recording", "on")
	}
	return b.add("-stream_recording", "off")
}
func (b *CommandBuilder) Pacifier(format string) *CommandBuilder {
	return b.add("-pacifier", format)
}
func (b *CommandBuilder) Commit() *CommandBuilder { return b.add("-commit") }
func (b *CommandBuilder) Eject(which string) *CommandBuilder {
	return b.add("-eject", which)
}

// Blanking and formatting
func (b *CommandBuilder) Blank(mode string) *CommandBuilder { return b.add("-blank", mode) }
func (b *CommandBuilder) Format(mode string) *CommandBuilder { return b.add("-format", mode) }

// Verification
func (b *CommandBuilder) CheckMedia(opts map[string]string) *CommandBuilder {
	args := []string{"-check_media"}
	for k, v := range opts {
		args = append(args, k+"="+v)
	}
	args = append(args, "--")
	return b.add(args...)
}
func (b *CommandBuilder) Compare(diskPath, isoPath string) *CommandBuilder {
	return b.add("-compare", diskPath, isoPath)
}

// Extraction
func (b *CommandBuilder) OsirroX(mode string) *CommandBuilder { return b.add("-osirrox", mode) }
func (b *CommandBuilder) Extract(isoPath, diskPath string) *CommandBuilder {
	return b.add("-extract", isoPath, diskPath)
}

// Error handling
func (b *CommandBuilder) AbortOn(severity string) *CommandBuilder {
	return b.add("-abort_on", severity)
}
func (b *CommandBuilder) ReportAbout(severity string) *CommandBuilder {
	return b.add("-report_about", severity)
}

// Build returns the final argument slice
func (b *CommandBuilder) Build() []string {
	return b.args
}
