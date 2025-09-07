<script lang="ts">
  export default {
    name: 'IntelligenceCenterWhite',
  }
</script>

<script setup lang="ts">
  import { Plus, Delete, Edit, Search } from '@element-plus/icons-vue'
  import {
    getInfoWhitePageApi,
    updateInfoWhiteApi,
    deleteInfoWhiteApi,
    exportTemplateInfoWhiteApi,
    infoWhiteImportInfoWhiteApi,
    infoWhiteExportInfoWhiteApi,
  } from '@/api-ecs/alert'
  import { QueryInfoCloud, InfoWhiteItem, AddInfoWhiteItem, InfoCustomOptions } from '~/src/types'
  import { formatTime } from '@/utils/time'
  import { FormInstance, Action } from 'element-plus'
  import { h, ref } from 'vue'
  import { ElDivider } from 'element-plus'
  import IntelligenceCenterWhiteUpload from './intelligence-center-white/intelligence-center-white-upload.vue'
  import { downloadFile } from '~/src/utils/download'

  const spacer = h(ElDivider, { direction: 'vertical' })
  const props = defineProps<{
    infoCustomOptions: { [key in 'threatLevel' | 'iocType' | 'reliable' | 'threatType']: InfoCustomOptions }
  }>()
  const $baseMessage: any = inject('$baseMessage')
  const queryForm = reactive<QueryInfoCloud>({
    pageNum: 1,
    pageSize: 12,
    searchStr: '',
  })
  const iocWhiteFormRef = ref<FormInstance>()
  const value3 = ref(true)
  const tableSelects = ref<number[]>([])
  const addIocWhitelistVisible = ref(false)
  const isSync = ref(false)
  const uploadisible = ref(false)
  const tableData = ref<InfoWhiteItem[]>([])
  const total = ref(0)
  const form = reactive<AddInfoWhiteItem>({
    whiteName: '',
    whiteType: '',
    remark: '',
  })
  const iocList = ref<InfoCustomOptions>([])
  const rules = {
    whiteName: [{ required: true, message: '请输入白名单内容', trigger: 'blur' }],
    whiteType: [{ required: true, message: '请选择白名单类型', trigger: 'blur' }],
  }
  const getInfoWhitePageHandle = async () => {
    const { data } = await getInfoWhitePageApi({ ...queryForm })
    total.value = data.total
    tableSelects.value = []
    tableData.value = data.records || []
    tableData.value.forEach((item: InfoWhiteItem) => {
      item['checked'] = false
    })
  }
  const handleSelectionChange = (id: number) => {
    if (!tableSelects.value.includes(id)) {
      tableSelects.value.push(id)
    } else {
      tableSelects.value = tableSelects.value.filter((item) => {
        return !tableSelects.value.includes(id)
      })
    }
  }

  const handleFormSave = () => {
    iocWhiteFormRef.value?.validate((valid) => {
      if (valid) {
        updateInfoWhiteApi(form).then(({ code }) => {
          if (code === 20) {
            $baseMessage(form.id ? '编辑白名单成功！' : '新增白名单成功！', 'success', 'vab-hey-message-success')
            addIocWhitelistVisible.value = false
            reloadHandle()
          }
        })
      }
    })
  }

  const whiteInfoHandle = async (ids: number[]) => {
    const { code } = await deleteInfoWhiteApi({ ids, deleteAll: false })
    if (code === 20) {
      $baseMessage('删除白名单成功！', 'success', 'vab-hey-message-success')
      reloadHandle()
    }
  }
  const deleteConfirmHandle = (item?: any) => {
    let arr: number[] = []
    item.id ? arr.push(item.id) : (arr = tableSelects.value)
    if (arr?.length == 0) return
    ElMessageBox.alert(arr.length > 1 ? '确认删除所选数据？' : '确认删除此数据？', '白名单IOC', {
      confirmButtonText: '确定',
      callback: (action: Action) => {
        action === 'confirm' && whiteInfoHandle(arr)
      },
    })
  }
  const editWhiteHandle = (data: InfoWhiteItem) => {
    const { id, whiteName, remark, whiteType } = data
    form.id = id
    form.whiteName = whiteName
    form.remark = remark
    form.whiteType = whiteType
    addIocWhitelistVisible.value = true
  }
  const clearAddDataHandle = () => {
    iocWhiteFormRef.value?.resetFields()
    form.id = undefined
  }
  const reloadHandle = () => {
    queryForm.pageNum = 1
    getInfoWhitePageHandle()
  }
  onMounted(() => {
    getInfoWhitePageHandle()
    iocList.value = props.infoCustomOptions.iocType
  })

  const showUploadisible = () => {
    uploadisible.value = true
  }

  const handleExport = async () => {
    const res = await infoWhiteExportInfoWhiteApi({ ids: tableSelects.value, searchStr: queryForm.searchStr })
    downloadFile(res, '情报白名单')
  }
