<script lang="ts">
  export default {
    // 激活页面
    name: 'Active',
  }
  const enum DeploymentType {
    Mysql = 1,
    InitUser = 2,
    Es = 3,
    Hugegraph = 4,
    Kafka = 5,
    Alarm = 6,
    Restart = 7,
  }
</script>
<script setup lang="ts">
  import { FormInstance, UploadInstance } from 'element-plus'
  import { getCheckActivateApi } from '~/src/api-ecs/public'
  import {
    getActivateKeyApi,
    toActivateApi,
    getDeploymentCheckStatusApi,
    deploymentCheckApi,
  } from '~/src/api-ecs/public'
  import { useSettingsStore } from '~/src/store/modules/settings'
  import { requireRules } from '~/src/utils/rules'
  const settingsStore = useSettingsStore()
  const showActiveContent = ref(true)
  const router = useRouter()
  const $baseMessage: any = inject('$baseMessage')
  const upload = ref<UploadInstance>()
  // 1代表mysql检测完毕 2代表初始化  3 是ES检测完成 4是Hugegraph 5是kafka检测完成
  const deploymentCheckStatus = ref(DeploymentType.Restart)
  const deploymentFormRef = ref<FormInstance>()
  const deploymentForm = ref()
  const deploymentFormRules = ref()
  const nextDisabled = ref(true)
  const deploymentLoading = ref(false)
  const deploymentData = {
    [DeploymentType.Mysql]: {
      ip: '',
      port: '',
      dbName: '',
      userName: '',
      password: '',
    },
    [DeploymentType.InitUser]: {
      loginName: '',
      nickName: '',
      mail: '',
      phone: '',
      password: '',
    },
    [DeploymentType.Es]: {
      ip: '',
      port: '',
      userName: '',
      password: '',
    },
    [DeploymentType.Hugegraph]: {
      ip: '',
      port: '',
    },
    [DeploymentType.Kafka]: {
      ip: '',
      port: '',
      userName: '',
      password: '',
    },
    [DeploymentType.Alarm]: {
      address: '',
    },
  }
  const deploymentRule = {
    [DeploymentType.Mysql]: {
      ip: requireRules,
      port: requireRules,
      dbName: requireRules,
      userName: requireRules,
      password: requireRules,
    },
    [DeploymentType.InitUser]: {
      loginName: requireRules,
      nickName: requireRules,
      password: requireRules,
    },
    [DeploymentType.Es]: {
      ip: requireRules,
      port: requireRules,
      userName: requireRules,
      password: requireRules,
    },
    [DeploymentType.Hugegraph]: {
      ip: requireRules,
      port: requireRules,
    },
    [DeploymentType.Kafka]: {
      ip: requireRules,
      port: requireRules,
      // userName: requireRules,
      // password: requireRules,
    },
    [DeploymentType.Alarm]: {
      address: requireRules,
    },
  }

  // 客户名称
  const customerName = ref('')
  // 激活码
  const machineCode = ref('')
  // 特征码
  const featureCode = ref(['', ''])
  const productKey = reactive(['', '', '', '', ''])

  // 获取机器码特征码
  const handleGetProductKey = async () => {
    const { data } = await getActivateKeyApi()
    customerName.value = data.customerName || ''
    featureCode.value = data.featureCode.split('-')
    machineCode.value = data.machineCode
  }

  // 提交成功
  const handleSuccess = ($val: any) => {
    const { code, msg } = $val
    if (code === 20) {
      $baseMessage($val.msg, 'success', 'vab-hey-message-success')
      handleCheckActivate(true)
    } else {
      $baseMessage(msg, 'error', 'vab-hey-message-error')
      deploymentLoading.value = false
    }
  }
  // 上传文件激活
  const submitUpload = () => {
    upload.value?.submit()
  }
  // 输入密钥激活
  const productKeyToActivate = () => {
    toActivateApi({
      customerName: customerName.value,
      machineCode: machineCode.value,
      featureCode: featureCode.value.join('-'),
      productKey: productKey.join('-'),
    }).then((res: any) => {
      if (res.code * 1 === 20) {
        $baseMessage(res.msg, 'success', 'vab-hey-message-success')
        showActiveContent.value = false
        deploymentLoading.value = false
        handleCheckActivate(true)
      }
    })
  }
  const nextDeploymentStatus = async () => {
    deploymentCheckStatus.value += 1
    nextDisabled.value = deploymentCheckStatus.value === DeploymentType.Alarm ? false : true
  }
  const submitDeploymentForm = async () => {
    deploymentLoading.value = true
    // 当初始化状态是6时，就是要重启了
    if (deploymentCheckStatus.value === DeploymentType.Restart) {
      try {
        await deploymentCheckApi(+deploymentCheckStatus.value, undefined)
        router.replace('/login')
      } finally {
        deploymentLoading.value = false
      }
    } else {
      // 如果不是就Alarm要执行相应的初始化操作
      deploymentFormRef.value?.validate(async (valid, fields) => {
        if (valid) {
          deploymentCheckApi(+deploymentCheckStatus.value, deploymentForm.value)
            .then(({ msg }) => {
              $baseMessage(msg, 'success', 'vab-hey-message-success')
              //如果是Alarm的话直接进入下一步，也就是重启
              if (deploymentCheckStatus.value === DeploymentType.Alarm) {
                deploymentCheckStatus.value += 1
              }
              nextDisabled.value = false
            })
            .finally(() => {
              deploymentLoading.value = false
            })
        } else {
          deploymentLoading.value = false
        }
      })
    }
  }
  const getDeploymentCheckStatus = async () => {
    const { data } = await getDeploymentCheckStatusApi()
    const status = +data.status
    if (status === DeploymentType.Restart) {
      router.replace('/login')
    } else {
      showActiveContent.value = false
      deploymentLoading.value = false
      deploymentCheckStatus.value = status + 1
      nextDisabled.value = deploymentCheckStatus.value === DeploymentType.Alarm ? false : true
    }
  }
  const handleCheckActivate = async (showErr = false) => {
    deploymentLoading.value = true
    const { data } = await getCheckActivateApi()
    const expUsed = data.expUsed * 1000
    const nowDate = new Date().getTime()
    if (expUsed > nowDate) {
      getDeploymentCheckStatus()
    } else {
      deploymentLoading.value = false
      if (showErr) {
        showActiveContent.value = true
        $baseMessage('激活文件/密钥已失效', 'error', 'vab-hey-message-error')
      }
    }
  }
  watch(
    () => deploymentCheckStatus.value,
    () => {
      if (deploymentCheckStatus.value === DeploymentType.Restart) return
      const status = deploymentCheckStatus.value
      deploymentForm.value = deploymentData[status]
      deploymentFormRules.value = deploymentRule[status]
    },
    {
      immediate: true,
    }
  )

  onMounted(() => {
    handleCheckActivate()
    handleGetProductKey()
  })
