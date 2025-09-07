<script lang="ts">
  export default {
    name: 'SiteHistory',
  }
</script>

<script setup lang="ts">
  import { getSiteHistoryApi, UpdateHistoryApi } from '@/api-ecs/site'

  import { SiteHistoryItem } from '@/types/index'
  import { h } from 'vue'
  import { ElMessageBox, ElInput } from 'element-plus'
  import dayjs from 'dayjs'
  import { useUserStore } from '@/store/modules/user'
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')
  const { username, userId } = useUserStore()
  interface Props {
    modelValue: boolean
    siteId: number | string
  }
  const props = withDefaults(defineProps<Props>(), {
    modelValue: false,
  })
  const emits = defineEmits<{
    (e: 'update:modelValue', visible: boolean): void
    (e: 'reload-search', data: any, hasTime: boolean): void
  }>()

  const visible = useVModel(props, 'modelValue', emits)
  const type = ref<'history' | 'favorites' | 'share' | 'shareto'>('history')
  const editNode = ref<Map<number, boolean>>(new Map())
  const queryForm = reactive({
    pageNum: 1,
    pageSize: 10,
  })
  const queryPage = reactive({
    total: 0,
    lodaing: true,
    listData: [] as SiteHistoryItem[],
  })
  const getNodeStatus = (id: number) => editNode.value.get(id)

  const tableNodeClickHandle = (id: number) => {
    if (type.value !== 'history') return
    editNode.value.set(id, true)
  }

  // const shareHandle = (row: SiteHistoryItem) => {
  //   const rNames = ref('')
  //   ElMessageBox({
  //     title: '分享给',
  //     confirmButtonText: '分享',
  //     message: () =>
  //       h(ElInput, {
  //         modelValue: rNames.value,
  //         type: 'textarea',
  //         placeholder: '请输入被分享人的姓名，多个被分享人使用英文逗号分隔。留空代表分享给所有人',
  //         rows: 5,
  //         resize: 'none',
  //         'onUpdate:modelValue': (val: string) => {
  //           rNames.value = val
  //         },
  //       }),
  //     beforeClose: async (action, instance, done) => {
  //       if (action === 'confirm') {
  //         instance.confirmButtonLoading = true
  //         const { code } = await toShareApi({
  //           rNames: rNames.value === '' ? '所有人' : rNames.value,
  //           workspaceId: props.spaceId,
  //           collectId: row.id!,
  //           sName: username,
  //         })
  //         $baseMessage('分享成功', 'success', 'vab-hey-message-success')
  //         done()
  //       } else {
  //         done()
  //       }
  //     },
  //   }).catch(() => {})
  // }
  // const cancelShareHandle = (row: SiteHistoryItem) => {
  //   $baseConfirm('你确定要取消分享', null, async () => {
  //     const { code } = await UpdateCollectShareApi(row.id!)
  //     $baseMessage('取消分享成功', 'success', 'vab-hey-message-success')
  //     getData()
  //   })
  // }

  // const shareClickHandle = (row: SiteHistoryItem) => {
  //   if (type.value !== 'share') return shareHandle(row)
  //   cancelShareHandle(row)
  // }

  const getHistory = async () => {
    queryPage.lodaing = true
    const params = { ...queryForm, siteId: props.siteId, collectStatus: type.value === 'history' ? false : true }
    const {
      data: { records, total },
    } = await getSiteHistoryApi(params)
    queryPage.listData = records
    queryPage.total = total
    queryPage.lodaing = false
  }
  // const getShare = async () => {
  //   queryPage.lodaing = true
  //   const params = type.value === 'share' ? { ...queryForm, sId: userId } : { ...queryForm, rName: username }
  //   const {
  //     data: { records, total },
  //   } = await getSkareHistoryApi(params)
  //   queryPage.listData = records
  //   queryPage.total = total
  //   queryPage.lodaing = false
  // }
  const updateCollect = async (row: SiteHistoryItem) => {
    const { code } = await UpdateHistoryApi({ ...row, collectStatus: !row.collectStatus })
    $baseMessage(row.collectStatus ? '已取消收藏' : '收藏成功', 'success', 'vab-hey-message-success')
    getData()
  }

  const updateRemarks = async (row: SiteHistoryItem) => {
    const { code } = await UpdateHistoryApi(row)
    $baseMessage('修改备注成功', 'success', 'vab-hey-message-success')
    getData()
  }

  const getData = () => {
    if (['history', 'favorites'].includes(type.value)) return getHistory()
    // getShare()
  }

  watchEffect(() => {
    getData()
  })
</script>

<template>
  <el-dialog v-model="visible" title="检索记录" width="1300px">
    <div v-loading="queryPage.lodaing">
      <el-button-group>
        <el-button :type="type === 'history' ? 'primary' : 'default'" @click="type = 'history'">回溯记录</el-button>
        <el-button :type="type === 'favorites' ? 'primary' : 'default'" @click="type = 'favorites'">
          指令收藏夹
        </el-button>
        <!-- <el-button :type="type === 'share' ? 'primary' : 'default'" @click="type = 'share'">我分享的</el-button>
        <el-button :type="type === 'shareto' ? 'primary' : 'default'" @click="type = 'shareto'">分享给我的</el-button> -->
      </el-button-group>
      <el-table border :data="queryPage.listData" style="margin-top: 20px; height: 550px">
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
              v-if="type === 'history'"
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
                icon="play-circle-line"
                style="font-size: 18px"
                @click="
                  () => {
                    emits('reload-search', row, false)
                    visible = false
                  }
                "
              />
            </span>
            <!-- 分享 -->
            <!-- <span
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
            </span> -->

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
</style>
