package goboy

import (
	"log"
//	"sort"
)

const (
	SCREEN_WIDTH = 160
	SCREEN_HEIGHT = 144
)

type palette struct {
	bg [4]uint8
	obj0 [4]uint8
	obj1 [4]uint8
}

type objdata struct {
	x, y int16
	tile, palette int
	yflip, xflip bool
	prio int
	num int
}

// type objdatalist [40]objdata

// TODO
// func (a objdatalist) Len() int { return len(a) }
// func (a objdatalist) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
// func (a objdatalist) Less(i, j int) bool {
// 	if a[i].x < a[j].x {
// 		return true
// 	}
// 	return a[i].num < a[j].num
// }

type gpu struct {
	vram [8192]uint8
	oam [SCREEN_WIDTH]uint8
	reg [0xBF]uint8
	od, odsorted [40]objdata
	scanrow [SCREEN_WIDTH]uint8
	Tilemap [512][8][8]uint8
	pal palette

	Screen [SCREEN_WIDTH*SCREEN_HEIGHT*4]uint8

	curline uint8
	curscan uint16
	linemode, modeclocks int

	yscrl, xscrl bool
	raster, ints uint8

	lcdon, bgon, objon, winon bool

	objsize int

	bgtilebase, bgmapbase, wintilebase uint16
}

func makeGPU() gpu {
	var g gpu
	for i := 0; i < 4; i += 1 {
		g.pal.bg[i] = 0xFF
		g.pal.obj0[i] = 0xFF
		g.pal.obj1[i] = 0xFF
	}

	log.Printf("gpu: initializing screen")
	for i := range g.Screen {
		g.Screen[i] = 255
	}

	g.linemode = 2

	for i := range g.od {
		g.od[i].x = -8
		g.od[i].y = -16
		g.od[i].num = i
	}

	g.bgtilebase = 0x0000
	g.bgmapbase = 0x1800
	g.wintilebase = 0x1800

	return g
}

func (g *gpu) Checkline() {
	g.modeclocks += Z80.R.M
	switch g.linemode {
		// hblank
		case 0:
			if g.modeclocks >= 51 {
				// end of hblank, last scanline: render
				if g.curline == 143 {
					g.linemode = 1
					MMU.If |= 1
				} else {
					g.linemode = 2
				}
				g.curline += 1
				g.curscan += 640
				g.modeclocks = 0
			}

		// vblank
		case 1:
			if g.modeclocks >= 114 {
				g.modeclocks = 0
				g.curline += 1
				if g.curline > 153 {
					g.curline = 0
					g.curscan = 0
					g.linemode = 2
				}
			}

		// oam read
		case 2:
			if g.modeclocks >= 20 {
				g.modeclocks = 0
				g.linemode = 3
			}

		// vram read
		case 3:
			// render scanline at end of allotted time
			if g.modeclocks >= 43 {
				g.modeclocks = 0
				g.linemode = 0
				if g.lcdon {
					if g.bgon {
						linebase := g.curscan
						var yscrl uint8 = 0
						if g.yscrl {
							yscrl = 1
						}
						xscrl := 0
						if g.xscrl {
							xscrl = 1
						}
						mapbase := g.bgmapbase + uint16(((g.curline + yscrl) >> 3) << 5)
						y := (g.curline + yscrl) & 7
						x := xscrl & 7
						t := uint16(xscrl >> 3) & 31

						if g.bgtilebase != 0 {
							tile := g.vram[mapbase + t]
							// TODO?
							//if tile < 128 {
						//		tile = 256 + tile
						//	}
							tilerow := g.Tilemap[tile][y]
							for w := SCREEN_WIDTH; w > 0; w -= 1 {
								g.scanrow[SCREEN_WIDTH-x-1] = tilerow[x]
								log.Printf("a: %x = %x\n", linebase+3, g.pal.bg[tilerow[x]])
								g.Screen[linebase + 3] = g.pal.bg[tilerow[x]]
								x += 1
								if x == 8 {
									t = (t + 1) & 31
									x = 0
									tile = g.vram[mapbase + t]
									// TODO?
									// if tile < 128 {
									// 	tile = 256 + tile
									// }
									tilerow = g.Tilemap[tile][y]
								}
								linebase += 4
							}
						} else {
							tilerow := g.Tilemap[g.vram[mapbase+t]][y]
							for w := SCREEN_WIDTH; w > 0; w -= 1 {
								g.scanrow[SCREEN_WIDTH-x-1] = tilerow[x]
								g.Screen[linebase + 3] = g.pal.bg[tilerow[x]]
								log.Printf("b: %x = %x\n", linebase+3, g.pal.bg[tilerow[x]])
								x += 1
								if x == 8 {
									t = (t + 1) & 31
									x = 0
									tilerow = g.Tilemap[g.vram[mapbase+t]][y]
								}
								linebase += 4
							}
						}
					}
	
					if g.objon {
						cnt := 0
						if g.objsize != 0 {
							// TODO
						} else {
							linebase := int16(g.curscan)
							for i := range g.od {
								obj := g.odsorted[i]
								var curline int16 = int16(g.curline)
								if obj.y <= curline && obj.y + 8 > curline {
									tilerow := g.Tilemap[obj.tile][curline - obj.y]
									if obj.yflip {
										tilerow = g.Tilemap[obj.tile][7 - (curline - obj.y)]
									}

									pal := g.pal.obj0
									if obj.palette != 0 {
										pal = g.pal.obj1
									}

									linebase = (curline * SCREEN_WIDTH + obj.x) * 4

									if obj.xflip {
										for x := int16(0); x < 8; x += 1 {
											if obj.x + x >= 0 && obj.x + x < SCREEN_WIDTH {
												if tilerow[7-x] != 0 && (obj.prio != 0 || g.scanrow[x] == 0) {
													g.Screen[linebase + 3] = pal[tilerow[7-x]]
												}
											}
											linebase += 4
										}
									} else {
										for x := int16(0); x < 8; x += 1 {
											if obj.x + x >= 0 && obj.x + x < SCREEN_WIDTH {
												if tilerow[x] != 0 && (obj.prio != 0 || g.scanrow[x] == 0) {
													g.Screen[linebase + 3] = pal[tilerow[x]]
												}
											}
											linebase += 4
										}
									}

									cnt += 1
									if cnt > 10 {
										break
									}
								}
							}
						}
					}
				}
			}
	}
}

