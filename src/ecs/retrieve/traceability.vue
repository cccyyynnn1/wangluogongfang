<script lang="ts">
  export default {
    name: 'RetrieveTraceability',
  }
</script>

<script setup lang="ts">
  import { getTetrieveTraceApi, getTetrieveTypeApi, getTetrieveLogApi } from '~/src/api-ecs/retrieve'
  import d3GraphForce2 from '@/components/d3-graph-force2.vue'
  import DbDetail from './components/detail-db.vue'
  import SmbDetail from './components/detail-smb.vue'
  import DnsDetail from './components/detail-dns.vue'
  import MailDetail from './components/detail-mail.vue'
  import FtpDetail from './components/detail-ftp.vue'
  import SessionInfo from '@/ecs/site/site-session/session-info.vue'
  import DetailOther from '@/ecs/retrieve/components/detail-other.vue'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import { useUserStore } from '@/store/modules/user'
  import { getAllDisPlaysFiledApi } from '@/api-ecs/public'
  import { TableColumnItemType } from '/#/store'
  import { formatNstime } from '@/utils/time'
  import dayjs from 'dayjs'
  const $baseMessage: any = inject('$baseMessage')
  interface Props {
    spaceId: number
  }
  const formatDate = (row: any, key: string) => {
    if (!row.key) return
    const time = row[key] / 1000000
    return dayjs(time).format('YYYY-MM-DD HH:mm:ss.SSS')
  }
  const { getTableColumn } = useUserStore()
  const detailDom: { [key: number]: any } = {
    2: DnsDetail,
    3: DbDetail,
    4: FtpDetail,
    5: SmbDetail,
    6: MailDetail,
    12: DetailOther,
    13: DetailOther,
    14: DetailOther,
    15: DetailOther,
    16: DetailOther,
    17: DetailOther,
    18: DetailOther,
    19: DetailOther,
    20: DetailOther,
    21: DetailOther,
  }
  const props = defineProps<Props>()
  const tableColumn = ref<TableColumnItemType[]>([])
  const indexTypeList = ref<{ name: string; value: number }[]>([])
  const indexType = ref<number | null>(null)
  const isFullscreen = ref(false)
  const chartLoading = ref(false)
  const graphRef = ref()
  const diggingVisible = ref(false)
  const showInfoData = ref()
  const detailVal = ref()
  const searchChartsData = reactive<{
    nodes: any[]
    links: any[]
  }>({
    nodes: [],
    links: [],
  })
  const queryForm = reactive({
    pageNum: 1,
    pageSize: 10,
    indexType: 0,
  })
  // 检索返回的数据
  const queryPage = reactive({
    total: 0,
    loading: true,
    listDate: [] as object[],
  })
  // HTTP检索单独使用！
  const sessionInfoVisible = ref(false)
  // 详情弹窗Visible 除HTTP检索使用
  const detailVisible = ref(false)
  // 获取表格序号
  const curIndex = computed(() => (queryForm.pageNum - 1) * queryForm.pageSize + 1)
  function echartClicckHandler(params: any) {
    diggingVisible.value = true
    showInfoData.value = params
  }

  const getTetrieveType = async () => {
    const { data } = await getTetrieveTypeApi(props.spaceId)
    indexTypeList.value = data || []
    queryForm.indexType = data[0]?.value || null
  }

  const getChartData = async () => {
    chartLoading.value = true
    try {
      const {
        data: { links, nodes },
      } = await getTetrieveTraceApi({
        isPathTrack: false,
        isTraceSource: true,
        isAssetVisit: false,
        workspaceId: props.spaceId,
      })
      if (!links || links.length === 0) {
        searchChartsData.links = []
        searchChartsData.nodes = []
        chartLoading.value = false
        return $baseMessage('暂无其他数据', 'warning', 'vab-hey-message-warning')
      }
      searchChartsData.links = links.map((i) => ({ ...i, source: i.clientIp, target: i.serverIp }))
      searchChartsData.nodes = nodes.map((i) => ({ ...i, id: i.ip }))
      chartLoading.value = false
    } catch (error) {
      searchChartsData.links = []
      searchChartsData.nodes = []
      chartLoading.value = false
    }
  }

  // 获取用户所有字段
  const getAllDisPlaysFiled = async () => {
    if (!queryForm.indexType) return
    const { data } = await getAllDisPlaysFiledApi()
    const userColumn = data[queryForm.indexType!] || []
    const allColum = getTableColumn(queryForm.indexType!)
    tableColumn.value = userColumn.map((key: number) => allColum.find((item) => item.id === key))
    getTetrieveLog()
  }
  function changeCellStyle(title: string) {
    if (['客户端口', '服务端口', '请求方式', '状态码'].includes(title)) {
      return '121'
    } else if (['客户端IP', '服务端IP', 'HOST', 'XFF'].includes(title)) {
      return '120'
    } else if (['请求时间', '响应时长'].includes(title)) {
      return '160'
    } else {
      return '150'
    }
  }
  const formatColumData = (columnKey: string, rowVal: any, indexType: number) => {
    if (rowVal === '' || rowVal === undefined || rowVal === null) return ''
    /**
     *  3:DB
     *  4:FTP
     *  5:SMB
     *  6:Email
     *  19:SSL
     **/
    const formatKey: { [key: number]: any } = {
      3: {
        result: {
          0: '成功',
          1: '失败',
        },
        operateType: {
          0: '登录',
          1: '操作',
        },
        model: {
          0: 'MySql',
          1: 'MSSQL',
          2: 'Oracle',
          3: 'Mongo',
          4: 'ElasticSearch',
          5: 'ClickHouse',
        },
      },
      4: {
        transModel: {
          0: '主动传输模式',
          1: '被动传输模式',
          2: '单端口模式',
        },
        result: {
          0: '成功',
          1: '失败',
        },
      },
      5: {
        cmd: {
          0: '创建文件',
          1: '创建文件夹',
          2: '读',
          3: '写',
          4: '删除文件',
          5: '删除文件夹',
        },
        result: {
          0: '成功',
          1: '失败',
        },
      },
      6: {
        model: {
          0: 'POP3',
          1: 'IMAP',
          2: 'SMTP',
        },
        operateType: {
          0: '登录',
          1: '邮件',
        },
        result: {
          0: '成功',
          1: '失败',
        },
        responseTime: 'ms',
      },
      19: {
        caType: {
          1: '可疑证书',
          2: '可信证书',
          3: '可疑中间证书',
          4: '可信中间证书',
        },
      },
    }

    if (columnKey === 'responseTime') return `${rowVal}/ms`
    try {
      const curKey = formatKey[indexType][columnKey]
      return curKey[rowVal]
    } catch (error) {
      return rowVal
    }
  }
  const getTetrieveLog = async () => {
    queryPage.loading = true
    const { type, clientIp, serverIp, ip } = showInfoData.value
    const {
      data: { records, total },
    } = await getTetrieveLogApi(
      type
        ? { ...queryForm, workspaceId: props.spaceId, clientIp, serverIp }
        : { ...queryForm, nodeIp: ip, workspaceId: props.spaceId }
    )
    queryPage.listDate = records || []
    queryPage.total = total || 0
    queryPage.loading = false
  }
  // 打开详情
  const handleInfo = (row: any) => {
    detailVal.value = { ...row, indexType: queryForm.indexType }
    if (queryForm.indexType === 1) return (sessionInfoVisible.value = true)
    detailVisible.value = true
  }
  watch(
    () => [queryForm.indexType, diggingVisible.value],
    () => {
      if (diggingVisible.value) {
        getAllDisPlaysFiled()
      }
    },
    {
      immediate: true,
    }
  )
  watch(
    () => props.spaceId,
    () => {
      getChartData()
    }
  )
  onMounted(() => {
    getChartData()
    getTetrieveType()
  })
