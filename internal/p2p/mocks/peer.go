// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	conn "github.com/tendermint/tendermint/internal/p2p/conn"
	log "github.com/tendermint/tendermint/libs/log"

	mock "github.com/stretchr/testify/mock"

	net "net"

	types "github.com/tendermint/tendermint/types"
)

// Peer is an autogenerated mock type for the Peer type
type Peer struct {
	mock.Mock
}

// CloseConn provides a mock function with given fields:
func (_m *Peer) CloseConn() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FlushStop provides a mock function with given fields:
func (_m *Peer) FlushStop() {
	_m.Called()
}

// Get provides a mock function with given fields: _a0
func (_m *Peer) Get(_a0 string) interface{} {
	ret := _m.Called(_a0)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// ID provides a mock function with given fields:
func (_m *Peer) ID() types.NodeID {
	ret := _m.Called()

	var r0 types.NodeID
	if rf, ok := ret.Get(0).(func() types.NodeID); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(types.NodeID)
	}

	return r0
}

// IsOutbound provides a mock function with given fields:
func (_m *Peer) IsOutbound() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsPersistent provides a mock function with given fields:
func (_m *Peer) IsPersistent() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsRunning provides a mock function with given fields:
func (_m *Peer) IsRunning() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NodeInfo provides a mock function with given fields:
func (_m *Peer) NodeInfo() types.NodeInfo {
	ret := _m.Called()

	var r0 types.NodeInfo
	if rf, ok := ret.Get(0).(func() types.NodeInfo); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(types.NodeInfo)
	}

	return r0
}

// OnReset provides a mock function with given fields:
func (_m *Peer) OnReset() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OnStart provides a mock function with given fields:
func (_m *Peer) OnStart() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OnStop provides a mock function with given fields:
func (_m *Peer) OnStop() {
	_m.Called()
}

// Quit provides a mock function with given fields:
func (_m *Peer) Quit() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// RemoteAddr provides a mock function with given fields:
func (_m *Peer) RemoteAddr() net.Addr {
	ret := _m.Called()

	var r0 net.Addr
	if rf, ok := ret.Get(0).(func() net.Addr); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(net.Addr)
		}
	}

	return r0
}

// RemoteIP provides a mock function with given fields:
func (_m *Peer) RemoteIP() net.IP {
	ret := _m.Called()

	var r0 net.IP
	if rf, ok := ret.Get(0).(func() net.IP); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(net.IP)
		}
	}

	return r0
}

// Reset provides a mock function with given fields:
func (_m *Peer) Reset() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Send provides a mock function with given fields: _a0, _a1
func (_m *Peer) Send(_a0 byte, _a1 []byte) bool {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(byte, []byte) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Set provides a mock function with given fields: _a0, _a1
func (_m *Peer) Set(_a0 string, _a1 interface{}) {
	_m.Called(_a0, _a1)
}

// SetLogger provides a mock function with given fields: _a0
func (_m *Peer) SetLogger(_a0 log.Logger) {
	_m.Called(_a0)
}

// SocketAddr provides a mock function with given fields:
func (_m *Peer) SocketAddr() *types.NetAddress {
	ret := _m.Called()

	var r0 *types.NetAddress
	if rf, ok := ret.Get(0).(func() *types.NetAddress); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.NetAddress)
		}
	}

	return r0
}

// Start provides a mock function with given fields:
func (_m *Peer) Start() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Status provides a mock function with given fields:
func (_m *Peer) Status() conn.ConnectionStatus {
	ret := _m.Called()

	var r0 conn.ConnectionStatus
	if rf, ok := ret.Get(0).(func() conn.ConnectionStatus); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(conn.ConnectionStatus)
	}

	return r0
}

// Stop provides a mock function with given fields:
func (_m *Peer) Stop() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// String provides a mock function with given fields:
func (_m *Peer) String() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// TrySend provides a mock function with given fields: _a0, _a1
func (_m *Peer) TrySend(_a0 byte, _a1 []byte) bool {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(byte, []byte) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Wait provides a mock function with given fields:
func (_m *Peer) Wait() {
	_m.Called()
}
