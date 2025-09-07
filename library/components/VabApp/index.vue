<script lang="ts" setup>
  import { enLocale, zhLocale } from '@/i18n/index'
  import { useSettingsStore } from '@/store/modules/settings'
  const { locale: language } = useI18n()
  const settingsStore = useSettingsStore()
  const { toolboxVisible, toolType, downloadVisible, dialogData, messageVisible, filePath } = storeToRefs(settingsStore)
  const locale = computed(() => (language.value === 'en' ? enLocale : zhLocale))
  const EcsToolboxs = defineAsyncComponent(() => import('@/ecs/toolbox/index.vue'))
  const EcsDownload = defineAsyncComponent(() => import('@/ecs/toolbox/download.vue'))
  const EcsMessage = defineAsyncComponent(() => import('@/ecs/toolbox/message.vue'))
</script>
<template>
  <el-config-provider
    :button="{
      autoInsertSpace: true,
    }"
    :locale="locale"
  >
    <ContextMenu />
    <router-view v-slot="{ Component }">
      <component :is="Component" />
    </router-view>
    <ecs-toolboxs v-if="toolboxVisible" v-model="toolboxVisible" :tool-type="toolType" />
    <ecs-download v-if="downloadVisible" v-model="downloadVisible" :file-name="filePath" :select-alert="dialogData" />
    <ecs-message v-if="messageVisible" v-model="messageVisible" />
  </el-config-provider>
</template>
