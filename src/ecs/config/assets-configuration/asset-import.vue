<script lang="ts">
  export default {
    name: 'AssetImport', //
  }
</script>

<script setup lang="ts">
  import AssetsKnown from '@/ecs/assets/components/assets-known.vue'
  import AssetsIpDetail from '@/ecs/assets//components/ip-detail.vue'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import { siteGetAllPageApi } from '~/src/api-ecs/site'
  import { getAllServerIpLabelApi } from '~/src/api-ecs/tagLib'
  import { AssetsIp_Detail } from '../../assets/type'

  const allSites = ref<any>([])

  const allTags = ref<any>([])

  const assetsKnownRef = ref()

  const ip_assets_visible = ref(false)

  const ip_assets_isEdit = ref(false)

  const current_row = ref()
  // 获取关联站点
  const getAllData = async () => {
    const { data } = await siteGetAllPageApi()
    allSites.value = data
    const res = await getAllServerIpLabelApi()
    allTags.value = res.data
  }

  onMounted(() => {
    getAllData
  })

  provide('reflashData', () => {
    try {
      setTimeout(() => {
        assetsKnownRef.value.getData()
      }, 0)
    } catch {
      console.log('error')
    }
  })

  provide(AssetsIp_Detail, {
    detailVisible: ip_assets_visible,
    isEdit: ip_assets_isEdit,
    currentRow: current_row,
  })

  function test(done: any) {
    ip_assets_isEdit.value = false
    done()
  }
</script>

<template>
  <div class="asset-import-container">
    <assets-known ref="assetsKnownRef" :all-sites="allSites" :all-tags="allTags" />
    <!-- IP资产详情 -->
    <vab-dialog
      v-model="ip_assets_visible"
      :before-close="test"
      destroy-on-close
      show-fullscreen
      title="IP资产详情"
      width="1175px"
    >
      <assets-ip-detail />
    </vab-dialog>
  </div>
</template>

<style scoped lang="scss"></style>
