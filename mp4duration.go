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

var printFileNames bool = false
var printTotal = false
var printMillis = false

func main() {
  app := &cli.App {
    Name: "mp4duration",
    Usage: "dumps the duration of mp4 or m4a file(s)",
    UsageText: "mp4duration [options] file...",
    Action: mp4Duration,
    HideHelp: true,
    Flags: []cli.Flag {
      &cli.BoolFlag {
        Name: "with-filename",
        Usage: "print the file names (default if more than one file)",
        Aliases: []string{"H"},
      },
      &cli.BoolFlag {
        Name: "no-filename",
        Usage: "don't print the file names (default if only one file)",
        Aliases: []string{"h"},
      },
      &cli.BoolFlag {
        Name: "millis",
        Usage: "include milliseconds",
        Aliases: []string{"m"},
      },
      &cli.BoolFlag {
        Name: "total",
        Usage: "show time as total, rather than 00:00:00",
        Aliases: []string{"t"},
      },
    },
  }
  app.Run(os.Args)
}

func mp4Duration(c *cli.Context) error {
  args := c.Args().Slice()
  if len(args) == 0 {
    cli.ShowAppHelp(c)
    return nil
  }
  printFileNames = setPrintFileNames(c, args)
  printTotal = c.Bool("total")
  printMillis = c.Bool("millis")
  // Find the longest arg and create a padding string, to make the ouput a bit prettier.
  maxLength := findMaxLength(args)
  padding := strings.Repeat(" ", maxLength)
  for _, arg := range(args) {
    if printFileNames {
      fmt.Printf("%s:%s %s\n", arg, padding[0:(maxLength-len(arg))], getDuration(arg))
    } else {
      fmt.Printf("%s\n", getDuration(arg))
    }
  }
  return nil
}

func setPrintFileNames(c *cli.Context, args []string) bool {
  if c.Bool("with-filename") {
    return true
  }
  if c.Bool("no-filename") {
    return false
  }
  return len(args) > 1
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
  // Get the number of units, rounded.
  totalMillis := ((units + timeUnit/2) * 1000)/ timeUnit
  if printTotal && printMillis {
    return fmt.Sprintf("%7d", totalMillis)
  }
  totalSecs := totalMillis / 1000
  millis := totalMillis % 1000
  mins := totalSecs / 60
  secs := totalSecs % 60
  hours := mins / 60
  mins = mins % 60
  if printMillis {
    return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, mins, secs, millis)
  }
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
