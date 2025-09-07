<script lang="ts">
  export default {
    name: 'SituationList', //
  }
</script>

<script setup lang="ts">
  import { ArrowLeft } from '@element-plus/icons-vue'
  import { setSystemConfigApi } from '~/src/api-ecs/system'
  import AesEncryptCBC from '~/src/utils/crypto'

  const emit = defineEmits<{
    (e: 'on-back-event', str: string): void
    (e: 'on-to-other', row: any): void
    (e: 'on-reflash'): void
  }>()

  const props = defineProps<{
    listDate: any
  }>()

  const listDates = ref<any>([])

  watch(
    () => props.listDate,
    () => {
      listDates.value = []
      const data = JSON.parse(JSON.stringify(props.listDate))
      data?.forEach((item: any) => {
        const obj = { id: undefined }
        obj.id = item.id
        if (item.value) {
          for (const key in item.value) {
            // @ts-ignore
            obj[key] = item.value[key]
          }
        }
        listDates.value.push(obj)
      })
    },
    { immediate: true, deep: true }
  )

  const handleChange = async (row: { id: number }, enable: boolean) => {
    const password = ref('')
    ElMessageBox({
      title: '提示',
      showCancelButton: true,
      customClass: 'need-password-message-box',
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      customStyle: {
        maxWidth: '500px',
      },
      message: () =>
        h('div', null, [
          h('div', { style: 'margin: 15px 0 5px 0;color:#55585b' }, '请输入敏感操作密码:(通过验证后方可进行操作)'),
          h(ElInput, {
            type: 'password',
            modelValue: password.value,
            placeholder: '请输入敏感操作密码',
            showPassword: true,
            style: 'margin-block: 10px',
            'onUpdate:modelValue': (val: string) => {
              password.value = val
            },
          }),
        ]),

      beforeClose: async (action, instance, done) => {
        if (action === 'confirm') {
          instance.confirmButtonLoading = true
          const formData = props.listDate
          const res = formData?.find((item: any, index: number) => {
            return item.id == row.id
          })
          res.value.enable = enable
          const { msg } = await setSystemConfigApi(
            { id: res.id, value: JSON.stringify(res.value) },
            { password: AesEncryptCBC(password.value) }
          )
          emit('on-reflash')
          ElMessage({ message: msg, type: 'success' })
          done()
        } else {
          done()
        }
      },
    }).catch(() => {})
  }

  const handleClick = () => {
    emit('on-back-event', 'chapter_1')
  }

  const handleConfig = (row: any) => {
    const formData = props.listDate
    const res = formData?.find((item: any, index: number) => {
      return item.id == row.id
    })
    emit('on-back-event', 'chapter_3')
    emit('on-to-other', res)
  }
</script>

<template>
  <div class="situation-list-container">
    <el-button :icon="ArrowLeft" @click="handleClick">返回</el-button>
    <el-table :border="true" class="field" :data="listDates" style="margin-top: 15px">
      <el-table-column align="center" fixed label="缩略图" prop="thumbnail" :resizable="true" show-overflow-tooltip>
        <template #default="{ row }">
          <el-image alt="" class="img" :src="row.thumbnail" />
        </template>
      </el-table-column>
      <el-table-column align="center" label="名称" prop="title" :resizable="true" show-overflow-tooltip />
      <el-table-column align="center" label="是否启用" prop="enable" :resizable="true" show-overflow-tooltip>
        <template #default="{ row }">
          <el-switch v-model="row.enable" @change="handleChange(row, row.enable)" />
        </template>
      </el-table-column>
      <el-table-column align="center" fixed="right" label="操作" :resizable="true" show-overflow-tooltip width="180">
        <template #default="{ row }">
          <el-button class="row_action" size="small" @click="handleConfig(row)">配置</el-button>
        </template>
      </el-table-column>
      <template #empty>
        <div style="height: 300px; line-height: 300px">暂无其他数据</div>
      </template>
    </el-table>
  </div>
</template>

<style scoped lang="scss">
  .field {
    .img {
      width: 128px !important;
      height: 72px !important;
      img {
        height: 100%;
        width: 100%;
      }
    }
  }
</style>
