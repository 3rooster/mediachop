# mediachop
## what dose mediachop do
mediachop is a simple server which can receive HLS or Dash streaming, and delivering.

# usage
## rules about streaming in chop
### concept
- event & stream
 We assume that every stream belongs to an event.
 An event can have many streams to publish.
- publish url & play url
 When stream publishing , it can only using HTTP METHOD PUT, which ffmpeg and other tools supported.
 Publish url will be in template `http://{$ip}:{$port}/{$event}/{$stream}/{$index.(m3u8|mpd)}`
 Play url will be in template `http://{$ip}:{$port}/{$event}/{$stream}/{$index.(m3u8|mpd)}`
