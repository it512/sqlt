package funcs

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

type UUID [16]byte

func (u UUID) String() string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x", u[:4], u[4:6], u[6:8], u[8:10], u[10:])
}

func RandomUUID() string {
	u, _ := newRandom()
	uuid := UUID{}
	copy(uuid[:], u)
	return uuid.String()
}

func TimeBaseUUID() string {
	u, _ := newTimeBase()
	uuid := UUID{}
	copy(uuid[:], u)
	return uuid.String()
}

func newRandom() ([]byte, error) {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid[:])
	if err != nil {
		return nil, err
	}
	if n != len(uuid) {
		return nil, errors.New("could not generate random bytes with 16 length")
	}

	SetVersion(uuid, VersionRandom)
	SetLayout(uuid, LayoutRFC4122)

	return uuid, nil
}

type Layout byte

const (
	LayoutInvalid Layout = iota
	LayoutNCS
	LayoutRFC4122
	LayoutMicrosoft
	LayoutFuture
)

func SetLayout(uuid []byte, layout Layout) {
	switch layout {
	case LayoutNCS:
		uuid[8] = (uuid[8] | 0x00) & 0x0f
	case LayoutRFC4122:
		uuid[8] = (uuid[8] | 0x80) & 0x8f
	case LayoutMicrosoft:
		uuid[8] = (uuid[8] | 0xc0) & 0xcf
	case LayoutFuture:
		uuid[8] = (uuid[8] | 0xe0) & 0xef
	default:
		panic("layout is invalid")
	}
}

type Version byte

const (
	VersionUnknown Version = iota
	VersionTimeBased
	VersionDCESecurity
	VersionNameBasedMD5
	VersionRandom
	VersionNameBasedSHA1
)

func SetVersion(uuid []byte, version Version) {
	switch version {
	case VersionTimeBased:
		uuid[6] = (uuid[6] | 0x10) & 0x1f
	case VersionDCESecurity:
		uuid[6] = (uuid[6] | 0x20) & 0x2f
	case VersionNameBasedMD5:
		uuid[6] = (uuid[6] | 0x30) & 0x3f
	case VersionRandom:
		uuid[6] = (uuid[6] | 0x40) & 0x4f
	case VersionNameBasedSHA1:
		uuid[6] = (uuid[6] | 0x50) & 0x5f
	default:
		panic("version is unknown")
	}
}

const (
	intervals = (2440587 - 2299160) * 86400 * 10000000
)

var (
	lastGenerated time.Time
	clockSequence uint16
	nodeID        []byte
	locker        sync.Mutex
)

func newTimeBase() ([]byte, error) {
	locker.Lock()
	defer locker.Unlock()

	uuid := make([]byte, 16)

	now := time.Now().UTC()
	timestamp := uint64(now.UnixNano()/100) + intervals
	if !now.After(lastGenerated) {
		clockSequence++
	} else {
		b := make([]byte, 2)
		_, err := rand.Read(b)
		if err != nil {
			return nil, err
		}
		clockSequence = uint16(int(b[0])<<8 | int(b[1]))
	}

	lastGenerated = now

	timeLow := uint32(timestamp & 0xffffffff)
	timeMiddle := uint16((timestamp >> 32) & 0xffff)
	timeHigh := uint16((timestamp >> 48) & 0xfff)

	binary.BigEndian.PutUint32(uuid[0:], timeLow)
	binary.BigEndian.PutUint16(uuid[4:], timeMiddle)
	binary.BigEndian.PutUint16(uuid[6:], timeHigh)
	binary.BigEndian.PutUint16(uuid[8:], clockSequence)

	if nodeID == nil {
		interfaces, err := net.Interfaces()
		if err != nil {
			return nil, err
		}

		for _, i := range interfaces {
			if len(i.HardwareAddr) >= 6 {
				nodeID = make([]byte, 6)
				copy(nodeID, i.HardwareAddr)
				break
			}
		}

		if nodeID == nil {
			nodeID = make([]byte, 6)
			_, err := rand.Read(nodeID)
			if err != nil {
				return nil, err
			}
		}
	}

	copy(uuid[10:], nodeID)

	SetVersion(uuid, VersionTimeBased)
	SetLayout(uuid, LayoutRFC4122)

	return uuid, nil
}
