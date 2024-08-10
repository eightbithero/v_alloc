## Memory Allocation Analysis Tool

This utility analyzes the distribution of RAM into occupied addresses. It outputs the results in a graphical representation, highlighting potential collisions (e.g., when a definition falls into a previously allocated range). 

### Features
- **Graphical Representation**: Visualizes the memory allocation to help identify and understand the distribution of occupied addresses.
- **Collision Detection**: Identifies and highlights possible collisions where new definitions overlap with previously allocated ranges.
- **Source Code Analysis**: Utilizes the source code, understanding the ca65 dialect, to perform its analysis.
- **VICE Labels Support**: Supports labels files in the VICE format for accurate and detailed memory mapping.

This tool is essential for developers working with memory allocation in systems using the ca65 assembler, providing a clear and detailed view of memory usage and potential issues.


### Example

In the following example, the tool analyzes the memory allocation of a project and outputs a representation of the occupied addresses.
For testing purposes, the tool uses a simple project https://github.com/jrs0ul/nes_survival by K. Lotužis
after building the project you will obtain labels file in VICE format, the tool can be run with the following command

```bash
go run . <path_to_the_labels>/Cold\ \&\ Starving\(USA\,\ Japan\).labels.txt <path_to_the_sourcecode>/src/ | sort | uniq 
```