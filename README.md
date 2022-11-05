# ClipAssist

note: this project integrates with the system clipboard and the notifications
system of your device, and as a result, requires CGO.

This tool serves as a framework for operating on selected text.
It provides functionality to run in the background of your system and listen in
on your clipboard operations.
When a piece of text is copied that matches a specified regex, a function can be
called against that piece of text.

For example, if you are copying a unix millisecond timestamp, it can
automatically send you a desktop notification with the timestamp represented in
a more human-readable fashion.
See the examples folder for a working example of this.

A sample systemd unit file to be placed into `/etc/systemd/user/clipassist.service`
is provided.
Be sure to change the path to the binary to match your home folder.
