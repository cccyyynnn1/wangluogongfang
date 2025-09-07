<script lang="ts">
  export default {
    name: 'RetrieveFavorites',
  }
</script>

<script setup lang="ts">
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import DbDetail from './components/detail-db.vue'
  import SmbDetail from './components/detail-smb.vue'
  import DnsDetail from './components/detail-dns.vue'
  import MailDetail from './components/detail-mail.vue'
  import FtpDetail from './components/detail-ftp.vue'
  import SessionInfo from '@/ecs/site/site-session/session-info.vue'
  import DetailOther from '@/ecs/retrieve/components/detail-other.vue'
  import { useUserStore } from '@/store/modules/user'
  import { getAllDisPlaysFiledApi } from '@/api-ecs/public'
  import DetailFile from './components/detail-file.vue'
  import { getCollectLogDataTypeApi, getLogPageApi, deletefavoritesLogPageApi } from '@/api-ecs/retrieve'
  import { useTableCopy, useCopy } from '@/utils'
  import { TableColumnItemType } from '/#/store'
  import dayjs from 'dayjs'
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')

  interface Props {
    spaceId: number
  }
  const props = defineProps<Props>()
  const multipleSelection = ref<number[]>([])
  const tableColumn = ref<TableColumnItemType[]>([])
  const userStore = useUserStore()
  const { getTableColumn } = userStore

  const indexTypeList = ref<{ name: string; value: number }[]>([])
  const userDisPlaysFiled = ref()
  // HTTP检索单独使用！
  const sessionInfoVisible = ref(false)
  // 详情弹窗Visible 除HTTP检索使用
  const detailVisible = ref(false)
  const detailType = ref(0)
  const detailVal = ref()
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
  // 获取表格序号
  const curIndex = computed(() => (queryForm.pageNum - 1) * queryForm.pageSize + 1)

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
    25: DetailOther,
    26: DetailOther,
    27: DetailFile,
    28: DetailOther,
    32: DetailOther,
  }
  function formatDate(row: any, key: string) {
    const time = row[key] / 1000000
    return dayjs(time).format('YYYY-MM-DD HH:mm:ss.SSS')
  }
  // 获取用户所有字段
  const getAllDisPlaysFiled = async () => {
    const { data } = await getAllDisPlaysFiledApi()
    userDisPlaysFiled.value = data
    formatColum(queryForm.indexType)
    queryPage.loading = false
  }
  function formatColum(type: number) {
    const columnData = getTableColumn(type)
    const userColumnData = userDisPlaysFiled.value[type] as number[]
    if (userColumnData?.length > 0) {
      tableColumn.value = userColumnData.map((key) =>
        columnData.find((item) => item.id === key)
      ) as TableColumnItemType[]
    } else {
      tableColumn.value = []
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

  const useTableCopyEvent = (row: any, column: any, cell: any, event: any) => {
    const target = column.property
    let value = undefined
    let label = ''
    if (target) {
      value = row[target]
      const res = tableColumn.value.find((item: any) => {
        return item.fieldNameEn == target
      })
      label = res!.fieldNameCn
    }
    useTableCopy(row, column, cell, event, [{ label: '复制', callback: useCopy, value: value }])
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
  // 打开详情
  const handleInfo = (row: any) => {
    detailVal.value = { ...row, indexType: detailType.value }
    if (detailType.value === 1) return (sessionInfoVisible.value = true)
    detailVisible.value = true
  }
  const getCollectLogDataTypeHandle = async () => {
    const { data } = await getCollectLogDataTypeApi(props?.spaceId)
    indexTypeList.value = data || []
    queryForm.indexType = data[0]?.value || 0
  }
  function handleSelectionChange(val: any[]) {
    multipleSelection.value = val.map((i) => i.id)
  }

  const deleteFavoritesLogPageHandle = (ids: number[]) => {
    $baseConfirm('你确定要删除当前项吗?', null, async () => {
      try {
        const { code } = await deletefavoritesLogPageApi(ids)
        $baseMessage('删除成功', 'success', 'vab-hey-message-success')
        getlogPageHandle()
      } catch (error) {
        console.log(error)
      }
    })
  }
  const getlogPageHandle = async () => {
    queryPage.loading = true
    const {
      data: { records, total },
    } = await getLogPageApi({ ...queryForm, workspaceId: props.spaceId })
    queryPage.listDate = records || []
    queryPage.total = total || 0
    queryPage.loading = false
    detailType.value = queryForm.indexType
  }

  watch(
    () => [queryForm.indexType, props.spaceId],
    () => {
      if (queryForm.indexType === 0) return (queryPage.listDate = [])
      getlogPageHandle()
      getAllDisPlaysFiled()
    },
    {
      immediate: true,
    }
  )
  onMounted(() => {
    getAllDisPlaysFiled()
    getCollectLogDataTypeHandle()
  })
</script>

<template>
  <div class="favorites-box">
    <el-empty v-if="queryForm.indexType === 0" class="vab-data-empty" description="暂无数据" />
    <template v-else>
      <div class="types">
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

        <el-button
          :disabled="multipleSelection.length === 0"
          type="danger"
          @click="() => deleteFavoritesLogPageHandle(multipleSelection)"
        >
          批量删除
        </el-button>
      </div>
      <el-table
        v-loading="queryPage.loading"
        class="favorites-table"
        :data="queryPage.listDate"
        style="margin-top: 20px; width: 100%"
        @cell-contextmenu="useTableCopyEvent"
        @selection-change="handleSelectionChange"
      >
        <el-table-column align="center" type="selection" width="55" />
        <el-table-column
          :align="'center'"
          fixed="left"
          :index="(index) => curIndex + index"
          label="序号"
          type="index"
          width="55"
        />
        <template v-for="item in tableColumn" :key="item.id">
          <el-table-column
            align="center"
            :label="item.fieldNameCn"
            :min-width="changeCellStyle(item.fieldNameCn)"
            :prop="item.fieldNameEn"
            resizable
            show-overflow-tooltip
          >
            <template
              #default="{ row }"
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
            >
              {{ formatColumData(item.fieldNameEn, row[item.fieldNameEn], queryForm.indexType) }}
            </template>
            <template v-else-if="item.fieldNameCn.includes('时间')" #default="{ row }">
              {{ row[item.fieldNameEn] && formatDate(row, item.fieldNameEn) }}
            </template>
          </el-table-column>
        </template>
        <el-table-column :align="'center'" fixed="right" label="操作" width="140">
          <template #default="{ row }">
            <el-button size="small" @click="handleInfo(row)">详情</el-button>
            <el-button size="small" @click="() => deleteFavoritesLogPageHandle([row.id])">删除</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty class="vab-data-empty" description="暂无数据" />
        </template>
      </el-table>
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
          @current-change="getlogPageHandle"
          @size-change="getlogPageHandle"
        />
      </el-space>
      <!-- 检索详情 -->
      <vab-dialog v-model="detailVisible" :close-on-click-modal="false" destroy-on-close title="详情" width="1125px">
        <component
          :is="detailDom[detailType]"
          v-if="detailType !== 1"
          :info-val="detailVal"
          :workspace-id="props.spaceId"
        />
      </vab-dialog>
    </template>
    <session-info
      v-if="sessionInfoVisible"
      :info-data="detailVal"
      :show-session-info="sessionInfoVisible"
      :workspace-id="props.spaceId"
      @on-close-event="sessionInfoVisible = false"
    />
  </div>
</template>

<style scoped lang="scss">
  .favorites-box {
    position: relative;
    .btbn {
      position: absolute;
      top: -30px;
      right: 0;
    }
  }
  .types {
    display: flex;
    justify-content: space-between;
  }
  .favorites-table {
    height: calc(100vh - 190px);
  }
</style>
