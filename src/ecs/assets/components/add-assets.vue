<script setup lang="ts">
  import { FormInstance } from 'element-plus'

  import { updateAssetsApi } from '~/src/api-ecs/assets'

  import { requireRules } from '~/src/utils/rules'

  const props = defineProps<{
    showPage: boolean
  }>()

  const formRef = ref<FormInstance>() // 表单实例

  const visible = ref(false) // 弹框显隐

  // 表单数据
  const formData = reactive({
    name: '',
    assetsIp: undefined,
    // assetsIpv6: undefined,
    hostName: undefined,
    inManageIp: undefined,
    outManageIp: undefined,
    macAddr: undefined,
    coreVersion: undefined,
    disk: undefined,
    cpu: undefined,
    ram: undefined,
    responsible: undefined,
    phone: undefined,
  })

  // 校验assetsIpv4和assetsIpv6的规则
  // const validateRule = (rule: any, value: any, callback: any) => {
  //   if (!formData.assetsIpv4 && !formData.assetsIpv6) {
  //     callback(new Error('IPV4和IPV6至少填写一项'))
  //   } else {
  //     formRef.value?.clearValidate(['assetsIpv4', 'assetsIpv6'])
  //     callback()
  //   }
  // }

  // 表单数据校验规则
  const rules = reactive({
    name: requireRules,
    assetsIp: requireRules,
    // assetsIpv4: [{ validator: validateRule, trigger: 'blur' }],
    // assetsIpv6: [{ validator: validateRule, trigger: 'blur' }],
  })

  const emit = defineEmits<{
    (e: 'on-close-event'): void
    (e: 'on-reflash'): void
  }>()

  // 关闭回调
  const handleClose = () => {
    emit('on-close-event')
  }

  // 提交表单
  const submitForm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
      if (valid) {
        const { msg } = await updateAssetsApi({ ...formData })
        setTimeout(() => {
          emit('on-reflash')
          ElMessage({ message: msg, type: 'success' })
          handleClose()
        }, 0)
      } else {
        console.log('error submit!', fields)
      }
    })
  }

  onMounted(() => {
    visible.value = props.showPage
  })
</script>

<script lang="ts">
  export default {
    name: 'AddAssets',
  }
</script>
<template>
  <div class="add-assets">
    <el-dialog v-model="visible" :before-close="handleClose" title="添加资产" width="1180px">
      <el-form ref="formRef" label-width="120px" :model="formData" :rules="rules">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="资产名称" prop="name">
              <el-input v-model="formData.name" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="资产IP" prop="assetsIp">
              <el-input v-model="formData.assetsIp" />
            </el-form-item>
          </el-col>
          <!-- <el-col :span="12">
            <el-form-item label="资产IVP6" prop="assetsIpv6">
              <el-input v-model="formData.assetsIpv6" placeholder="IPV4和IPV6至少填写一项" />
            </el-form-item>
          </el-col> -->
          <el-col :span="12">
            <el-form-item label="主机名" prop="hostName">
              <el-input v-model="formData.hostName" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="带内管理IP" prop="inManageIp">
              <el-input v-model="formData.inManageIp" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="带外管理IP" prop="outManageIp">
              <el-input v-model="formData.outManageIp" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="MAC地址" prop="macAddr">
              <el-input v-model="formData.macAddr" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="内核版本" prop="coreVersion">
              <el-input v-model="formData.coreVersion" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="磁盘" prop="disk">
              <el-input v-model="formData.disk" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="CPU" prop="cpu">
              <el-input v-model="formData.cpu" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="内存" prop="ram">
              <el-input v-model="formData.ram" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="负责人" prop="responsible">
              <el-input v-model="formData.responsible" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="负责人电话" prop="phone">
              <el-input v-model="formData.phone" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :offset="20" :span="4" style="text-align: right">
            <el-button type="primary" @click="submitForm(formRef)">确认</el-button>
            <el-button @click="handleClose">取消</el-button>
          </el-col>
        </el-row>
      </el-form>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss"></style>
