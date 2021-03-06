package ble

import (
	"log"
	"github.com/godbus/dbus"
)

const (
	deviceInterface = "org.bluez.Device1"
	interfacesAdded = "org.freedesktop.DBus.ObjectManager.InterfacesAdded"
)

// The Device type corresponds to the org.bluez.Device1 interface.
// See bluez/doc/devicet-api.txt
type Device interface {
	BaseObject

	UUIDs() []string
	Connected() bool
	Paired() bool

	Connect() error
	Disconnect() error
	Pair() error

	Address() string
	Rssi() int16
	Adapter() dbus.ObjectPath
}

func (conn *Connection) matchDevice(matching predicate) (Device, error) {
	return conn.findObject(deviceInterface, matching)
}

func (conn *Connection) matchDevices(matching predicate) ([]Device, error) {
	arr, err := conn.findObjects(deviceInterface, matching)
	var devices = []Device{}

	for _, d := range(arr) {
		devices = append(devices, d)
	}

	return devices, err
}

// GetDevice finds a Device in the object cache matching the given UUIDs.
func (conn *Connection) GetDevice(uuids ...string) (Device, error) {
	return conn.matchDevice(func(device *blob) bool {
		return uuidsInclude(device.UUIDs(), uuids)
	})
}

func (conn *Connection) GetDevices(uuids ...string) ([]Device, error) {
	return conn.matchDevices(func(device *blob) bool {
		return uuidsInclude(device.UUIDs(), uuids)
	})
}

func uuidsInclude(advertised []string, uuids []string) bool {
	for _, u := range uuids {
		if !ValidUUID(u) {
			log.Printf("invalid UUID %s", u)
			return false
		}
		if !stringArrayContains(advertised, u) {
			return false
		}
	}
	return true
}

// GetDeviceByName finds a Device in the object cache with the given name.
func (conn *Connection) GetDeviceByName(name string) (Device, error) {
	return conn.matchDevice(func(device *blob) bool {
		return device.Name() == name
	})
}

func (device *blob) UUIDs() []string {
	return device.properties["UUIDs"].Value().([]string)
}

func (device *blob) Connected() bool {
	return device.properties["Connected"].Value().(bool)
}

func (device *blob) Address() string {
	return device.properties["Address"].Value().(string)
}

func (device *blob) Rssi() int16 {
	rssi:= device.properties["RSSI"].Value()
	if rssi == nil {
		return -1
	}
	return rssi.(int16)
}

func (device *blob) Adapter() dbus.ObjectPath {
	return device.properties["Adapter"].Value().(dbus.ObjectPath)
}

func (device *blob) Paired() bool {
	return device.properties["Paired"].Value().(bool)
}

func (device *blob) Connect() error {
	log.Printf("%s: connecting", device.Name())
	return device.call("Connect")
}

func (device *blob) Disconnect() error {
	log.Printf("%s: disconnecting", device.Name())
	return device.call("Disconnect")
}

func (device *blob) Pair() error {
	log.Printf("%s: pairing", device.Name())
	return device.call("Pair")
}

func stringArrayContains(a []string, str string) bool {
	for _, s := range a {
		if s == str {
			return true
		}
	}
	return false
}
