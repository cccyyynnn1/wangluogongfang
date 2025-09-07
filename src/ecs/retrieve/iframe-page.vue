<script lang="ts">
  export default {
    name: 'IframePage', //
  }
</script>

<script setup lang="ts">
  import { getLogin29443Api } from '~/src/api-ecs/public'

  const jsessionIds = ref()

  let token: string | null = null

  const loading = ref(false)

  const iframeRef = ref()

  const login29443Api = async () => {
    loading.value = true
    try {
      const { data } = await getLogin29443Api()
      if (data) {
        token = data.probeToken
        renderIframe()
      }
    } catch (error) {
      console.error(error)
    }
  }

  const renderIframe = () => {
    localStorage.setItem('Authorization', token || '')
    jsessionIds.value = true
  }

  const handleLoad = () => {
    loading.value = false
    const iframe = iframeRef.value
    if (iframe) {
      const hash = iframe.contentWindow?.location.hash
      if (hash === '#/login') {
        loading.value = false
        jsessionIds.value = false
        setTimeout(() => {
          renderIframe()
        }, 0)
      }
    }
  }

  onMounted(() => {
    login29443Api()
  })
</script>

<template>
  <div v-loading="loading" class="iframe-container">
    <iframe
      v-if="jsessionIds"
      ref="iframeRef"
      frameborder="0"
      :src="`/proxy/#/probeAnalysisMenu`"
      style="width: 100%; height: 100%; border: none"
      @load="handleLoad"
    ></iframe>
  </div>
</template>

<style scoped lang="scss">
  .iframe-container {
    overflow: hidden;
    width: 100%;
    height: 100%;
    display: flex;
  }
</style>
