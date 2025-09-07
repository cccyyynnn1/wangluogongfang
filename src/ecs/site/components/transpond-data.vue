<script lang="ts">
  export default {
    name: 'TranspondData',
  }
</script>
<script setup lang="ts">
  import { kafkagetAllApi } from '~/src/api-ecs/kafka'

  import { requireRules } from '~/src/utils/rules'

  import { forwordRuleUpdateType } from '@/types'

  import { FormInstance } from 'element-plus'

  import { forwordRuleUpdateApi } from '@/api-ecs/forword-rule'

  const props = defineProps<{
    showTranspond: boolean
    siteSessionId: number | string
  }>()

  const visible = ref(false)

  // 表单数据
  const formData: forwordRuleUpdateType = reactive({
    name: '',
    description: '',
    kafkaConfigId: '',
    forwardTopic: '',
    filters: [],
    matchTopic: 1,
    siteSessionId: 0,
  })

  const rules = reactive({
    name: requireRules,
    kafkaConfigId: requireRules,
    forwardTopic: requireRules,
  })

  const formRef = ref<FormInstance>() // 表单实例

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

  const options = ref()

  onMounted(() => {
    visible.value = props.showTranspond
    formData.siteSessionId = props.siteSessionId
    getAllKafka()
  })

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

  const emit = defineEmits<{
    (e: 'on-closeEvent', val: boolean): void
  }>()

  // 确定转发
  const submit = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
      if (valid) {
        const { msg } = await forwordRuleUpdateApi({ ...formData })
        ElMessage({ message: msg, type: 'success' })
        handleClose()
      } else {
        console.log('error submit!', fields)
      }
    })
  }

  const handleClose = () => {
    emit('on-closeEvent', false)
  }
</script>

<template>
  <div class="transpond-data">
    <el-dialog v-model="visible" :before-close="handleClose" title="转发" width="650px">
      <el-form ref="formRef" label-position="right" label-width="120px" :model="formData" :rules="rules">
        <el-form-item label="名称：" prop="name">
          <el-input v-model="formData.name" />
        </el-form-item>
        <el-form-item label="描述：" prop="description">
          <el-input v-model="formData.description" />
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
        <el-form-item v-if="false" label="转发自段：">
          <div class="filed">
            <el-checkbox-group>
              <el-row>
                <el-col v-for="(item, index) in fieldList" :key="index" :span="3">
                  <el-checkbox :label="item" />
                </el-col>
              </el-row>
            </el-checkbox-group>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="submit(formRef)">确认</el-button>
          <el-button @click="handleClose">取消</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  .transpond-data {
    .filed {
      width: 100%;
      height: 300px;
      padding: 10px 20px;
      border: 1px solid var(--el-border-color);
    }
  }
</style>
