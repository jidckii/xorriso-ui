// Пакет udev предоставляет netlink listener для udev событий блочных устройств.
// Используется для отслеживания вставки/извлечения оптических дисков без polling.
package udev

import (
	"context"
	"strings"

	"golang.org/x/sys/unix"
)

// Event представляет udev-событие ядра
type Event struct {
	Action    string // "add", "remove", "change"
	DevPath   string // "/devices/.../block/sr0"
	Subsystem string // "block"
	DevName   string // "sr0"
	DevType   string // "disk"
}

// Monitor слушает netlink udev события ядра
type Monitor struct {
	fd     int
	stopCh chan struct{}
}

// New создаёт новый Monitor, открывая AF_NETLINK socket для NETLINK_KOBJECT_UEVENT.
func New() (*Monitor, error) {
	fd, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_DGRAM, unix.NETLINK_KOBJECT_UEVENT)
	if err != nil {
		return nil, err
	}

	addr := &unix.SockaddrNetlink{
		Family: unix.AF_NETLINK,
		Groups: 1, // kernel events group
	}

	if err := unix.Bind(fd, addr); err != nil {
		_ = unix.Close(fd)
		return nil, err
	}

	return &Monitor{
		fd:     fd,
		stopCh: make(chan struct{}),
	}, nil
}

// Listen читает udev пакеты и отправляет отфильтрованные события в канал.
// Фильтрация: SUBSYSTEM=block и DEVNAME начинается с "sr".
// Канал закрывается при отмене контекста или вызове Close().
func (m *Monitor) Listen(ctx context.Context) (<-chan Event, error) {
	ch := make(chan Event, 16)

	// Устанавливаем таймаут чтения, чтобы периодически проверять отмену контекста
	tv := unix.Timeval{Sec: 1}
	if err := unix.SetsockoptTimeval(m.fd, unix.SOL_SOCKET, unix.SO_RCVTIMEO, &tv); err != nil {
		close(ch)
		return nil, err
	}

	go func() {
		defer close(ch)
		buf := make([]byte, 4096)

		for {
			select {
			case <-ctx.Done():
				return
			case <-m.stopCh:
				return
			default:
			}

			n, err := unix.Read(m.fd, buf)
			if err != nil {
				// Таймаут — нормальная ситуация, продолжаем
				if err == unix.EAGAIN || err == unix.EINTR {
					continue
				}
				// Сокет закрыт
				return
			}
			if n <= 0 {
				continue
			}

			ev := ParseUevent(buf[:n])
			if ev == nil {
				continue
			}

			// Фильтруем: только block-устройства sr*
			if ev.Subsystem != "block" || !strings.HasPrefix(ev.DevName, "sr") {
				continue
			}

			select {
			case ch <- *ev:
			case <-ctx.Done():
				return
			case <-m.stopCh:
				return
			}
		}
	}()

	return ch, nil
}

// Close закрывает netlink socket и останавливает горутину Listen.
func (m *Monitor) Close() error {
	select {
	case <-m.stopCh:
		// Уже закрыт
	default:
		close(m.stopCh)
	}
	return unix.Close(m.fd)
}

// ParseUevent парсит raw netlink uevent пакет.
// Формат: первая строка "action@devpath\0", затем KEY=VALUE\0 пары.
// Возвращает nil если пакет слишком короткий или невалидный.
func ParseUevent(data []byte) *Event {
	if len(data) < 4 {
		return nil
	}

	// Разбиваем на null-terminated строки
	parts := splitNullTerminated(data)
	if len(parts) == 0 {
		return nil
	}

	// Первая строка: action@devpath
	header := parts[0]
	before, after, ok := strings.Cut(header, "@")
	if !ok {
		return nil
	}

	ev := &Event{
		Action:  before,
		DevPath: after,
	}

	// Парсим KEY=VALUE пары
	for _, part := range parts[1:] {
		before, after, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		key := before
		val := after

		switch key {
		case "SUBSYSTEM":
			ev.Subsystem = val
		case "DEVNAME":
			ev.DevName = val
		case "DEVTYPE":
			ev.DevType = val
		}
	}

	return ev
}

// splitNullTerminated разбивает байты по null-байтам, пропуская пустые строки.
func splitNullTerminated(data []byte) []string {
	var result []string
	start := 0
	for i, b := range data {
		if b == 0 {
			if i > start {
				result = append(result, string(data[start:i]))
			}
			start = i + 1
		}
	}
	// Обработка последнего фрагмента без завершающего null
	if start < len(data) {
		result = append(result, string(data[start:]))
	}
	return result
}
