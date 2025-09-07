<script lang="ts">
  export default {
    name: 'DrawerMemorise', // 记忆库
  }
</script>

<script setup lang="ts">
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import { FormInstance } from 'element-plus'
  import { requireRules } from '~/src/utils/rules'
  // const props = defineProps<{
  //   alldata: any
  // }>()
  const emits = defineEmits<{
    (e: 'on-close'): void
  }>()

  const dialogVisible = ref(false)
  const multipleSelection = ref<{ id: number }[]>([])
  const dataList = reactive([
    { id: 0, sourceIp: '172.168.2.234', clientIp: '168.123.212.212，168.' },
    { id: 1, sourceIp: '172.168.2.234', clientIp: '168.123.212.212，168.' },
    { id: 2, sourceIp: '172.168.2.234', clientIp: '168.123.212.212，168.' },
    { id: 3, sourceIp: '172.168.2.234', clientIp: '168.123.212.212，168.' },
    { id: 4, sourceIp: '172.168.2.234', clientIp: '168.123.212.212，168.' },
  ])
  const mode = ref('新增')

  const formRef = ref<FormInstance>()
  const formData = reactive({
    type: 0,
    sourceIp: '',
    clientIp: '',
  })
  const rules = {
    type: requireRules,
    sourceIp: requireRules,
    clientIp: requireRules,
  }
  const selectList = ref([
    { value: 0, label: '一对一（“一个源IP”连线“一个目的IP“）' },
    { value: 1, label: '一对多（“一个源IP”连线“多个目的IP“）' },
    { value: 2, label: '多对一（“多个源IP”连线“一个目的IP“）' },
  ])

  const handleUpdate = () => {
    dialogVisible.value = true
    mode.value = '新增'
  }
  const handleDelete = (arr: number[]) => {}
  const handleEdit = (row: number[]) => {}
  const handleSelectionChange = (val: { id: number }[]) => {
    multipleSelection.value = val
  }
  // 提交表单
  const hanldeSave = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
      if (valid) {
        console.log(1)
      } else {
        console.log('error submit!', fields)
      }
    })
  }
</script>

<template>
  <div class="drawer-memorise">
    <div class="drawer-memorise-top">
      <div class="left">
        <el-button plain @click="emits('on-close')">
          <el-icon><Back /></el-icon>
          返回
        </el-button>
        <el-divider direction="vertical" />
        <div class="title">记忆库</div>
      </div>
      <div class="right">
        <el-button color="#6954F0" @click="handleUpdate">新增</el-button>
        <el-button color="#6954F0">导入</el-button>
        <el-button color="#6954F0">导出</el-button>
        <el-dropdown style="margin-left: 10px">
          <span class="el-dropdown-link">
            <el-button type="danger">
              批量删除
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </el-button>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleDelete(multipleSelection.map((item) => item.id))">
                删除选中
              </el-dropdown-item>
              <el-dropdown-item @click="handleDelete([])">删除所有</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
    <el-table class-name="memorise-table" :data="dataList" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" />
      <el-table-column label="序号" type="index" width="75" />
      <el-table-column label="源IP" prop="sourceIp" />
      <el-table-column label="目的IP" prop="clientIp" />
      <el-table-column align="center" fixed="right" label="操作" width="140">
        <template #default="{ row }">
          <el-button size="small" @click="() => handleEdit(row)">编辑</el-button>
          <el-button size="small" @click="() => handleDelete([row.id])">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <vab-dialog
      v-model="dialogVisible"
      align-center
      :close-on-click-modal="false"
      destroy-on-close
      title="新增"
      width="870px"
    >
      <el-form ref="formRef" label-position="top" :model="formData" :rules="rules">
        <el-form-item label="连线类型：" prop="type">
          <el-select v-model="formData.type" clearable collapse-tags style="width: 100%">
            <el-option v-for="item in selectList" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <div class="row_1">
          <div class="row_1_inside">
            <el-form-item label="输入源IP：" prop="sourceIp">
              <el-input v-model="formData.sourceIp" resize="none" style="width: 370px" type="textarea" />
            </el-form-item>
          </div>
          <div class="row_1_arrow"></div>
          <div class="row_1_outside">
            <el-form-item label="输入目的IP：" prop="clientIp">
              <el-input v-model="formData.clientIp" resize="none" style="width: 370px" type="textarea" />
            </el-form-item>
          </div>
        </div>
      </el-form>
      <template #footer>
        <el-button :auto-insert-space="false" type="primary" @click="hanldeSave(formRef)">确认</el-button>
        <el-button :auto-insert-space="false" @click="dialogVisible = false">取消</el-button>
      </template>
    </vab-dialog>
  </div>
</template>

<style scoped lang="scss">
  :deep() {
    .el-textarea {
      height: 400px;
      .el-textarea__inner {
        height: 100%;
      }
    }
  }
  .drawer-memorise {
    position: absolute;
    height: 100%;
    right: 0;
    top: 0;
    padding: 10px 15px;
    background-color: #fff;
    border-left: 1px solid #e9e6f9;
    width: 1030px;
    z-index: 99;
    .drawer-memorise-top {
      display: flex;
      align-items: center;
      justify-content: space-between;
      width: 100%;
      margin-top: 5px;
      .left {
        display: flex;
        align-items: center;
      }
      .title {
        font-weight: 500;
        font-size: 18px;
        color: #303133;
      }
    }
  }
  .memorise-table {
    margin-top: 15px;
  }
  .row_1 {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;

    .row_1_inside,
    .row_1_outside {
      width: 370px;
    }
    .row_1_arrow {
      width: 42px;
      height: 0px;
      border: 1px solid #aea8d4;
      position: relative;
      margin-left: -4px;
      &::after {
        position: absolute;
        content: '';
        top: -8px;
        right: -17px;
        width: 0;
        height: 0;
        border-width: 8px;
        border-style: solid;
        border-color: transparent transparent transparent #8781b3;
      }
    }
  }
</style>
