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
Package logger provides logging mechanism for SLAV projects.

It is a completely new logger package, as there we didn't find any ready one
that meets our requirements:
	* have syslog-like levels;
	* do not panic with highest priority log;
	* have structured log entities (to allow automated parsing).

Basics

The use of logger is pretty simple.

First you will need a Logger structure that will: filter, serialize and write log entities for you.
You can create your own Logger with NewLogger function ...
	log = logger.NewLogger()
	log.Warning("Some warning.")
or even easier use the default one, that is already created fo you.
	logger.Warning("Some warning.")
Then you will be able to call Logger's methods for logging messages with proper level, e.g. Warning
as you can see in the above examples.

Logger structure can be setup by defining threshold and different backends.
Log entities can be also more compex containing properties or information about an error.
All of these enhanced features will be described below.

Levels

There are 8 log levels defined - exactly matching syslog levels:
	* EmergLevel - used when system is unusable;
	* AlertLevel - used when action must be taken immediately;
	* CritLevel - used when critical conditions occur;
	* ErrLevel - used when error conditions occur;
	* WarningLevel - used when warning conditions occur;
	* NoticeLevel - used when normal, but significant, conditions occur;
	* InfoLevel - used for logging informational message;
	* DebugLevel - used for logging debug-level message.
To log an entity with one of these levels simply use dedicated method of Logger or existing
log message entity, e.g.
	logger.Notice("You were notified!")
	logger.Noticef("You were notified about %s!", "usage of logger")
For every log level there are 2 methods available: normal and formatted that can be used
the same way fmt.Print and fmt.Printf are.

You can also use generic methods: Log, Logf which require level as the first parameter.

Notice that all of functions described in this paragraph complete creation of log message entity
and causes Logger structure to process the entity. It will be more clear after reading about
Entry structure and creation of structured logs.

Threshold

Every Logger structure has a Threshold value that can be set with SetThreshold method
and read with Threshold method.
All log message entities that are processed by Logger are verified, if they pass threshold test.
It checks if log message level is higher than threshold, e.g.
If Logger's threshold is set to Warning level then Debug, Info and Notice messages
won't be logged, because their levels is lower, but Warinig, Error, Critical, Alert
and Emergency logs will be processed fine.

Entry

Every log message entity is kept in Entry structure.
It contains following elements:
	* Level - level of log message entity defined by completing function;
	* Message - string part of log message entity defined by completing function;
	* Properties - a key-value map of log message entities;
	* Timestamp - time of log creation - set during log entity processing;
	* CallContext - source code context - set during log entity processing.
The structure should not be modified directly by the user.
There are several methods to operate on this structure.

Properties

WithProperty and WithProperties methods insert new or update existing key-value properties
of the log message entity, e.g.
	logger.WithProperty("ID", 17).Info("New object created.")
	logger.WithProperties(logger.Properties{
		name:		"John",
		lastname:	"Covalsky",
	}).Critical("Object deceased.")
There is also a special property for logging errors than can be added by WithError function:
	result, err := doStuff()
	if err != nil {
		logger.WithErrror(err).Error("Getting things done failed.")
	}

If these methods are run on existing Entry structure, they modify and return
the Entry structure.

If these methods are run on Logger structure, they create and return new Entry structure with
defined properties.

CallContext and Timestamp are collected automatically when completed log message entity
is processed.

CallContext

Every log message entity get source code context during processing. The context contains:
	* Path - path the source file from which a log was created;
	* File - filename of the source file;
	* Line - line number;
	* Package - package name;
	* Type - name of type for which method was defined (in case of methods only);
	* Function - name of function or method.
The full context is available during processing of the log message entity, but which information
is written to logs depends on Logger (Serializers used in Backends to be more precise).

So both collection of call context and writing it is pretty automated. However there are
situations when some kind of auxiliary helper functions log an error. In such case
the code author's intention is probably to log place where helper function is used
instead of specific line in helper function. In such case an IncDepth method can be used to
change the default call stack depth and get call context of the caller, e.g.
	func helper(check Stuff) {
		if !check.OK() {
			logger.IncDep(1).Error("Stuff check failed v1.")
			logger.Error("Stuff check failed v2.")	// call context of v2 will point here.
		}
	}
	func stuffMaker() {
		stuff := NewStuff()
		helper(stuff)	// call context of v1 will point here.
	}

Processing log messages

Every log message entity is processed after calling one of Log, Logf, Debug, Debugf, Info,
Infof, ... logging functions.

Processing of an entity consist of following steps:
 
1) Verification of threshold test. If test fails, the log entity is dropped.

2) Completing Entry structure by adding timestamp and call context information.

3) Passing an Entry structure to every Backend registered in Logger and continuing processing
in every backend.

Backends

Backends are the customizable parts of the logger that allow filtering logs, defining the way
they are formatted and choosing the destination where they are finally written.

Logger can have multiple backends registered. Every backend works independently, so e.g.
filtering an Entry by one of them does not affect processing the log entitty in another one.
Backends are identified with name (string), so adding a backend with name already used, will
replace the old one.

Multiple backends with different filters can be used for logging specific entities into additional
files, logs, network locations, etc. for example all security logs or network logs containing
some special property can be passed to specific files.

Backends can be dynamically added or removed from Logger with following functions:
	* AddBackend - add (or replace) a single backend;
	* RemoveBackend - remove a single backend;
	* RemoveAllBackends - clear all backends collection from Logger.
After removing all backends, you should probably add at least one backend, as it will make
no sense to have no backends, as you logger won't be able to log anything at all.

Every backend consists of 3 elements:
	* Filter - for choosing which entities should be handled by the Backend;
	* Serializer - for marshaling Entry structure into []byte;
	* Writer - for saving/sending entities.

Filter

Filter's role is to verify if log message entity should be logged by a backend.
It is an interface that requires implementation of a single method:
	Verify(*Entry) (bool, error)

There is a trivial Filter implementation: FilterPassAll, that accepts all log message entities.

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
destinations, e.g. syslog that require this information.

There are 3 example implementations of this interface:
	* WriterFile - that saves log entities into files;
	* WriterStderr - that prints logs to standard error output;
	* WriterSyslog - that logs to system logger using log/syslog package.
See their constructors for more customized usage.

Default

To make usage of logger package easier, the default Logger is already initialized and
ready to be used in simple scenarios. It has threshold set to Info level, one backend
with FilterPassAll, SerializerText and WriterStderr.

All global package functions operate on the default Logger.

The default logger can be set to customized solution with SetDefault function, e.g.
	func() initializeMyLogger(filename string, level string) error {
		lvl, err := StringToLevel(level)
		if err != nil {
			return err
		}
		myLogger := NewLogger()
		myLogger.SetThreshold(lvl)
		myLogger.AddBackend("myBackend", Backend{
			Filter:     NewFilterPassAll(),
			Serializer: NewSerializerJSON(),
			Writer:     NewWriterFile(filename, 0444),
		})
		logger.SetDefault(myLogger)
		return nil
	}
*/
package logger
