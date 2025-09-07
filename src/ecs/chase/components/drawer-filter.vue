<script lang="ts">
  export default {
    name: 'DrawerFilter', // 筛选
  }
</script>

<script setup lang="ts">
  const props = defineProps<{
    alldata: any
  }>()
  const emits = defineEmits<{
    (e: 'on-close'): void
    (e: 'on-change-event', data: any): void
  }>()

  const radioValue = ref('疑似沦陷')

  const count = ref('0/0')

  const filterData = ref()

  onMounted(() => {
    filterData.value = props.alldata
  })

  const handleChange = (mode: string) => {
    emits('on-change-event', mode)
  }

  watch(
    () => filterData.value,
    () => {
      const keys = Object.keys(filterData.value)
      let index = 0
      keys.forEach((item) => {
        if (filterData.value[item]) {
          index++
        }
      })
      count.value = `${index}/${keys.length}`
      console.log(count.value)
    },
    { deep: true }
  )
</script>

<template>
  <div class="drawer-info sketchpad-drawer">
    <div class="sketchpad-drawer-top">
      <div>
        <span style="font-size: 18px; color: #1e1842; margin-right: 6px">筛选</span>
        <span style="font-size: 14px; color: #4d3fa2">已选：{{ count }}</span>
      </div>
      <el-icon @click="emits('on-close')"><Close /></el-icon>
    </div>
    <!-- <div class="my-title"><span>场景筛选</span></div>
    <div class="warp-item">
      <div class="left-text">威胁专题</div>
      <div class="right-img"><img alt="" :src="require('@/assets/chase/close.svg')" /></div>
    </div>
    <div class="warp-item">
      <div class="left-text">IOC信息</div>
      <div class="right-img"><img alt="" :src="require('@/assets/chase/open.svg')" /></div>
    </div>
    <el-divider />
    <div class="my-title"><span>威胁等级</span></div>
    <div class="warp-item">
      <div class="left-text">危急</div>
      <div class="right-img"><img alt="" :src="require('@/assets/chase/close.svg')" /></div>
    </div>
    <div class="warp-item">
      <div class="left-text">高危</div>
      <div class="right-img"><img alt="" :src="require('@/assets/chase/open.svg')" /></div>
    </div>
    <div class="warp-item">
      <div class="left-text">中危</div>
      <div class="right-img"><img alt="" :src="require('@/assets/chase/close.svg')" /></div>
    </div>
    <div class="warp-item">
      <div class="left-text">低危</div>
      <div class="right-img"><img alt="" :src="require('@/assets/chase/open.svg')" /></div>
    </div>
    <el-divider /> -->
    <div class="my-title"><span>信息展示</span></div>
    <div class="warp-item">
      <div class="left-text">IP地址</div>
      <div v-if="filterData?.ip" class="right-img" @click="handleChange('ip')">
        <img alt="" :src="require('@/assets/chase/open.svg')" />
      </div>
      <div v-else class="right-img" @click="handleChange('ip')">
        <img alt="" :src="require('@/assets/chase/close.svg')" />
      </div>
    </div>
    <div class="warp-item">
      <div class="left-text">资产名称</div>
      <div v-if="filterData?.assetName" class="right-img" @click="handleChange('assetName')">
        <img alt="" :src="require('@/assets/chase/open.svg')" />
      </div>
      <div v-else class="right-img" @click="handleChange('assetName')">
        <img alt="" :src="require('@/assets/chase/close.svg')" />
      </div>
    </div>
    <!-- <el-divider />
    <div class="my-title">
      <span>外网攻击源IP</span>
      <el-radio-group v-model="radioValue" class="my-title-radio">
        <el-radio label="全部" />
        <el-radio label="疑似沦陷" />
      </el-radio-group>
    </div>
    <div class="warp-item">
      <el-icon><CaretRight /></el-icon>
      <div class="left-text">192.168.0.0/16（23）</div>
      <div class="right-img"><img alt="" :src="require('@/assets/chase/close.svg')" /></div>
    </div>
    <div class="warp-item">
      <el-icon><CaretBottom /></el-icon>
      <div class="left-text">192.168.1.0/24（23）</div>
      <div class="right-img"><img alt="" :src="require('@/assets/chase/open.svg')" /></div>
    </div>
    <div class="warp-item">
      <div style="width: 14px"></div>
      <div class="left-text">192.168.1.0/24（23）</div>
      <div class="right-img"><img alt="" :src="require('@/assets/chase/close.svg')" /></div>
    </div>
    <div class="warp-item">
      <div style="width: 14px"></div>
      <div class="left-text">192.168.1.0/24（71）</div>
      <div class="right-img"><img alt="" :src="require('@/assets/chase/open.svg')" /></div>
    </div>
    <div class="warp-item">
      <el-icon><CaretRight /></el-icon>
      <div class="left-text">192.168.0.0/16（23）</div>
      <div class="right-img"><img alt="" :src="require('@/assets/chase/open.svg')" /></div>
    </div> -->
  </div>
</template>

<style scoped lang="scss">
  :deep() {
    .el-divider--horizontal {
      margin: 10px 0 !important;
    }

    .my-title-radio {
      float: right;
    }
    .el-radio {
      height: 24px;
      margin-left: 15px;
      margin-right: 0;
    }

    .left-text {
      width: 180px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      word-break: break-all;
      word-wrap: break-word;
      font-size: 14px;
      color: #918da5;
    }
    .right-img {
      width: calc(100% - 180px);
      white-space: normal;
      word-break: break-all;
      word-wrap: break-word;
      font-size: 14px;
      color: #2b2742;
    }

    img {
      width: 20px;
      cursor: pointer;
      float: right;
    }
  }
</style>
