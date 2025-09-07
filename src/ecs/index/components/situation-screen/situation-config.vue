<script lang="ts">
  export default {
    name: 'SituationConfig', //
  }
</script>

<script setup lang="ts">
  import { ref } from 'vue'
  import { FormInstance, UploadUserFile, UploadProps, UploadInstance } from 'element-plus'
  import AddItem from './add-item.vue'
  import ImportItem from './import-item.vue'
  import { requireRules } from '~/src/utils/rules'
  import { setSystemConfigApi } from '~/src/api-ecs/system'
  import { getToken } from '@/utils/token'
  import { Plus } from '@element-plus/icons-vue'
  import { stationListApi, stationDelApi, stationExportApi, firewallApi, flowDeviceApi } from '~/src/api-ecs/situation'
  import { downloadFile } from '~/src/utils/download'
  import AesEncryptCBC from '~/src/utils/crypto'

  const $baseMessage: any = inject('$baseMessage')

  const $baseConfirm: any = inject('$baseConfirm')

  const props = defineProps<{
    remark: any
  }>()

  const emit = defineEmits<{
    (e: 'on-back-event', str: string): void
    (e: 'on-reflash'): void
  }>()

  const upload = ref<UploadInstance>()

  const formRef = ref<FormInstance>() // 表单实例

  const isLoading = ref(false) // 按钮动画

  const listLoading = ref(false) // 表格动画

  const visible = ref(false) // 新增编辑

  const showImp = ref(false) // 导出

  const currentItem = ref()

  const meta = ref('add')

  const listDate = ref()

  const flowDeviceData = ref()

  const firewallData = ref()

  // 表单数据
  const formData = reactive<{
    thumbnail: string
    title: string
    totalMileage?: string
    totalService?: string
    totalGantry?: string
    enable: boolean
  }>({
    thumbnail: '',
    title: '',
    totalMileage: undefined,
    totalService: undefined,
    totalGantry: undefined,
    enable: true,
  })

  const queryData = reactive({
    pageNum: 1,
    pageSize: 10,
  })

  // 表单数据校验规则
  const rules = ref<any>({
    thumbnail: [
      {
        required: true,
        message: '请选择上传文件',
        trigger: 'change',
      },
    ],
    title: requireRules,
  })

  const logUrl = ref<UploadUserFile[]>([])

  const mode = ref()

  const curItem = ref()

  const initData = async () => {
    curItem.value = props.remark
    mode.value = curItem.value.key
    const res1 = await flowDeviceApi({ page: 1, limit: 1000, query: JSON.stringify({}) })
    flowDeviceData.value = res1.data.list
    const res2 = await firewallApi({ page: 1, limit: 1000, query: JSON.stringify({}) })
    firewallData.value = res2.data.list
    for (const key in curItem.value.value) {
      // @ts-ignore
      formData[key] = curItem.value.value[key]
    }
    const path = curItem.value.value.thumbnail
    const arr = path.split('/')
    if (curItem.value.value.thumbnail) {
      logUrl.value = [{ name: arr[arr.length - 1], url: curItem.value.value.thumbnail }]
    }
  }

  const initRules = () => {
    if (mode.value == 'screen_highway_liaoning') {
      rules.value = {
        thumbnail: [
          {
            required: true,
            message: '请选择上传文件',
            trigger: 'change',
          },
        ],
        title: requireRules,
        totalMileage: requireRules,
        totalService: requireRules,
        totalGantry: requireRules,
      }
    }
  }
  watch(
    () => mode.value,
    () => {
      initRules()
    },
    { immediate: true }
  )

  // 获取表格序号
  const curIndex = computed(() => (queryData.pageNum - 1) * queryData.pageSize + 1)

  const initTable = async () => {
    listLoading.value = true
    try {
      const { data } = await stationListApi(queryData)
      total.value = data.total
      listDate.value = data.records
    } finally {
      listLoading.value = false
    }
  }

  const total = ref(0)

  onMounted(() => {
    initData()
    initRules()
    initTable()
  })

  // 页容量改变
  const handleSizeChange = (val: number) => {
    queryData.pageSize = val
    initTable()
  }

  // 页面改变
  const handleCurrentChange = (val: number) => {
    queryData.pageNum = val
    initTable()
  }

  const handleDelete = (row: any) => {
    $baseConfirm('你确定要删除吗', null, async () => {
      const arr = []
      arr.push(row.id)
      const { msg } = await stationDelApi({ ids: arr })
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      await initTable()
    })
  }

  const handleDdit = (row: any) => {
    visible.value = true
    meta.value = 'edit'
    currentItem.value = row
  }

  const handleAdd = () => {
    visible.value = true
    meta.value = 'add'
  }

  const handleExport = async () => {
    listLoading.value = true
    try {
      const res = await stationExportApi()
      // console.log(res)
      // $baseMessage('', 'success', 'vab-hey-message-success')
      downloadFile(res, '大屏数据')
    } finally {
      listLoading.value = false
    }
  }

  const handleImport = () => {
    showImp.value = true
  }

  const handleClick = () => {
    emit('on-back-event', 'chapter_2')
  }

  const beforeAvatarUpload: UploadProps['beforeUpload'] = (rawFile) => {
    if (rawFile.size / 1024 / 1024 > 100) {
      ElMessage.error('上传图片大小不能超过100MB!')
      return false
    }
    return true
  }

  const handleAvatarSuccess: UploadProps['onSuccess'] = (response, uploadFile) => {
    if (response.code === 20) {
      formData.thumbnail = response.data
    } else {
      $baseMessage(response.msg, 'error', 'vab-hey-message-error')
    }
  }

  const handleExceed: UploadProps['onExceed'] = (files, uploadFiles) => {
    upload.value!.clearFiles()
  }

  const changeXdrasUuid = (id: string) => {
    const data = flowDeviceData.value && JSON.parse(JSON.stringify(flowDeviceData.value))
    // @ts-ignore
    const res = data?.find((item: any) => {
      return item.id == id
    })
    return res?.ip
  }

  // 提交表单
  const submitForm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
      if (valid) {
        const password = ref('')
        ElMessageBox({
          title: '提示',
          showCancelButton: true,
          customClass: 'need-password-message-box',
          confirmButtonText: '确认',
          cancelButtonText: '取消',
          customStyle: {
            maxWidth: '500px',
          },
          message: () =>
            h('div', null, [
              h('div', { style: 'margin: 15px 0 5px 0;color:#55585b' }, '请输入敏感操作密码:(通过验证后方可进行操作)'),
              h(ElInput, {
                type: 'password',
                modelValue: password.value,
                placeholder: '请输入敏感操作密码',
                showPassword: true,
                style: 'margin-block: 10px',
                'onUpdate:modelValue': (val: string) => {
                  password.value = val
                },
              }),
            ]),

          beforeClose: async (action, instance, done) => {
            if (action === 'confirm') {
              instance.confirmButtonLoading = true
              isLoading.value = true
              try {
                const { msg } = await setSystemConfigApi(
                  {
                    id: curItem.value.id,
                    value: JSON.stringify({ ...formData }),
                  },
                  { password: AesEncryptCBC(password.value) }
                )
                ElMessage({ message: msg, type: 'success' })
                isLoading.value = false
                emit('on-reflash')
                handleClick()
                done()
              } catch (error: any) {
                isLoading.value = false
              }
            } else {
              done()
            }
          },
        }).catch(() => {})
      } else {
        console.log('error submit!', fields)
      }
    })
  }
