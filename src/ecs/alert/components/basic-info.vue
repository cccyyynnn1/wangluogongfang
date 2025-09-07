<script setup lang="ts">
  import { formatNstime, formatTime } from '@/utils/time'
  import { levelKey } from '../data/index'
  import { downloadLogPacket } from '@/utils/download'
  import { useSettingsStore } from '@/store/modules/settings'
  import AlertDetailLog from './alert-detail-log.vue'
  import { useUserStore } from '@/store/modules/user'
  import { getAllDisPlaysFiledApi } from '@/api-ecs/public'
  import { TableColumnItemType } from '/#/store'
  import dayjs from 'dayjs'
  // @ts-ignore
  import Whitelist from './white-list/index.vue'
  import { getAssetsPreviewListApi } from '~/src/api-ecs/assets-preview'
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')
  const settingsStore = useSettingsStore()
  const router = useRouter()
  const props = defineProps<{
    infoVal: any
    type?: string
  }>()
  const alertIndexType = 10
  const loading = ref(false)
  const userStore = useUserStore()
  const { getTableColumn } = userStore
  const alertColumns = getTableColumn(alertIndexType)
  const tableColumn = ref<TableColumnItemType[]>([])
  function getLevel(level: number) {
    return levelKey[level]
  }

  const infoValData = ref()

  const show = ref(false)

  async function formatColum() {
    const { data: userDisPlaysFiled } = await getAllDisPlaysFiledApi()
    const userColumnData = userDisPlaysFiled[alertIndexType] as number[]
    if (userColumnData?.length > 0) {
      tableColumn.value = userColumnData.map((key) =>
        alertColumns.find((item) => item.id === key)
      ) as TableColumnItemType[]
    } else {
      tableColumn.value = []
    }
  }

  const handleDownloadExcel = () => {
    const tHeader: string[] = []
    const filterVal: any = []
    try {
      tableColumn.value.forEach(({ fieldNameCn, fieldNameEn }) => {
        tHeader.push(fieldNameCn)
        filterVal.push(fieldNameEn)
      })
      loading.value = true
      import('@/utils/excel').then((excel) => {
        const data = formatJson(filterVal, [props.infoVal])
        excel.export_json_to_excel({
          header: tHeader,
          data,
          filename: `告警日志-${formatNstime(props.infoVal?.startTimeNs)}`,
          autoWidth: true,
          bookType: 'xlsx',
        })
        loading.value = false
      })
    } catch (error) {
      console.log(error)
      loading.value = false
    }
  }
  function formatJson(filterVal: any, jsonData: any) {
    return jsonData.map((v: any) =>
      filterVal.map((j: any) => {
        return formatExcelData(v, j)
      })
    )
  }

  function formatExcelData(val: any, key: string) {
    switch (key) {
      case 'startTimeNs':
      case 'warnTime':
        return formatNstime(val[key])
      default:
        return val[key] || ''
    }
  }

  async function download() {
    try {
      const {
        attackIp: clientIp,
        sourcePort: clientPort,
        victimIp: serverIp,
        targetPort: serverPort,
        startTimeNs,
      } = props.infoVal
      downloadLogPacket({ clientIp, clientPort, serverIp, serverPort, requestTimeNs: startTimeNs }, 2)
    } catch (error) {
      console.error(error)
    }
  }

  const addWhiteList = async () => {
    show.value = true
    // $baseConfirm('确认加入白名单', null, async () => {
    //   const { ruleId, clientIp, serverIp } = props.infoVal
    //   const { msg } = await getUpdateWarnWhiteApi({ ruleId, clientIp, serverIp })
    //   $baseMessage(msg, 'success', 'vab-hey-message-success')
    // })
  }

  const fullFlowSurveyHandle = async () => {
    const { clientIp, serverIp, sourcePort, targetPort, startTimeNs } = props.infoVal
    const time = startTimeNs / 1000000
    const sqlStr = `源ip = "${clientIp}" and 源端口 = "${sourcePort}" and 目的ip = "${serverIp}" and 目的端口 = "${targetPort}"`
    const indexType = 1
    const start = formatTime(dayjs(time).subtract(15, 'minute'))
    const end = formatTime(dayjs(time).add(15, 'minute'))
    const resolveRouter = router.resolve({
      path: '/retrieve/index',
      query: {
        info: encodeURIComponent(
          JSON.stringify({
            sql: sqlStr,
            indexType,
            timeRanges: [start, end],
          })
        ),
      },
    })
    window.open(resolveRouter.href, '_blank')
  }
  onMounted(() => {
    formatColum()
    infoValData.value = props.infoVal
    getAssetName()
  })

  const clientIPAssetName = ref()
  const serverIPAssetName = ref()
  const getAssetName = async () => {
    const { clientIp, serverIp } = props.infoVal
    const arr: string[] = []
    if (clientIp) {
      arr.push(clientIp)
    }
    if (serverIp) {
      arr.push(clientIp)
    }
    const {
      data: { records },
    } = await getAssetsPreviewListApi({
      pageNum: 1,
      pageSize: 10,
      searchIps: arr,
      dataSourceId: [] as number[],
      otherDataSource: false,
      searchStr: '',
    })
    if (clientIp) {
      clientIPAssetName.value = records[0]?.businessName || records[0]?.appName
    }
    if (serverIp) {
      serverIPAssetName.value = records[records.length - 1]?.businessName || records[records.length - 1]?.appName
    }
  }
  defineExpose({
    getAssetName,
  })
