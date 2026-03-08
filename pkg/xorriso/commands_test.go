package xorriso

import (
	"slices"
	"strings"
	"testing"
)

func assertArgs(t *testing.T, got, want []string) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Errorf("аргументы не совпадают\n  got:  %v\n  want: %v", got, want)
	}
}

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()
	if cmd == nil {
		t.Fatal("NewCommand() вернул nil")
	}
	args := cmd.Build()
	if len(args) != 0 {
		t.Errorf("Build() должен вернуть пустой слайс, получили: %v", args)
	}
}

func TestPktOutput(t *testing.T) {
	assertArgs(t, NewCommand().PktOutput().Build(), []string{"-pkt_output", "on"})
}

func TestDevice(t *testing.T) {
	assertArgs(t, NewCommand().Device("/dev/sr0").Build(), []string{"-dev", "/dev/sr0"})
}

func TestDeviceLinks(t *testing.T) {
	assertArgs(t, NewCommand().Devices().Build(), []string{"-device_links"})
}

func TestVolumeID(t *testing.T) {
	assertArgs(t, NewCommand().VolumeID("MY_DISC").Build(), []string{"-volid", "MY_DISC"})
}

func TestRockRidge_On(t *testing.T) {
	assertArgs(t, NewCommand().RockRidge(true).Build(), []string{"-rockridge", "on"})
}

func TestRockRidge_Off(t *testing.T) {
	assertArgs(t, NewCommand().RockRidge(false).Build(), []string{"-rockridge", "off"})
}

func TestJoliet_On(t *testing.T) {
	assertArgs(t, NewCommand().Joliet(true).Build(), []string{"-joliet", "on"})
}

func TestJoliet_Off(t *testing.T) {
	assertArgs(t, NewCommand().Joliet(false).Build(), []string{"-joliet", "off"})
}

func TestHFSPlus_On(t *testing.T) {
	assertArgs(t, NewCommand().HFSPlus(true).Build(), []string{"-hfsplus", "on"})
}

func TestHFSPlus_Off(t *testing.T) {
	assertArgs(t, NewCommand().HFSPlus(false).Build(), []string{"-hfsplus", "off"})
}

func TestZisofs_On(t *testing.T) {
	assertArgs(t, NewCommand().Zisofs(true).Build(), []string{"-zisofs", "by_magic"})
}

func TestZisofs_Off(t *testing.T) {
	assertArgs(t, NewCommand().Zisofs(false).Build(), []string{"-zisofs", "off"})
}

func TestISOLevel(t *testing.T) {
	assertArgs(t, NewCommand().ISOLevel(3).Build(), []string{"-iso_level", "3"})
}

func TestMap(t *testing.T) {
	assertArgs(t, NewCommand().Map("/home/user/data", "/data").Build(),
		[]string{"-map", "/home/user/data", "/data"})
}

func TestCheckMedia_WithOpts(t *testing.T) {
	opts := map[string]string{
		"use":     "outdev",
		"min_lba": "0",
	}
	args := NewCommand().CheckMedia(opts).Build()

	if len(args) < 2 {
		t.Fatalf("слишком мало аргументов: %v", args)
	}
	if args[0] != "-check_media" {
		t.Errorf("первый аргумент должен быть -check_media, получили: %s", args[0])
	}
	if args[len(args)-1] != "--" {
		t.Errorf("последний аргумент должен быть --, получили: %s", args[len(args)-1])
	}

	// Проверяем наличие всех опций в виде key=value
	middle := strings.Join(args[1:len(args)-1], " ")
	for k, v := range opts {
		kv := k + "=" + v
		if !strings.Contains(middle, kv) {
			t.Errorf("опция %s не найдена в аргументах: %v", kv, args)
		}
	}
}

func TestCheckMedia_NilOpts(t *testing.T) {
	assertArgs(t, NewCommand().CheckMedia(nil).Build(), []string{"-check_media", "--"})
}

func TestFluentChain(t *testing.T) {
	args := NewCommand().
		Device("/dev/sr0").
		VolumeID("BACKUP").
		RockRidge(true).
		Map("/home/user/docs", "/docs").
		Commit().
		Build()

	expected := []string{
		"-dev", "/dev/sr0",
		"-volid", "BACKUP",
		"-rockridge", "on",
		"-map", "/home/user/docs", "/docs",
		"-commit",
	}
	assertArgs(t, args, expected)
}

func TestDummy_On(t *testing.T) {
	assertArgs(t, NewCommand().Dummy(true).Build(), []string{"-dummy", "on"})
}

func TestDummy_Off(t *testing.T) {
	assertArgs(t, NewCommand().Dummy(false).Build(), []string{"-dummy", "off"})
}

func TestClose_On(t *testing.T) {
	assertArgs(t, NewCommand().Close(true).Build(), []string{"-close", "on"})
}

func TestClose_Off(t *testing.T) {
	assertArgs(t, NewCommand().Close(false).Build(), []string{"-close", "off"})
}

func TestWriteSpeed(t *testing.T) {
	assertArgs(t, NewCommand().WriteSpeed("4x").Build(), []string{"-speed", "4x"})
}

func TestPadding(t *testing.T) {
	assertArgs(t, NewCommand().Padding(300).Build(), []string{"-padding", "300k"})
}

func TestBlank(t *testing.T) {
	assertArgs(t, NewCommand().Blank("fast").Build(), []string{"-blank", "fast"})
}

func TestFormat(t *testing.T) {
	assertArgs(t, NewCommand().Format("full").Build(), []string{"-format", "full"})
}

func TestStdioOutDevice(t *testing.T) {
	assertArgs(t, NewCommand().StdioOutDevice("/tmp/output.iso").Build(),
		[]string{"-outdev", "stdio:/tmp/output.iso"})
}
