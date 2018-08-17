package ble

import "github.com/godbus/dbus"

const (
	// GATTMTU is the maximum size of a write to a GATT characteristic.
	GATTMTU = 20

	serviceInterface        = "org.bluez.GattService1"
	characteristicInterface = "org.bluez.GattCharacteristic1"
	descriptorInterface     = "org.bluez.GattDescriptor1"
)

func (conn *Connection) findGattObject(iface string, uuid string) (*blob, error) {
	return conn.findObject(iface, func(desc *blob) bool {
		return desc.UUID() == uuid
	})
}

// GattHandle is the interface satisfied by GATT handles.
type GattHandle interface {
	BaseObject

	UUID() string
}

// UUID returns the handle's UUID
func (handle *blob) UUID() string {
	return handle.properties["UUID"].Value().(string)
}

// Service corresponds to the org.bluez.GattService1 interface.
// See bluez/doc/gatt-api.txt
type Service interface {
	GattHandle
	Primary() bool
	Device() dbus.ObjectPath
	Includes() []dbus.ObjectPath
}

func (handle *blob) Primary() bool {
	return handle.properties["Primary"].Value().(bool)
}

func (handle *blob) Device() dbus.ObjectPath {
	return handle.properties["Device"].Value().(dbus.ObjectPath)
}

func (handle *blob) Includes() []dbus.ObjectPath {
	return handle.properties["Includes"].Value().([]dbus.ObjectPath)
}

// GetService finds a Service with the given UUID.
func (conn *Connection) GetService(uuid string) (Service, error) {
	return conn.findGattObject(serviceInterface, uuid)
}

// ReadWriteHandle is the interface satisfied by GATT objects
// that provide ReadValue and WriteValue operations.
type ReadWriteHandle interface {
	GattHandle

	ReadValue() ([]byte, error)
	WriteValue([]byte) error
}

// ReadValue reads the handle's value.
func (handle *blob) ReadValue() ([]byte, error) {
	var data []byte
	err := handle.callv("ReadValue", Properties{}).Store(&data)
	return data, err
}

// WriteValue writes a value to the handle.
func (handle *blob) WriteValue(data []byte) error {
	return handle.call("WriteValue", data, Properties{})
}

// NotifyHandler represents a function that handles notifications.
type NotifyHandler func([]byte)

// Characteristic corresponds to the org.bluez.GattCharacteristic1 interface.
// See bluez/doc/gatt-api.txt
type Characteristic interface {
	ReadWriteHandle

	Notifying() bool

	StartNotify() error
	StopNotify() error

	HandleNotify(NotifyHandler) error
}

// GetCharacteristic finds a Characteristic with the given UUID.
func (conn *Connection) GetCharacteristic(uuid string) (Characteristic, error) {
	return conn.findGattObject(characteristicInterface, uuid)
}

// Notifying returns whether or not a Characteristic is notifying.
func (handle *blob) Notifying() bool {
	return handle.properties["Notifying"].Value().(bool)
}

// StartNotify starts notifying.
func (handle *blob) StartNotify() error {
	return handle.call("StartNotify")
}

// StartNotify stops notifying.
func (handle *blob) StopNotify() error {
	return handle.call("StopNotify")
}

// Descriptor corresponds to the org.bluez.GattDescriptor1 interface.
// See bluez/doc/gatt-api.txt
type Descriptor interface {
	ReadWriteHandle
}

// GetDescriptor finds a Descriptor with the given UUID.
func (conn *Connection) GetDescriptor(uuid string) (Descriptor, error) {
	return conn.findGattObject(descriptorInterface, uuid)
}
