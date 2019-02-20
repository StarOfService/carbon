## Carbon installation
1. Download the latest version for your platform from the project releases page: https://github.com/StarOfService/carbon/releases
2. Rename the binary file:
    * to `carbon` for Unix systems
    * to `carbon.exe` for Windows systems
3. Put the binary file to the folder from PATH environment variable:
    * by default `/usr/local/bin` is recommended for Unix systems. Another option is to create a folder specific for your user and add it to PATH environment variable at `.bash_profile` or `.bashrc`
    * Windows doesn't provide a default folder for user executable files. The simplest (but not the best) option is to locate the `carbon.exe` at C:/Windows/System32. Or you can create a custom folder for executable files and add it to PATH environment variable. You can find a PATH environment variable extension process for different versions of Windows here: https://www.computerhope.com/issues/ch000549.htm

We understand that the installation process is complicated (especially for Windows users), and we're going to improve this in the future.