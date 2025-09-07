<script lang="ts">
  export default {
    name: 'LevelRule', // 规则列表
  }
</script>

<script setup lang="ts">
  import { DelLevelRulePageAPI, GetLevelRulePageAPI } from '~/src/api-ecs/retrieve'

  import EditRule from './edit-rule.vue'

  import { formatNstime } from '@/utils/time'

  import { useUserStore } from '@/store/modules/user'

  import { IndexTypeTpye } from '~/src/types'

  const props = defineProps<{
    isShow: boolean
  }>()

  const userStore = useUserStore()

  const { getTableColumn, getAllIndexType, username } = userStore

  const $baseConfirm: any = inject('$baseConfirm')

  const $baseMessage: any = inject('$baseMessage')

  const tagType = ref<IndexTypeTpye[]>(getAllIndexType())

  const total = ref(0) // 总条数

  const visible = ref(false) // 显隐

  const showPage = ref(false) // 显隐

  let listDate = reactive<object[]>([]) // 表格数据

  const itemData = ref()

  const listLoading = ref(false) // 是否加载

  // 检索参数
  const queryData = reactive({
    pageNum: 1,
    pageSize: 10,
  })

  const mode = ref()

  // 添加
  const handleAdd = () => {
    showPage.value = true
    mode.value = 'add'
    itemData.value = undefined
  }

  // 获取表格序号
  const curIndex = computed(() => (queryData.pageNum - 1) * queryData.pageSize + 1)

  // 改变页面容量
  function handleSizeChange(params: number) {
    queryData.pageSize = params
    getData()
  }

  // 改变页面
  function handleCurrentChange(params: number) {
    queryData.pageNum = params
    getData()
  }

  const emit = defineEmits<{
    (e: 'on-closeEvent', val: boolean): void
  }>()

  const handleClose = () => {
    emit('on-closeEvent', false)
  }

  let delList: [] = [] // 删除的数组
  let selectList: [] = []

  onMounted(() => {
    visible.value = props.isShow
  })

  // 多选项改变
  const setSelectRows = (e: any) => {
    delList = []
    e.forEach((item: any) => {
      // @ts-ignore
      delList.push(item.id)
    })
  }

  // 删除
  const handleDelete = ({ ...row }) => {
    if (row.row) {
      delList = []
      // @ts-ignore
      delList.push(row.row.id)
    }
    if (delList.length > 0 && !row.deleteAll) {
      const ids = delList.join(',')
      const flag = selectList.some((item: { hasChildren: boolean }) => {
        return item.hasChildren
      })
      const message = flag ? '删除项中存在API，确定删除' : '您确定要删除所选项吗'
      $baseConfirm(message, null, async () => {
        const { msg } = await DelLevelRulePageAPI({ ids: delList, deleteAll: false })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        delList = []
        await getData()
      })
    }
    if (!row.row && row.deleteAll) {
      $baseConfirm('你确定要删除所有数据吗', null, async () => {
        const { msg } = await DelLevelRulePageAPI({ deleteAll: true })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        await getData()
      })
    }
  }

  const handlerEeit = (row: any) => {
    showPage.value = true
    itemData.value = row
    mode.value = 'edit'
  }

  const getData = async () => {
    try {
      listLoading.value = true
      const { data } = await GetLevelRulePageAPI({ ...queryData })
      listDate = data.records
      total.value = data.total
    } finally {
      listLoading.value = false
    }
  }

  const changeData = (str: string) => {
    const res: IndexTypeTpye[] = tagType.value?.filter((item: any) => {
      return item.value == str
    })
    return res[0].label
  }

  const getThreatLevel = (type: number) => {
    const obj: {
      [key: string]: string
    } = {
      严重: 'high',
      一般: 'mid',
      普通: 'low',
    }
    return obj[type] || 'default'
  }

  const handlerExChange = (str: number) => {
    if (str == 0) {
      return '普通'
    } else if (str == 1) {
      return '一般'
    } else if (str == 2) {
      return '严重'
    }
  }

  onMounted(() => {
    getData()
  })
</script>

<template>
  <div class="level-rule">
    <el-dialog v-model="visible" :before-close="handleClose" destroy-on-close title="规则列表页" width="1180px">
      <div class="btn">
        <el-row :gutter="20">
          <el-button type="primary" @click="handleAdd">添加</el-button>
          <!-- <el-button type="danger" @click="handleDelete">批量删除</el-button> -->
          <el-dropdown style="margin-left: 10px">
            <span class="el-dropdown-link">
              <el-button type="danger">
                批量删除
                <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </el-button>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleDelete">删除选中</el-dropdown-item>
                <el-dropdown-item @click="(e) => handleDelete({ row: false, deleteAll: true })">
                  删除所有
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </el-row>
      </div>
      <el-table
        v-loading="listLoading"
        align="center"
        border
        :data="listDate"
        style="width: 100%; height: 450px"
        @selection-change="setSelectRows"
      >
        <el-table-column show-overflow-tooltip type="selection" />
        <el-table-column align="center" :index="(index) => curIndex + index" label="序号" type="index" width="70" />
        <el-table-column align="center" label="规则名称" prop="ruleName" />
        <el-table-column align="center" label="等级" prop="level">
          <template #default="{ row }">
            <span :class="['alert_tag', getThreatLevel(row.level)]">{{ row.level }}</span>
          </template>
        </el-table-column>
        <el-table-column align="center" label="协议" prop="indexType">
          <template #default="{ row }">
            {{ changeData(row.indexType) }}
          </template>
        </el-table-column>
        <el-table-column :align="'center'" label="更新时间" prop="updateTime" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.updateTime ? formatNstime(row.updateTime, false) : row.createTime }}
          </template>
        </el-table-column>
        <el-table-column align="center" fixed="right" label="操作" width="150">
          <template #default="{ row }">
            <el-button class="row_action" size="small" @click="handlerEeit(row)">编辑</el-button>
            <el-button v-permissions="['Admin']" class="row_action" size="small" @click="handleDelete({ row })">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="queryData.pageNum"
        v-model:page-size="queryData.pageSize"
        background
        class="known_pagination"
        layout="total, sizes, prev, pager, next, jumper"
        :page-num-sizes="[10, 20, 30]"
        :total="total"
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
      />
    </el-dialog>
    <EditRule
      v-if="showPage"
      :cur-data="itemData"
      :is-show="showPage"
      :remark="mode"
      @on-closeEvent="showPage = false"
      @on-reflash="getData"
    />
  </div>
</template>

<style scoped lang="scss">
  $criticalColor: #ff0202;
  $lowColor: #1b81fe;
  $midColor: #f1b04d;
  $highColor: #fa6d15;
  $defaultColor: #909399;
  .alert_tag {
    display: inline-block;
    width: 55px;
    line-height: 30px;
    border-radius: 5px;
    color: #fff;
    height: 30px;
    &.critical {
      background-color: rgba($criticalColor, 0.8);
    }
    &.low {
      background-color: rgba($lowColor, 0.8);
    }
    &.mid {
      background-color: rgba($midColor, 0.8);
    }
    &.high {
      background-color: rgba($highColor, 0.8);
    }
    &.default {
      background-color: #fff;
      color: inherit;
    }
  }
  .btn {
    margin: 0 0 20px 10px;
  }
</style>
