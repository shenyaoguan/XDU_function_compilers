package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"strconv"

	"compilers/lexer"
	"compilers/parser"
)

type Drawing struct {
	origin        image.Point
	scale         image.Point
	rotation      float64
	width, height int
}

func NewDrawing(width, height int) *Drawing {
	return &Drawing{
		origin:   image.Point{X: width / 2, Y: height / 2}, // Initial center point
		scale:    image.Point{X: 1, Y: 1},                  // Initial scale 1:1
		rotation: 0,                                        // Initial rotation angle 0
		width:    width,
		height:   height,
	}
}

// Set origin
func (d *Drawing) SetOrigin(x, y int) {
	d.origin = image.Point{X: x, Y: y}
}

// Set scale
func (d *Drawing) SetScale(x, y int) {
	d.scale = image.Point{X: x, Y: y}
}

// Set rotation
func (d *Drawing) SetRotation(angle float64) {
	d.rotation = angle
}

// Rotate point
func (d *Drawing) RotatePoint(p image.Point) image.Point {
	// Convert to radians
	angle := d.rotation * math.Pi / 180
	sin, cos := math.Sin(angle), math.Cos(angle)

	// Rotate coordinates
	newX := cos*float64(p.X-d.origin.X) - sin*float64(p.Y-d.origin.Y) + float64(d.origin.X)
	newY := sin*float64(p.X-d.origin.X) + cos*float64(p.Y-d.origin.Y) + float64(d.origin.Y)

	return image.Point{X: int(newX), Y: int(newY)}
}

// Draw the image
func (d *Drawing) Draw() {
	// Create a blank image
	img := image.NewRGBA(image.Rect(0, 0, d.width, d.height))

	// Set background color to white
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Define the start and end points of the line
	start := image.Point{X: d.origin.X - 50, Y: d.origin.Y - 50}
	end := image.Point{X: d.origin.X + 50, Y: d.origin.Y + 50}

	// Apply scale to the points
	start.X = d.origin.X + (start.X-d.origin.X)*d.scale.X
	start.Y = d.origin.Y + (start.Y-d.origin.Y)*d.scale.Y
	end.X = d.origin.X + (end.X-d.origin.X)*d.scale.X
	end.Y = d.origin.Y + (end.Y-d.origin.Y)*d.scale.Y

	// Apply rotation to the points
	start = d.RotatePoint(start)
	end = d.RotatePoint(end)

	// Draw the line between the start and end points
	drawLine(img, start, end, color.Black)

	// Create and save the output file
	outFile, err := os.Create("output.png")
	if err != nil {
		fmt.Println("Error creating image file:", err)
		return
	}
	defer outFile.Close()

	// Output the image to the file
	if err := png.Encode(outFile, img); err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}

	fmt.Println("Drawing saved as output.png")
}

// drawLine draws a line between two points on the image
func drawLine(img *image.RGBA, p1, p2 image.Point, col color.Color) {
	dx := abs(p2.X - p1.X)
	dy := abs(p2.Y - p1.Y)
	sx := -1
	if p1.X < p2.X {
		sx = 1
	}
	sy := -1
	if p1.Y < p2.Y {
		sy = 1
	}
	err := dx - dy

	for {
		img.Set(p1.X, p1.Y, col)
		if p1.X == p2.X && p1.Y == p2.Y {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			p1.X += sx
		}
		if e2 < dx {
			err += dx
			p1.Y += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	// Assume we receive the program code from the command line or other input source
	program := "ORIGIN IS (100, 200); SCALE IS (2, 2); ROT IS 45;"

	// Initialize the lexer
	l := lexer.New(program)
	p := parser.New(l)

	// Parse the program
	statements := p.ParseProgram()

	// Create a drawing object
	drawing := NewDrawing(500, 500)

	// Execute the corresponding operations based on the parsed statements
	for _, stmt := range statements {
		switch stmt := stmt.(type) {
		case *parser.OriginStatement:
			// Set origin
			x, _ := strconv.Atoi(stmt.X.(*parser.ConstantExpression).Value)
			y, _ := strconv.Atoi(stmt.Y.(*parser.ConstantExpression).Value)
			drawing.SetOrigin(x, y)

		case *parser.ScaleStatement:
			// Set scale
			x, _ := strconv.Atoi(stmt.X.(*parser.ConstantExpression).Value)
			y, _ := strconv.Atoi(stmt.Y.(*parser.ConstantExpression).Value)
			drawing.SetScale(x, y)

		case *parser.RotStatement:
			// Set rotation
			angle, _ := strconv.ParseFloat(stmt.Angle.(*parser.ConstantExpression).Value, 64)
			drawing.SetRotation(angle)
		}
	}

	// Draw the image
	drawing.Draw()
}
