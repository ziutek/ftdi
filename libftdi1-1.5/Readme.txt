
Date of compilation: 19-July-2020
Xiaofan Chen <xiaofanc AT gmail dot com)

Please use libftdi mailing list for support. I do not
reply to technical questions through private email.

To use libftdi1 under Windows, you need to use libusb 
and supported drivers.

The copyright information about libftdi1 and libusb are 
inside the copyright directory.

The html_doc directory is the doxygen generated HTML document 
for libftdi1-1.5. libusb is at the following URL.
  http://libusb.info/

The bin directory contains the dlls and example programs for 
libftdi1-1.5 and libusb 1.0.23 Windows. 

The source code libftdi1-1.5 and libusb-1.0.23 release for 
the build are also included in the src directory.

Tools used to build this package:
  MSYS2 up to date as of 19-July-2020
  CMake 3.18.0: http://www.cmake.org/

MSYS2 has CMake as well but I had some problems with it so
I do not use it. Python bindings are not built and you can
try pyftdi or pylibftdi.
