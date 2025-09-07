<script setup lang="ts">
  import { UploadFilled } from '@element-plus/icons-vue'
  import { genFileId } from 'element-plus'
  import type { UploadUserFile, UploadProps, FormInstance } from 'element-plus'
  import { dockerCopyApi } from '~/src/api-ecs/algorithms'
  import { requireAndNumberRules, requireRules } from '~/src/utils/rules'

  const props = defineProps<{
    currentItem: any
  }>()

  const emit = defineEmits<{
    (e: 'closeEvent'): void
  }>()

  const formRef = ref<FormInstance>()

  const state = reactive<{
    container: string
    permission: number | undefined
    target: string
    filedata: UploadUserFile[]
  }>({
    container: '',
    permission: undefined,
    filedata: [],
    target: '',
  })
  const rules = reactive({
    container: requireRules,
    target: requireRules,
    permission: requireAndNumberRules,
    filedata: [
      {
        required: true,
        message: '请选择上传文件',
        trigger: 'change',
      },
    ],
  })
  const submitForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.validate(async (valid) => {
      if (valid) {
        const { msg } = await dockerCopyApi(
          { container: props.currentItem.Id, permission: state.permission, target: state.target },
          { filedata: state.filedata[0].raw }
        )
        ElMessage({ message: msg, type: 'success' })
        emit('closeEvent')
      } else {
        console.log('error submit!')
        return false
      }
    })
  }

  onMounted(() => {
    initData()
  })

  // 初始化数据
  const initData = () => {
    if (props.currentItem.Names.length > 0) {
      state.container = props.currentItem.Names[0]
    }
  }

  const resetForm = (formEl: FormInstance | undefined) => {
    emit('closeEvent')
    if (!formEl) return
    formEl.resetFields()
  }

  const handleExceed: UploadProps['onExceed'] = (files) => {
    state.filedata = files
  }
</script>

<script lang="ts">
  export default {
    name: 'ExamplesDetailUpload',
  }
</script>

<template>
  <el-form ref="formRef" class="examplesDetailUpload-form" label-width="100px" :model="state" :rules="rules">
    <el-form-item label="容器名称" prop="container">
      <el-input v-model="state.container" disabled />
    </el-form-item>
    <el-form-item label="目标路径" prop="target">
      <el-input v-model="state.target" />
    </el-form-item>
    <el-form-item label="文件权限" prop="permission">
      <el-input v-model="state.permission" />
    </el-form-item>
    <el-form-item label="上传文件" prop="filedata">
      <el-upload
        v-model:file-list="state.filedata"
        :auto-upload="false"
        class="upload-demo"
        drag
        :limit="1"
        :on-exceed="handleExceed"
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          将文件拖到此处，或
          <em>点击上传</em>
        </div>
      </el-upload>
    </el-form-item>
    <el-form-item label="&nbsp">
      <el-button type="primary" @click="submitForm(formRef)">确认</el-button>
      <el-button @click="resetForm(formRef)">取消</el-button>
    </el-form-item>
  </el-form>
</template>

<style scoped lang="scss">
  .examplesDetailUpload-form {
    width: 700px;
    margin: 0 auto;
  }

  .upload-demo {
    width: 300px;
  }
</style>
