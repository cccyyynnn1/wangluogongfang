<script lang="ts">
  export default {
    name: 'RetrieveHistory',
  }
</script>

<script setup lang="ts">
  import {
    getRetrieveHistoryApi,
    getSkareHistoryApi,
    UpdateHistoryApi,
    UpdateCollectShareApi,
    toShareApi,
    deleteCollectApi,
    getAllUserNameAPI,
  } from '@/api-ecs/retrieve'

  import { RetrieveHistoryItem } from '@/types/index'
  import { h } from 'vue'
  import { ElMessageBox, ElInput, ElSelect } from 'element-plus'
  import dayjs from 'dayjs'
  import { useUserStore } from '@/store/modules/user'
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')
  const { username, userId } = useUserStore()
  const nameList = ref([])
  interface Props {
    modelValue: boolean
    spaceId: number
    active?: 'history' | 'favorites' | 'share' | 'shareto'
  }
  const props = withDefaults(defineProps<Props>(), {
    modelValue: false,
    active: 'history',
  })
  const emits = defineEmits<{
    (e: 'update:modelValue', visible: boolean): void
    (e: 'reload-search', data: any, hasTime: boolean): void
  }>()
  // 自定义时间
  const timeDate = ref()
  const visible = useVModel(props, 'modelValue', emits)
  const type = ref<'history' | 'favorites' | 'share' | 'shareto'>('history')

  const editNode = ref<Map<number, boolean>>(new Map())
  const queryForm = reactive({
    pageNum: 1,
    pageSize: 10,
    searchStr: '',
  })
  const queryPage = reactive({
    total: 0,
    lodaing: true,
    listData: [] as RetrieveHistoryItem[],
  })
  const getNodeStatus = (id: number) => editNode.value.get(id)

  const tableNodeClickHandle = (id: number) => {
    if (type.value !== 'history' && type.value !== 'favorites') return
    editNode.value.set(id, true)
  }

  const shareHandle = (row: RetrieveHistoryItem) => {
    const rNames = ref<string[]>([])
    const remarks = ref(row.remarks)
    ElMessageBox({
      title: '分享给',
      confirmButtonText: '分享',
      customStyle: {
        maxWidth: '500px',
      },
      message: () =>
        h('div', null, [
          h(
            ElSelect,
            {
              modelValue: rNames.value,
              placeholder: '请选择被分享人',
              multiple: true,
              class: 'abc',
              style: { width: '100%' },
              'onUpdate:modelValue': (val: string[]) => {
                rNames.value = val
              },
            },
            () =>
              nameList.value.map((item: string) => {
                return h(ElSelect.Option, {
                  key: item,
                  label: item,
                  value: item,
                })
              })
          ),
          h(ElInput, {
            modelValue: remarks.value,
            type: 'textarea',
            placeholder: '请输入收藏备注',
            rows: 5,
            style: 'margin-top: 10px',
            resize: 'none',
            'onUpdate:modelValue': (val: string) => {
              remarks.value = val
            },
          }),
        ]),

      beforeClose: async (action, instance, done) => {
        if (action === 'confirm') {
          instance.confirmButtonLoading = true
          const { code } = await toShareApi({
            rNames: rNames.value.length === 0 ? '所有人' : rNames.value.join(','),
            workspaceId: props.spaceId,
            collectId: row.id!,
            sName: username,
            remarks: remarks.value,
          })
          $baseMessage('分享成功', 'success', 'vab-hey-message-success')
          done()
        } else {
          done()
        }
      },
    }).catch(() => {})
  }
  const cancelShareHandle = (row: RetrieveHistoryItem) => {
    $baseConfirm('你确定要取消分享', null, async () => {
      const { code } = await UpdateCollectShareApi(row.id!)
      $baseMessage('取消分享成功', 'success', 'vab-hey-message-success')
      getData()
    })
  }

  const shareClickHandle = (row: RetrieveHistoryItem) => {
    if (type.value !== 'share') return shareHandle(row)
    cancelShareHandle(row)
  }

  const deleteClickHandle = async (row: RetrieveHistoryItem) => {
    const { msg } = await deleteCollectApi([row.id!], username)
    $baseMessage('删除分享成功', 'success', 'vab-hey-message-success')
    getData()
  }
  const getHistory = async () => {
    queryPage.lodaing = true
    const [searchStTime, searchEdTime] = timeDate.value || [null, null]
    const params = {
      ...queryForm,
      searchStTime,
      searchEdTime,
      workspaceId: props.spaceId,
      collectStatus: type.value === 'history' ? false : true,
    }
    const {
      data: { records, total },
    } = await getRetrieveHistoryApi(params)
    queryPage.listData = records
    queryPage.total = total
    queryPage.lodaing = false
  }
  const getShare = async () => {
    queryPage.lodaing = true
    const params = type.value === 'share' ? { ...queryForm, sId: userId } : { ...queryForm, rName: username }
    const {
      data: { records, total },
    } = await getSkareHistoryApi(params)
    queryPage.listData = records
    queryPage.total = total
    queryPage.lodaing = false
  }
  const updateCollect = async (row: RetrieveHistoryItem) => {
    const { code } = await UpdateHistoryApi({ ...row, collectStatus: !row.collectStatus })
    $baseMessage(row.collectStatus ? '已取消收藏' : '收藏成功', 'success', 'vab-hey-message-success')
    getData()
  }

  const updateRemarks = async (row: RetrieveHistoryItem) => {
    const { code } = await UpdateHistoryApi(row)
    $baseMessage('修改备注成功', 'success', 'vab-hey-message-success')
    getData()
  }

  const getData = () => {
    if (['history', 'favorites'].includes(type.value)) return getHistory()
    getShare()
  }
  onMounted(() => {
    type.value = props.active
    getAllUserNameAPI().then((res) => {
      nameList.value = res.data || []
    })
  })
  watch(
    type,
    () => {
      queryForm.searchStr = ''
      timeDate.value = null
      setTimeout(() => {
        getData()
      }, 100)
    },
    {
      immediate: true,
    }
  )
