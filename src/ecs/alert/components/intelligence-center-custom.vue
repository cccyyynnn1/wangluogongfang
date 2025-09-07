<script lang="ts">
  export default {
    name: 'IntelligenceCenterCustom',
  }
</script>

<script setup lang="ts">
  import '@/assets/js/particles.min'
  import { Plus, Search } from '@element-plus/icons-vue'
  import {
    getInfoCustomPageApi,
    updateInfoCustomPageApi,
    deleteInfoCustomApi,
    updateInfoCustomStatusApi,
    exportInfoCustomApi,
    importInfoCustomApi,
    infoCustomExportTempApi,
    applyRulesApi,
  } from '@/api-ecs/alert'
  import { QueryInfoCloud, InfoCustomItem, InfoCustomOptions, AddInfoCloud } from '~/src/types'
  import { formatTime } from '@/utils/time'
  import { Action, FormInstance, FormItemRule } from 'element-plus'
  import { downloadFile } from '~/src/utils/download'
  import IntelligenceCenterWhiteUpload from './intelligence-center-white/intelligence-center-white-upload.vue'
  type IOCSubType = 'IP地址' | '域名' | 'URI' | '端口' | 'MD5'
  type IOCRulesType =
    | 'threatName'
    | 'threatType'
    | 'threatLevel'
    | 'reliable'
    | 'ip'
    | 'domain'
    | 'uri'
    | 'md5'
    | 'port'

  const DICTIONARY: {
    [key in IOCSubType]: {
      key: 'ip' | 'domain' | 'uri' | 'md5' | 'port'
      rules: FormItemRule[]
    }
  } = {
    IP地址: {
      key: 'ip',
      rules: [
        {
          validator: validateIP,
          trigger: 'blur',
        },
      ],
    },
    域名: {
      key: 'domain',
      rules: [
        {
          required: true,
          message: '请输入域名',
          trigger: 'blur',
        },
        {
          pattern: /^([0-9a-zA-Z-]{1,}\.)+([a-zA-Z]{2,})$/,
          message: '请输入合法域名',
          trigger: 'blur',
        },
      ],
    },
    端口: {
      key: 'port',
      rules: [
        {
          required: true,
          message: '请输入端口',
          trigger: 'blur',
        },
        {
          pattern: /^([0-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{4}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$/,
          message: '请输入合法端口',
          trigger: 'blur',
        },
      ],
    },
    URI: {
      key: 'uri',
      rules: [
        {
          required: true,
          message: '请输入URI',
          trigger: 'blur',
        },
      ],
    },
    MD5: {
      key: 'md5',
      rules: [
        {
          required: true,
          message: '请输入MD5',
          trigger: 'blur',
        },
        {
          pattern: /^[a-fA-F0-9]{16}$|^[a-fA-F0-9]{32}$/,
          message: '请输入合法MD5',
          trigger: 'blur',
        },
      ],
    },
  }
  const props = defineProps<{
    infoCustomOptions: { [key in 'threatLevel' | 'iocType' | 'reliable' | 'threatType']: InfoCustomOptions }
  }>()

  const $baseMessage: any = inject('$baseMessage')
  const queryForm = reactive<QueryInfoCloud>({
    pageNum: 1,
    pageSize: 100,
    searchStr: '',
  })
  const uploadFile = ref()
  const uploadVisible = ref(false)
  const tableSelects = ref<InfoCustomItem[]>([])
  const addIocVisible = ref(false)
  const iocFormRef = ref<FormInstance>()
  const tableData = ref<InfoCustomItem[]>([])
  const total = ref(0)
  const iocOptions = ref()
  const form = reactive<AddInfoCloud>({
    iocType: '1',
    threatName: '',
    threatType: '',
    threatLevel: '',
    reliable: '',
  })
  const rules = reactive<Partial<{ [key in IOCRulesType]: FormItemRule[] }>>({
    threatName: [{ required: true, message: '请输入威胁名称', trigger: 'blur' }],
    threatType: [{ required: true, message: '请选择威胁类型', trigger: 'blur' }],
    threatLevel: [{ required: true, message: '请选择威胁等级', trigger: 'blur' }],
    reliable: [{ required: true, message: '请选择置信度', trigger: 'blur' }],
  })
  function validateIP(rule: any, value: any, callback: any) {
    const IPV4 =
      /^((\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5])\.){3}(\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5])(?::(?:[0-9]|[1-9][0-9]{1,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5]))?$/
    const IPV6 =
      /(^(?:(?:(?:[0-9A-Fa-f]{1,4}:){7}[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){6}:[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){5}:([0-9A-Fa-f]{1,4}:)?[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){4}:([0-9A-Fa-f]{1,4}:){0,2}[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){3}:([0-9A-Fa-f]{1,4}:){0,3}[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){2}:([0-9A-Fa-f]{1,4}:){0,4}[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){6}((\b((25[0-5])|(1\d{2})|(2[0-4]\d)|(\d{1,2}))\b)\.){3}(\b((25[0-5])|(1\d{2})|(2[0-4]\d)|(\d{1,2}))\b))|(([0-9A-Fa-f]{1,4}:){0,5}:((\b((25[0-5])|(1\d{2})|(2[0-4]\d)|(\d{1,2}))\b)\.){3}(\b((25[0-5])|(1\d{2})|(2[0-4]\d)|(\d{1,2}))\b))|(::([0-9A-Fa-f]{1,4}:){0,5}((\b((25[0-5])|(1\d{2})|(2[0-4]\d)|(\d{1,2}))\b)\.){3}(\b((25[0-5])|(1\d{2})|(2[0-4]\d)|(\d{1,2}))\b))|([0-9A-Fa-f]{1,4}::([0-9A-Fa-f]{1,4}:){0,5}[0-9A-Fa-f]{1,4})|(::([0-9A-Fa-f]{1,4}:){0,6}[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){1,7}:))$)|(^\[(?:(?:(?:[0-9A-Fa-f]{1,4}:){7}[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){6}:[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){5}:([0-9A-Fa-f]{1,4}:)?[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){4}:([0-9A-Fa-f]{1,4}:){0,2}[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){3}:([0-9A-Fa-f]{1,4}:){0,3}[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){2}:([0-9A-Fa-f]{1,4}:){0,4}[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){6}((\b((25[0-5])|(1\d{2})|(2[0-4]\d)|(\d{1,2}))\b)\.){3}(\b((25[0-5])|(1\d{2})|(2[0-4]\d)|(\d{1,2}))\b))|(([0-9A-Fa-f]{1,4}:){0,5}:((\b((25[0-5])|(1\d{2})|(2[0-4]\d)|(\d{1,2}))\b)\.){3}(\b((25[0-5])|(1\d{2})|(2[0-4]\d)|(\d{1,2}))\b))|(::([0-9A-Fa-f]{1,4}:){0,5}((\b((25[0-5])|(1\d{2})|(2[0-4]\d)|(\d{1,2}))\b)\.){3}(\b((25[0-5])|(1\d{2})|(2[0-4]\d)|(\d{1,2}))\b))|([0-9A-Fa-f]{1,4}::([0-9A-Fa-f]{1,4}:){0,5}[0-9A-Fa-f]{1,4})|(::([0-9A-Fa-f]{1,4}:){0,6}[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){1,7}:))\](?::(?:[0-9]|[1-9][0-9]{1,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5]))?$)/i

    if (IPV4.test(value) || IPV6.test(value)) {
      return callback()
    }
    if (!value) {
      return callback(new Error('请输入IP地址'))
    }
    if (!IPV4.test(value)) {
      callback(new Error('请输入合法IPV4地址'))
    }
    if (!IPV6.test(value)) {
      callback(new Error('请输入合法IPV6地址'))
    }
  }
  const getInfoCustomPageHandle = async () => {
    const { data } = await getInfoCustomPageApi(queryForm)
    total.value = data.total
    tableData.value = data.records || []
  }

  const handleSelectionChange = (val: InfoCustomItem[]) => {
    tableSelects.value = val
  }
  const handleFormSave = () => {
    iocFormRef.value?.validate((valid) => {
      if (valid) {
        updateInfoCustomPageApi(form).then((res) => {
          if (res.code === 20) {
            $baseMessage('新增简易IOC成功！', 'success', 'vab-hey-message-success')
            reloadHandle()
            addIocVisible.value = false
          }
        })
      }
    })
  }
  const clearAddDataHandle = () => {
    iocFormRef.value?.resetFields()
    form.id = undefined
  }
  const deleteConfirmHandle = (ids: number[]) => {
    if (!ids?.length) return
    ElMessageBox.alert(ids.length > 1 ? '确认删除所选数据？' : '确认删除此数据？', '自定义情报', {
      confirmButtonText: '确定',
      callback: (action: Action) => {
        action === 'confirm' && customInfoHandle(ids)
      },
    })
  }
  const editWhiteHandle = (data: InfoCustomItem) => {
    Object.assign(form, data)
    addIocVisible.value = true
  }
  const customInfoHandle = async (ids: number[]) => {
    const { code } = await deleteInfoCustomApi({ ids, deleteAll: false })
    if (code === 20) {
      $baseMessage('删除自定义情报成功！', 'success', 'vab-hey-message-success')
      reloadHandle()
    }
  }

  const updateInfoCustomStatusHandle = async (ids: number[], enable: any) => {
    if (!ids?.length) return
    const { code } = await updateInfoCustomStatusApi({ ids, enable })
    if (code === 20) {
      $baseMessage('修改状态成功！', 'success', 'vab-hey-message-success')
      reloadHandle()
    }
  }
  const reloadHandle = () => {
    queryForm.pageNum = 1
    getInfoCustomPageHandle()
  }
  const exportInfoCustomHandle = async () => {
    const data = await exportInfoCustomApi({ searchStr: queryForm.searchStr || '' })
    downloadFile(data, '自定义情报')
  }
  const setFormOrRuls = (iocType: string, isSetForm = true) => {
    const curOptin = props.infoCustomOptions.iocType.find((item) => item.dictValue == iocType)
    if (!curOptin) return
    const { dictLabel, dictValue } = curOptin
    const labels = dictLabel.split(':')
    labels.forEach((label) => {
      const curForm = DICTIONARY[label as IOCSubType]
      const key = curForm.key as 'ip' | 'domain' | 'uri' | 'md5' | 'port'
      if (isSetForm) {
        form[key] = form[key] || ''
        rules[key] = curForm.rules
      } else {
        delete form[key]
        delete rules[key]
      }
    })
  }
  // 根据iocType获取相应IOC下的Form值和Ruls
  const getCustomIocTypeVal = (newVal: string, oldVal?: string) => {
    if (oldVal) setFormOrRuls(oldVal, false)
    setFormOrRuls(newVal)
  }
  const getConfidence = (code: string) => {
    return ['低', '中', '高'].findIndex((i) => code === i) + 1
  }

  const handleToApply = async () => {
    const { msg } = await applyRulesApi()
    $baseMessage(msg, 'success', 'vab-hey-message-success')
  }

  // const handleExceed = (files: File[]) => {
  //   console.log(files, 'files')
  // }

  // const handleUploadSuccess = (response: any) => {
  //   $baseMessage(response.msg, 'success', 'vab-hey-message-success')
  //   uploadFile.value = undefined
  // }
  // const hanldeFileUpload = () => {
  //   console.log(uploadFile.value, 'uploadFile.value')
  //   return importInfoCustomApi({ file: uploadFile.value[0].raw })
  // }
  watch(
    () => form.iocType,
    (newVal, oldVal) => {
      getCustomIocTypeVal(newVal, oldVal)
    }
  )
  onMounted(() => {
    form.reliable = props.infoCustomOptions.reliable[0].dictValue
    form.threatType = props.infoCustomOptions.threatType[0].dictValue
    form.threatLevel = props.infoCustomOptions.threatLevel[0].dictValue
    iocOptions.value = {
      reliable: props.infoCustomOptions.reliable.reduce(
        (obj: any, item) => Object.assign(obj, { [item.dictValue]: item.dictLabel }),
        {}
      ),
      threatLevel: props.infoCustomOptions.threatLevel.reduce(
        (obj: any, item) => Object.assign(obj, { [item.dictValue]: item.dictLabel }),
        {}
      ),
    }
    getCustomIocTypeVal(props.infoCustomOptions.iocType[0].dictValue)
    getInfoCustomPageHandle()
  })
