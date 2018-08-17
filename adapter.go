package ble

import (
	"log"
	"time"
)

const (
	adapterInterface = "org.bluez.Adapter1"
)

// The Adapter type corresponds to the org.bluez.Adapter1 interface.
// See bluez/doc/adapter-api.txt
//
// StartDiscovery starts discovery on the adapter.
//
// StopDiscovery stops discovery on the adapter.
//
// RemoveDevice removes the specified device and its pairing information.
//
// SetDiscoveryFilter sets the discovery filter to require
// LE transport and the given UUIDs.
//
// Discover performs discovery for a device with the given UUIDs,
// for at most the specified timeout, or indefinitely if timeout is 0.
// See also the Discover method of the ObjectCache type.
type Adapter interface {
	BaseObject

	StartDiscovery() error
	StopDiscovery() error
	RemoveDevice(Device) error
	SetDiscoveryFilter(uuids ...string) error

	Discover(timeout time.Duration, uuids ...string) error
}

// GetAdapter finds an Adapter in the object cache and returns it.
func (conn *Connection) GetAdapter() (Adapter, error) {
	// Get Default Adapter
	adapters, err := conn.GetAdapters()
	if err != nil {
		return nil, err
	}

	return adapters[0], nil
}

func (conn *Connection) GetAdapters() ([]Adapter, error) {
	adapters, err := conn.findObjects(adapterInterface, func(_ *blob) bool { return true })
	var r_adapters = []Adapter{}
	for _, adapter := range(adapters) {
		r_adapters = append(r_adapters, adapter)
	}
	return r_adapters, err
}

func (adapter *blob) StartDiscovery() error {
	log.Printf("%s: starting discovery", adapter.Name())
	return adapter.call("StartDiscovery")
}

func (adapter *blob) StopDiscovery() error {
	log.Printf("%s: stopping discovery", adapter.Name())
	return adapter.call("StopDiscovery")
}

func (adapter *blob) RemoveDevice(device Device) error {
	log.Printf("%s: removing device %s", adapter.Name(), device.Name())
	return adapter.call("RemoveDevice", device.Path())
}