</script>

<template>
  <el-dialog v-model="visible" title="检索记录" width="1300px">
    <div v-loading="queryPage.lodaing">
      <el-button-group>
        <el-button :type="type === 'history' ? 'primary' : 'default'" @click="type = 'history'">回溯记录</el-button>
        <el-button :type="type === 'favorites' ? 'primary' : 'default'" @click="type = 'favorites'">
          指令收藏夹
        </el-button>
        <el-button :type="type === 'share' ? 'primary' : 'default'" @click="type = 'share'">我分享的</el-button>
        <el-button :type="type === 'shareto' ? 'primary' : 'default'" @click="type = 'shareto'">分享给我的</el-button>
      </el-button-group>
      <div style="float: right; display: flex">
        <vab-date-time-picker
          v-if="['history', 'favorites'].includes(type)"
          v-model="timeDate"
          :allow-clear="true"
          style="width: 350px"
        />
        <el-input
          v-model="queryForm.searchStr"
          clearable
          placeholder="模糊匹配"
          style="width: 350px; margin-inline: 20px"
        />
        <el-button type="primary" @click="getData()">检索</el-button>
      </div>

      <el-table :data="queryPage.listData" style="margin-top: 20px; height: 550px">
        <el-table-column
          v-if="['history', 'favorites'].includes(type)"
          :formatter="({ createTime }) => dayjs(createTime).format('YYYY-MM-DD HH:mm:ss')"
          label="检索时间"
          prop="createTime"
          width="180"
        />
        <el-table-column label="数据类型" prop="dataType" width="120" />
        <el-table-column label="检索时间范围" prop="address" width="180">
          <template #default="{ row }">
            <div>
              {{ row.searchStTime }}
              <br />
              {{ row.searchEdTime }}
            </div>
          </template>
        </el-table-column>
        <el-table-column label="检索条件" prop="searchSql" show-overflow-tooltip />
        <el-table-column label="备注" prop="remarks">
          <template #default="{ row }">
            <el-input
              v-if="['history', 'favorites'].includes(type)"
              v-model="row.remarks"
              :readonly="!getNodeStatus(row.id)"
              @click="() => tableNodeClickHandle(row.id)"
              @keydown.enter="() => updateRemarks(row)"
            />
          </template>
        </el-table-column>
        <el-table-column v-if="type === 'shareto'" label="分享人" prop="sName" />
        <el-table-column v-if="type === 'share'" label="被分享人" prop="rNames" />
        <el-table-column align="center" fixed="right" label="操作" width="170">
          <template #default="{ row }">
            <span class="btn">
              <vab-icon
                icon="bubble-chart-line"
                style="font-size: 18px"
                @click="
                  () => {
                    emits('reload-search', row, false)
                    visible = false
                  }
                "
              />
            </span>
            <span class="btn">
              <vab-icon
                icon="play-circle-line"
                style="font-size: 18px"
                @click="
                  () => {
                    emits('reload-search', row, true)
                    visible = false
                  }
                "
              />
            </span>
            <!-- 分享 -->
            <span
              v-if="['history', 'favorites', 'share'].includes(type)"
              class="btn"
              @click="() => shareClickHandle(row)"
            >
              <template v-if="type !== 'share'">
                <vab-icon icon="share-line" style="font-size: 18px" />
              </template>
              <template v-else>
                <vab-icon icon="share-fill" style="font-size: 18px; color: var(--el-color-warning)" />
              </template>
            </span>

            <!-- 收藏 -->
            <span v-if="['history', 'favorites'].includes(type)" class="btn" @click="() => updateCollect(row)">
              <template v-if="type === 'history'">
                <vab-icon
                  v-if="row.collectStatus"
                  icon="star-fill"
                  style="font-size: 18px; color: var(--el-color-warning)"
                />
                <vab-icon v-else icon="star-line" style="font-size: 18px" />
              </template>
              <template v-else>
                <vab-icon icon="star-fill" style="font-size: 18px; color: var(--el-color-warning)" />
              </template>
            </span>

            <!-- 删除 -->
            <span v-if="['shareto'].includes(type)" class="btn" @click="() => deleteClickHandle(row)">
              <vab-icon icon="close-circle-line" style="font-size: 18px" />
            </span>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="queryForm.pageNum"
        v-model:page-size="queryForm.pageSize"
        background
        layout=" sizes, prev, pager, next, jumper"
        :page-sizes="[10, 20, 30]"
        :total="queryPage.total"
        @current-change="getData"
        @size-change="getData"
      />
    </div>
  </el-dialog>
</template>

<style scoped lang="scss">
  .btn {
    cursor: pointer;
    margin: 0 8px;
  }
  :deep(.el-popper) {
    max-width: fit-content !important;
  }
</style>
