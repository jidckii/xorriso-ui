package udev

import (
	"testing"
)

// buildUeventPacket собирает тестовый netlink пакет из header и KEY=VALUE пар
func buildUeventPacket(header string, kvPairs ...string) []byte {
	var buf []byte
	buf = append(buf, []byte(header)...)
	buf = append(buf, 0)
	for _, kv := range kvPairs {
		buf = append(buf, []byte(kv)...)
		buf = append(buf, 0)
	}
	return buf
}

func TestParseUevent_ChangeEventSr0(t *testing.T) {
	data := buildUeventPacket(
		"change@/devices/pci0000:00/0000:00:1f.2/ata1/host0/target0:0:0/0:0:0:0/block/sr0",
		"ACTION=change",
		"DEVPATH=/devices/pci0000:00/0000:00:1f.2/ata1/host0/target0:0:0/0:0:0:0/block/sr0",
		"SUBSYSTEM=block",
		"DEVNAME=sr0",
		"DEVTYPE=disk",
	)

	ev := ParseUevent(data)
	if ev == nil {
		t.Fatal("ParseUevent вернул nil для валидного пакета")
	}
	if ev.Action != "change" {
		t.Errorf("Action = %q, want %q", ev.Action, "change")
	}
	if ev.DevPath != "/devices/pci0000:00/0000:00:1f.2/ata1/host0/target0:0:0/0:0:0:0/block/sr0" {
		t.Errorf("DevPath = %q, unexpected", ev.DevPath)
	}
	if ev.Subsystem != "block" {
		t.Errorf("Subsystem = %q, want %q", ev.Subsystem, "block")
	}
	if ev.DevName != "sr0" {
		t.Errorf("DevName = %q, want %q", ev.DevName, "sr0")
	}
	if ev.DevType != "disk" {
		t.Errorf("DevType = %q, want %q", ev.DevType, "disk")
	}
}

func TestParseUevent_AddEventSr1(t *testing.T) {
	data := buildUeventPacket(
		"add@/devices/pci0000:00/block/sr1",
		"ACTION=add",
		"SUBSYSTEM=block",
		"DEVNAME=sr1",
		"DEVTYPE=disk",
	)

	ev := ParseUevent(data)
	if ev == nil {
		t.Fatal("ParseUevent вернул nil для add event")
	}
	if ev.Action != "add" {
		t.Errorf("Action = %q, want %q", ev.Action, "add")
	}
	if ev.DevName != "sr1" {
		t.Errorf("DevName = %q, want %q", ev.DevName, "sr1")
	}
}

func TestParseUevent_NoSubsystem(t *testing.T) {
	// Пакет без SUBSYSTEM — должен вернуть Event с пустым Subsystem
	// (фильтрация происходит в Listen, не в ParseUevent)
	data := buildUeventPacket(
		"change@/devices/some/path",
		"DEVNAME=sda",
		"DEVTYPE=disk",
	)

	ev := ParseUevent(data)
	if ev == nil {
		t.Fatal("ParseUevent вернул nil для пакета без SUBSYSTEM")
	}
	if ev.Subsystem != "" {
		t.Errorf("Subsystem = %q, want empty", ev.Subsystem)
	}
	if ev.DevName != "sda" {
		t.Errorf("DevName = %q, want %q", ev.DevName, "sda")
	}
}

func TestParseUevent_EmptyPacket(t *testing.T) {
	ev := ParseUevent([]byte{})
	if ev != nil {
		t.Errorf("ParseUevent([]) = %+v, want nil", ev)
	}
}

func TestParseUevent_ShortPacket(t *testing.T) {
	ev := ParseUevent([]byte{0x01, 0x02})
	if ev != nil {
		t.Errorf("ParseUevent(short) = %+v, want nil", ev)
	}
}

func TestParseUevent_NoAtSign(t *testing.T) {
	// Header без @ — невалидный формат
	data := buildUeventPacket(
		"invalidheader",
		"SUBSYSTEM=block",
	)

	ev := ParseUevent(data)
	if ev != nil {
		t.Errorf("ParseUevent(no @) = %+v, want nil", ev)
	}
}

func TestParseUevent_GarbageData(t *testing.T) {
	// Случайные байты без null-терминаторов
	garbage := []byte("thisisgarbagedatawithoutanynulls")
	ev := ParseUevent(garbage)
	// Нет '@' в строке — должен вернуть nil
	if ev != nil {
		t.Errorf("ParseUevent(garbage) = %+v, want nil", ev)
	}
}

func TestParseUevent_GarbageWithAtSign(t *testing.T) {
	// Мусор с '@' но без null-терминаторов — должен парсить header
	garbage := []byte("change@/some/path")
	ev := ParseUevent(garbage)
	if ev == nil {
		t.Fatal("ParseUevent вернул nil для данных с @")
	}
	if ev.Action != "change" {
		t.Errorf("Action = %q, want %q", ev.Action, "change")
	}
	if ev.DevPath != "/some/path" {
		t.Errorf("DevPath = %q, want %q", ev.DevPath, "/some/path")
	}
}