</script>

<template>
  <div class="situation-config-container">
    <div class="top-title">态势大屏配置</div>
    <div class="content">
      <div class="config">
        <el-form ref="formRef" class="examplesDetailUpload-form" label-width="150px" :model="formData" :rules="rules">
          <div class="default">
            <el-row>
              <el-col :span="12">
                <el-form-item label="态势大屏名称" prop="title">
                  <el-input v-model="formData.title" clearable style="width: 520px" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item v-if="mode == 'screen_highway_liaoning'" label="mileage" prop="totalMileage">
                  <el-input v-model="formData.totalMileage" clearable style="width: 520px" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item v-if="mode == 'screen_highway_liaoning'" label="service数量" prop="totalService">
                  <el-input v-model="formData.totalService" clearable style="width: 520px" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item v-if="mode == 'screen_highway_liaoning'" label="gantry数量" prop="totalGantry">
                  <el-input v-model="formData.totalGantry" clearable style="width: 520px" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item label="是否启用">
              <el-switch v-model="formData.enable" />
            </el-form-item>
          </div>
          <el-form-item label="上传文件" prop="filedata">
            <el-upload
              ref="upload"
              v-model:file-list="logUrl"
              accept="image/jpg,image/png,image/jpeg"
              action="/v3/ecsPlatform/public/uploadThumbnail"
              :auto-upload="true"
              :before-upload="beforeAvatarUpload"
              class="upload-demo"
              drag
              :headers="{ Cookies: `EcsSessionId=${getToken()}` }"
              :limit="1"
              list-type="picture"
              :multiple="false"
              name="file"
              :on-exceed="handleExceed"
              :on-remove="() => (formData.thumbnail = '')"
              :on-success="handleAvatarSuccess"
            >
              <el-icon class="el-icon--upload"><upload-filled /></el-icon>
              <div class="el-upload__text">
                <span style="color: #0d88fe">点击上传缩略图</span>
                <br />
                只能上传png/jpg/jpeg
              </div>
            </el-upload>
          </el-form-item>
          <el-form-item label="&nbsp" style="margin-top: 80px">
            <el-button :loading="isLoading" type="primary" @click="submitForm(formRef)">保存</el-button>
            <el-button @click="handleClick">返回</el-button>
          </el-form-item>
        </el-form>
      </div>
      <div v-if="mode == 'screen_highway_liaoning'" class="bottom-table">
        <div class="table-top">
          <div class="title">关联信息</div>
          <div class="btns">
            <el-space alignment="flex-end" :size="10">
              <el-button :icon="Plus" type="primary" @click="handleAdd">添加</el-button>
              <el-button type="primary" @click="handleImport">导入</el-button>
              <el-button type="primary" @click="handleExport">导出</el-button>
            </el-space>
          </div>
        </div>
        <el-table v-loading="listLoading" :border="true" class="field" :data="listDate" style="margin-top: 15px">
          <el-table-column :align="'center'" label="序号" width="65">
            <template #default="{ $index }">
              <span>{{ curIndex + $index }}</span>
            </template>
          </el-table-column>
          <el-table-column align="center" label="名称" prop="stationName" :resizable="true" show-overflow-tooltip />
          <el-table-column align="center" label="探针" prop="xdrasUuid" :resizable="true" show-overflow-tooltip>
            <template #default="{ row }">
              <span>{{ changeXdrasUuid(row.xdrasUuid) }}</span>
            </template>
          </el-table-column>
          <el-table-column align="center" label="防火墙" prop="firewallIps" :resizable="true" show-overflow-tooltip>
            <!-- <template #default="{ row }">
              <span>{{ changeFirewallIps(row.firewallIps) }}</span>
            </template> -->
          </el-table-column>
          <el-table-column align="center" label="资产IP" prop="assetIp" :resizable="true" show-overflow-tooltip />
          <el-table-column align="center" label="防御类型" prop="defenseType" :resizable="true" show-overflow-tooltip>
            <template #default="{ row }">
              <span>{{ row.defenseType == 1 ? '检测' : '防御' }}</span>
            </template>
          </el-table-column>
          <el-table-column align="center" label="X轴" prop="eastLongitude" :resizable="true" show-overflow-tooltip />
          <el-table-column align="center" label="Y轴" prop="northernLatitude" :resizable="true" show-overflow-tooltip />
          <el-table-column
            align="center"
            fixed="right"
            label="操作"
            :resizable="true"
            show-overflow-tooltip
            width="180"
          >
            <template #default="{ row }">
              <el-button class="row_action" size="small" @click="handleDdit(row)">编辑</el-button>
              <el-button class="row_action" size="small" @click="handleDelete(row)">删除</el-button>
            </template>
          </el-table-column>
          <template #empty>
            <div style="height: 300px; line-height: 300px">暂无其他数据</div>
          </template>
        </el-table>
        <el-pagination
          v-model:current-page="queryData.pageNum"
          v-model:page-size="queryData.pageSize"
          background
          class="site_pagination"
          layout="total, sizes, prev, pager, next, jumper"
          :page-sizes="[10, 20, 30]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
    <AddItem
      v-if="visible"
      :curitem="currentItem"
      :firewall="firewallData"
      :flow-device="flowDeviceData"
      :mode="meta"
      :show-edit-site="visible"
      @on-closeEvent="visible = false"
      @on-reflash="initTable()"
    />
    <ImportItem v-if="showImp" :show-upload="showImp" @on-close-event="showImp = false" @on-reflash="initTable()" />
  </div>
