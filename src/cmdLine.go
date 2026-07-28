/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "os"
  "io"
  "flag"
  "strings"
)

type Args struct {
  Help        bool
  DryRun      bool
  Wait        bool
  ConfigPath  string
}

func parseArgs() (Args) {
  var args Args
  
  options := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
  options.SetOutput(io.Discard)
  
  options.BoolVar(&args.Help, "help", false, "Show list of all arguments.")
  options.BoolVar(&args.DryRun, "dry-run", false, "Program will exit before starting the executable.")
  options.BoolVar(&args.Wait, "wait", false, "Program will wait for the executable to terminate before exiting.")
  options.StringVar(&args.ConfigPath, "config", "launcher.json", "File path to the config file to use.")

  if err := options.Parse(os.Args[1:]); err != nil {
    var output strings.Builder
    output.WriteString("Failed to parse command-line arguments: " + err.Error())
    output.WriteString("\n\nUsage:\n\n")
    options.SetOutput(&output)
    options.PrintDefaults()
    panic("Launcher", output.String())
  }
  
  if args.Help {
    var output strings.Builder
    output.WriteString("Usage:\n\n")
    options.SetOutput(&output)
    options.PrintDefaults()
    informAndExit("Launcher", output.String())
  }
  
  return args
}
