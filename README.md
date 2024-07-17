# go-qmc5883
A Golang library for the qmc5883 compass

# Const
Inclination: 68° 5' 
Variation: 15.8

# SPEC
https://cdn-shop.adafruit.com/datasheets/HMC5883L_3-Axis_Digital_Compass_IC.pdf

# BUILD
for raspberry pi
`env GOOS=linux GOARCH=arm GOARM=5 go build`

# DEBUG
i2cdetect -y 1

// get a value
i2cget -y 1 0x1e 0x01

// set to continue
i2cset 1 0x1e 0x02 0x00

// status
i2cget -y 1 0x1e 0x09
