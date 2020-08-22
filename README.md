# Kast
An open source streaming stick with a few tricks up its sleeves

alternatively,

An extensible open-source daemon for displaying arbitrary media on the big screen.

## Goals / Guiding Principles
* To achieve a user experience on par with the Google Chromecast
* To preserve the privacy and respect the attention of the user
* To provide an extensible platform for developing lightweight big screen experiences
  * Note: The "Chromecast Use-case" is considered the first class use-case

## Guiding Principles
* Less UI is better UI

## Features
* HTTP API that allows playing of any youtube-dl supported URLs
* Random pictures when idle

## Design
In its current implementation, Kast acts as a mediator between multiple processes (modules) which share exclusive access to the screen. The HTTP API allows local clients to invoke functions from each module, which may decide to pre-empt the current module (the default of which being the slideshow module). Once the running module is terminated, the new module is run. When any module terminates without being pre-empted, the slideshow module is automatically loaded. In addition to the slideshow module, the media module allows any youtube-dl supported media to be streamed.
