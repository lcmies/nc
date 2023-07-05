<script setup>
import { toRefs, ref, onMounted, onUnmounted, computed } from 'vue'

const props = defineProps({ fname: String })
const { fname } = toRefs(props)
const rawContent = ref("")
const pos = ref(0)
const autorefresh = ref(false)

const content = computed(() => {
  let text = rawContent.value.replaceAll(/\r\n/g, "\n")
  text = text.replaceAll(/[^\r\n]*\r/g, "")
  return text
})

let intervalId = 0

function fetchContent() {
  return fetch(`/data/${fname.value}`, {
    headers: {
      "Range": `bytes=${pos.value}-`
    }
  }).then((res) => {
    if (res.status < 300) {
      pos.value = parseInt(res.headers.get("Content-Range").match(/\/([0-9]+)$/)[1])
      return res.text().then((txt) => rawContent.value += txt)
    }
  })
}

function refreshContent() {
  if (!autorefresh.value) return
  fetchContent().then(() => {
    const elm = document.querySelector(".modal-body")
    if (!elm) return
    elm.scrollTo({
      top: elm.scrollHeight,
      left: 0,
      behavior: "instant"
    })
  })
}

onMounted(() => {
  fetchContent()
  intervalId = setInterval(refreshContent, 1000)
})

onUnmounted(() => {
  autorefresh.value = false
  rawContent.value = ""
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