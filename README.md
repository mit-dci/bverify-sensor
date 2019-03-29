# bverify-sensor
Small integrated example for a sensor reader that commits to a bverify server

This project consists of two parts:

- [Client](client/) is running on a board that is connected to a sensor. It will witness the sensor readings to the b_verify server. *Any* sensor that can be read from bash will be eligible for this. Once the client receives confirmation of a witnessing, it will output the proof to a folder that can be read in by the web frontend
- [Web](web/) runs on a server and reads from a folder of sensor data. Each sensor has its own subdirectory and in there the files from the client are dumped. It will read all the exported statements and provide a nice webfront to browse them and present a QR code to verify the witnessing with the [mobile app](https://github.com/mit-dci/bverify-mobile)