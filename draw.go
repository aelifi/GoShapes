package main

import (
  "fmt"
  //"math"
  "os"
  "errors"
)

//HEADER
//Program written by Aleksandr Elifirenko aelifi2, 12/1/2023
//Project 3: (Go)
//
//This program has two intefaces implemented: geometry for shapes and screen for drawing the shapes. The 
//shapes implemented are triangle, rectangle and circle. Colors for the shapes are stored in a map. So when
//the shapes are drawn, you can write out the ppm file the cordinates of the shapes. 
//

// Constants for in a map
var red Color = 0
var green Color = 1
var blue Color = 2
var yellow Color = 3
var orange Color = 4
var purple Color = 5
var brown Color = 6
var black Color = 7
var white Color= 8

//RGB struct for each matrix index
type RGB struct {
  R int
  G int
  B int
}

//Color variable for in map call
type Color int

//Map for all the color values
var colorsMap = map[Color]RGB {
  0 : {255, 0, 0},
  1 : {0, 255, 0},
  2 : {0, 0, 255},
  3 : {255, 255, 0},
  4 : {255, 164, 0},
  5 : {128, 0, 128},
  6 : {165, 42, 42},
  7 : {0, 0, 0},
  8 : {255, 255, 255},
}

//Geometry interface implementation with the functions in it
type geometry interface {
  draw(scn screen) (err error)
  shape() (s string)
}

//Point struct for x and y
type Point struct {
  x, y int
}

//Rectangle struct for lower and upper points with its color
type Rectangle struct {
  ll, ur Point
  c Color
}

//Circle struct with center point, radius and color
type Circle struct {
  cp Point
  r int
  c Color
}

//Triangle struct with three angle points and color
type Triangle struct {
  pt0, pt1, pt2 Point
  c Color
}

//Variables for each error
var outOfBoundsErr = errors.New("geometry of bounds")
var colorUnknownErr = errors.New("color unknown")

//Function checks if the color is right
func colorUnknown(c Color) bool {
  if _, ok := colorsMap[c]; ok { //Check if the color is in the map
    return false
  }
  return true
}

//Function to check out of bounds
func outOfBounds(p Point, scn screen) bool {
  maxX, maxY := scn.getMaxXY() //Retrieve each cord
  if(p.x < 0 || p.x >= maxX || p.y < 0 || p.y >= maxY) { //Check the bounds of cords
    return true
  }
  return false
}

//next three functions return the correct representation string of each shape
func (T Triangle) shape() string {
  return "Triangle"
}

func (R Rectangle) shape() string {
  return "Triangle"
}

func (C Circle) shape() string {
  return "Triangle"
}

//  https://gabrielgambetta.com/computer-graphics-from-scratch/07-filled-triangles.html
func interpolate (l0, d0, l1, d1 int) (values []int) {
  a := float64(d1 - d0) / float64(l1 - l0)
  d  := float64(d0)

  count := l1-l0+1
  for ; count>0; count-- {
    values = append(values, int(d))
    d = d+a
  }
  return
}

//  https://gabrielgambetta.com/computer-graphics-from-scratch/07-filled-triangles.html
func (tri Triangle) draw(scn screen) (err error) {
  if outOfBounds(tri.pt0,scn) || outOfBounds(tri.pt1,scn)  || outOfBounds(tri.pt2,scn) {
    return outOfBoundsErr
  }
  if colorUnknown(tri.c) {
    return colorUnknownErr
  }

  y0 := tri.pt0.y
  y1 := tri.pt1.y
  y2 := tri.pt2.y

  // Sort the points so that y0 <= y1 <= y2
  if y1 < y0 { tri.pt1, tri.pt0 = tri.pt0, tri.pt1 }
  if y2 < y0 { tri.pt2, tri.pt0 = tri.pt0, tri.pt2 }
  if y2 < y1 { tri.pt2, tri.pt1 = tri.pt1, tri.pt2 }

  x0,y0,x1,y1,x2,y2 := tri.pt0.x, tri.pt0.y, tri.pt1.x, tri.pt1.y, tri.pt2.x, tri.pt2.y

  x01 := interpolate(y0, x0, y1, x1)
  x12 := interpolate(y1, x1, y2, x2)
  x02 := interpolate(y0, x0, y2, x2)

  // Concatenate the short sides

  x012 := append(x01[:len(x01)-1],  x12...)

  // Determine which is left and which is right
  var x_left, x_right []int
  m := len(x012) / 2
  if x02[m] < x012[m] {
    x_left = x02
    x_right = x012
  } else {
    x_left = x012
    x_right = x02
  }

  // Draw the horizontal segments
  for y := y0; y<= y2; y++  {
    for x := x_left[y - y0]; x <=x_right[y - y0]; x++ {
      scn.drawPixel(x, y, tri.c)
    }
  }
  return
}

