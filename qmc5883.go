package goqmc5883

import (
	"encoding/binary"
	"log"
	"math"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

type Magnetometer struct {
	Device *i2c.Dev
	i2cbus i2c.BusCloser
	minX   int16
	maxX   int16
	minY   int16
	maxY   int16
	minZ   int16
	maxZ   int16
}

const COM_ADDR = 0x1e

// qmc5883 regs
const (
	CRA  = 0x00
	CRB  = 0x01
	MR   = 0x02
	XMSB = 0x03
	XLSB = 0x04
	ZMSB = 0x05
	ZLSB = 0x06
	YMSB = 0x07
	YLSB = 0x08
)

// New creates a new Magnotometer struct
func New() *Magnetometer {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	i2cbus, err := i2creg.Open("")
	if err != nil {
		log.Fatal("fail to open:", err)
	}

	var mag Magnetometer
	mag.i2cbus = i2cbus
	mag.Device = &i2c.Dev{Addr: COM_ADDR, Bus: i2cbus}

	x, y, z, err := mag.GetXYZ()
	if err != nil {
		log.Fatal("fail to get init measure:", err)
	}
	mag.maxX = x
	mag.minX = x
	mag.maxY = y
	mag.minY = y
	mag.minZ = z
	mag.maxZ = z

	return &mag
}

func (m *Magnetometer) Close() {
	defer m.i2cbus.Close()
}

func (m *Magnetometer) GetXYZ() (int16, int16, int16, error) {
	xm, err := m.ReadData([]byte{XMSB})
	if err != nil {
		return 0, 0, 0, err
	}
	xl, err := m.ReadData([]byte{XLSB})
	if err != nil {
		return 0, 0, 0, err
	}
	zm, err := m.ReadData([]byte{ZMSB})
	if err != nil {
		return 0, 0, 0, err
	}
	zl, err := m.ReadData([]byte{ZLSB})
	if err != nil {
		return 0, 0, 0, err
	}
	ym, err := m.ReadData([]byte{YMSB})
	if err != nil {
		return 0, 0, 0, err
	}
	yl, err := m.ReadData([]byte{YLSB})
	if err != nil {
		return 0, 0, 0, err
	}
	x := calcTwoC(xm[0], xl[0])
	z := calcTwoC(zm[0], zl[0])
	y := calcTwoC(ym[0], yl[0])

	return x, y, z, nil
}

func (m *Magnetometer) ReadData(writeBuf []byte) ([]byte, error) {
	// Send a command 0x10 and expect a 1 byte reply.
	write := writeBuf
	read := make([]byte, 1)
	if err := m.Device.Tx(write, read); err != nil {
		return nil, err
	}
	return read, nil
}

func calcTwoC(msb, lsb byte) int16 {
	val := int16(binary.LittleEndian.Uint16([]byte{lsb, msb}))
	return val
}

func (m *Magnetometer) GetTrueHeadingAzimuth() (int, error) {
	x, y, z, err := m.GetXYZ()
	if err != nil {
		return 0, err
	}
	log.Printf("x: %d y: %d z: %d\n", x, y, z)

	if x > m.maxX {
		m.maxX = x
	}
	if x < m.minX {
		m.minX = x
	}

	if y > m.maxY {
		m.maxY = y
	}
	if y < m.minY {
		m.minY = y
	}

	if z > m.maxZ {
		m.maxZ = z
	}
	if z < m.minZ {
		m.minZ = z
	}

	offsetX := (m.minX + m.maxX) / 2
	offsetY := (m.minY + m.maxY) / 2
	offsetZ := (m.minZ + m.maxZ) / 2
	x = x - offsetX
	y = y - offsetY
	z = z - offsetZ

	x = x * -1

	log.Printf("maxx: %d minx: %d\n", m.maxX, m.minX)
	log.Printf("maxy: %d miny: %d\n", m.maxY, m.minY)
	log.Printf("maxz: %d minz: %d\n", m.maxZ, m.minZ)
	log.Printf("norm: x: %d y: %d z: %d\n", x, y, z)

	azimuth := math.Atan2(float64(y), float64(x)) * (180 / math.Pi)

	res := int(azimuth) % 360
	if res <= 0 {
		res = res + 360
	}

	log.Printf("magnetic heading: %d", res)
	res += 15

	res = int(res) % 360
	if res <= 0 {
		res = res + 360
	}
	log.Printf("true heading: %d", res)

	return res, nil
}
