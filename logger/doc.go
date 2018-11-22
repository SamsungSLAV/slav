/*
 *  Copyright (c) 2018 Samsung Electronics Co., Ltd All Rights Reserved
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License
 */

/*
Package logger provides logging mechanism for SLAV projects. It was created, as we didn't find a
logger fitting our requirements:

* have syslog-like levels;

* do not panic with highest priority log;

* have structured log entities (to allow automated parsing).

Basics

The usage of logger is really simple:
	if data, err := someFunction(ctx); err != nil{
		logger.WithError(err).WithProperty("context", ctx).Warning("Some warning.")
	}

Below, you can see output from logger with default settings. You can adjust it to your liking (see
Backend section).
	[0.000093] [WAR] [filename.go:11] "Some warning." {context:someFuncInput;error:"error msg";}
	 ^timestamp ^lvl  ^call context    ^message       ^properties            ^error

Where:

* Timestamp - time of log creation - set during log entity processing;

* Level - level of log message entity - defined during call to a logging function;

* CallContext - source code context - set during log entity processing;

* Message - string part of log message entity - defined during call to a logging function;

* Properties - a key-value map of log message entities - defined by calls of: WithProperty,
WithProperties and WithError functions.

Usage

The simplest use-case would be to use default logger (residing in the logger package):
	logger.Warning("Some warning.")

You can also create your instance of logger and customize it:
	// Create new instance of Logger.
	log := logger.NewLogger()

	// Set threshold to omit logs below Notice level.
	log.SetThreshold(logger.NoticeLevel)

	// Register custom backend.
	log.AddBackend("myBackend", Backend{
		Filter:     NewFilterPassAll(),
		Serializer: NewSerializerJSON(),
		Writer:     NewWriterFile(filename, 0444),
	})

	// Set the logger as default.
	logger.SetDefault(log)

Logging entries

There are 8 log levels defined - exactly matching syslog levels. Please refer to documentation of
type Level for consts names.

To log an entity with one of these levels simply use dedicated method of Logger. For every log
level there are 2 methods available: normal and formatted that can be used the same way fmt.Print
and fmt.Printf are:
	logger.Notice("You were notified!")
	logger.Noticef("You were notified about %s!", "usage of logger")

You can also use generic methods: Log, Logf which require level as the first parameter.

Adding more information

To make parsing of logger messages easier, you should use logger.Properties. Properties can be
serialized separately and easily parsed. WithProperty and WithProperties methods insert new or
update existing key-value properties of the log message entity, e.g.
	logger.WithProperty("ID", 17).Info("New object created.")
	logger.WithProperties(logger.Properties{
		name:		"John",
		lastname:	"Covalsky",
	}).Critical("Object deceased.")
There is also a special property for logging errors than can be added by WithError function:
	result, err := doStuff()
	if err != nil {
		logger.WithError(err).Error("Getting things done failed.")
	}

Every log message entity gets CallContext during processing, containing:

* Path - path the source file from which a log was created;

* File - filename of the source file;

* Line - line number;

* Package - package name;

* Type - name of type for which method was defined (in case of methods only);

* Function - name of function or method.

It is up to Serializer used in Backends (described further down this document) which information is
written to logs.

There are situations when some kind of auxiliary helper functions log an error. In such case your
intention is probably to have context of calling the helper function rather than helper function
itself. IncDepth method can be used to change the default call stack depth and get call context of
the caller, e.g.
	...
	15	func exitOnErr(msg string, err error) {
	16		if err != nil {
	17			logger.IncDepth(1).WithError(err).Error(msg + " IncDep(1)")
	18			logger.WithError(err).Error(msg)
	19			os.Exit(1)
	20		}
	21	}
	22
	23	func main() {
	24		obj, err := NewObj()
	25		exitOnErr("Failed to create object.", err)
	26 }
Will give you following (when default settings are used):
	[0.000099] [ERR] [yourapp.go:25] "Failed to create object. IncDep(1)" {error:"error msg";}
	[0.000209] [ERR] [yourapp.go:18] "Failed to create object." {error:"error msg";}

If methods from this paragraph are run on an existing Entry structure, they modify and return it.
If they are run on Logger structure, they create and return new Entry structure with defined
properties.

Processing log messages

Every log message entity is processed after calling one of Log, Logf, Debug, Debugf, Info,
Infof, ... logging functions.

Processing of an entity consist of following steps:

1) Verification of threshold. If it fails, the log entity is dropped.

2) Adding timestamp and call context.

3) Passing an Entry structure to every Backend registered in Logger and continuing processing
in every backend.

Backends

Backends are customizable parts of logger that allow filtering logs, defining the way
they are formatted and choosing the destination where they are finally written.

Logger can have multiple backends registered. Every backend works independently, so e.g. filtering
an Entry by one of them does not affect processing the log entitty in another one. Backends are
identified with name (string), so adding new backend with a name that is already used, will replace
old backend.

Multiple backends with different filters can be used for logging specific entities into additional
files, logs, network locations, etc. For example: all security logs or network logs containing
some special property can be passed to specific files.

Backends can be dynamically added or removed from Logger with following functions:

* AddBackend - add (or replace) a single backend;

* RemoveBackend - remove a single backend;

* RemoveAllBackends - clear all backends collection from Logger.

After removing all backends, you should add at least one, as your logger won't be able to log
anything at all.

Every backend consists of 3 elements:

* Filter - for choosing which entities should be handled by the Backend;

* Serializer - for marshalling Entry structure into []byte;

* Writer - for saving/sending entities.

Filter

Filter's role is to verify if log message entity should be logged by a backend.
It is an interface that requires implementation of a single method:
	Verify(*Entry) (bool, error)

There is a FilterPassAll implementation of Filter implementation. It accepts all log message
entities.

Serializer

Serializer's role is to marshal Entry structure to a slice of bytes, so it can be written
by Writer.
It is an interface that requires implementation of a single method:
	Serialize(*Entry) ([]byte, error)

There are 2 example implementations of this interface:

* SerializerJSON - that uses JSON format for Entry serialization;

* SerializerText - that is intended to produce human-readable from of logs for consoles
or log files.

Both of them are configurable. Please see fields' descriptions of structures defining them
for details.

Writer

Writer's role is to save/send serialized log message entity.
It is an interface that requires implementation of a single method:
	Write(level Level, p []byte) (n int, err error)
which is very similiar to io.Writer interface, but requiring a log level as there are some
destinations (e.g. syslog) that require this information.

There are 3 example implementations of this interface:

* WriterFile - that saves log entities into files;

* WriterStderr - that prints logs to standard error output;

* WriterSyslog - that logs to system logger using log/syslog package.

See their constructors for more customized usage.

*/
package logger
