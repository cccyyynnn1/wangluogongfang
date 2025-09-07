<script lang="ts">
  export default {
    name: 'AddKafka',
  }
</script>
<script setup lang="ts">
  import { requireRules } from '~/src/utils/rules'

  import { FormInstance } from 'element-plus'

  import { kafkaSaveOrUpdateApi, CheckKafkaApi } from '~/src/api-ecs/kafka'

  import { kafkaSaveOrUpdateType } from '@/types'

  // 表单数据
  const formData = reactive<kafkaSaveOrUpdateType>({
    name: '',
    description: '',
    addressPorts: '',
    host: '',
    acount: '',
    password: '',
  })

  const formRef = ref<FormInstance>() // 表单实例
  const isDisabled = ref(true)

  // 表单数据校验
  const rules = reactive({
    name: requireRules,
    addressPorts: requireRules,
    host: requireRules,
    acount: requireRules,
    password: requireRules,
  })

  const props = defineProps<{
    showDrawer: boolean
    mode: string
    currentData: kafkaSaveOrUpdateType
  }>()

  const emit = defineEmits<{
    (e: 'on-closeEvent'): void
    (e: 'on-reflash'): void
  }>()

  const drawer = ref<boolean>()

  onMounted(() => {
    initData()
  })

  // 回显数据
  const initData = () => {
    drawer.value = props.showDrawer
    if (props.mode == 'edit') {
      const data = JSON.parse(JSON.stringify(props.currentData))
      Object.keys(formData).forEach((item) => {
        // @ts-ignore
        formData[item] = data[item]
      })
      if (formData.addressPorts) {
        const arr = formData.addressPorts.split(':')
        formData.addressPorts = arr[0]
        formData.host = arr[1]
      }
    }
  }

  // 保存
  const confirmClick = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
      if (valid) {
        const { msg } = await kafkaSaveOrUpdateApi({
          ...formData,
          id: props.mode == 'edit' ? props.currentData.id : '',
          addressPorts: `${formData.addressPorts}:${formData.host}`,
        })
        emit('on-reflash')
        ElMessage({ message: msg, type: 'success' })
        closeEvent()
      } else {
        console.log('error submit!', fields)
      }
    })
  }

  const handleConfirm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
      if (valid) {
        const { msg, code } = await CheckKafkaApi({
          ip: formData.host,
          post: formData.addressPorts,
          userName: formData.acount,
          password: formData.password,
        })
        ElMessage({ message: msg, type: 'success' })
        if (code == 20) isDisabled.value = false
      } else {
        console.log('error submit!', fields)
      }
    })
  }

  const closeEvent = () => {
    emit('on-closeEvent')
  }
</script>

<template>
  <div class="add-kafka">
    <el-drawer v-model="drawer" size="35%" @close="closeEvent">
      <template #header>
        <h4>{{ props.mode == 'edit' ? '编辑Kafka' : '新建Kafka' }}</h4>
      </template>
      <template #default>
        <el-form ref="formRef" label-position="right" label-width="160px" :model="formData" :rules="rules">
          <div class="my-title"><span>基本信息</span></div>
          <el-form-item label="Kafka名称：" prop="name">
            <el-input v-model="formData.name" />
          </el-form-item>
          <el-form-item label="描述：" prop="description">
            <el-input v-model="formData.description" />
          </el-form-item>
          <el-form-item label="地址：" prop="addressPorts">
            <el-input v-model="formData.addressPorts" />
          </el-form-item>
          <el-form-item label="端口：" prop="host">
            <el-input v-model="formData.host" />
          </el-form-item>
          <div class="my-title"><span>鉴权信息</span></div>
          <el-form-item label="账户：" prop="acount">
            <el-input v-model="formData.acount" />
          </el-form-item>
          <el-form-item label="密码：" prop="password">
            <el-input v-model="formData.password" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleConfirm(formRef)">连通性测试</el-button>
            <el-button :disabled="isDisabled" type="primary" @click="confirmClick(formRef)">保存</el-button>
            <el-button @click="closeEvent">取消</el-button>
          </el-form-item>
        </el-form>
      </template>
    </el-drawer>
  </div>
</template>

<style lang="scss" scoped>
  .add-kafka {
    padding-right: 40px;

    .my-title {
      margin-bottom: 30px;

      &::before {
        display: inline-block;
        width: 3px;
        height: 10px;
        margin: 0 6px 0 30px;
        content: '';
        background: #1890ff;
      }
    }
  }
</style>