</script>

<script lang="ts">
  export default {
    name: 'BasicInfo',
  }
</script>

<template>
  <div v-loading="loading">
    <el-space v-if="!type" alignment="flex-end" :size="10">
      <el-button type="primary" @click="fullFlowSurveyHandle">全流量调查</el-button>
      <el-button type="primary" @click="addWhiteList">加入白名单</el-button>
      <el-button type="primary" @click="download">PCAP保存</el-button>
      <el-button type="primary" @click="handleDownloadExcel">下载告警详情</el-button>
    </el-space>
    <div v-else class="alert-header">
      <img class="form_icon" src="@/assets/alert_images/alert.svg" />
      <div>
        <span>告警设备IP:</span>
        <div>{{ infoVal?.warnDeviceIp || '&nbsp;&nbsp;-&nbsp;&nbsp;' }}</div>
      </div>
      <div>
        <span>状态码:</span>
        {{ infoVal?.responseStatusCode ?? '&nbsp;&nbsp;-&nbsp;&nbsp;' }}
      </div>
      <div>
        <span>XFF:</span>
        {{ infoVal?.xff || '&nbsp;&nbsp;-&nbsp;&nbsp;' }}
      </div>
      <div>
        <span>攻击结果:</span>
        {{ infoVal?.attackResult || '&nbsp;&nbsp;-&nbsp;&nbsp;' }}
      </div>
      <!-- <div>
        <span>攻击描述:</span>
        {{ infoVal?.attackDesc || '&nbsp;&nbsp;-&nbsp;&nbsp;' }}
      </div> -->
    </div>

    <alert-detail-log :log-val="JSON.parse(infoVal.sourceData)" />
  </div>
  <Whitelist v-if="show" :cur-data="infoValData" :is-show="show" @on-closeEvent="show = false" />
</template>

<style scoped lang="scss">
  .alert-header {
    width: 100%;
    height: 60px;
    background: #fff5f5;
    border-radius: 4px;
    margin-bottom: 20px;
    display: flex;
    align-items: center;
    img {
      margin-inline: 20px;
    }
    & > div {
      font-weight: 400;
      font-size: 13px;
      color: #4a4759;
      margin-right: 40px;
      width: 200px;
      &:nth-of-type() {
        flex: 1;
        margin-right: 0;
      }
      span {
        font-weight: 500;
        color: #342e58;
      }
      div {
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        display: inline-block;
        max-width: 130px;
        vertical-align: bottom;
      }
    }
  }
  .el-descriptions {
    margin-bottom: 15px;
    clear: both;
    :deep() {
      .el-descriptions__label {
        width: 172px;
      }

      .el-descriptions__content {
        width: 305px;
      }
    }
  }

  .el-space--horizontal {
    float: right;
    margin: 15px -10px 20px 0;
  }

  .my-item {
    background: #f8fbff;
    padding: 20px;
    line-height: 20px;
    overflow: hidden;
    word-break: break-all;
    // white-space: nowrap;
    word-wrap: break-word;
    box-shadow: 0 0 0 1px var(--el-input-border-color, var(--el-border-color)) inset;
    border-radius: var(--el-input-border-radius, var(--el-border-radius-base));
    transition: var(--el-transition-box-shadow);
    white-space: pre-wrap;

    &:hover {
      cursor: text;
    }
  }
</style>
