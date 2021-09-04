package rj32

import (
	"fmt"

	"github.com/rj45/rj32/emu/data"
)

// BusHandler handles a bus transaction
type BusHandler interface {
	HandleBus(bus Bus) Bus
}

// BusHandlerFunc implements BusHandler with a func
type BusHandlerFunc func(bus Bus) Bus

func (fn BusHandlerFunc) HandleBus(bus Bus) Bus {
	return fn(bus)
}

// Device is a memory mapped device
type Device struct {
	Addr, Size int
	Handler    BusHandler
}

// MemMap is a BusHandler for a list of Devices
type MemMap []Device

func (mm MemMap) HandleBus(bus Bus) Bus {
	addr := bus.Address()
	found := false
	for _, dev := range mm {
		if addr >= dev.Addr && addr < (dev.Addr+dev.Size) {
			if found && !bus.WE() && bus.Ack() {
				panic(fmt.Sprintf("Bus conflict on address %x", bus.Address()))
			}
			found = true

			bus = dev.Handler.HandleBus(bus)
		}
	}
	return bus
}

// AddrOffset is a BusHandler that subtracts Offset from
// the address before passing it off to a sub-handler
type AddrOffset struct {
	Offset  int
	Handler BusHandler
}

func (off *AddrOffset) HandleBus(bus Bus) Bus {
	addr := bus.Address()
	return off.Handler.
		HandleBus(bus.SetAddress(addr - off.Offset)).
		SetAddress(addr + off.Offset)
}

// RAM is full read/write memory
type RAM struct {
	*data.Memory
}

func (ram *RAM) HandleBus(bus Bus) Bus {
	addr := bus.Address()

	if bus.WE() {
		ram.Write(addr, uint16(bus.Data()))
		return bus.SetAck(true)
	}

	return bus.
		SetAck(true).
		SetData(int(ram.Read(addr)))
}

// ROM ignores (but acks) writes but handles reads
type ROM struct {
	*data.Memory
}

func (rom *ROM) HandleBus(bus Bus) Bus {
	addr := bus.Address()

	if bus.WE() {
		return bus.SetAck(true)
	}

	return bus.
		SetAck(true).
		SetData(int(rom.Read(addr)))
}

// ShadowMem writes to a memory but does not ack the
// request, expecting another memory to do that.
type ShadowMem struct {
	*data.Memory
}

func (mem *ShadowMem) HandleBus(bus Bus) Bus {
	addr := bus.Address()

	if bus.WE() {
		mem.Write(addr, uint16(bus.Data()))
		return bus
	}

	return bus
}
