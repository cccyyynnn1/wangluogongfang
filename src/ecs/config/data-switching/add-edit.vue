<script lang="ts">
  export default {
    name: 'AddEdit',
  }
</script>
<script setup lang="ts">
  import { requireRules } from '~/src/utils/rules'

  import { FormInstance } from 'element-plus'

  import { transpondSaveOrUpdateType, IndexTypeTpye } from '@/types'

  import { transpondSaveOrUpdateApi } from '~/src/api-ecs/transpond'

  import { kafkagetAllApi } from '~/src/api-ecs/kafka'

  import { useUserStore } from '@/store/modules/user'

  const { getAllIndexType } = useUserStore()

  const tagType = ref<IndexTypeTpye[]>(getAllIndexType())

  const props = defineProps<{
    showEdit: boolean
    mode: string
    currentData?: transpondSaveOrUpdateType
  }>()

  const emit = defineEmits<{
    (e: 'on-closeEvent', val: boolean): void
    (e: 'on-reflash'): void
  }>()

  const visible = ref(false)

  const mode = ref('')

  const formRef = ref<FormInstance>() // 表单实例

  // 表单数据
  const formData = reactive<transpondSaveOrUpdateType>({
    name: '',
    description: '',
    matchTopic: 1,
    kafkaConfigId: null,
    forwardTopic: '',
    filters: [],
    forwardFields: '',
    siteSessionId: null,
    id: null,
  })

  // 表单数据校验规则
  const rules = reactive({
    name: requireRules,
    matchTopic: requireRules,
    kafkaConfigId: requireRules,
    forwardTopic: requireRules,
  })

  const options = ref()

  // 字段
  const fieldList = [
    '请求时间',
    '请求时间',
    '请求时间',
    '请求时间',
    '请求时间',
    '请求时间',
    '请求时间',
    'URL',
    'URL',
    'URL',
    'URL',
    'URL',
    'URL',
    'URL',
    '分析设备ID',
    '分析设备ID',
    '分析设备ID',
    '分析设备ID',
    '分析设备ID',
    '分析设备ID',
    '分析设备ID',
    '请求接收协议',
    '请求接收协议',
    '请求接收协议',
    '请求接收协议',
    '请求接收协议',
    '请求接收协议',
    '请求接收协议',
    '请求其他数据',
    '请求其他数据',
    '请求其他数据',
    '请求其他数据',
    '请求其他数据',
    '请求其他数据',
    '请求其他数据',
    '用户代理',
    '用户代理',
    '用户代理',
    '用户代理',
    '用户代理',
    '用户代理',
    '用户代理',
    '响应时间',
    '响应时间',
    '响应时间',
    '响应时间',
    '响应时间',
    '响应时间',
    '响应时间',
    '响应其他数据',
    '响应其他数据',
    '响应其他数据',
    '响应其他数据',
    '响应其他数据',
    '响应其他数据',
    '响应其他数据',
    '测试新增字段',
    '测试新增字段',
    '测试新增字段',
    '测试新增字段',
    '测试新增字段',
    '测试新增字段',
    '测试新增字段',
    '测试字段',
    '测试字段',
    '测试字段',
    '测试字段',
    '测试字段',
    '测试字段',
    '测试字段',
  ]

  onMounted(() => {
    initData()
    getAllKafka()
  })

  // 回显数据
  const initData = () => {
    visible.value = props.showEdit
    mode.value = props.mode
    if (props.mode == 'edit') {
      const data = JSON.parse(JSON.stringify(props.currentData))
      Object.keys(formData).forEach((item) => {
        // @ts-ignore
        if (item == 'matchTopic') {
          formData[item] = Number(data[item])
        } else if (item == 'kafkaConfigId') {
          formData[item] = String(data[item])
        } else {
          // @ts-ignore
          formData[item] = data[item]
        }
      })
    }
  }

  // 保存
  const confirmClick = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
      if (valid) {
        const { msg } = await transpondSaveOrUpdateApi({ ...formData })
        emit('on-reflash')
        ElMessage({ message: msg, type: 'success' })
        handleClose()
      } else {
        console.log('error submit!', fields)
      }
    })
  }

  // 获取所有kafka配置信息
  const getAllKafka = async () => {
    const { data } = await kafkagetAllApi()
    options.value = []
    if (data) {
      Object.keys(JSON.parse(data)).forEach((item) => {
        const obj = {
          label: '',
          value: '',
        }
        obj.value = item
        obj.label = JSON.parse(data)[item]
        options.value.push(obj)
      })
    }
  }

  // 重置表单
  const resetForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.resetFields()
  }

  const handleClose = () => {
    resetForm(formRef.value)
    emit('on-closeEvent', false)
  }
