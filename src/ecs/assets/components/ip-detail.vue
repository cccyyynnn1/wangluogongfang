<script setup lang="ts">
  import { getKeyVarApi } from '@/api-ecs/retrieve'
  import ServiceVisits from './ip-detail-service-visits.vue'

  import { AssetsIp_Detail, IpDetailTabsValue } from '../type'

  import { injectStrict } from '@/utils/inject'

  import { FormInstance } from 'element-plus'

  import { updateAssetsApi } from '@/api-ecs/assets'

  const { detailVisible, isEdit, currentRow } = injectStrict(AssetsIp_Detail)

  const callback = inject('reflashData')

  const detailTabsValue = ref<IpDetailTabsValue>('serviceVisits')

  const formRef = ref<FormInstance>() // 表单实例
  const formData = reactive({
    name: '',
    assetsIp: '',
    assetsGroupId: '',
    assetsTypeIds: '',
    assetsTagIds: '',
    hostName: '',
    inManageIp: '',
    outManageIp: '',
    macAddr: '',
    sys: '',
    coreVersion: '',
    disk: '',
    cpu: '',
    ram: '',
    responsible: '',
    phone: '',
  })

  const idEditForm = ref(false)

  const submitForm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
      if (valid) {
        const { msg } = await updateAssetsApi(
          // @ts-ignore
          { ...formData, id: currentRow.value.id }
        )
        ElMessage({ message: msg, type: 'success' })
        // @ts-ignore
        callback()
        detailVisible.value = false
      } else {
        console.log('error submit!', fields)
      }
    })
  }
  const initData = () => {
    // console.log(currentRow.value)
    Object.keys(formData).forEach((item) => {
      // @ts-ignore
      formData[item] = currentRow.value[item]
    })
  }
  onMounted(() => {})
  watchEffect(() => {
    idEditForm.value = isEdit.value
    initData()
  })
</script>

<script lang="ts">
  export default {
    name: 'AssetsIpDetail',
  }
</script>

<template>
  <div class="ip-detail">
    <div>
      <el-tabs v-model="detailTabsValue">
        <el-tab-pane label="服务访问" name="serviceVisits">
          <service-visits />
        </el-tab-pane>
        <el-tab-pane :class="{ 'detail-tab-pane': idEditForm }" label="信息" name="overview">
          <el-form ref="formRef" :disabled="!idEditForm" label-width="120px" :model="formData">
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
            <el-row v-if="idEditForm">
              <el-col :offset="20" :span="4" style="text-align: right">
                <el-button type="primary" @click="submitForm(formRef)">确认</el-button>
                <el-button @click="idEditForm = false">取消</el-button>
              </el-col>
            </el-row>
          </el-form>
          <el-icon class="detail-form-edit" @click="idEditForm = true">
            <Edit />
          </el-icon>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<style scoped lang="scss">
  #pane-overview {
    position: relative;
    margin-top: 10px;

    .detail-form-edit {
      position: absolute;
      right: 0;
      top: 0;
      font-size: 20px;
      cursor: pointer;
    }

    &.detail-tab-pane {
      background: #f8fafd;

      .detail-form-edit {
        display: none;
      }
    }

    .el-form {
      padding: 42px 54px 40px 0;
    }
  }
</style>
