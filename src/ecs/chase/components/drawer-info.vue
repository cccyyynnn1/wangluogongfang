<script lang="ts">
  export default {
    name: 'DrawerAlert', // 告警列表
  }
</script>

<script setup lang="ts">
  const props = defineProps<{
    alldata: any
  }>()

  const emits = defineEmits<{
    (e: 'on-close'): void
  }>()

  const baseInfo = ref()

  /**
   *
   * @param groupId 资产组id  1:流量指纹 2:规则指纹 3:三方指纹
   */
  const getFingerprintLabel = (groupId: 1 | 2 | 3) => {
    const group = baseInfo.value?.labels?.find((item: any) => item.groupId === groupId)
    return group ? group.labelList : []
  }

  const formatHostVoList = (data: any) => {
    let str = ''
    data.forEach((item: any) => {
      str += item.apis.join(',')
    })
    return str || '-'
  }

  watch(
    () => props.alldata,
    () => {
      baseInfo.value = props.alldata
    },
    { immediate: true }
  )
</script>

<template>
  <div class="drawer-info sketchpad-drawer">
    <div class="sketchpad-drawer-top">
      <div class="title">{{ baseInfo?.indexIp }}</div>
      <el-icon @click="emits('on-close')"><Close /></el-icon>
    </div>
    <div class="sketchpad-drawer-top-tip">
      <span>来源：</span>
      <span>{{ baseInfo?.datasource || baseInfo?.netTypeStr || '-' }}</span>
      <!-- <span style="margin-left: 30px">告警：</span>
      <span>23条</span> -->
    </div>
    <div class="sketchpad-drawer-content">
      <div class="my-title"><span>基本信息</span></div>
      <div class="tip-item">
        <div class="left">业务名称</div>
        <div class="right">{{ baseInfo?.businessName || '-' }}</div>
      </div>
      <div class="tip-item">
        <div class="left">应用名称</div>
        <div class="right">{{ baseInfo?.appName || '-' }}</div>
      </div>
      <div class="tip-item">
        <div class="left">应用类型</div>
        <div class="right">{{ baseInfo?.appType || '-' }}</div>
      </div>
      <div class="tip-item">
        <div class="left">数据中心</div>
        <div class="right">{{ baseInfo?.datasource || '-' }}</div>
      </div>
      <div class="tip-item">
        <div class="left">业务IPV4</div>
        <div class="right">{{ baseInfo?.ipv4 || '-' }}</div>
      </div>
      <div class="tip-item">
        <div class="left">业务IPV6</div>
        <div class="right">{{ baseInfo?.ipv6 || '-' }}</div>
      </div>
      <div class="tip-item">
        <div class="left">网络分区</div>
        <div class="right">{{ baseInfo?.netPartation || '-' }}</div>
      </div>
      <div class="tip-item">
        <div class="left">系统版本</div>
        <div class="right">{{ baseInfo?.sysVersion || '-' }}</div>
      </div>
      <div class="tip-item">
        <div class="left">中间库版本</div>
        <div class="right">{{ baseInfo?.middlewareVersion || '-' }}</div>
      </div>
      <div class="tip-item">
        <div class="left">数据库版本</div>
        <div class="right">{{ baseInfo?.dbVersion || '-' }}</div>
      </div>
      <div class="tip-item">
        <div class="left">联系方式</div>
        <div class="right">{{ baseInfo?.phone || '-' }}</div>
      </div>
      <div v-if="baseInfo.labels?.length > 0" class="my-title"><span>资产标签</span></div>
      <div class="labelBox foldBox">
        <div v-for="(item, index) in baseInfo.labels" :key="index" class="table-label">
          <span>{{ item.groupName }}</span>
        </div>
      </div>
      <div v-if="baseInfo.labels?.length > 0" class="my-title"><span>流量指纹</span></div>
      <div class="labelBox foldBox">
        <div v-for="(item, index) in getFingerprintLabel(1)" :key="index" class="table-label">
          <span>{{ item.labelName }}</span>
        </div>
      </div>
      <div class="my-title">
        <span>资产信息</span>
        <div class="tip-item">
          <div class="left">HOST信息</div>
          <div class="right">{{ baseInfo?.host || '-' }}</div>
        </div>
        <div class="tip-item">
          <div class="left">服务端口</div>
          <div class="right">{{ (baseInfo?.portList?.length > 0 && baseInfo?.portList?.join(',')) || '-' }}</div>
        </div>
        <div class="tip-item">
          <div class="left">业务端口</div>
          <div class="right">{{ '-' }}</div>
        </div>
        <div class="tip-item">
          <div class="left">API接口</div>
          <div class="right">
            {{ (baseInfo.hostVoList?.length > 0 && formatHostVoList(baseInfo.hostVoList)) || '-' }}
          </div>
        </div>
      </div>
      <!-- <div class="my-title"><span>IOC情报信息</span></div>
      <div class="tip-item">
        <div class="left">业务名称</div>
        <div class="right">public/plugins/aquaqpublic/plugins/aquaqpublic/plugins/aquaq</div>
      </div>
      <div class="tip-item">
        <div class="left">应用名称</div>
        <div class="right">网络特洛伊木马</div>
      </div> -->
    </div>
  </div>
</template>

<style scoped lang="scss">
  .drawer-info {
    left: 0;
    border-right: 1px solid #e9e6f9;
    border-left: 0 solid #fff !important;
    .sketchpad-drawer-top-tip {
      color: #5e6169;
      margin-top: 5px;
      line-height: 20px;
    }
    .sketchpad-drawer-content {
      width: 100%;
      height: calc(100% - 60px);
      overflow-y: auto;
      &::-webkit-scrollbar {
        width: 0;
        height: 0;
      }
    }

    .labelBox {
      display: flex;
      &.foldBox {
        display: block;
        .table-label {
          display: inline-block;
        }
      }
      .table-label {
        padding: 0 10px;
        height: 30px;
        border-radius: 4px;
        line-height: 30px;
        background: #f5f4ff;
        margin-right: 10px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }
</style>
