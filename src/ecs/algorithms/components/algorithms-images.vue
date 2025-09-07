<script setup lang="ts">
  import VabDialog from '@/plugins/VabDialog/index.vue'

  import AlgorithmsUploadImage from './upload-image.vue'

  import { Search, Plus } from '@element-plus/icons-vue'

  import { algorithmsListType } from '~/src/types'

  import VueEvent from '@/data/event'

  import { useTableCopy } from '@/utils'

  import {
    getAlgorithmsListApi,
    uploadContainersApi,
    deleteContainersApi,
    createContainersApi,
    startContainersApi,
    getAlgorithmsContainersApi,
    stopContainersApi,
    saveUpdateContainersApi,
  } from '~/src/api-ecs/algorithms'

  const $baseConfirm: any = inject('$baseConfirm')

  const $baseMessage: any = inject('$baseMessage')

  const uploadImaageVisible = ref(false)

  const isLoading = ref(false)

  const createContainerVisible = ref(false)

  const runParam = ref('')

  const currentRowId = ref()

  const listDate = ref<any[]>([]) // 表格数据

  // 检索参数
  const queryData = reactive<algorithmsListType>({
    pageNum: 1,
    pageSize: 10,
    algorName: undefined,
    imageName: undefined,
  })

  const total = ref(0) // 总条数

  const listLoading = ref(false) // 是否加载

  const upload_image_ref = ref<InstanceType<typeof AlgorithmsUploadImage>>()

  // 获取表格序号
  const curIndex = computed(() => (queryData.pageNum - 1) * queryData.pageSize + 1)

  const containersList = ref<any>([]) // 所有容器列表

  const containersRunningList = ref<any>([]) // 所有启动容器列表

  // function confirmEvent() {
  //   console.log('confirmEvent')
  // }

  // 获取容器数据
  const getContainersList = async () => {
    const res = await getAlgorithmsContainersApi({ all: true, size: false })
    containersList.value = []
    containersRunningList.value = []
    // @ts-ignore
    res.forEach((item: any) => {
      containersList.value.push(item.Id)
      if (item.State == 'running') {
        containersRunningList.value.push(item.Id)
      }
    })
    // console.log(containersList.value)
    // console.log(containersRunningList.value)
  }

  // 获取数据
  const getData = async () => {
    listLoading.value = true
    const { data } = await getAlgorithmsListApi({ ...queryData })
    total.value = data.total
    listDate.value = data.records
    listLoading.value = false
  }

  // 删除
  const handleDelete = ({ ...row }) => {
    if (row.row) {
      $baseConfirm('你确定要删除当前项吗', null, async () => {
        const { msg } = await deleteContainersApi({ id: row.row.id })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        await getData()
      })
    }
  }

  function onUploadSubmit() {
    // @ts-ignore
    upload_image_ref.value.uploadImageRef.validate(async (valid) => {
      if (valid) {
        isLoading.value = true
        // @ts-ignore
        const { algorName, description, upload } = upload_image_ref.value.state
        // @ts-ignore
        const { data, msg, code } = await uploadContainersApi({ algorName, description }, { file: upload[0].raw })
        const url = data.fileUrl as string
        const imageSize = data.imageSize as number
        msg || ElMessage({ message: msg })
        if (code == 20) {
          const res = await saveUpdateContainersApi({ algorName, description }, url, imageSize)
          ElMessage({ message: res.msg })
        }
        isLoading.value = false
        uploadImaageVisible.value = false
        getData()
      } else {
        console.log('error submit!')
        isLoading.value = false
        uploadImaageVisible.value = false
        return false
      }
    })
  }

  // 改变页面容量
  function handleSizeChange(params: number) {
    queryData.pageSize = params
    getData()
  }

  // 改变页面
  function handleCurrentChange(params: number) {
    queryData.pageNum = params
    getData()
  }

  function onUploadCancel() {
    uploadImaageVisible.value = false
  }

  // 开启容器
  const startContainer = async (row: any) => {
    const { msg } = await startContainersApi({ id: currentRowId.value })
    msg || ElMessage({ message: msg })
    // createContainerVisible.value = false
    getContainersList()
    VueEvent.emit('containerRefresh')
  }

  // 关闭容器
  const stopContainer = async (row: any) => {
    const { msg } = await stopContainersApi({ id: currentRowId.value })
    msg || ElMessage({ message: msg })
    getContainersList()
    // createContainerVisible.value = false
  }

  // 打开创建容器页面
  const openContainer = (row: any) => {
    currentRowId.value = row.id
    createContainerVisible.value = true
  }

  // 创建容器
  const createContainer = async () => {
    const { msg } = await createContainersApi({ id: currentRowId.value, runParam: runParam.value })
    msg || ElMessage({ message: msg })
  }

  onMounted(() => {
    getContainersList()
    getData()
  })
