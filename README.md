coin is a Golang repo. 
to run, you only need Golang 1.23

coin's current implementation is a program you can run, and connect multiple terminals to with the command 'telnet localhost 8080'.
you may create your own user name, and a menu will allow you to create or join chat rooms. 
within a chat room, you can get a stream of messages from other clients attached to that specific chat room.

planned features:
give a chat room a password, only clients with the password may enter.
host this with AWS, so that users can remotely connect.
log errors.
log user connections & messages once this chat room closes.
