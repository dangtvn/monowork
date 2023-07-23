<script lang="ts">
  import bg from './assets/darius-anton-coffee-01.gif';

  let currentSong = 'Loading...';
</script>

<div class="container">
  <div class="background" style="background-image: url({bg});" />
  <!-- <div class="overlay" style="background-image: url({overlay});"></div> -->
  <!-- <div class="song-name">{{ currentSong }}</div> -->
  <div class="song-name">[Click on the screen to play]</div>
</div>

<style>
  .container {
    width: 100%;
    height: 100%;
    position: relative;
  }

  .background {
    width: 100%;
    height: 100%;
    position: absolute;
    background: transparent no-repeat center center;
    background-size: cover;
  }

  .background {
    opacity: 0.25;
    z-index: 1;
  }

  .song-name {
    position: absolute;
    bottom: 20px;
    left: 20px;
    color: #fff;
    text-shadow: 0 0 7px #fff, 0 0 10px #fff, 0 0 21px #fff, 0 0 42px #0fa,
      0 0 82px #0fa, 0 0 92px #0fa, 0 0 102px #0fa, 0 0 151px #0fa;
    font-family: monospace;
    font-size: 16px;
  }
</style>
