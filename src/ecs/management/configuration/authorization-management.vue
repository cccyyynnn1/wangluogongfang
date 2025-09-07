<script lang="ts">
  export default {
    name: 'AuthorizationManagement',
  }
</script>

<script setup lang="ts">
  import { getAuthorization, getActivateKeyApi, toActivateApi } from '~/src/api-ecs/public'
  import type { UploadInstance } from 'element-plus'
  import { formatNstime } from '~/src/utils/time'
  const $baseMessage: any = inject('$baseMessage')

  // 客户名称
  const customerName = ref('')
  // 激活码
  const machineCode = ref('')
  // 特征码
  const featureCode = ref(['', ''])
  const productKey = reactive(['', '', '', '', ''])

  const dialogVisible = ref(false) // 弹框显隐

  // 提交成功
  const handleSuccess = ($val: any) => {
    const { code, msg } = $val
    if (code === 20) {
      $baseMessage($val.msg, 'success', 'vab-hey-message-success')
      handleGetMachineCode()
    } else {
      $baseMessage(msg, 'error', 'vab-hey-message-error')
      upload.value?.clearFiles()
    }
    dialogVisible.value = false
  }

  const upload = ref<UploadInstance>()

  const authorization = ref() // 激活码

  // 上传
  const submitUpload = () => {
    upload.value!.submit()
  }

  // 重新激活
  const resetActive = () => {
    dialogVisible.value = true
  }

  // 获取机器码特征码
  const handleGetProductKey = async () => {
    const { data } = await getActivateKeyApi()
    customerName.value = data.customerName || ''
    featureCode.value = data.featureCode.split('-')
    machineCode.value = data.machineCode
  }
  const productKeyToActivate = () => {
    toActivateApi({
      customerName: customerName.value,
      machineCode: machineCode.value,
      featureCode: featureCode.value.join('-'),
      productKey: productKey.join('-'),
    }).then((res: any) => {
      if (res.code * 1 === 20) {
        $baseMessage(res.msg, 'success', 'vab-hey-message-success')
      }
    })
  }
  // 获取激活code
  const handleGetMachineCode = async () => {
    const { data } = await getAuthorization()
    authorization.value = data
  }
  onMounted(() => {
    handleGetMachineCode()
    handleGetProductKey()
  })
</script>

<template>
  <div class="authorization-management-container">
    <div class="content">
      <vab-query-form>
        <vab-query-form-left-panel :span="12">
          <h3>授权管理</h3>
        </vab-query-form-left-panel>
      </vab-query-form>
      <el-row>
        <el-col :span="8" />
        <el-col :span="8">
          <el-form ref="formRef" label-position="left" label-width="100px">
            <el-form-item label="授权单位" prop="title">
              {{ authorization?.name }}
            </el-form-item>
            <el-form-item label="授权时间" prop="title">
              {{ formatNstime(authorization?.authTime * 1000, false) }}
            </el-form-item>
            <el-form-item label="服务到期时间" prop="title">
              {{ formatNstime(authorization?.expService * 1000, false) }}
            </el-form-item>
            <el-form-item label="使用到期时间" prop="title">
              {{ formatNstime(authorization?.expUsed * 1000, false) }}
            </el-form-item>
            <el-form-item label="机器码" prop="title">
              {{ authorization?.machineCode }}
            </el-form-item>
          </el-form>
        </el-col>
        <el-col :span="8" />
      </el-row>
      <div class="btn">
        <el-button type="primary" @click="resetActive">重新激活</el-button>
      </div>
    </div>
    <el-dialog v-model="dialogVisible" title="上传新许可" width="35%">
      <el-tabs class="active-tabs" :model-value="'file'">
        <el-tab-pane label="文件激活" name="file">
          <el-upload
            ref="upload"
            action="/v3/eht/public/fileActivation"
            :auto-upload="false"
            drag
            :limit="1"
            name="filedata"
            :on-success="handleSuccess"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              将文件拖到此处，或
              <em>点击上传</em>
              <br />
              只能上传bin文件，且不超过500b
            </div>
          </el-upload>
          <div style="margin: 20px 0 10px; text-align: center">
            <el-button type="primary" @click="submitUpload">激活</el-button>
          </div>
        </el-tab-pane>
        <el-tab-pane label="序列号激活" name="productKey">
          <el-form class="activeForm" label-position="right" label-width="100px">
            <el-form-item label="客户名称:">
              <el-input v-model="customerName" />
            </el-form-item>
            <el-form-item label="机器码">
              <el-input readonly :value="machineCode" />
            </el-form-item>
            <el-form-item label="特征码">
              <el-space>
                <el-input readonly :value="featureCode[0]" />
                <span>-</span>
                <el-input readonly :value="featureCode[1]" />
              </el-space>
            </el-form-item>
            <el-form-item label="产品密钥">
              <el-space>
                <el-input v-model="productKey[0]" maxlength="5" />
                <span>-</span>
                <el-input v-model="productKey[1]" maxlength="5" />
                <span>-</span>
                <el-input v-model="productKey[2]" maxlength="5" />
                <span>-</span>
                <el-input v-model="productKey[3]" maxlength="5" />
                <span>-</span>
                <el-input v-model="productKey[4]" maxlength="5" />
              </el-space>
            </el-form-item>
            <el-form-item label=" ">
              <el-button type="primary" @click="productKeyToActivate">激活设备</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>
  </div>
</template>

<style lang="scss" scoped>
  .authorization-management-container {
    :deep() {
      .avatar-uploader {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;

        .el-upload__tip {
          margin: 20px 0;
        }
      }
    }
    .content {
      h3 {
        margin-block: 0 0.5em;
      }
      .avatar-uploader .avatar {
        display: block;
        width: 128px;
        height: 128px;
      }

      .btn {
        display: flex;
        align-items: center;
        justify-content: center;
        margin-top: 20px;
      }
    }
  }
</style>
<style>
  .avatar-uploader .el-upload {
    position: relative;
    overflow: hidden;
    cursor: pointer;
    border: 1px dashed var(--el-border-color);
    border-radius: 6px;
    transition: var(--el-transition-duration-fast);
  }

  .avatar-uploader .el-upload:hover {
    border-color: var(--el-color-primary);
  }

  .el-icon.avatar-uploader-icon {
    width: 128px;
    height: 68px;
    font-size: 28px;
    color: #8c939d;
    text-align: center;
  }
</style>
