<script lang="ts">
  export default {
    name: 'AddItem',
  }
</script>
<script setup lang="ts">
  import { requireRules } from '@/utils/rules'

  import type { FormInstance } from 'element-plus'

  import { ElMessage } from 'element-plus'

  import { stationAddApi, stationUpdate } from '~/src/api-ecs/situation'

  const props = defineProps<{
    showEditSite: boolean
    curitem: any
    mode: string
    flowDevice: any
    firewall: any
  }>()

  const emit = defineEmits<{
    (e: 'on-closeEvent', val: boolean): void
    (e: 'on-reflash'): void
  }>()

  const visible = ref(false)

  const formRef = ref<FormInstance>()

  const modes = ref('')

  const title = ref('')

  const flowDeviceData = ref()

  const firewallData = ref()

  // 表单数据
  const formData = reactive<any>({
    id: undefined,
    assetIp: '',
    defenseType: undefined,
    eastLongitude: undefined,
    firewallIps: '',
    northernLatitude: undefined,
    stationName: '',
    xdrasUuid: '',
  })

  // 表单数据校验
  const rules = reactive({
    // assetIp: requireRules,
    defenseType: requireRules,
    eastLongitude: requireRules,
    // firewallIps: requireRules,
    northernLatitude: requireRules,
    stationName: requireRules,
    // xdrasUuid: requireRules,
  })

  // 提交
  const submitForm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
      if (valid) {
        const url = modes.value == 'edit' ? stationUpdate : stationAddApi
        // console.log(formData)
        const { msg } = await url({
          ...formData,
          eastLongitude: +formData.eastLongitude,
          northernLatitude: +formData.northernLatitude,
        })
        ElMessage({ message: msg, type: 'success' })
        emit('on-reflash')
        handleClose()
      } else {
        console.log('error submit!', fields)
      }
    })
  }

  // 重置表单
  const resetForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.resetFields()
  }

  // 关闭
  const handleClose = () => {
    resetForm(formRef.value)
    emit('on-closeEvent', false)
  }

  const initData = async () => {
    visible.value = props.showEditSite
    modes.value = props.mode
    flowDeviceData.value = props.flowDevice
    firewallData.value = props.firewall
    if (modes.value == 'edit') {
      title.value = '编辑'
      for (const key in formData) {
        // @ts-ignore
        formData[key as keyof typeof formData] = props.curitem[key as keyof typeof formData]
      }
    } else {
      title.value = '新增'
      formRef?.value?.resetFields()
    }
  }

  onMounted(() => {
    initData()
  })
</script>

<template>
  <div class="add-item">
    <el-dialog v-model="visible" :before-close="handleClose" :title="title" width="600px">
      <el-form ref="formRef" label-position="right" label-width="110px" :model="formData" :rules="rules">
        <el-form-item label="名称" prop="stationName">
          <el-input v-model="formData.stationName" />
        </el-form-item>
        <el-form-item label="探针" prop="xdrasUuid">
          <el-select v-model="formData.xdrasUuid" style="width: 100%">
            <template v-for="item in flowDeviceData" :key="item.id">
              <el-option :label="item?.ip" :value="item?.id" />
            </template>
          </el-select>
          <!-- <el-input v-model="" /> -->
        </el-form-item>
        <el-form-item v-if="firewallData" label="防火墙" prop="firewallIps">
          <el-select v-model="formData.firewallIps" style="width: 100%">
            <el-option
              v-for="item in firewallData"
              :key="item.firewallIp"
              :label="item.firewallIp"
              :value="item.firewallIp"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="资产IP" prop="assetIp">
          <el-input
            v-model="formData.assetIp"
            placeholder="支持单IP地址、IP子网、IP段，多个地址以英文逗号分隔。示例：192.168.0.1,192.168.1.0/24,192.168.2.0-192.168.2.100"
            resize="none"
            :rows="4"
            type="textarea"
          />
        </el-form-item>
        <el-form-item label="防御类型" prop="defenseType">
          <el-select v-model="formData.defenseType" style="width: 100%">
            <el-option key="blue-black" label="检测" :value="1" />
            <el-option key="blue-white" label="防御" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="X轴" prop="eastLongitude">
          <el-input v-model="formData.eastLongitude" />
        </el-form-item>
        <el-form-item label="Y轴" prop="northernLatitude">
          <el-input v-model="formData.northernLatitude" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="submitForm(formRef)">保存</el-button>
          <el-button @click="handleClose">取消</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  .space {
    width: 100%;
    min-height: 150px;
    padding: 10px;
    border: 1px solid var(--el-border-color);
    border-radius: 2px;
  }

  :deep(.el-tag) {
    margin-left: 10px;
  }
</style>