</script>

<template>
  <div class="add-edit">
    <el-dialog v-model="visible" :before-close="handleClose" :title="mode === 'add' ? '添加' : '编辑'" width="650px">
      <el-form ref="formRef" label-position="right" label-width="120px" :model="formData" :rules="rules">
        <el-form-item label="名称：" prop="name">
          <el-input v-model="formData.name" />
        </el-form-item>
        <el-form-item label="描述：" prop="description">
          <el-input v-model="formData.description" />
        </el-form-item>
        <el-form-item label="日志类型：" prop="matchTopic">
          <el-select v-model="formData.matchTopic" class="m-2" style="width: 100%">
            <el-option v-for="item in tagType" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <!-- <el-row :gutter="20">
          <el-col :span="12"> -->
        <el-form-item label="KAFKA服务：" prop="kafkaConfigId">
          <el-select v-model="formData.kafkaConfigId" class="m-2" style="width: 100%">
            <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <!-- </el-col>
          <el-col :span="12"> -->
        <el-form-item label="TOPIC：" prop="forwardTopic">
          <el-input v-model="formData.forwardTopic" />
        </el-form-item>
        <!-- </el-col>
        </el-row> -->
        <!-- <el-form-item label="匹配规则：" prop="filters">
                                                          <div class="regulation">
                                                            <template v-for="(item, index) in formData.filters" :key="item.id">
                                                              <el-row class="row" :gutter="20">
                                                                <el-col :span="7">
                                                                  <el-select v-model="item.source" class="m-2" style="width: 100%">
                                                                    <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
                                                                  </el-select>
                                                                </el-col>
                                                                <el-col :span="7">
                                                                  <el-select v-model="item.relation" class="m-2" style="width: 100%">
                                                                    <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
                                                                  </el-select>
                                                                </el-col>
                                                                <el-col :span="7">
                                                                  <el-input v-model="item.host" />
                                                                </el-col>
                                                                <el-col :span="3">
                                                                  <div class="right">
                                                                    <el-icon v-if="formData.filters.length - 1 === index" class="icon">
                                                                      <CirclePlus />
                                                                    </el-icon>
                                                                    <el-icon class="icon">
                                                                      <CircleClose />
                                                                    </el-icon>
                                                                  </div>
                                                                </el-col>
                                                              </el-row>
                                                            </template>
                                                          </div>
                                                        </el-form-item> -->
        <!-- <el-form-item label="转发字段：">
                                                          <div class="filed">
                                                            <el-checkbox-group v-model="formData.forwardFields">
                                                              <el-row>
                                                                <el-col v-for="(item, index) in fieldList" :key="index" :span="3">
                                                                  <el-checkbox :label="item" />
                                                                </el-col>
                                                              </el-row>
                                                            </el-checkbox-group>
                                                          </div>
                                                        </el-form-item> -->
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="confirmClick(formRef)">确定</el-button>
          <el-button type="primary" @click="handleClose">取消</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  .filed {
    width: 100%;
    height: 300px;
    padding: 10px 20px;
    border: 1px solid var(--el-border-color);
  }

  .regulation {
    width: 100%;
    padding: 20px;
    border: 1px solid var(--el-border-color);
  }

  .right {
    height: 100%;
    display: flex;
    align-items: center;
  }

  .row {
    margin-bottom: 20px;

    &:last-child {
      margin-bottom: 0px;
    }
  }

  .icon {
    font-size: 20px;
    margin-right: 6px;
    color: #44a2fd;

    &:nth-child(2) {
      color: #bbbbbb;
    }

    &:hover {
      cursor: pointer;
    }
  }
</style>
