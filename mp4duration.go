package main

import (
  "encoding/binary"
  "fmt"
  "io"
  "os"
  "strings"
  "github.com/urfave/cli/v2"
  "github.com/brothertoad/btu"
  "github.com/brothertoad/mp4atom"
)

func main() {
  app := &cli.App {
    Name: "mp4duration",
    Usage: "dumps the duration of mp4 or m4a file(s)",
    Action: mp4Duration,
  }
  app.Run(os.Args)
}

func mp4Duration(c *cli.Context) error {
  args := c.Args().Slice()
  if len(args) == 0 {
    cli.ShowAppHelp(c)
    return nil
  }
  // Find the longest arg and create a padding string, to make the ouput a bit prettier.
  maxLength := findMaxLength(args)
  padding := strings.Repeat(" ", maxLength)
  for _, arg := range(args) {
    fmt.Printf("%s:%s %s\n", arg, padding[0:(maxLength-len(arg))], getDuration(arg))
  }
  return nil
}

func getDuration(path string) string {
  file := btu.OpenFile(path)
  defer file.Close()
  mp4atom.FindAtomPath(file, "moov/mvhd")
  // Skip 12 bytes, then read the timeUnits and time, both of which are 32 bit integers
  // in network byte order.
  file.Seek(12, io.SeekCurrent)
  b := make([]byte, 4)
  file.Read(b)
  timeUnit := binary.BigEndian.Uint32(b)
  file.Read(b)
  units := binary.BigEndian.Uint32(b)
  return formatDuration(units, timeUnit)
}

func formatDuration(units, timeUnit uint32) string {
  // Get the number of seconds, rounded.
  totalSecs := (units + timeUnit/2) / timeUnit
  mins := totalSecs / 60
  secs := totalSecs % 60
  hours := mins / 60
  mins = mins % 60
  return fmt.Sprintf("%02d:%02d:%02d", hours, mins, secs)
}

func findMaxLength(args []string) int {
  max := 0
  for _, s := range(args) {
    if len(s) > max {
      max = len(s)
    }
  }
  return max
}
