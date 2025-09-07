<script lang="ts">
  export default {
    name: 'SnmpConfig',
  }
</script>
<script setup lang="ts">
  const formData = reactive({
    name: '',
    worknode: '',
    address: '',
    post: '',
    OID: '',
    versions: '',
    credential: '',
    decoder: '',
    code: '',
    description: '',
    timeOut: '',
    pollingTime: '',
    maxLog: '',
    analysisNum: '',
    readNum: '',
    pollingEffect: '',
  })

  const roleOptions = [
    {
      value: 'Option1',
      label: 'Option1',
    },
  ]

  // 取消
  const cancelClick = () => {}

  // 保存
  const confirmClick = () => {}
</script>

<template>
  <div class="snmp-config">
    <el-form
      ref="formRef"
      label-position="right"
      label-width="140px"
      :model="formData"
    >
      <div class="my-title">基本设置</div>
      <el-row>
        <el-col :span="5" />
        <el-col :span="12">
          <el-form-item label="名称：" prop="name">
            <el-input v-model="formData.name" />
          </el-form-item>
          <el-form-item label="工作节点：" prop="worknode">
            <el-input v-model="formData.worknode" />
          </el-form-item>
          <el-form-item label="地址：" prop="address">
            <el-input v-model="formData.address" />
          </el-form-item>
          <el-form-item label="端口：" prop="post">
            <el-input v-model="formData.post" />
          </el-form-item>
          <el-form-item label="OID：" prop="OID">
            <el-input v-model="formData.OID" />
          </el-form-item>
          <el-form-item label="版本：" prop="versions">
            <el-select
              v-model="formData.versions"
              class="m-2"
              placeholder="选择类型"
            >
              <el-option
                v-for="item in roleOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="凭据：" prop="credential">
            <el-input v-model="formData.credential" />
          </el-form-item>
          <el-form-item label="解码器：" prop="decoder">
            <el-select
              v-model="formData.decoder"
              class="m-2"
              placeholder="选择类型"
            >
              <el-option
                v-for="item in roleOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="编码：" prop="code">
            <el-select
              v-model="formData.code"
              class="m-2"
              placeholder="选择类型"
            >
              <el-option
                v-for="item in roleOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="描述：" prop="description">
            <el-input v-model="formData.description" />
          </el-form-item>
        </el-col>
        <el-col :span="9" />
      </el-row>
      <div class="my-title">高级设置</div>
      <el-row>
        <el-col :span="5" />
        <el-col :span="12">
          <el-form-item label="超时时间（ms）：" prop="timeOut">
            <el-input v-model="formData.timeOut" />
          </el-form-item>
          <el-form-item label="轮询时间（ms）：" prop="pollingTime">
            <el-input v-model="formData.pollingTime" />
          </el-form-item>
          <el-form-item label="日志处理上限：" prop="maxLog">
            <el-input v-model="formData.maxLog" />
          </el-form-item>
          <el-form-item label="解析并发个数：" prop="analysisNum">
            <el-input v-model="formData.analysisNum" />
          </el-form-item>
          <el-form-item label="读取并发个数：" prop="readNum">
            <el-input v-model="formData.readNum" />
          </el-form-item>
          <el-form-item label="轮询生效期：" prop="pollingEffect">
            <el-input v-model="formData.pollingEffect" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="confirmClick">保存</el-button>
            <el-button @click="cancelClick">取消</el-button>
          </el-form-item>
        </el-col>
        <el-col :span="9" />
      </el-row>
    </el-form>
  </div>
</template>

<style lang="scss" scoped>
  .snmp-config {
    padding: 10px 20px;

    .my-title {
      margin: 30px 0;
    }
  }
</style>
