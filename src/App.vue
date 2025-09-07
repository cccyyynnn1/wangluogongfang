<template>
  <vab-app />
</template>
<script setup lang="ts">
  let toeknPulse = undefined
  const toeknPulseHandle = () => {
    toeknPulse = setTimeout(async () => {
      await fetch('/v3/ecsPlatform/public/version')
      toeknPulseHandle()
    }, 1000 * 60)
  }
  toeknPulseHandle()
  onUnmounted(() => {
    clearTimeout(toeknPulse!)
  })
</script>
