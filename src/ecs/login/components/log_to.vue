<script lang="ts">
  export default {
    name: 'LogTo', //
  }
</script>

<script setup lang="ts">
  import type { FormInstance } from 'element-plus'
  import { LoginOtherType } from '@/types/index'
  import sha1 from 'sha1'
  import { requireRules } from '~/src/utils/rules'
  const formRef = ref<FormInstance>()

  const props = defineProps<{
    visible: boolean
    mode: number
  }>()

  const emit = defineEmits<{
    (e: 'on-closeEvent'): void
    (e: 'on-loginEvent', form: LoginOtherType): void
  }>()

  const form = reactive<LoginOtherType>({
    username: '',
    password: '',
    authType: 1,
    loginType: '1',
  })

  const dialogFormVisible = ref(false)

  const loading = ref(false)

  const rules = {
    username: requireRules,
    password: requireRules,
  }

  // 登录
  const handleConfirm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl!.validate(async (valid) => {
      if (valid) {
        emit('on-loginEvent', form)
      } else {
        return
      }
    })
  }

  const handleClose = () => {
    emit('on-closeEvent')
  }

  onMounted(() => {
    dialogFormVisible.value = props.visible
    form.authType = props.mode
  })
</script>

<template>
  <div class="log_to">
    <el-dialog
      v-model="dialogFormVisible"
      center
      :close-on-click-modal="false"
      title="登录"
      width="500"
      @close="handleClose"
    >
      <el-form ref="formRef" label-width="70px" :model="form" :rules="rules">
        <el-form-item label="账号：" prop="username">
          <el-input v-model="form.username" autocomplete="off" />
        </el-form-item>
        <el-form-item label="密码：" prop="password">
          <el-input v-model="form.password" autocomplete="off" type="password" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <div class="my_btns">
            <el-button class="login-btn" :loading="loading" type="primary" @click="handleConfirm(formRef)">
              登录
            </el-button>
          </div>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  .my_btns {
    .login-btn {
      background: var(--el-color-primary);
      &:hover {
        opacity: 0.9 !important;
      }
    }
  }
</style>
