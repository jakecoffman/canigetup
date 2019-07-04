<template>
  <div id="app">
    <button @click="on()">On</button>
    <button @click="off()">Off</button>
    <button @click="toggle()">Toggle</button>
  </div>
</template>

<script>

export default {
  methods: {
    async on() {
      const response = await fetch("/api/on", {method: "POST"})
      const json = await response.json()
      console.log(json)
    },
    async off() {
      const response = await fetch("/api/off", {method: "POST"})
      const json = await response.json()
      console.log(json)
    },
    async toggle() {
      const response = await fetch("/api/toggle", {method: "POST"})
      const json = await response.json()
      console.log(json)
    }
  }
}
</script>

<style>
  button {
    background: white;
    border: 1px black solid;
    padding: 1rem;
    margin: 1rem;
  }
</style>