func (g *gpu) UpdateTile(addr uint16, value uint8) {
	if (addr & 0x1) != 0 {
		addr -= 1
	}
	saddr := addr

	tile := (addr >> 4) & 511
	y := (addr >> 1) & 7
	for x := 0; x < 8; x += 1 {
		sx := byte(1 << (7 - uint(x)))

		t := 0
		if g.vram[saddr] & sx != 0 {
			t = 1
		}
		if g.vram[saddr+1] & sx != 0 {
			t |= 2
		}
		g.Tilemap[tile][y][x] = byte(t)
	}
}

func (g *gpu) UpdateOAM(addr uint16, value uint8) {
	addr -= 0xFE00

	obj := addr >> 2
	if obj < 40 {
		switch addr & 3 {
			case 0: g.od[obj].y = int16(value) - 16
			case 1: g.od[obj].x = int16(value) - 8
			case 2:
				g.od[obj].tile = int(value)
				if g.objsize != 0 {
					g.od[obj].tile = int(value) & 0xFE
				}
			case 3:
				g.od[obj].palette = 0
				if (value & 0x10) != 0{
					g.od[obj].palette = 1
				}
				g.od[obj].xflip = false
				if (value & 0x20) != 0{
					g.od[obj].xflip = true
				}
				g.od[obj].yflip = false
				if (value & 0x40) != 0{
					g.od[obj].yflip = true
				}
				g.od[obj].prio = 0
				if (value & 0x80) != 0{
					g.od[obj].prio = 1
				}
		}
	}

	for i := range g.od {
		g.odsorted[i] = g.od[i]
	}
	// TODO
	// sort.Sort(g.odsorted)
}

func (g gpu) ReadByte(addr uint16) uint8 {
	gaddr := addr - 0xFF40
	switch gaddr {
		case 0:
			value := 0
			if g.lcdon { value |= 0x80 }
			if g.bgtilebase == 0x0000 { value |= 0x10 }
			if g.bgmapbase == 0x1C00 { value |= 0x08 }
			if g.objsize != 0 { value |= 0x04 }
			if g.objon { value |= 0x02 }
			if g.bgon { value |= 0x01 }
			return byte(value & 0xFF)
		case 1:
			value := g.linemode
			if g.curline == g.raster { value |= 0x4 }
			return uint8(value & 0xFF)
		case 2:
			if g.yscrl { return 1 }
			return 0
		case 3:
			if g.xscrl { return 1 }
			return 0
		case 4: return byte(g.curline)
		case 5: return g.raster
	}
	return g.reg[gaddr]
}

func (g *gpu) WriteByte(addr uint16, value uint8) {
	gaddr := addr - 0xFF40
	g.reg[gaddr] = value
	switch gaddr {
		case 0:
			g.lcdon = false
			if (value & 0x80) != 0 { g.lcdon = true }
			g.bgtilebase = 0x0800
			if (value & 0x10) != 0 { g.bgtilebase = 0x0000 }
			g.bgmapbase = 0x1800
			if (value & 0x08) != 0 { g.bgmapbase = 0x1C00 }
			g.objsize = 0
			if (value & 0x04) != 0 { g.objsize = 1 }
			g.objon = false
			if (value & 0x02) != 0 { g.objon = true }
			g.bgon = false
			if (value & 0x01) != 0 { g.bgon = true }
		case 2:
			g.yscrl = false
			if value != 0 { g.yscrl = true }
		case 3:
			g.xscrl = false
			if value != 0 { g.xscrl = true }
		case 5: g.raster = value
		case 6:
			// OAM DMA
			for i := uint16(0); i < SCREEN_WIDTH; i += 1 {
				v := MMU.ReadByte(uint16(value << 8) + i)
				g.oam[i] = v
				g.UpdateOAM(uint16(0xFE00) + i, v)
			}
		case 7:
			// BG palette mapping
			for i := uint(0); i < 4; i += 1 {
				switch (value >> (i*2)) & 3 {
					case 0: g.pal.bg[i] = 255
					case 1: g.pal.bg[i] = 192
					case 2: g.pal.bg[i] = 96
					case 3: g.pal.bg[i] = 0
				}
			}

		case 8:
			// obj0 palette mapping
			for i := uint(0); i < 4; i += 1 {
				switch (value >> (i*2)) & 3 {
					case 0: g.pal.obj0[i] = 255
					case 1: g.pal.obj0[i] = 192
					case 2: g.pal.obj0[i] = 96
					case 3: g.pal.obj0[i] = 0
				}
			}
		case 9:
			// obj1 palette mapping
			for i := uint(0); i < 4; i += 1 {
				switch (value >> (i*2)) & 3 {
					case 0: g.pal.obj1[i] = 255
					case 1: g.pal.obj1[i] = 192
					case 2: g.pal.obj1[i] = 96
					case 3: g.pal.obj1[i] = 0
				}
			}
	}
}
