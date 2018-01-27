// Copyright (C) 2018 Lee H <lee@leeh.uk>
// Licensed under the BSD 2-Clause License as found in LICENSE.txt

/*
Package clicommand provides CLI applications with subcommand/api-style interfaces and option/parameter handling

The clicommand library makes the creation of Go CLI applications using a subcommand
interface easier.  The subcommand interface is structured as a parent/child tree so
the application can mimic an api, with edges of the tree running custom Handler
functions and the tree providing a structured way of grouping commands, attaching
option arguments and finding additional parameters.

The tree itself operates under a set of rules.  Each parent Command object within the
tree may have any number of children.  Each child Command object within the tree
has a single parent.  A child Command object within the tree can have its own
children, except when it has a Handler function.

This allows building a CLI application which can mimic an API, e.g.:
  ./clicommand                         // parent, has children
  ./clicommand http                    // child of clicommand, has children itself
  ./clicommand http get => Handler()   // child of clicommand->http, calls Handler() when
                                       // run.  Cannot have children.
  ./clicommand http post => Handler()  // child of clicommand->http, calls Handler() when
                                       // run.  Cannot have children.
*/
package clicommand
