<script>
  let audioContext = new AudioContext();
  let audioBuffer = null;
  let audioSource = null;

  let ws = new WebSocket('ws://localhost:4444/stream.mp3');

  ws.binaryType = 'arraybuffer';

  ws.onmessage = function(event) {
    // Create an audio buffer from the received ArrayBuffer
    audioContext.decodeAudioData(event.data, function(buffer) {
      audioBuffer = buffer;
    });
  };

  function playMusic() {
    // Create a new AudioBufferSourceNode and connect it to the AudioContext
    audioSource = audioContext.createBufferSource();
    audioSource.buffer = audioBuffer;
    audioSource.connect(audioContext.destination);

    // Start playing the audio buffer
    audioSource.start();
  }

  function stopMusic() {
    // Stop playing the audio buffer
    audioSource.stop();
  }
</script>

<button on:click={playMusic}>Play Music</button>
<button on:click={stopMusic}>Stop Music</button>