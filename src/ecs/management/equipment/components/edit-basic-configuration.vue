<script setup lang="ts">
  import { FormInstance } from 'element-plus'
  import { json } from 'node:stream/consumers'
  import { settingVarUpdateVarApi } from '~/src/api-ecs/equipment'
  import AesEncryptCBC from '~/src/utils/crypto'

  import { requireRules } from '~/src/utils/rules'

  const formRef = ref<FormInstance>() // 表单实例

  const props = defineProps<{
    mode: string
    showPage: boolean
    currentData: object
  }>()

  const title = ref('添加') // 标题

  const mode = ref() // 标题

  const visible = ref(false) // 弹框显隐

  // 表单数据
  const formData = reactive({
    value: {
      addr: '',
      account: '',
      password: '',
      key: '',
      model: '1',
      port: '',
    },
    key: '',
    name: '',
    id: '' as string | number,
  })

  const formData1 = reactive({
    value: [
      {
        addr: '',
        account: '',
        password: '',
        key: '',
        model: '1',
        port: '',
      },
    ],
    key: '',
    name: '',
    id: '' as string | number,
  })

  // 表单数据校验规则
  const rules = reactive({
    addr: requireRules,
    port: requireRules,
  })

  const emit = defineEmits<{
    (e: 'on-close-event'): void
    (e: 'on-reflash'): void
  }>()

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
              try {
                let id = 0
                let value = ''
                if (mode.value === 'ecs_kafka') {
                  id = formData1.id as number
                  value = JSON.stringify(formData1.value)
                } else {
                  id = formData.id as number
                  value = JSON.stringify(formData.value)
                }
                const { msg } = await settingVarUpdateVarApi({ id, value }, { password: AesEncryptCBC(password.value) })
                ElMessage({ message: msg, type: 'success' })
                handleClose()
                emit('on-reflash')
                done()
              } catch (error) {
                instance.confirmButtonLoading = false
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

  // 关闭回调
  const handleClose = () => {
    emit('on-close-event')
  }

  const initData = () => {
    mode.value = props.mode
    visible.value = props.showPage
    Object.keys(formData).forEach((item) => {
      if (mode.value === 'ecs_kafka') {
        // @ts-ignore
        formData1[item] = JSON.parse(JSON.stringify(props.currentData[mode.value][item]))
      } else {
        // @ts-ignore
        formData[item] = JSON.parse(JSON.stringify(props.currentData[mode.value][item]))
      }
      if (item == 'value') {
        if (mode.value === 'ecs_kafka') {
          formData1[item][0].model = formData1[item][0].model ? formData1[item][0].model : '1'
        } else {
          formData[item].model = formData[item].model ? formData[item].model : '1'
        }
      }
    })
    switch (mode.value) {
      case 'yunche_ras':
        title.value = '编辑数据节点'
        break
      case 'ecs_kafka':
        title.value = 'KAFKA编辑'
        break
      case 'ecs_storage':
        title.value = '存储库编辑'
        break
      case 'ecs_neo4j':
        title.value = '图数据库编辑'
        break

      default:
        break
    }
  }

  onMounted(() => {
    initData()
  })
</script>

<script lang="ts">
  export default {
    name: 'EditBasicConfiguration', // 编辑基础配置
  }
</script>
<template>
  <div class="edit-basic-configuration">
    <el-dialog v-model="visible" :before-close="handleClose" :title="title" width="600px">
      <el-form v-if="mode != 'ecs_kafka'" ref="formRef" label-width="100px" :model="formData.value" :rules="rules">
        <el-form-item label="配置类型：" prop="model">
          <el-radio-group v-model="formData.value.model" class="ml-4">
            <el-radio label="1">默认配置</el-radio>
            <el-radio label="2">自定义配置</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="IP：" prop="addr">
          <el-input v-if="formData.value.model === '1'" :disabled="formData.value.model === '1'" show-word-limit />
          <el-input v-else v-model="formData.value.addr" show-word-limit />
        </el-form-item>
        <el-form-item v-if="mode != 'yunche_ras'" label="端口：" prop="port">
          <el-input v-if="formData.value.model === '1'" disabled show-word-limit />
          <el-input v-else v-model="formData.value.port" show-word-limit />
        </el-form-item>
        <el-form-item v-if="mode != 'yunche_ras'" label="账号：" prop="account">
          <el-input v-if="formData.value.model === '1'" disabled show-word-limit />
          <el-input v-else v-model="formData.value.account" show-word-limit />
        </el-form-item>
        <el-form-item v-if="mode != 'yunche_ras'" label="密码：" prop="password">
          <el-input v-if="formData.value.model === '1'" disabled show-word-limit />
          <el-input v-else v-model="formData.value.password" show-word-limit />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submitForm(formRef)">保存</el-button>
        </el-form-item>
      </el-form>
      <el-form v-else ref="formRef" label-width="100px" :model="formData1.value[0]" :rules="rules">
        <el-form-item label="配置类型：" prop="model">
          <el-radio-group v-model="formData1.value[0].model" class="ml-4">
            <el-radio label="1">默认配置</el-radio>
            <el-radio label="2">自定义配置</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="IP：" prop="addr">
          <el-input v-if="formData1.value[0].model === '1'" disabled show-word-limit />
          <el-input v-else v-model="formData1.value[0].addr" show-word-limit />
        </el-form-item>
        <el-form-item label="端口：" prop="port">
          <el-input v-if="formData1.value[0].model === '1'" disabled show-word-limit />
          <el-input v-else v-model="formData1.value[0].port" show-word-limit />
        </el-form-item>
        <el-form-item label="账号：" prop="account">
          <el-input v-if="formData1.value[0].model === '1'" disabled show-word-limit />
          <el-input v-else v-model="formData1.value[0].account" show-word-limit />
        </el-form-item>
        <el-form-item label="密码：" prop="password">
          <el-input v-if="formData1.value[0].model === '1'" disabled show-word-limit />
          <el-input v-else v-model="formData1.value[0].password" show-word-limit />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submitForm(formRef)">保存</el-button>
        </el-form-item>
      </el-form>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss"></style>
