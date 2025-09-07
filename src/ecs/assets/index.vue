<script setup lang="ts">
  // import AssetsKnown from './components/assets-known.vue'

  // import AssetsUnknown from './components/assets-unknown.vue'

  import NetworkPartition from './components/network-partition.vue'

  import AssetsSite from './components/assets-site.vue'

  import AssetsUNSite from './components/assets-unsite.vue'

  import ServiceChain from './components/service-chain.vue'

  import AssetVisits from './components/asset-visits.vue'
  import AssetOverview from './components/asset-overview/index.vue'
  import VabDialog from '@/plugins/VabDialog/index.vue'

  import AssetsIpDetail from './components/ip-detail.vue'

  import { AssetsType, AssetsIp_Detail } from './type'

  import { getAssetsSetApi, updateAssetsSetApi } from '~/src/api-ecs/assets'

  import { siteGetAllPageApi } from '@/api-ecs/site'

  import { getAllServerIpLabelApi } from '~/src/api-ecs/tagLib'

  const route = useRoute()

  const $baseMessage: any = inject('$baseMessage')

  const activeName = ref<AssetsType>('asset-overview')
  const assetsKnownRef = ref()

  const assetsUnKnownRef = ref()

  const current_row = ref()

  const ip_assets_visible = ref(false)

  const ip_assets_isEdit = ref(false)

  const ip_assets_intranetScope_visible = ref(false)

  const allSites = ref<any>([])

  const allTags = ref<any>([])

  const assetsSetData = reactive({
    range: '',
  })

  function test(done: any) {
    ip_assets_isEdit.value = false
    done()
  }

  // 获取资产网段
  const getAssetsSet = async () => {
    ip_assets_intranetScope_visible.value = true
    const { data } = await getAssetsSetApi()
    assetsSetData.range = data
  }

  // 修改资产网段
  const saveData = async () => {
    const { msg } = await updateAssetsSetApi({ range: assetsSetData.range })
    $baseMessage(msg, 'success', 'vab-hey-message-success')
    ip_assets_intranetScope_visible.value = false
  }

  provide(AssetsIp_Detail, {
    detailVisible: ip_assets_visible,
    isEdit: ip_assets_isEdit,
    currentRow: current_row,
  })

  provide('reflashData', () => {
    try {
      setTimeout(() => {
        assetsKnownRef.value.getData()
        assetsUnKnownRef.value.getData()
      }, 0)
    } catch {
      console.log('error')
    }
  })
  // 获取关联站点
  const getAllData = async () => {
    const { data } = await siteGetAllPageApi()
    allSites.value = data
    const res = await getAllServerIpLabelApi()
    allTags.value = res.data
  }

  onMounted(() => {
    getAllData()
    const type = route.query.params as AssetsType
    const isOverview = ['asset-overview', 'known', 'unknow'].includes(type)
    activeName.value = isOverview ? 'asset-overview' : type
  })
</script>

<script lang="ts">
  export default {
    name: 'Assets',
  }
</script>

<template>
  <div class="assets-container">
    <div>
      <!-- v-model="activeName"  -->
      <el-tabs v-model="activeName" class="demo-tabs">
        <el-tab-pane label="资产总览" name="asset-overview">
          <asset-overview />
        </el-tab-pane>
        <el-tab-pane label="资产访问" lazy name="assetVisits">
          <AssetVisits :all-sites="allSites" :all-tags="allTags" />
        </el-tab-pane>
        <!-- <el-tab-pane label="已知IP资产" lazy name="known">
          <assets-known ref="assetsKnownRef" :all-sites="allSites" :all-tags="allTags" />
        </el-tab-pane>
        <el-tab-pane label="未知IP资产" lazy name="unKnown">
          <assets-unknown ref="assetsUnKnownRef" :all-sites="allSites" :all-tags="allTags" />
        </el-tab-pane> -->
        <el-tab-pane label="已知站点" lazy name="site">
          <assets-site v-if="activeName === 'site'" :key="0" :module="1" />
        </el-tab-pane>
        <el-tab-pane label="未知站点" lazy name="nosite">
          <AssetsUNSite v-if="activeName === 'nosite'" :key="1" />
        </el-tab-pane>
        <el-tab-pane label="网络分区" lazy name="network">
          <network-partition />
        </el-tab-pane>
        <!-- <el-tab-pane label="业务链梳理" lazy name="service-chain">
          <service-chain />
        </el-tab-pane> -->
      </el-tabs>
      <!-- 已知IP资产内网范围 -->
      <el-button v-if="activeName === 'known'" class="assets-setting" link @click="getAssetsSet">
        <el-icon>
          <Setting />
        </el-icon>
        &nbsp;设置
      </el-button>
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
      <!-- 内网IP范围 -->
      <vab-dialog v-model="ip_assets_intranetScope_visible" destroy-on-close title="内网IP范围" width="475px">
        <el-input
          v-model="assetsSetData.range"
          :autosize="{ minRows: 4 }"
          placeholder="Please input"
          resize="none"
          type="textarea"
        />
        <template #footer>
          <el-button type="primary" @click="saveData">保存</el-button>
          <el-button @click="ip_assets_intranetScope_visible = false">取消</el-button>
        </template>
      </vab-dialog>
    </div>
  </div>
</template>

<style scoped lang="scss">
  :deep() {
    .el-scrollbar {
      height: calc(100vh - 410px);
    }
  }
  .assets-container {
    position: relative;
    height: calc(100vh - 20px);
    overflow: hidden;

    .assets-setting {
      position: absolute;
      top: 20px;
      right: 20px;
    }
  }
</style>
