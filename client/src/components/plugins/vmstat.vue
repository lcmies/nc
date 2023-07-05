<script setup>
import { toRefs, ref, computed } from 'vue'
import Chart from 'chart.js/auto'
import { Line } from 'vue-chartjs'

const props = defineProps({ content: String })
const { content } = toRefs(props)

const points = ref("60")

function chartOptions(title) {
  return {
    responsive: true,
    aspectRatio: 3,
    maintainAspectRatio: true,
    animation: false,
    elements: {
      point: { radius: 0 },
      line: { borderWidth: 1 }
    },
    plugins: {
      legend: { display: false },
      title: { display: true, text: title }
    }
  }
}

function padArray(arr, size) {
  return [...Array(size - arr.length)].map(() => null).concat(arr)
}

const vmstatData = computed(() => {
  let cpu = []
  let freemem = []

  content.value.split(/\n/).forEach((line) => {
    line = line.trim()

    if (line.match(/^procs/)) return
    if (line.match(/^r/)) return

    const data = line.split(/\s+/)
    cpu.push(100 - parseInt(data[14]))
    freemem.push(parseInt(data[3]) + parseInt(data[4]) + parseInt(data[5]))
  })

  return {
    cpu: cpu,
    freemem: freemem
  }
})

const cpuData = computed(() => {
  return {
    labels: [...Array(parseInt(points.value))].map(() => ""),
    datasets: [
      {data: padArray(vmstatData.value.cpu.slice(-parseInt(points.value)), parseInt(points.value))}
    ]
  }
})

const memData = computed(() => {
  return {
    labels: [...Array(parseInt(points.value))].map(() => ""),
    datasets: [
      {data: padArray(vmstatData.value.freemem.slice(-parseInt(points.value)), parseInt(points.value))}
    ]
  }
})

</script>

<template lang="pug">
.d-flex.align-items-center
  small.pe-2 Points
  select.form-select.form-select-sm(v-model="points")
    option(value="30") 30
    option(value="60") 60
    option(value="120") 120
    option(value="300") 300
    option(value="600") 600
    option(value="1200") 1200
Line(:data="cpuData" :options="chartOptions('CPU')")
Line(:data="memData" :options="chartOptions('Free MEM')")

</template>