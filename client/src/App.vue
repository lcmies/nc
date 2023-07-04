<script setup>
import { ref, watch, onMounted } from 'vue'
import dayjs from 'dayjs'
import Reader from './components/reader.vue'

const limits = [10, 20, 50, 100]

const page = ref(1)
const limit = ref(20)
const rows = ref([])
const fname = ref("")

watch(limit, () => {
  getRows()
})

function getRows() {
  fetch(`/api/search?page=${page.value}&limit=${limit.value}`)
    .then((res) => res.json())
    .then((r) => rows.value = r)
}

function prevPage() {
  if (page.value == 1) return
  page.value -= 1
  getRows()
}

function nextPage() {
  page.value += 1
  getRows()
}

onMounted(() => {
  getRows()
})
</script>

<template lang="pug">
nav.navbar.fixed-top.bg-body-tertiary
  .container-fluid
    a.navbar-brand(href='#') NC
    .d-flex
      button.btn(v-if="page > 1" @click='prevPage') ＜
      select.form-select(v-model="limit")
        option(v-for="l in limits" :value="l") {{ l }}
      button.btn(@click='nextPage') ＞
      button.btn.btn-outline-primary(@click='getRows') ↺

main.container
  table.table.align-middle
    thead
      tr
        th Time
        th Host
        th

    tbody
      tr(v-for="row in rows")
        td.font-monospace {{ dayjs.unix(row.start).format("YYYY-MM-DD HH:mm:ss") }}
        td.font-monospace {{ row.host }}
        td
          button.btn.btn-outline-success(@click='fname=row.fname') READ
  
Reader(v-if='fname' :fname='fname' @close='fname=""')
</template>
