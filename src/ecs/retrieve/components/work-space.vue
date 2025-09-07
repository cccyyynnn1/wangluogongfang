<script lang="ts">
  export default {
    name: 'RetrieveWorkSpace',
  }
</script>

<script setup lang="ts">
  import { getAllSpaceApi, saveOrUpdateSpaceApi, deleteSpaceApi, emptySpaceApi } from '@/api-ecs/retrieve'
  import { WorkerSpaceItem } from '@/types'
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')
  const visible = ref(false) // 显隐
  const modle = ref('add')
  const curSpace = reactive({
    id: 1,
    spaceName: '',
  })

  const props = defineProps<{
    showLink?: boolean
  }>()
  const spaceList = ref<WorkerSpaceItem[]>([])
  // 删除
  const handleDelete = (row: WorkerSpaceItem) => {
    $baseConfirm('你确定要删除当前项吗', null, async () => {
      const { code } = await deleteSpaceApi(row.id)
      const spaceId = localStorage.getItem('ecs-space')
      if (spaceId === row.id.toString()) {
        localStorage.setItem('ecs-space', '')
      }
      $baseMessage('删除工作空间成功', 'success', 'vab-hey-message-success')
      getAllSpace()
    })
  }
  // 删除
  const handleClear = (row: WorkerSpaceItem) => {
    $baseConfirm('你确定要清空该工作空间的数据吗？', null, async () => {
      const { code } = await emptySpaceApi(row.id)
      $baseMessage('清空工作空间成功', 'success', 'vab-hey-message-success')
      getAllSpace()
    })
  }

  const getAllSpace = async () => {
    try {
      const { data } = await getAllSpaceApi()
      spaceList.value = data
      const spaceId = localStorage.getItem('ecs-space')
      curSpace.id = spaceId ? +spaceId : spaceList.value[0].id
    } catch (error) {
      console.log(error)
    }
  }
  const showModel = (modleVal: string) => {
    modle.value = modleVal
    visible.value = true
  }

  const saveSpace = async () => {
    if (!curSpace.spaceName) return $baseMessage('工作空间名称不能为空', 'warning', 'vab-hey-message-warning')
    const { code } = await saveOrUpdateSpaceApi(curSpace.spaceName)
    $baseMessage('新建成功', 'success', 'vab-hey-message-success')
    getAllSpace()
    visible.value = false
  }
  const spaceChange = (id: number) => {
    localStorage.setItem('ecs-space', id.toString())
    $baseMessage('工作空间切换成功', 'success', 'vab-hey-message-success')
  }
  onMounted(() => {
    getAllSpace()
  })
  defineExpose({
    showModel,
    spaceId: toRef(curSpace, 'id'),
  })
</script>

<template>
  <div class="favorites-action">
    <el-select v-if="showLink" placeholder="请选择链路" style="margin-right: 10px" />
    <el-select v-model="curSpace.id" @change="spaceChange">
      <el-option v-for="item in spaceList" :key="item.id" :label="item.spaceName" :value="item.id" />
    </el-select>
    <el-divider direction="vertical" />
    <el-tooltip content="添加工作空间" effect="dark" placement="top">
      <el-icon :size="20" style="vertical-align: middle; color: #b3b9c8" @click="showModel('add')">
        <CirclePlus />
      </el-icon>
    </el-tooltip>
    <el-divider direction="vertical" />
    <el-tooltip content="编辑工作空间" effect="dark" placement="top">
      <el-icon :size="20" style="vertical-align: middle; color: #b3b9c8" @click="showModel('edit')"><Edit /></el-icon>
    </el-tooltip>
  </div>

  <el-dialog
    v-model="visible"
    :title="modle === 'add' ? '新建工作空间' : '管理工作空间'"
    :width="modle === 'add' ? '450px' : '600px'"
    @closed="() => (curSpace.spaceName = '')"
  >
    <template v-if="modle === 'add'">
      <p>
        名称：
        <el-input v-model="curSpace.spaceName" style="text-align: right; width: 344px" />
      </p>
      <p style="text-align: right; margin-right: 20px; margin-top: 10%">
        <el-button type="primary" @click="saveSpace">新建</el-button>
      </p>
    </template>
    <template v-else>
      <el-table border :data="spaceList" style="height: 400px">
        <el-table-column label="工作空间名称" prop="spaceName" />
        <el-table-column align="center" fixed="right" label="操作" width="140">
          <template #default="{ row }">
            <el-button :disabled="row.default" size="small" @click="() => handleDelete(row)">删除</el-button>
            <el-button size="small" @click="() => handleClear(row)">清空</el-button>
          </template>
        </el-table-column>
      </el-table>
    </template>
  </el-dialog>
</template>

<style scoped lang="scss"></style>