</script>

<template>
  <div class="info-white-container">
    <div style="margin-bottom: 14px">
      <el-input
        v-model="queryForm.searchStr"
        placeholder="请输入域名、IP、URL、文件MD5等IOC进行查询"
        :prefix-icon="Search"
        style="width: calc(100% - 100px); border-radius: 2px"
      />
      <el-button
        :auto-insert-space="false"
        style="width: 90px; height: 36px; margin-left: 10px"
        type="primary"
        @click="getInfoWhitePageHandle"
      >
        检索
      </el-button>
    </div>
    <el-space :size="10" style="width: 100%; justify-content: space-between" wrap>
      <div>
        <el-button
          :auto-insert-space="false"
          :icon="Plus"
          type="primary"
          @click="
            () => {
              isSync = false
              addIocWhitelistVisible = true
            }
          "
        >
          新增
        </el-button>
        <el-button plain @click="deleteConfirmHandle">批量删除</el-button>
      </div>
      <div style="margin-right: -10px">
        <!-- <el-button
          plain
          @click="
            () => {
              isSync = true
              addIocWhitelistVisible = true
            }
          "
        >
          同步配置
        </el-button> -->
        <el-button plain @click="showUploadisible">批量导入</el-button>
        <el-button plain @click="handleExport">批量导出</el-button>
      </div>
    </el-space>
    <div class="white-list">
      <el-row :gutter="10">
        <el-col v-for="item in tableData" :key="item.id" :lg="6" :md="12" :sm="24" :xl="6" :xs="24">
          <el-card shadow="never">
            <template #header>
              <div class="card-header">
                <img alt="icon" :src="require('@/assets/alert_images/white.png')" />
                <div>
                  {{ item.whiteName }}
                </div>
                <el-checkbox v-model="item.checked" size="large" @change="() => handleSelectionChange(item.id)" />
              </div>
            </template>
            <template #default>
              <div class="card-body">
                <ul>
                  <li>
                    <label>类型</label>
                    <span>{{ item.whiteTypeStr }}</span>
                  </li>
                  <li>
                    <label>创建人</label>
                    <span>{{ item.createUserStr }}</span>
                  </li>
                  <li>
                    <label>备注</label>
                    <span style="line-height: 20px; margin-top: 4px">{{ item.remark }}</span>
                  </li>
                </ul>
              </div>
              <el-space class="card-footer" :spacer="spacer">
                <el-button :icon="Edit" size="small" text @click="editWhiteHandle(item)">编辑</el-button>
                <el-button :icon="Delete" size="small" text @click="deleteConfirmHandle(item)">删除</el-button>
              </el-space>
            </template>
          </el-card>
        </el-col>
      </el-row>
    </div>
    <el-pagination
      v-model:current-page="queryForm.pageNum"
      v-model:page-size="queryForm.pageSize"
      background
      class="pagination"
      layout="total, sizes, prev, pager, next, jumper"
      :page-sizes="[12, 24, 36, 48]"
      :total="total"
      @current-change="getInfoWhitePageHandle"
      @size-change="reloadHandle"
    />
    <el-dialog
      v-model="addIocWhitelistVisible"
      class="ioc-whitelistdialog"
      :title="isSync ? '同步配置' : form.id ? '编辑白名单' : '新增白名单'"
      width="765px"
      @closed="clearAddDataHandle"
    >
      <template v-if="isSync">
        <div class="tips">
          <el-alert
            :closable="false"
            show-icon
            title="开启推送开关，系统会自动同步自定义情报至采集设备，采集设备需要配置有效的免密登录密钥。"
            type="warning"
          />
          <el-form class="ioc-form" label-width="128px" :model="form">
            <el-form-item label="情报自动推送开关:">
              <el-switch v-model="value3" active-text="开" inactive-text="关" inline-prompt />
            </el-form-item>
          </el-form>
        </div>
      </template>
      <el-form v-else ref="iocWhiteFormRef" class="ioc-form" label-width="140px" :model="form" :rules="rules">
        <el-form-item label="白名单:" prop="whiteName" required>
          <el-input v-model="form.whiteName" />
        </el-form-item>
        <el-form-item label="白名单类型:" prop="whiteType" required>
          <el-select v-model="form.whiteType">
            <el-option
              v-for="option in iocList"
              :key="option.dictValue"
              :label="option.dictLabel"
              :value="option.dictValue"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="备注:" prop="remark">
          <el-input v-model="form.remark" :autosize="false" :rows="3" type="textarea" />
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button v-if="isSync" :auto-insert-space="false" type="primary" @click="addIocWhitelistVisible = false">
            确定
          </el-button>
          <el-button v-else :auto-insert-space="false" type="primary" @click="handleFormSave">保存</el-button>
          <el-button :auto-insert-space="false" @click="addIocWhitelistVisible = false">取消</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
  <IntelligenceCenterWhiteUpload
    v-model:visible="uploadisible"
    :download-temp="exportTemplateInfoWhiteApi"
    :import-assets-fnc="infoWhiteImportInfoWhiteApi"
    title="导入情报白名单"
    @reflash="getInfoWhitePageHandle"
  />
