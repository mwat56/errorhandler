# ErrorHandler

## Purpose

The out-of-the-box `Go` web-server send's _plain text error messages_ whenever an `HTTP` errror occurs. This package provides a simple facility to send whatever `HTML` page you like for error-pages.

## Installation

You can use `Go` to install this package for you:

    go get github.com/mwat56/go-errorhandler

## Usage

This package defines the `TErrorPager` interface which requires just one function or method:

    GetErrorPage(aData []byte, aStatus int) []byte

The arguments are:

* `aData` is the orignal error text.

* `aStatus` is the error number of the actual HTTP error.

If the function's return value is _empty_ then `aData` is sent as is to the remote user, otherwise the function result is sent thus delighting the user with your customised page.

Once you've implemented such a function you call the package's `Wrap()` function:

    Wrap(aHandler http.Handler, aPager TErrorPager) http.Handler

The arguments are:

* `aHandler` is your original HTTP handler which will be wrapped by this packe and used internally.

* `aPager` is the provider of error message pages as discussed above.

In the directory `_demo` there is the file `demo.go` which shows the bare minimum of how to integrate this package with your web-server.

## Licence

    Copyright (C) 2019  M.Watermann, 10247 Berlin, FRG
                All rights reserved
            EMail : <support@mwat.de>

> This program is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 3 of the License, or (at your option) any later version.
>
> This software is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
>
> You should have received a copy of the GNU General Public License along with this program.  If not, see the [GNU General Public License](http://www.gnu.org/licenses/gpl.html) for details.
