package chip8

import (
	"fmt"
	"time"
)

type Chip8 struct {
	Memory   [4096]byte
	V        [16]byte
	I        uint16
	PC       uint16
	Stack    [16]uint16
	SP       uint8
	Delay    byte
	Sound    byte
	Keys     [16]byte
	DrawFlag bool
	Screen   [32][64]bool // 64x32 pixel screen
}

func New() *Chip8 {
	return &Chip8{
		PC: 0x200,
	}
}

func (c *Chip8) LoadROM(bytes []byte) {
	// todo: annotate statements and operations
	var rom [4096]byte
	mem := make([]byte, 0x1000)
	copy(mem[0x200:], bytes)
	copy(rom[:], mem)
	c.Memory = rom
}

func (c *Chip8) Run() {
	for {
		if c.PC >= 4096 {
			fmt.Println("Program counter out of bounds")
			break
		}

		firstByte := c.Memory[c.PC]
		secondByte := c.Memory[c.PC+1]

		firstNibble := firstByte >> 4     // Extract upper 4 bits of firstByte
		secondNibble := firstByte & 0x0F  // Extract lower 4 bits of firstByte
		thirdNibble := secondByte >> 4    // Extract upper 4 bits of secondByte
		fourthNibble := secondByte & 0x0F // Extract lower 4 bits of secondByte

		NN := secondByte // 8-bit immediate value (secondByte as-is)
		NNN := uint16(firstByte&0x0F)<<8 | uint16(secondByte)

		c.PC += 2

		switch firstNibble {
		case 0x0:
			switch secondNibble {
			case 0x0: // CLS
				if thirdNibble == 0x0 && fourthNibble == 0xE {
					fmt.Println("CLS")
				}
			default:
				fmt.Println("Unhandled 0x00 instruction")
			}
		case 0x1: // JMP
			// fmt.Println("JMP")
			c.PC = NNN
		case 0x6: // set register VX
			fmt.Println("Set Register VX")
			c.V[secondNibble] = NN
		case 0x7: // Add value NN to VX
			fmt.Println("Add the value NN to VX")
			c.V[secondNibble] += NN
		case 0xA: // set register I
			fmt.Println("Set Register I")
			c.I = NNN
		case 0xD: // DRW
			x := c.V[secondNibble]
			y := c.V[thirdNibble]
			height := fourthNibble
			fmt.Printf("DRW at (%d, %d) with height %d\n", x, y, height)
			c.draw(x, y, height)
		default:
			fmt.Printf("Unknown instruction %x, %x\n", firstByte, secondByte)
		}

		time.Sleep(time.Millisecond * 17) // ~60Hz
	}
}

// todo: work with specification to implement fully
func (c *Chip8) draw(x, y byte, height byte) {
	screenX := x % 64 // Wrap around the screen width
	screenY := y % 32 // Wrap around the screen height
	c.V[0xF] = 0x0    // set VF to zero

	for row := byte(0); row < height; row++ {
		spriteRow := c.Memory[c.I+uint16(row)]
		x := screenX
		y := screenY

		for idx := 0; idx < 8; idx++ {
			if screenX >= 64 {
				continue
			}
			pixel := (spriteRow>>(7-idx))&1 != 0
			c.Screen[y+row][x] = pixel
			x++
		}
	}

	c.Render()
}

func (c *Chip8) Render() {
	fmt.Print("\033[H\033[2J")
	// Print the screen as a grid of '█' and ' ' for pixels on and off
	for y := 0; y < 32; y++ { // todo: put resolution in constants
		for x := 0; x < 64; x++ {
			if c.Screen[y][x] {
				fmt.Print("█") // Pixel on
			} else {
				fmt.Print(" ") // Pixel off
			}
		}
		fmt.Println() // New line after each row
	}
}
