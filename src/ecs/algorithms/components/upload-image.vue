<script setup lang="ts">
  import { UploadFilled } from '@element-plus/icons-vue'

  import type { UploadUserFile, UploadProps, FormInstance } from 'element-plus'

  const formRef = ref<FormInstance>()

  const state = reactive<{
    algorName: string
    // runParam: string
    description?: string
    upload: UploadUserFile[]
  }>({
    algorName: '',
    // runParam: '',
    upload: [],
    description: undefined,
  })

  const rules = reactive({
    algorName: [
      {
        required: true,
        message: '请输入算法名称',
        trigger: 'blur',
      },
    ],
    upload: [
      {
        required: true,
        message: '请选择上传文件',
        trigger: 'change',
      },
    ],
  })

  const submitForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.validate((valid) => {
      if (valid) {
        console.log('submit!!')
      } else {
        console.log('error submit!')
        return false
      }
    })
  }

  const handleExceed: UploadProps['onExceed'] = (files) => {
    state.upload = files
  }

  defineExpose({
    uploadImageRef: formRef,
    state,
  })
</script>

<script lang="ts">
  export default {
    name: 'AlgorithmsUploadImage',
  }
</script>

<template>
  <el-form ref="formRef" class="examplesDetailUpload-form" label-width="100px" :model="state" :rules="rules">
    <el-form-item label="算法名称" prop="algorName">
      <el-input v-model="state.algorName" />
    </el-form-item>
    <el-form-item label="描述" prop="description">
      <el-input v-model="state.description" />
    </el-form-item>
    <!-- <el-form-item label="启动参数" prop="runParam">
      <el-input v-model="state.runParam" />
    </el-form-item> -->
    <el-form-item label="上传文件" prop="upload">
      <el-upload
        v-model:file-list="state.upload"
        accept=".tar,"
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
          <br />
          只能上传tar文件
        </div>
      </el-upload>
    </el-form-item>
  </el-form>
</template>

<style scoped lang="scss">
  .upload-demo {
    width: 550px;
  }
</style>
