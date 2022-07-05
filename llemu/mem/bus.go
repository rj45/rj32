package mem

type Bus struct {
	Address uint32
	Data    uint32
	Mask    uint32
	WE      bool
	Ack     bool
}

// BusHandler handles a bus transaction
type BusHandler interface {
	HandleBus(bus *Bus)
}

// BusHandlerFunc implements BusHandler with a func
type BusHandlerFunc func(bus *Bus)

func (fn BusHandlerFunc) HandleBus(bus *Bus) {
	fn(bus)
}

// Device is a memory mapped device
type Device struct {
	Address, Size uint32
	Handler       BusHandler
}

// MemMap is a BusHandler for a list of Devices
type MemMap []Device

func (mm MemMap) HandleBus(bus *Bus) {
	for _, dev := range mm {
		if bus.Address >= dev.Address && bus.Address < (dev.Address+dev.Size) {
			dev.Handler.HandleBus(bus)

			if bus.Ack {
				break
			}
		}
	}
}

// AddrOffset is a BusHandler that subtracts Offset from
// the address before passing it off to a sub-handler
type AddrOffset struct {
	Offset  uint32
	Handler BusHandler
}

func (off *AddrOffset) HandleBus(bus *Bus) {
	oldAddress := bus.Address
	bus.Address -= off.Offset
	off.Handler.HandleBus(bus)
	bus.Address = oldAddress
}