//Function to draw the rectangle
func (rect Rectangle) draw(scn screen) (err error) {
  if colorUnknown(rect.c) { //Check color
    return colorUnknownErr
  }

  if outOfBounds(rect.ll, scn) || outOfBounds(rect.ur, scn) { //Check bounds
    return outOfBoundsErr
  }

  for x := rect.ll.x; x <= rect.ur.x; x++ { //Draw horizontal lines
    for y := rect.ll.y; y <= rect.ur.y; y++ { //Draw vertical lines
      scn.drawPixel(x, y, rect.c)
    }
  }

  return nil
}

//Function to draw circle
func (cir Circle) draw(scn screen) (err error) {
  maxX, maxY := scn.getMaxXY() //Get max x and y coordinates

  for x := 0; x < maxX; x++ { //Loop though x
    for y := 0; y < maxY; y++ { //Loop though y
      distance := (x-cir.cp.x)*(x-cir.cp.x) + (y-cir.cp.y)*(y-cir.cp.y) //Calculate the distance

      if distance <= cir.r*cir.r { //if distance is withing the r^2 aka withing the area then draw the pixel
        scn.drawPixel(x, y, cir.c)
      }
    }
  }
  return
}

//Screen interface with all its functions
type screen interface {
  initialize(maxX, maxY int)
  getMaxXY() (maxX, maxY int)
  drawPixel(x, y int, c Color) (err error)
  getPixel(x, y int) (c Color, err error)
  clearScreen()
  screenShot(f string) (err error)
}

//Display struct with the max cords and color matrix
type Display struct {
  maxX, maxY int
  matrix [][]Color
}

var display Display //variable for the display with Display type to test in main

//initialize function
func (display *Display) initialize(maxX, maxY int) {
  display.maxX = maxX //set maxX
  display.maxY = maxY //set maxY

  //Two loops to make the screen with x as rows and y as columns
  display.matrix = make([][]Color, maxX)
  for i := range display.matrix {
    display.matrix[i] = make([]Color, maxY)
    for j := range display.matrix[i] {
      display.matrix[i][j] = white
    }
  }
}

//Just gets the maxX and maxY
func (display *Display) getMaxXY() (int, int) {
  return display.maxX, display.maxY
}

//fucntion Draws a pixel
func (display *Display) drawPixel(x, y int, c Color) error {
  if x < 0 || x >= display.maxX || y < 0 || y >= display.maxY { //Check the bounds
    return errors.New("Out of bounds")
  }
  display.matrix[x][y] = c //set the matrix aka the pixel in the display to a color
  return nil
}

//Function returns a pixel
func (display *Display) getPixel(x, y int) (Color, error) {
  if x < 0 || x >= display.maxX || y < 0 || y >= display.maxY { //CHeck the bounds
    return -1, errors.New("Out of bounds")
  }
  return display.matrix[y][x], nil //return the pixel aka the matrix space with given cords and the nil aka no error
}

//Functions to clear the screen
func (display *Display) clearScreen() {
  //Loop the whole display matrix
  for i := range display.matrix {
    for j := range display.matrix[i] {
      display.matrix[i][j] = white //set every matrix index aka pixel to white
    }
  }
}

//Function to print the diplay
func (display *Display) screenShot(f string) (err error) {
  file, err := os.Create(f + ".ppm") //Create the file
  if err != nil { //Check for error creating
    return err
  }

  defer file.Close()

  header := fmt.Sprintf("P3\n%d %d\n255\n", display.maxX, display.maxY) //Create the PPM header maxX aka max width, maxY aka max height and 255 max color value
  file.WriteString(header) //Write the header string

  //Loop through the whole matrix
  for i := range display.matrix {
    for j := range display.matrix[i] {
      col, err := display.getPixel(j, i) //Get the pixel
      if err != nil {
        return err
      }
      rgb := colorsMap[col] //Get the RGB value of the pixel
      file.WriteString(fmt.Sprintf("%d %d %d ", rgb.R, rgb.G, rgb.B)) //Print the RGB value to the file
    }
    file.WriteString("\n") //Make a new line
  }
  return nil
}