</template>

<style scoped lang="scss">
  .info-white-container {
    height: 100%;
    padding-inline: 35px;
    :deep() {
      .el-input__inner {
        height: 34px;
      }
    }
    .white-list {
      margin-top: 14px;
      height: calc(100vh - 220px);
      overflow-y: auto;
      &::-webkit-scrollbar {
        width: 0;
        height: 0;
      }
      .card-header {
        display: flex;
        height: 80px;
        align-items: center;
        padding-inline: 15px;
        background-color: #fafaff;
        color: #342e58;
        font-size: 16px;
        font-weight: 500;
        img {
          width: 54px;
          height: 54px;
          border-radius: 20px;
        }
        div {
          overflow: hidden;
          display: -webkit-box;
          -webkit-box-orient: vertical;
          line-clamp: 2;
          -webkit-line-clamp: 2;
          white-space: pre-wrap;
          word-break: break-all;
          word-wrap: break-word;
          text-overflow: ellipsis;
          margin-inline: 10px;
          flex: 1;
        }

        // :deep(.el-checkbox) {
        // }
      }
      .card-footer {
        width: 100%;
        height: 44px;
        border-top: 1px solid #f2f0ff;
      }
      .card-body {
        padding: 0 15px;
        height: 120px;
        display: flex;
        align-items: center;
        ul {
          display: flex;
          flex-direction: column;
          padding-left: 0;
          margin-bottom: 0;
          margin-top: 0;
          li {
            display: flex;
            line-height: 30px;
            font-weight: 400;
            font-size: 14px;
            color: #342e58;
            span {
              flex: 1;
              overflow: hidden;
              display: -webkit-box;
              -webkit-box-orient: vertical;
              line-clamp: 2;
              -webkit-line-clamp: 2;
              white-space: pre-wrap;
              word-break: break-all;
              word-wrap: break-word;
              text-overflow: ellipsis;
            }
            label {
              width: 50px;
              display: inline-block;
              text-align: justify;
              text-align-last: justify;
              color: #9d9baa;
              margin-inline-end: 8px;
              &::after {
                content: ':';
              }
            }
          }
        }
      }
      :deep() {
        .el-card {
          border-radius: 16px;
          border-color: #f2f0ff;
          .el-card__header {
            padding: 0 !important;
            border-color: #f2f0ff;
          }
        }
        .el-divider {
          border-color: #f2f0ff;
        }
        .el-card__body {
          padding: 0;
        }
        .el-space__item {
          flex: 1;
          .el-button {
            font-weight: 400;
            font-size: 14px;
            color: #9d9baa;
            .el-icon {
              font-size: 14px;
            }
            &:hover,
            &:focus {
              background: transparent;
            }
          }
        }
      }
    }
    .tips {
      width: 625px;
      margin: 0 auto;
      .el-alert {
        padding: 14px 20px;
        background: #fff3f3;
        border-color: #ff4340;
        color: #ff4340;
        :deep() {
          .el-icon,
          .el-icon svg {
            color: #ff4340;
            width: 28px;
            height: 28px;
          }
          .el-alert__title {
            font-weight: 400;
            font-size: 14px;
          }
        }
      }
    }
    :deep() {
      .el-select {
        width: 100%;
      }
    }
  }
</style>
