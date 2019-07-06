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
    <h3>Status</h3>
    <button class="main-btn" @click="On = true; put()" :class="{selected: On}">On</button>
    <button class="main-btn" @click="On = false; put()" :class="{selected: !On}">Off</button>
    <button class="main-btn" @click="On = !On; put()">Toggle</button>

    <h3>Schedules</h3>
    <div v-if="!Schedules || !Schedules.length">
      No schedules
    </div>
    <div v-for="(schedule, index) of Schedules" :key="schedule.At" class="mb-1">
      <span>{{schedule.At}}</span>
      <span class="mx-1">Turn <span v-if="schedule.On">on</span><span v-else>off</span></span>
      <button @click="remove(index)">Remove</button>
    </div>

    <h3>Add a new schedule</h3>
    <div class="add-schedule">
      <input type="time" v-model="newScheduleAt">
      <button @click="newScheduleOn = true" :class="{selected: newScheduleOn}">On</button>
      <button @click="newScheduleOn = false"  :class="{selected: newScheduleOn === false}">Off</button>
      <button @click="add()">Add</button>
    </div>
  </div>
</template>

<script>

export default {
  data() {
    return {
      On: false,
      Schedules: [],
      error: null,
      loading: true,
      newScheduleAt: null,
      newScheduleOn: null
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
    },
    async add() {
      if (this.newScheduleAt === null) {
        return
      }
      if (this.newScheduleOn === null) {
        return
      }

      this.loading = true
      if (this.Schedules === null) {
        this.Schedules = []
      }
      this.Schedules.push({At: this.newScheduleAt, On: this.newScheduleOn})
      const response = await fetch("/api/state", {method: "PUT", body: this.state()})
      await this.handleResponse(response)
      this.newScheduleAt = null
      this.newScheduleOn = null
    },
    async remove(index) {
      this.loading = true
      this.Schedules.splice(index, 1)
      const response = await fetch("/api/state", {method: "PUT", body: this.state()})
      await this.handleResponse(response)
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
    background: #000;
    color: #a3a3a3;
    font-size: 16pt;
  }
  #app {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  .error {
    background: #9c0000;
    width: 100%;
    padding: 1rem;
  }
  button {
    background: black;
    color: #a3a3a3;
    border: 1px #a3a3a3 solid;
    font-size: 16pt;
  }
  input {
    background: black;
    color: #a3a3a3;
    border: 1px #a3a3a3 solid;
    font-size: 16pt;
  }

  .main-btn {
    width: 80%;
    padding: 1rem;
    margin-bottom: 1rem;
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

  .status {
    margin-top: 2rem;
  }

  .add-schedule {
    display: flex;
  }

  .selected {
    background: #a3a3a3;
    color: black;
  }

  .mx-1 {
    margin-left: 1rem;
    margin-right: 1rem;
  }

  .mb-1 {
    margin-bottom: 0.5rem;
  }
</style>
