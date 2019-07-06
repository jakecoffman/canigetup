<template>
  <div id="app">
    <transition name="fade">
      <div class="loader" v-if="loading"><div></div></div>
    </transition>
    <div class="error" v-if="error === true">
      Error connecting to Pi
    </div>
    <div class="error" v-if="error && error !== true">
      Error: {{error}}
    </div>
    <div>
      Currently <span v-if="On">ON</span><span v-else>OFF</span>
    </div>
    <button @click="On = true; put()">On</button>
    <button @click="On = false; put()">Off</button>
    <button @click="On = !On; put()">Toggle</button>
  </div>
</template>

<script>

export default {
  data() {
    return {
      On: false,
      Schedules: [],
      error: null,
      loading: true
    }
  },
  methods: {
    async put() {
      this.loading = true
      const response = await fetch("/api/state", {method: "PUT", body: this.state()})
      await this.handleResponse(response)
    },
    async handleResponse(response) {
      let json
      try {
        json = await response.json()
      } catch (e) {
        // eslint-disable-next-line no-console
        console.error(e)
      }
      if (response.status !== 200) {
        this.error = true
        if (json) {
          this.error = json
        }
      } else {
        this.On = json.On
        this.Schedules = json.Schedules
      }
      this.loading = false
    },
    state() {
      return JSON.stringify({
        On: this.On,
        Schedules: this.Schedules
      })
    }
  },
  async created() {
    this.loading = true
    const response = await fetch("/api/state")
    await this.handleResponse(response)
  }
}
</script>

<style>
  html, body {
    margin: 0;
    padding: 0;
  }
  #app {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  .error {
    background: red;
    color: white;
    width: 100%;
    padding: 1rem;
  }
  button {
    width: 80%;
    background: white;
    border: 1px black solid;
    padding: 1rem;
    margin: 1rem;
  }

  .loader {
    position: fixed;
    left: 50%;
    top: 50%;
    transform: rotate(45deg) scale(5) translate(-50%, -50%);
    transform-origin: 40px 32px;
  }
  .loader div {
    top: 23px;
    left: 19px;
    position: absolute;
    width: 26px;
    height: 26px;
    background: red;
    animation: loader 1.2s infinite cubic-bezier(0.215, 0.61, 0.355, 1);
  }
  .loader div:after,
  .loader div:before {
    content: " ";
    position: absolute;
    display: block;
    width: 26px;
    height: 26px;
    background: red;
  }
  .loader div:before {
    left: -17px;
    border-radius: 50% 0 0 50%;
  }
  .loader div:after {
    top: -17px;
    border-radius: 50% 50% 0 0;
  }
  @keyframes loader {
    0% {
      transform: scale(0.95);
    }
    5% {
      transform: scale(1.1);
    }
    39% {
      transform: scale(0.85);
    }
    45% {
      transform: scale(1);
    }
    60% {
      transform: scale(0.95);
    }
    100% {
      transform: scale(0.9);
    }
  }

  .fade-enter-active, .fade-leave-active {
    transition: opacity 1s;
  }
  .fade-enter, .fade-leave-to {
    opacity: 0;
  }

</style>
