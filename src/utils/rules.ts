import type { FormItemRule } from 'element-plus'
//  手机号
const mobileRules: FormItemRule[] = [
  // 必填
  { required: true, message: '请输入手机号' },
  // 手机号
  { pattern: /^1[3-9]\d{9}$/, message: '手机号不合法' },
]

//  密码框
const passwordRules: FormItemRule[] = [
  // 必填
  { required: true, message: '请输入密码' },
  // 密码：8~24
  { pattern: /^\w{8,24}$/, message: '密码不合法' },
]

// 必须项
const requireRules: FormItemRule[] = [
  // 必填
  { required: true, message: '该选项为必填项' },
]

// 必须项且为数字
const requireAndNumberRules: FormItemRule[] = [
  { pattern: /^\d+$/, message: '请输入数字', trigger: 'blur' },
  {
    required: true,
    message: '该选项为必填项',
    trigger: 'blur',
  },
]

export { mobileRules, passwordRules, requireRules, requireAndNumberRules }
