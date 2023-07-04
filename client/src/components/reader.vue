<script setup>
import { toRefs, ref, onMounted, onUnmounted, watch } from 'vue'

const props = defineProps({ fname: String })
const { fname } = toRefs(props)
const content = ref("")
const pos = ref(0)
const autorefresh = ref(false)

let intervalId = 0

function fetchContent() {
  fetch(`/data/${fname.value}`, {
    headers: {
      "Range": `bytes=${pos.value}-`
    }
  }).then((res) => {
    if (res.status < 300) {
      pos.value = parseInt(res.headers.get("Content-Range").match(/\/([0-9]+)$/)[1])
      res.text().then((txt) => content.value += txt)
    }
  })
}

function refreshContent() {
  if (!autorefresh.value) return
  fetchContent()
}

onMounted(() => {
  fetchContent()
  intervalId = setInterval(refreshContent, 1000)
})

onUnmounted(() => {
  autorefresh.value = false
  content.value = ""
  pos.value = 0
  clearInterval(intervalId)
})


</script>

<template lang="pug">
.modal.d-block
  .modal-dialog.modal-dialog-scrollable.modal-xl
    .modal-content
      .modal-header
        span
          input.btn-check#autofresh(type='checkbox' :value='autorefresh')
          label.btn.btn-outline-primary(for='autofresh' @click='autorefresh=!autorefresh') ↺
        button.btn-close(@click='$emit("close")')
      .modal-body
        pre {{ content }}

</template>

<style scoped lang="scss">
.modal {
    background-color: rgba(255, 255, 255, 0.25);
}
</style>