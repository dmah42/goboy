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
	bg [4]byte
	obj0 [4]byte
	obj1 [4]byte
}

type objdata struct {
	x, y, tile, palette int
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
	vram, oam, reg []byte
	od, odsorted [40]objdata
	scanrow [SCREEN_WIDTH]byte
	tilemap [512][8][8]byte
	pal palette

	scrn [SCREEN_WIDTH*SCREEN_HEIGHT*4]byte

	curline int
	curscan, linemode, modeclocks int

	yscrl, xscrl bool
	raster, ints byte

	lcdon, bgon, objon, winon bool

	objsize int

	bgtilebase, bgmapbase, wintilebase int
}

func makeGPU() gpu {
	var g gpu
	for i := 0; i < 4; i += 1 {
		g.pal.bg[i] = 0xFF
		g.pal.obj0[i] = 0xFF
		g.pal.obj1[i] = 0xFF
	}

	log.Printf("gpu: initializing screen")
	for i := range g.scrn {
		g.scrn[i] = 255
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
					// TODO: render
					log.Println("********** ")
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
						yscrl := 0
						if g.yscrl {
							yscrl = 1
						}
						xscrl := 0
						if g.xscrl {
							xscrl = 1
						}
						mapbase := g.bgmapbase + ((((g.curline + yscrl) & 0xFF) >> 3) << 5)
						y := (g.curline + yscrl) & 7
						x := xscrl & 7
						t := (xscrl >> 3) & 31

						if g.bgtilebase != 0 {
							tile := g.vram[mapbase + t]
							// TODO?
							//if tile < 128 {
						//		tile = 256 + tile
						//	}
							tilerow := g.tilemap[tile][y]
							for w := SCREEN_WIDTH; w > 0; w -= 1 {
								g.scanrow[SCREEN_WIDTH-x] = tilerow[x]
								g.scrn[linebase + 3] = g.pal.bg[tilerow[x]]
								x += 1
								if x == 8 {
									t = (t + 1) & 31
									x = 0
									tile = g.vram[mapbase + t]
									// TODO?
									// if tile < 128 {
									// 	tile = 256 + tile
									// }
									tilerow = g.tilemap[tile][y]
								}
								linebase += 4
							}
						} else {
							tilerow := g.tilemap[g.vram[mapbase+t]][y]
							for w := SCREEN_WIDTH; w > 0; w -= 1 {
								g.scanrow[SCREEN_WIDTH-x] = tilerow[x]
								g.scrn[linebase + 3] = g.pal.bg[tilerow[x]]
								x += 1
								if x == 8 {
									t = (t + 1) & 31
									x = 0
									tilerow = g.tilemap[g.vram[mapbase+t]][y]
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
							linebase := g.curscan
							for i := range g.od {
								obj := g.odsorted[i]
								if obj.y <= g.curline && obj.y + 8 > g.curline {
									tilerow := g.tilemap[obj.tile][g.curline - obj.y]
									if obj.yflip {
										tilerow = g.tilemap[obj.tile][7 - (g.curline - obj.y)]
									}

									pal := g.pal.obj0
									if obj.palette != 0 {
										pal = g.pal.obj1
									}

									linebase = (g.curline * SCREEN_WIDTH + obj.x) * 4
									// TODO: if/else xflip

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

func (g *gpu) UpdateTile(addr int, value byte) {
	// TODO
}

func (g *gpu) UpdateOAM(addr int, value byte) {
	// TODO
	for i := range g.od {
		g.odsorted[i] = g.od[i]
	}
	// TODO
	// sort.Sort(g.odsorted)
}

func (g gpu) ReadByte(addr int) byte {
	gaddr := addr - 0xFF40
	switch gaddr {
		case 0:
			// TODO
		case 1:
			// TODO
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

func (g *gpu) WriteByte(addr int, value byte) {
	gaddr := addr - 0xFF40
	g.reg[gaddr] = value
	switch gaddr {
		case 0:
			// TODO
		case 2:
			g.yscrl = false
			if value != 0 {
				g.yscrl = true
			}
		case 3:
			g.xscrl = false
			if value != 0 {
				g.xscrl = true
			}
		case 5: g.raster = value
		case 6:
			// TODO: OAM DMA
		case 7:
			// TODO: BG palette mapping

		case 8:
			// TODO: obj0 palette mapping
		case 9:
			// TODO: obj1 palette mapping
	}
}
