# bverify-sensor client

The client is very simple. Once you've built the binary, ensure that the `sensor.sh` file is modified to execute whatever's needed to read out your sensor and print its result to `STDOUT`. The sensor client will execute this command, read the output, and witness the output to the b_verify server. Once the server has confirmed the reading's been witnessed, it will repeat.