</template>

<style scoped lang="scss">
  .top-title {
    background-color: #fff;
    height: 80px;
    height: 28px;
    font-size: 20px;
    font-weight: 500;
    color: #303133;
    line-height: 28px;
  }
  .situation-config-container {
    .content {
      height: calc(100vh - 150px);
      overflow-y: auto;
      &::-webkit-scrollbar {
        width: 0;
        height: 0;
      }
    }
  }

  .title {
    height: 28px;
    font-size: 16px;
    font-weight: 500;
    color: #303133;
    line-height: 28px;
  }
  .config {
    margin-top: 25px;

    :deep() {
      .upload-demo {
        // padding: 0;
        width: 318px;
        height: 150px;
        border-radius: 2px;
        border: 1px solid #dadfea;
      }

      .el-upload-list {
        height: 60px;
        width: 320px;
      }
      .el-upload-list__item {
        height: 100%;
        width: 100%;
      }
      .el-upload-dragger {
        padding: 0;
        width: 320px;
        height: 150px;
        .el-icon--upload {
          margin-top: 15px;
          height: 40px;
          width: 54px;
        }
        .el-upload__text {
          margin-top: -15px;
        }
      }
    }
  }

  .bottom-table {
    padding-bottom: 40px;
    .table-top {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
  }
</style>
