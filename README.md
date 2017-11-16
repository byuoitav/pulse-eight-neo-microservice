# pulse-eight-neo-microservice
[![CircleCI](https://img.shields.io/circleci/project/byuoitav/pulse-eight-neo-microservice.svg)](https://circleci.com/gh/byuoitav/pulse-eight-neo-microservice) [![Apache 2 License](https://img.shields.io/hexpm/l/plug.svg)](https://raw.githubusercontent.com/byuoitav/pulse-eight-neo-microservice/master/LICENSE)

[![View in Swagger](http://jessemillar.github.io/view-in-swagger-button/button.svg)](http://byuoitav.github.io/swagger-ui/?url=https://raw.githubusercontent.com/byuoitav/pulse-eight-neo/master/swagger.json)

## Audio Troubleshooting	
If the audio does not play over a display, try changing the EDID (Extended Display Identification Data) for the port the input display is connected to. 

1. Open the Neo's web UI (just type the Neo's IP address in a browser)
2. Click on ' Cloud connection'.  An option that reads 'Admin' will appear at the bottom.
3. Return to the video switching matrix screen. Each input should have a little settings icon on it.
4. Click on the icon choose an EDID profile
5. Despite the deceptive ballon that reads "settings saved," the device must be rebooted for changes to take effect
6. Poof! It's done