</script>

<script lang="ts">
  export default {
    name: 'AlgorithmsImages',
  }
</script>

<template>
  <el-space alignment="start" :size="0" :style="{ justifyContent: 'space-between', width: '100%' }">
    <el-form class="demo-form-inline" inline :model="queryData">
      <el-form-item label="算法名称">
        <el-input v-model="queryData.algorName" clearable />
      </el-form-item>
      <el-form-item label="镜像名称">
        <el-input v-model="queryData.imageName" clearable />
      </el-form-item>
      <el-form-item>
        <el-button :icon="Search" :loading="listLoading" type="primary" @click="getData">检索</el-button>
      </el-form-item>
    </el-form>
    <el-button :icon="Plus" type="primary" @click="uploadImaageVisible = true">上传镜像</el-button>
  </el-space>
  <el-table v-loading="listLoading" align="center" border :data="listDate" @cell-contextmenu="useTableCopy">
    <el-table-column label="序号" type="index" width="100">
      <template #default="{ $index }">{{ curIndex + $index }}</template>
    </el-table-column>
    <el-table-column align="center" label="算法名称" prop="algorName" />
    <el-table-column align="center" label="镜像名称" prop="imageName" width="200" />
    <el-table-column align="center" label="描述" prop="description" />
    <el-table-column align="center" label="镜像大小" prop="imageSize" width="150" />
    <el-table-column align="center" label="操作" width="250">
      <template #default="{ row }">
        <el-button
          :disabled="row.containerId || containersList.includes(row.containerId)"
          size="small"
          @click="openContainer(row)"
        >
          创建容器
        </el-button>
        <el-button
          v-if="row.containerId || containersRunningList.includes(row.containerId)"
          size="small"
          @click="stopContainer(row)"
        >
          暂停容器
        </el-button>
        <el-button v-else link size="small" type="warning" @click="startContainer(row)">启动容器</el-button>
        <el-button link size="small" type="danger" @click="handleDelete({ row })">删除镜像</el-button>
      </template>
    </el-table-column>
  </el-table>
  <el-pagination
    v-model:current-page="queryData.pageNum"
    v-model:page-size="queryData.pageSize"
    background
    class="known_pagination"
    layout="total, sizes, prev, pager, next, jumper"
    :page-sizes="[10, 20, 30]"
    :total="total"
    @current-change="handleCurrentChange"
    @size-change="handleSizeChange"
  />
  <!-- 创建容器 -->
  <vab-dialog v-model="createContainerVisible" class="startImage" destroy-on-close title="创建容器" width="650px">
    <el-input
      v-model="runParam"
      placeholder="这里把上传的时候的参数读取出来，方便修改"
      resize="none"
      :rows="5"
      type="textarea"
    />
    <template #footer>
      <el-button type="primary" @click="createContainer">创建容器</el-button>
    </template>
  </vab-dialog>
  <!-- 上传镜像 -->
  <vab-dialog v-model="uploadImaageVisible" destroy-on-close title="上传镜像" width="650px">
    <AlgorithmsUploadImage ref="upload_image_ref" />
    <template #footer>
      <el-button :loading="isLoading" type="primary" @click="onUploadSubmit">确认</el-button>
      <el-button @click="onUploadCancel">取消</el-button>
    </template>
  </vab-dialog>
</template>

<style scoped lang="scss">
  :deep() {
    .startImage .el-dialog__footer {
      text-align: center;
    }
  }
</style>
