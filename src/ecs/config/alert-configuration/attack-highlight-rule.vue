<script lang="ts">
  export default {
    name: 'AttackHighlightRule',
  }
</script>

<script setup lang="ts">
  import { FormInstance } from 'element-plus'
  import { updateAttackHighlightConfigApi, updateAttackHighlightWhiteApi } from '@/api-ecs/attack-characterization'
  import { AttackCharacterizationType } from '@/types/index'
  const $baseMessage: any = inject('$baseMessage')

  const props = defineProps<{
    attackHighlightData: Partial<AttackCharacterizationType> & { type?: 'highLightConfig' | 'highLightWhite' }
  }>()
  const emits = defineEmits<{
    (e: 'ecs-dialog-close'): void
    (e: 'ecs-dialog-confirm'): void
  }>()
  const attackRuleForm = ref<FormInstance>()
  const attackRuleData: AttackCharacterizationType & { type: 'highLightConfig' | 'highLightWhite' } = reactive(
    Object.assign({ content: '', remark: '', type: 'highLightConfig', scope: 'all' }, props.attackHighlightData)
  )
  const attackFormRules = reactive({
    content: [
      { required: true, message: '请输入高亮特征', trigger: 'blur' },
      { max: 35, message: '高亮特征最多35位', trigger: 'change' },
    ],
    type: { required: true, trigger: 'change' },
    scope: { required: true, trigger: 'change' },
  })

  const attackRuleFormSubmit = async () => {
    attackRuleForm.value?.validate(async (valid) => {
      if (!valid) return
      const ajax = { highLightConfig: updateAttackHighlightConfigApi, highLightWhite: updateAttackHighlightWhiteApi }
      const initType = props.attackHighlightData.id ? props.attackHighlightData.type : attackRuleData.type
      const { msg } = await ajax[initType as 'highLightConfig' | 'highLightWhite'](attackRuleData)
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      emits('ecs-dialog-confirm')
    })
  }
  onMounted(() => {
    props.attackHighlightData.content && attackRuleForm.value?.validate()
  })
</script>

<template>
  <div class="attack-rule-container">
    <el-form
      ref="attackRuleForm"
      label-position="top"
      :model="attackRuleData"
      :rules="attackFormRules"
      validate-on-rule-change
    >
      <el-form-item label="高亮字段" prop="content">
        <el-input v-model="attackRuleData.content" placeholder="请输入高亮字段" />
      </el-form-item>
      <el-form-item label="添加到" prop="type">
        <el-radio-group v-model="attackRuleData.type">
          <el-radio label="highLightConfig">高亮特征库</el-radio>
          <el-radio label="highLightWhite">特征库白名单</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="作用域" prop="scope">
        <el-select v-model="attackRuleData.scope" style="width: 100%">
          <el-option label="全部" value="all" />
          <el-option label="请求头" value="requestHeader" />
          <el-option label="请求体" value="requestPayload" />
          <el-option label="响应头" value="responseHeader" />
          <el-option label="响应体" value="responsePayload" />
        </el-select>
      </el-form-item>
      <el-form-item label="备注" prop="remark">
        <el-input
          v-model="attackRuleData.remark"
          :autosize="{ minRows: 4, maxRows: 4 }"
          placeholder="请输入备注"
          resize="none"
          type="textarea"
        />
      </el-form-item>
    </el-form>
    <div class="attack-rule-footer">
      <el-button :auto-insert-space="false" type="primary" @click="attackRuleFormSubmit">确定</el-button>
      <el-button :auto-insert-space="false" @click="() => emits('ecs-dialog-close')">取消</el-button>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .attack-rule-container {
    .attack-rule-footer {
      text-align: right;
      .el-button {
        border-radius: 4px;
      }
    }
  }
</style>
