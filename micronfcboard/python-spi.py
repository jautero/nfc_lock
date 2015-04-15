import glob, serial, sys
device=glob.glob("/dev/tty.usbmodem*")
if not device:
  print "Device not found. Try reseting micronfcboard"
  sys.exit(1)
ser=serial.serial_for_url(device[0])
ser.write(bytearray([ 1 << 7 | 0x37 << 1,0]))
print "%02x" % ord(ser.read(2)[1])