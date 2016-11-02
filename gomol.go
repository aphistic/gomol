// Package gomol is the GO Multi-Output Logger, a structured logging library supporting
// multiple outputs at once. Gomol grew from a desire to have a structured logging
// library that could write to any number of outputs while also keeping a small
// in-band footprint.
//
// Gomol has a few basic concepts and most should be familiar to those who have used
// other logging libraries in the past.  There are multiple logging levels and the
// ability to limit the levels that are logged.
//
// In order to provide maximum flexibility gomol has the concept of a base for logging
// functionality, represented by the Base struct.  The Base can have zero or more Logger
// instances added to it.  A Logger is an implementation of a way to display or store
// log messages.  For example, there is a Logger for logging to the console, one for
// logging to Graylog and others. Once Loggers are added to the Base and initialized
// any messages logged using that Base will be sent to all the loggers added.
//
// In most use cases you will probably not need to create your own Base and can just
// use the default one created by gomol on startup.  To use the default Base, simply
// call the functions in the root of the gomol package.
package gomol
