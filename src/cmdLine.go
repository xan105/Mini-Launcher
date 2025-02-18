/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "os"
  "flag"
)

type Args struct {
  Help        bool
  DryRun      bool
  ConfigPath  string
}

func parseArgs() (Args) {
  var args Args
  
  flag.BoolVar(&args.Help, "help", false, "Show list of all arguments.")
  flag.BoolVar(&args.DryRun, "dry-run", false, "Program will exit before starting the executable.")
  flag.StringVar(&args.ConfigPath, "config", "launcher.json", "File path to the config file to use.")
  flag.Parse()
  
  if args.Help {
    alert(
      "Launcher",
      "--config filePath\n" +
      "File path to the config file to use.\n" +
      "--dry-run\n" +
      "Program will exit before starting the executable.\n" +
      "\n" +
      "--help\n" +
      "Show list of all arguments\n",
    )
    os.Exit(0)
  }
  return args
}

