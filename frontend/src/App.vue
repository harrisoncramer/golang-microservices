<template>
  <div class="test">
    <div class="buttons">
      <button @click="handleTestBroker">Test Broker</button>
      <button @click="handleTestAuth">Test Auth</button>
      <button @click="handleTestLogger">Test Logger</button>
      <button @click="handleTestMail">Test Mail**</button>
    </div>
    <div class="content" v-if="!loading">
      <div v-if="s">
        <b>Sent:</b> <pre>{{ s }}</pre>
      </div>
      <div v-if="r">
        <b>Response:</b> <pre>{{ r }}</pre>
      </div>
      <div v-if="e">
        <b>Error:</b> <pre>{{ e }}</pre>
      </div>
    </div>
    <div v-else>
      Loading...
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const r = ref<string>()
const s = ref<string>()
const e = ref<string>()
const loading = ref<boolean>(false)

async function reset () {
  r.value = ''
  s.value = ''
  e.value = ''
  loading.value = false
}

async function handleTestBroker () {
  reset()
  loading.value = true
  const body = { method: 'POST' }

  s.value = JSON.stringify(body, undefined, 4)

  try {
    /* Normally all requests to the broker go through /handle, but for 
    testing to see that it's up we can use the /test route */
    const response = await fetch("http://localhost:8080/test", body)
    const json = await response.json()
    r.value = JSON.stringify(json, undefined, 4)
  } catch (err) {
    e.value = JSON.stringify(err, Object.getOwnPropertyNames(err))
  } finally {
    loading.value = false
  }
}

async function handleTestLogger () {
  reset()
  loading.value = true
  const payload = { action: 'log', log: { name: 'event', data: 'Some log!'} }
  const headers = new Headers()
  headers.append('Content-Type', 'application/json')

  const body = {
    method: 'POST',
    body: JSON.stringify(payload),
    headers,
  }

  s.value = JSON.stringify(body, undefined, 4)

  try {
    const response = await fetch("http://localhost:8080/handle", body)
    const json = await response.json()
    r.value = JSON.stringify(json, undefined, 4)
  } catch (err) {
    e.value = JSON.stringify(err, Object.getOwnPropertyNames(err))
  } finally {
    loading.value = false
  }
}

async function handleTestAuth () {
  reset()
  loading.value = true
  const payload = { action: 'auth', auth: { email: 'admin@example.com', password: 'verysecret'} }
  const headers = new Headers()
  headers.append('Content-Type', 'application/json')

  const body = {
    method: 'POST',
    body: JSON.stringify(payload),
    headers,
  }

  s.value = JSON.stringify(body, null, 4)

  try {
    const response = await fetch("http://localhost:8080/handle", body)
    const json = await response.json()
    r.value = JSON.stringify(json, null, 4)
  } catch (err) {
    e.value = JSON.stringify(err, Object.getOwnPropertyNames(err))
  } finally {
    loading.value = false
  }
}

async function handleTestMail () {
  reset()
  loading.value = true
  const payload = { action: 'mail', mail: { from: 'me@example.com', to: 'you@there.com', subject: 'Test subject', message: 'Hello world!' } }
  const headers = new Headers()
  headers.append('Content-Type', 'application/json')

  const body = {
    method: 'POST',
    body: JSON.stringify(payload),
    headers,
  }

  s.value = JSON.stringify(body, undefined, 4)

  try {
    const response = await fetch("http://localhost:8080/handle", body)
    const json = await response.json()
    r.value = JSON.stringify(json, undefined, 4)
  } catch (err) {
    e.value = JSON.stringify(err, Object.getOwnPropertyNames(err))
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>

.test {
  padding: 32px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: start;
}

.buttons {
  display: flex;
  justify-content: start;
  width: 100%;
  gap: 1em;
  padding: 16px;
  margin-bottom: 36px;
}

.content { 
  width: 100%;
}
</style>