</script>

<template>
  <div class="info-custom-container">
    <div style="margin-bottom: 14px">
      <el-input
        v-model="queryForm.searchStr"
        placeholder="请输入域名、IP、URL、文件MD5等IOC进行查询"
        :prefix-icon="Search"
        size="large"
        style="width: calc(100% - 100px); border-radius: 2px"
      />
      <el-button
        :auto-insert-space="false"
        style="width: 90px; height: 36px; margin-left: 10px"
        type="primary"
        @click="getInfoCustomPageHandle"
      >
        检索
      </el-button>
    </div>
    <el-space :size="10" style="width: 100%; justify-content: space-between" wrap>
      <div>
        <el-button :auto-insert-space="false" :icon="Plus" type="primary" @click="addIocVisible = true">新增</el-button>
        <el-button
          plain
          @click="
            () =>
              updateInfoCustomStatusHandle(
                tableSelects.map((i) => i.id),
                1
              )
          "
        >
          批量启用
        </el-button>
        <el-button
          plain
          @click="
            () =>
              updateInfoCustomStatusHandle(
                tableSelects.map((i) => i.id),
                0
              )
          "
        >
          批量停用
        </el-button>
        <el-button plain @click="() => deleteConfirmHandle(tableSelects.map((i) => i.id))">批量删除</el-button>
        <el-button plain @click="handleToApply">应用情报</el-button>
      </div>
      <div style="margin-right: -10px">
        <el-button plain @click="uploadVisible = true">批量导入</el-button>
        <el-button plain @click="exportInfoCustomHandle">批量导出</el-button>
      </div>
    </el-space>
    <el-table
      :data="tableData"
      style="width: 100%; height: calc(100% - 156px); margin-top: 14px"
      @selection-change="handleSelectionChange"
    >
      <el-table-column align="center" fixed="left" type="selection" width="55" />
      <el-table-column align="center" fixed="left" label="情报编号" prop="infoNum" show-overflow-tooltip width="120" />
      <el-table-column align="center" label="IOC" prop="iocStr" show-overflow-tooltip width="120" />
      <el-table-column align="center" label="IOC类型" prop="iocTypeStr" />
      <el-table-column align="center" label="所属组织" prop="organize" show-overflow-tooltip width="120" />
      <el-table-column align="center" label="置信度" prop="reliable">
        <template #default="{ row }">
          <el-rate
            :colors="['#67C23A', '#67C23A', '#67C23A']"
            disabled
            disabled-void-color="#C7C6D4"
            :max="3"
            :model-value="getConfidence(row.reliableStr)"
          />
        </template>
      </el-table-column>
      <el-table-column align="center" label="威胁名称" prop="threatName" show-overflow-tooltip width="120" />
      <el-table-column align="center" label="威胁类型" prop="threatTypeStr" show-overflow-tooltip width="120" />
      <el-table-column align="center" label="威胁等级" prop="threatLevel" width="90">
        <template #default="{ row }">
          <span :class="['level', `level-${row.threatLevel}`]">
            <el-icon><WarnTriangleFilled /></el-icon>
            {{ iocOptions?.threatLevel[row.threatLevel] }}
          </span>
        </template>
      </el-table-column>
      <el-table-column align="center" label="威胁简介" prop="threatDesc" show-overflow-tooltip width="140" />
      <el-table-column align="center" label="启用状态" prop="enable" width="100">
        <template #default="{ row }">
          <el-switch
            v-model="row.enable"
            active-text="是"
            :active-value="1"
            inactive-text="否"
            :inactive-value="0"
            inline-prompt
            @change="(enable) => updateInfoCustomStatusHandle([row.id], enable)"
          />
        </template>
      </el-table-column>
      <el-table-column align="center" label="情报来源" prop="origin" show-overflow-tooltip width="120" />
      <el-table-column
        align="center"
        :formatter="({ createTime }) => formatTime(createTime, 'YYYY-MM-DD HH:mm:ss.SSS')"
        label="创建时间"
        prop="createTime"
        show-overflow-tooltip
        width="120"
      />
      <el-table-column
        align="center"
        :formatter="({ updateTime }) => formatTime(updateTime, 'YYYY-MM-DD HH:mm:ss.SSS')"
        label="更新时间"
        prop="updateTime"
        show-overflow-tooltip
        width="120"
      />
      <el-table-column align="center" fixed="right" label="操作" width="135px">
        <template #default="{ row }">
          <el-button plain size="small" @click="editWhiteHandle(row)">编辑</el-button>
          <el-button plain size="small" @click="deleteConfirmHandle([row.id])">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      v-model:current-page="queryForm.pageNum"
      v-model:page-size="queryForm.pageSize"
      background
      class="pagination"
      layout="total, sizes, prev, pager, next, jumper"
      :page-sizes="[100, 150, 200, 300]"
      :total="total"
      @current-change="getInfoCustomPageHandle"
      @size-change="getInfoCustomPageHandle"
    />
    <el-dialog
      v-model="addIocVisible"
      class="ioc-dialog"
      destroy-on-close
      :title="form.id ? '编辑简易IOC' : '新增简易IOC'"
      width="765px"
      @closed="clearAddDataHandle"
    >
      <el-form ref="iocFormRef" class="ioc-form" label-width="140px" :model="form" :rules="rules">
        <el-form-item label="IOC类型:" prop="iocType" required>
          <el-select v-model="form.iocType">
            <el-option
              v-for="option in infoCustomOptions.iocType"
              :key="option.dictValue"
              :label="option.dictLabel"
              :value="option.dictValue"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-if="['1', '5', '7', '9'].includes(form.iocType)" label="IP地址:" prop="ip" required>
          <el-input v-model="form.ip" />
        </el-form-item>
        <el-form-item v-if="['2', '4', '6', '8'].includes(form.iocType)" label="域名:" prop="domain" required>
          <el-input v-model="form.domain" />
        </el-form-item>
        <el-form-item v-if="['4', '5', '8', '9'].includes(form.iocType)" label="端口:" prop="port" required>
          <el-input v-model="form.port" />
        </el-form-item>
        <el-form-item v-if="['4', '5', '6', '7'].includes(form.iocType)" label="URI:" prop="uri" required>
          <el-input v-model="form.uri" />
        </el-form-item>
        <el-form-item v-if="['3'].includes(form.iocType)" label="MD5:" prop="md5">
          <el-input v-model="form.md5" />
        </el-form-item>
        <el-form-item label="所属组织:" prop="organize">
          <el-input v-model="form.organize" />
        </el-form-item>
        <el-form-item label="置信度:" prop="reliable" required>
          <el-select v-model="form.reliable">
            <el-option
              v-for="option in infoCustomOptions.reliable"
              :key="option.dictValue"
              :label="option.dictLabel"
              :value="option.dictValue"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="情报来源:" prop="origin">
          <el-input v-model="form.origin" />
        </el-form-item>
        <el-form-item label="威胁名称:" prop="threatName" required>
          <el-input v-model="form.threatName" />
        </el-form-item>
        <el-form-item label="威胁类型:" prop="threatType" required>
          <el-select v-model="form.threatType">
            <el-option
              v-for="option in infoCustomOptions.threatType"
              :key="option.dictValue"
              :label="option.dictLabel"
              :value="option.dictValue"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="威胁等级:" required>
          <el-select v-model="form.threatLevel" prop="threatLevel">
            <el-option
              v-for="option in infoCustomOptions.threatLevel"
              :key="option.dictValue"
              :label="option.dictLabel"
              :value="option.dictValue"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="威胁简介:" prop="threatDesc">
          <el-input v-model="form.threatDesc" :autosize="false" :rows="3" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button :auto-insert-space="false" type="primary" @click="handleFormSave">保存</el-button>
          <el-button :auto-insert-space="false" @click="addIocVisible = false">取消</el-button>
        </div>
      </template>
    </el-dialog>

    <IntelligenceCenterWhiteUpload
      v-model:visible="uploadVisible"
      :download-temp="infoCustomExportTempApi"
      :import-assets-fnc="importInfoCustomApi"
      title="导入自定义情报"
      @reflash="getInfoCustomPageHandle"
    />
  </div>
</template>

<style scoped lang="scss">
  .info-custom-container {
    height: 100%;
    padding-inline: 35px;
    .upload-demo {
      margin: 20px;
    }
    :deep() {
      .el-input__inner {
        height: 34px;
      }
    }
    :deep() {
      .level {
        display: inline-block;
        width: 55px;
        padding: 2px 3px;
        border-radius: 5px;
        font-size: 14px;
        background-color: #9d9aba;
        color: #fff;
        .el-icon {
          margin-right: -4px;
          vertical-align: -3px;
          font-size: 17px;
        }
        &-critical {
          background: #ca0a08;
        }
        &-high {
          background: #ff2927;
        }
        &-medium {
          background: #ff7212;
        }
        &-low {
          background: #ffbe36;
        }
      }
      .ioc-dialog {
        .el-select {
          width: 100%;
        }
      }
    }
  }
</style>
