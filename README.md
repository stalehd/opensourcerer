# The Open Sourcerer

Utility to apply open source license to a Go source repository. No cleverness
so if you decide to add the license multiple times it will be messy. It does
check for a LICENSE file before modifying the source code but that's it.

Currently only the Apache 2.0 license is supported.

* Adds a LICENSE file at the root
* Adds license header to all of the source files