</script>

<template>
  <div class="active">
    <div class="active-box" v-if="showActiveContent" v-loading="deploymentLoading">
      <el-tabs class="active-tabs" :model-value="'file'">
        <el-tab-pane label="文件激活" name="file">
          <el-form label-width="100px" label-position="right" class="activeForm">
            <el-form-item label="机器码:">
              <el-input :value="machineCode" readonly />
            </el-form-item>
            <el-form-item label="">
              <el-upload
                class="upload-demo"
                ref="upload"
                action="/v3/eht/public/fileActivation"
                name="filedata"
                :limit="1"
                :on-success="handleSuccess"
                :before-upload="
                  () => {
                    deploymentLoading = true
                    return true
                  }
                "
                :auto-upload="false"
              >
                <template #trigger>
                  <el-button type="primary" @click="upload?.clearFiles()">选取文件</el-button>
                </template>
                <el-button style="margin-left: 10px" class="ml-3" type="success" @click="submitUpload">
                  激活设备
                </el-button>
                <template #tip>
                  <div class="el-upload__tip">只能上传bin文件，且不超过500b</div>
                </template>
              </el-upload>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        <el-tab-pane label="序列号激活" name="productKey">
          <el-form label-width="100px" label-position="right" class="activeForm">
            <el-form-item label="客户名称:">
              <el-input v-model="customerName" />
            </el-form-item>
            <el-form-item label="机器码">
              <el-input :value="machineCode" readonly />
            </el-form-item>
            <el-form-item label="特征码">
              <el-space>
                <el-input :value="featureCode[0]" readonly />
                <span>-</span>
                <el-input :value="featureCode[1]" readonly />
              </el-space>
            </el-form-item>
            <el-form-item label="产品密钥">
              <el-space>
                <el-input v-model="productKey[0]" maxlength="5" />
                <span>-</span>
                <el-input v-model="productKey[1]" maxlength="5" />
                <span>-</span>
                <el-input v-model="productKey[2]" maxlength="5" />
                <span>-</span>
                <el-input v-model="productKey[3]" maxlength="5" />
                <span>-</span>
                <el-input v-model="productKey[4]" maxlength="5" />
              </el-space>
            </el-form-item>
            <el-form-item label=" ">
              <el-button type="success" @click="productKeyToActivate">激活设备</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </div>
    <div class="deployment-box" v-else>
      <header>第三方组件检测</header>
      <div class="deployment-nav" v-if="deploymentCheckStatus < DeploymentType.Restart">
        <el-timeline style="max-width: 600px">
          <el-timeline-item
            center
            :hide-timestamp="true"
            icon="CircleCheckFilled"
            :class="{ isActive: deploymentCheckStatus === DeploymentType.Mysql }"
          >
            <div class="navBar">Mysql</div>
          </el-timeline-item>
          <el-timeline-item
            center
            :hide-timestamp="true"
            icon="CircleCheckFilled"
            :class="{ isActive: deploymentCheckStatus === DeploymentType.InitUser }"
          >
            <div class="navBar">初始化登录人</div>
          </el-timeline-item>
          <el-timeline-item
            center
            :hide-timestamp="true"
            icon="CircleCheckFilled"
            :class="{ isActive: deploymentCheckStatus === DeploymentType.Es }"
          >
            <div class="navBar">ES</div>
          </el-timeline-item>
          <el-timeline-item
            center
            :hide-timestamp="true"
            icon="CircleCheckFilled"
            :class="{ isActive: deploymentCheckStatus === DeploymentType.Hugegraph }"
          >
            <div class="navBar">Hugegraph</div>
          </el-timeline-item>
          <el-timeline-item
            center
            :hide-timestamp="true"
            icon="CircleCheckFilled"
            :class="{ isActive: deploymentCheckStatus === DeploymentType.Kafka }"
          >
            <div class="navBar">Kafka</div>
          </el-timeline-item>
          <el-timeline-item
            center
            :hide-timestamp="true"
            icon="CircleCheckFilled"
            :class="{ isActive: deploymentCheckStatus === DeploymentType.Alarm }"
          >
            <div class="navBar">威胁检测</div>
          </el-timeline-item>
        </el-timeline>
      </div>
      <div
        class="deployment-form"
        :class="{ 'deployment-form-only': deploymentCheckStatus === DeploymentType.Restart }"
      >
        <el-form
          ref="deploymentFormRef"
          :model="deploymentForm"
          :rules="deploymentFormRules"
          label-position="top"
          :validate-on-rule-change="false"
        >
          <template v-if="deploymentCheckStatus === DeploymentType.Mysql">
            <el-form-item label="地址" prop="ip">
              <el-input v-model="deploymentForm.ip" />
            </el-form-item>
            <el-form-item label="端口" prop="port">
              <el-input v-model="deploymentForm.port" />
            </el-form-item>
            <el-form-item label="数据库名称" prop="dbName">
              <el-input v-model="deploymentForm.dbName" />
            </el-form-item>
            <el-form-item label="用户名" prop="userName">
              <el-input v-model="deploymentForm.userName" />
            </el-form-item>
            <el-form-item label="密码" prop="password">
              <el-input v-model="deploymentForm.password" type="password" show-password />
            </el-form-item>
          </template>
          <template v-if="deploymentCheckStatus === DeploymentType.InitUser">
            <el-form-item label="用户姓名" prop="loginName">
              <el-input v-model="deploymentForm.loginName" />
            </el-form-item>
            <el-form-item label="真实姓名" prop="nickName">
              <el-input v-model="deploymentForm.nickName" />
            </el-form-item>
            <el-form-item label="邮箱地址" prop="mail">
              <el-input v-model="deploymentForm.mail" />
            </el-form-item>
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="deploymentForm.phone" />
            </el-form-item>
            <el-form-item label="密码" prop="password">
              <el-input v-model="deploymentForm.password" type="password" show-password />
            </el-form-item>
          </template>
          <template v-if="deploymentCheckStatus === DeploymentType.Es">
            <el-form-item label="地址" prop="ip">
              <el-input v-model="deploymentForm.ip" />
            </el-form-item>
            <el-form-item label="端口号" prop="port">
              <el-input v-model="deploymentForm.port" />
            </el-form-item>
            <el-form-item label="用户名" prop="userName">
              <el-input v-model="deploymentForm.userName" />
            </el-form-item>
            <el-form-item label="密码" prop="password">
              <el-input v-model="deploymentForm.password" type="password" show-password />
            </el-form-item>
          </template>
          <template v-if="deploymentCheckStatus === DeploymentType.Hugegraph">
            <el-form-item label="地址" prop="ip">
              <el-input v-model="deploymentForm.ip" />
            </el-form-item>
            <el-form-item label="端口号" prop="port">
              <el-input v-model="deploymentForm.port" />
            </el-form-item>
          </template>
          <template v-if="deploymentCheckStatus === DeploymentType.Kafka">
            <el-form-item label="地址" prop="ip">
              <el-input v-model="deploymentForm.ip" />
            </el-form-item>
            <el-form-item label="端口号" prop="port">
              <el-input v-model="deploymentForm.port" />
            </el-form-item>
            <el-form-item label="用户名" prop="userName">
              <el-input v-model="deploymentForm.userName" />
            </el-form-item>
            <el-form-item label="密码" prop="password">
              <el-input v-model="deploymentForm.password" type="password" show-password />
            </el-form-item>
          </template>
          <template v-if="deploymentCheckStatus === DeploymentType.Alarm">
            <el-form-item label="地址" prop="address">
              <template #label>
                地址
                <el-tooltip content="多个地址使用逗号进行分隔" placement="right">
                  <vab-icon icon="question-line" />
                </el-tooltip>
              </template>
              <el-input v-model="deploymentForm.address" />
            </el-form-item>
          </template>
          <template v-if="deploymentCheckStatus === DeploymentType.Restart">
            <div class="restart">
              <el-icon><CircleCheckFilled /></el-icon>
              <h1>检测完成</h1>
              <p>第三方检测已通过，请点击下方按钮重启系统</p>
              <el-button
                type="primary"
                class="submit-btn deployment-btn"
                :loading="deploymentLoading"
                @click="submitDeploymentForm"
              >
                重启系统
              </el-button>
            </div>
          </template>
        </el-form>
        <footer v-if="deploymentCheckStatus < DeploymentType.Restart">
          <span>提示：检测成功后可进行下一步</span>
          <el-button
            type="primary"
            class="next-btn deployment-btn"
            :disabled="nextDisabled"
            @click="nextDeploymentStatus"
            v-if="deploymentCheckStatus <= DeploymentType.Alarm"
          >
            下一步
          </el-button>
          <el-button
            :loading="deploymentLoading"
            type="primary"
            class="submit-btn deployment-btn"
            @click="submitDeploymentForm"
          >
            检测
          </el-button>
        </footer>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .active {
    height: 100vh;
    background: url('~@/assets/login_images/background.jpg') center center fixed no-repeat;
    background-size: cover;

    .active-box {
      width: 700px;
      height: 500px;
      background: linear-gradient(120deg, rgba(68, 100, 245, 0.3) 0%, rgba(47, 48, 174, 0.2) 100%);
      backdrop-filter: blur(3px);
      overflow: hidden;
      border-radius: 30px;
      position: fixed;
      top: calc(50% - 250px);
      right: 6vw;
      // padding: 100px;
      .active-tabs {
        padding: 30px 40px;
      }
      :deep() {
        .el-tabs__item {
          color: #a8aebb;
        }

        .el-tabs__active-bar {
          background-color: #fff;
        }
        .el-tabs__nav-wrap::after {
          background-color: #a8aebb;
        }
        .el-form-item__label,
        .el-upload__tip,
        .el-tabs__item.is-active {
          color: #fff;
        }
        .el-tabs__header,
        .el-form-item {
          margin-bottom: 30px;
        }
        .el-space {
          span {
            color: #fff;
            margin-inline: 8px;
          }
          .el-space__item {
            margin-right: 0 !important;
          }
        }
        .upload-demo {
          width: 100%;

          .el-upload-list__item {
            background-color: var(--el-fill-color-light);
          }
        }
      }
    }
    .deployment-box {
      width: 860px;
      position: fixed;
      height: 680px;
      background: #ffffff;
      border-radius: 10px;
      top: calc(50% - 340px);
      right: 6vw;
      overflow: hidden;
      .restart {
        text-align: center;
        margin-top: 66px;
        .el-icon {
          color: #3cba6b;
          font-size: 68px;
        }
        h1 {
          font-weight: 500;
          font-size: 28px;
          color: #303133;
        }
        p {
          font-weight: 400;
          font-size: 14px;
          color: #606266;
        }
        .el-button {
          background: #6954f0;
        }
      }
      header {
        height: 75px;
        line-height: 75px;
        border-bottom: 1px solid #eeeef0;
        font-weight: 500;
        font-size: 18px;
        color: #303133;
        text-indent: 40px;
      }
      .deployment-nav {
        width: 254px;
        height: 604px;
        border-right: 1px solid #eeeef0;
        float: left;
        padding-top: 30px;

        .nav-item {
          height: 42px;
          width: 184px;
          margin-left: 40px;
        }

        :deep() {
          .el-timeline-item__wrapper {
            cursor: pointer;
          }
          .el-timeline-item__icon {
            font-size: 16px;
            color: #d2cfdf;
          }
          .el-timeline-item__tail {
            border-left-style: dashed;
          }
          .el-timeline-item__node {
            top: 14px;
          }
          .el-timeline-item__content {
            width: fit-content;
            font-family: PingFangSC, PingFang SC;
            font-size: 14px;
            font-weight: 500;
            color: #606266;
          }
          .navBar {
            position: relative;
            padding-inline: 14px;
            width: 150px;
            height: 42px;
            line-height: 42px;
            border-radius: 4px;
            transition: all 0.3s;
          }
          .isActive {
            .el-timeline-item__node {
              background: var(--el-color-white);
              border-color: #6954f0;
              border-style: solid;
              border-width: 3px;
            }
            .el-timeline-item__icon {
              color: #6954f0;
            }
            .navBar {
              color: #6954f0;
              background-color: #efeefe;
              &::before {
                opacity: 1;
              }
            }
          }
        }
      }
      .deployment-form {
        margin-left: 255px;
        height: 604px;
        margin-right: 40px;
        padding: 30px 0 30px 30px;
        &.deployment-form-only {
          margin-inline: 0;
        }
        footer {
          position: absolute;
          bottom: 40px;
          width: 535px;
          span {
            font-weight: 400;
            font-size: 14px;
            line-height: 36px;
            color: #a9acb3;
          }
          .deployment-btn {
            height: 36px;
            font-size: 14px;
            border: 0;
            width: 88px;
            float: right;
            &.submit-btn {
              background: #6954f0;
            }
            &:hover {
              opacity: 0.9;
            }
            &.next-btn {
              background: #6954f0;
              color: #fff;
              margin-left: 15px;
            }
            &.is-disabled {
              color: #74737d;
              background: #eeedf5;
            }
          }
        }
      }
    }
  }
</style>
