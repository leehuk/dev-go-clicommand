// Copyright (C) 2018 Lee H <lee@leeh.uk>
// Licensed under the BSD 2-Clause License as found in LICENSE.txt

package clicommand

// A Handler represents a function to be called when a particular Command
// is selected, or a callback is required.  The Handler function must be
// defined as:
//   func myHandlerFunc(data *Data) error {
//     // your code goes here
//   }
//
// The Data parameter contains a pointer to the selected Command, a map
// of supplied options and an array of the supplied parameters.
//
// If the Handler function encounters an error it should return this as
// an error and it will automatically be sent to stderr.  The Handler
// function should return nil on success.
type Handler func(*Data) (err error)
