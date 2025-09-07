<script lang="ts">
  export default {
    name: 'AddRole',
  }
</script>

<script setup lang="ts">
  import { useUserStore } from '@/store/modules/user'
  import { SysRole } from '@/types'
  import { ElTree } from 'element-plus'
  import { editSystemRoleApi } from '@/api-ecs/system'
  const $baseMessage: any = inject('$baseMessage')
  const props = defineProps<{
    showName: string
    dialogVal?: SysRole
  }>()
  const userStore = useUserStore()
  const emit = defineEmits<{
    (e: 'on-saveData', load: boolean): void
  }>()

  const treeRef = ref<InstanceType<typeof ElTree>>()
  const activeName = ref('addRole') // tabs选中项
  const loading = ref(false)
  const formData = ref<any>({ roleName: '', enable: 1, menuIds: [] as number[] }) // 表单数据
  const defaultProps = {
    children: 'children',
    label: ({ meta }: any) => {
      return meta.title
    },
    disabled: 'disabled',
  }

  // 保存
  const saveData = async () => {
    try {
      loading.value = true
      const selectKeys = treeRef.value?.getCheckedKeys(false) as number[]
      const loginMenuIdList = treeRef.value?.getHalfCheckedKeys() as number[]
      const { enable, roleName } = formData.value
      const selectedMenus = [...selectKeys, ...loginMenuIdList]
      if (selectedMenus.length === 0) return $baseMessage('请为该角色配置功能权限', 'error', 'vab-hey-message-error')
      await editSystemRoleApi({
        enable,
        menuIds: selectKeys,
        roleName,
        id: formData.value.id,
        loginMenuIdList: selectedMenus,
      })
      $baseMessage(formData.value.id ? '修改角色成功' : '新增角色成功', 'success', 'vab-hey-message-success')
      emit('on-saveData', true)
    } finally {
      loading.value = false
    }
  }
  const getId = (item: any) => {
    return [item.id, item.children.map(getId)]
  }
  watchEffect(() => {
    if (!props.dialogVal) return
    const { cloned } = useCloned(props.dialogVal)
    formData.value = cloned.value
  })
</script>

<template>
  <div class="add-role">
    <el-tabs v-model="activeName">
      <el-tab-pane :label="props.showName === 'add' ? '添加角色' : '编辑角色'" name="addRole">
        <el-row>
          <el-col :offset="4" :span="13">
            <el-form class="role-form" label-position="right" label-width="220px">
              <el-form-item label="角色名称" required>
                <el-input v-model="formData.roleName" maxlength="32" minlength="8" />
              </el-form-item>
              <el-form-item label="状态">
                <el-radio-group v-model="formData.enable">
                  <el-radio :label="1">启用</el-radio>
                  <el-radio :label="0">停用</el-radio>
                </el-radio-group>
              </el-form-item>
              <el-form-item label="功能权限" required>
                <el-tree
                  ref="treeRef"
                  :data="userStore.systemMenus"
                  :default-checked-keys="formData.menuIds"
                  node-key="id"
                  :props="defaultProps"
                  show-checkbox
                />
              </el-form-item>
              <el-form-item class="btn">
                <el-button :loading="loading" type="primary" @click="saveData">保存</el-button>
                <el-button @click="emit('on-saveData', false)">取消</el-button>
              </el-form-item>
            </el-form>
          </el-col>
        </el-row>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<style lang="scss" scoped>
  .add-role {
    .role-form {
      margin-top: 20px;
    }
  }
  :deep() {
    .el-tree {
      width: 600px;
      height: 500px;
      overflow: auto;
    }
  }
</style>