</script>

<template>
  <div v-loading="chartLoading" class="retrieveTraceability">
    <d3-graph-force2 ref="graphRef" :graph-data="searchChartsData" is-traceability @click="echartClicckHandler" />
    <vab-dialog
      v-model="diggingVisible"
      align-center
      :close-on-click-modal="false"
      destroy-on-close
      title="深度挖掘"
      width="1375px"
    >
      <el-button-group>
        <el-button
          v-for="type in indexTypeList"
          :key="type.value"
          :type="queryForm.indexType === type.value ? 'primary' : 'default'"
          @click="queryForm.indexType = type.value"
        >
          {{ type.name }}
        </el-button>
      </el-button-group>
      <div v-loading="queryPage.loading" style="min-height: 480px">
        <el-table :data="queryPage.listDate" style="margin-top: 20px; width: 100%">
          <el-table-column align="center" type="selection" width="55" />
          <el-table-column
            :align="'center'"
            fixed="left"
            :index="(index) => curIndex + index"
            label="序号"
            type="index"
            width="55"
          />

          <el-table-column
            v-for="item in tableColumn"
            :key="item.id"
            align="center"
            :label="item.fieldNameCn"
            :min-width="changeCellStyle(item.fieldNameCn)"
            :prop="item.fieldNameEn"
            resizable
            show-overflow-tooltip
          >
            <template
              v-if="
                [
                  'caType',
                  'cmd',
                  'result',
                  'model',
                  'operateType',
                  'transModel',
                  'operateType',
                  'responseTime',
                ].includes(item.fieldNameEn)
              "
              #default="{ row }"
            >
              {{ formatColumData(item.fieldNameEn, row[item.fieldNameEn], queryForm.indexType) }}
            </template>
            <template v-else-if="item.fieldNameCn.includes('时间')" #default="{ row }">
              {{ formatDate(row, item.fieldNameEn) }}
            </template>
          </el-table-column>
          <el-table-column :align="'center'" fixed="right" label="操作" width="140">
            <template #default="{ row }">
              <el-button size="small" @click="handleInfo(row)">详情</el-button>
            </template>
          </el-table-column>
          <template #empty>
            <el-empty class="vab-data-empty" description="暂无数据" />
          </template>
        </el-table>
      </div>
      <el-space v-if="queryPage.total > 0" :size="10" style="width: 100%; justify-content: end" wrap>
        <div>
          <template v-if="queryPage.total > 1000">根据系统配置，仅显示前1000条</template>
          <template v-else>共 {{ queryPage.total }} 条</template>
        </div>
        <el-pagination
          v-model:current-page="queryForm.pageNum"
          v-model:page-size="queryForm.pageSize"
          background
          layout="sizes, prev, pager, next, jumper"
          :page-sizes="[10, 20, 30, 40, 50, 100]"
          style="margin-bottom: 20px"
          :total="queryPage.total > 1000 ? 1000 : queryPage.total"
          @current-change="getAllDisPlaysFiled"
          @size-change="getAllDisPlaysFiled"
        />
      </el-space>
    </vab-dialog>
    <!-- 检索详情 -->
    <vab-dialog
      v-model="detailVisible"
      :close-on-click-modal="false"
      destroy-on-close
      :loading="queryPage.loading"
      title="详情"
      width="1125px"
    >
      <component :is="detailDom[queryForm.indexType]" v-if="queryForm.indexType !== 1" :info-val="detailVal" />
    </vab-dialog>
    <session-info
      v-if="sessionInfoVisible"
      :info-data="detailVal"
      :show-session-info="sessionInfoVisible"
      @on-close-event="sessionInfoVisible = false"
    />
  </div>
</template>

<style scoped lang="scss">
  .retrieveTraceability {
    height: 100%;
    background-color: #f3f9ff;
  }
</style>
