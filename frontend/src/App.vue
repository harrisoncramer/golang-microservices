<template>
  <div class="test broker-test">
    <button @click="handleTestBroker">Test Broker</button>
    <div>
      Response: {{ brokerResponse }}
    </div>
    <div>
      Error: {{ brokerErr ?? "No error" }}
    </div>
    <div v-if="brokerResponseCount">
      API Calls: {{ brokerResponseCount }}
    </div>
  </div>
</template>

<style scoped>
.test {
  gap: 1em;
}
.broker-test {
  display: flex;
  align-items: center;
  justify-content: start;
}
</style>

<script setup lang="ts">

import { ref } from 'vue'

const brokerResponse = ref<string>()
const brokerErr = ref<string>()
const brokerResponseCount = ref<number>(0)

async function handleTestBroker () {
  const body = { method: 'POST' }

  try {
    brokerResponseCount.value += 1
    const response = await fetch("http://localhost:8080", body)
    const json = await response.json()
    brokerResponse.value = json

    if(json.error) {
      brokerErr.value = json.message
    } else {
      brokerResponse.value = json.message
    }
  } catch (err) {
    brokerErr.value = 'Could not reach broker'
  }
}

</script>
