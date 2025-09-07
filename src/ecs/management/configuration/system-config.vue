<script lang="ts">
  export default {
    name: 'SystemConfig',
  }
</script>
<script setup lang="ts">
  import AuthenticationConfig from './components/authentication-config.vue'

  import StorageDuration from './components/storage-duration.vue'

  import SnmpConfig from './components/snmp-config.vue'

  import TimeConfig from './components/time-config.vue'

  import LadpConfig from './components/ladp-config.vue'

  import radiusConfig from './components/radius-config.vue'

  import ssoConfig from './components/sso-config.vue'

  import safeConfig from './components/safe-config.vue'

  import Notification from './components/notification.vue'

  import MessageConfig from './components/message-config.vue'

  import SyslogConfig from './components/syslog-config.vue'

  import ExternalStorage from './components/external-storage.vue'

  import MailConfig from './components/mail-config.vue'

  import ModelConfig from './components/model-config.vue'

  import DateInfo from './components/date-info.vue'

  import { getSystemConfigApi } from '~/src/api-ecs/system'
  const route = useRoute()
  const activeName = ref((route.query.params as string) || 'syslogConfig') // tabs选中项

  type objType = {
    [key: string]: any
  }

  const obj: objType = {
    ladp: 'ecs_ldap',
    mail: 'ecs_mail',
    sso: 'ecs_sso',
    radius: 'ecs_radius',
    message: 1,
    storage: 2,
    syslog: 'ecs_syslog',
    notification: 5,
    ecsDataTime: 'ecs_data_time',
    moduleConfig: 'module_config',
  }

  const configData = reactive<objType>({
    syslog: undefined,
    ladp: undefined,
    sso: undefined,
    radius: undefined,
    mail: undefined,
    message: undefined,
    storage: undefined,
    notification: undefined,
    ecsDataTime: undefined,
    moduleConfig: undefined,
  })

  const KEYS = 'ecs_syslog,ecs_mail,ecs_ldap,ecs_sso,ecs_radius,ecs_data_time,module_config'
  const getData = async () => {
    const { data } = await getSystemConfigApi({ keys: KEYS })
    for (const key in configData) {
      configData[key] = data[obj[key]]
    }
  }

  onMounted(() => {
    getData()
  })
</script>

<template>
  <div class="system-config-container">
    <el-tabs v-model="activeName">
      <el-tab-pane label="时间设置" name="timeConfig">
        <time-config />
      </el-tab-pane>
      <el-tab-pane label="AI配置" name="modelConfig">
        <model-config :all-data="configData.moduleConfig?.value" />
      </el-tab-pane>
      <el-tab-pane label="SYSLOG服务设置" name="syslogConfig">
        <syslog-config :id="configData.syslog?.id" :syslog="configData.syslog?.value" />
      </el-tab-pane>
      <!-- <el-tab-pane label="KAFKA设置" name="kafkaConfig">
              <kafka-config />
            </el-tab-pane> -->
      <!-- <el-tab-pane label="SNMP设置" name="snmpConfig">
        <snmp-config />
      </el-tab-pane> -->
      <el-tab-pane label="邮件设置" name="mailConfig">
        <mail-config :id="configData.mail?.id" :mail="configData.mail?.value" />
      </el-tab-pane>
      <el-tab-pane label="LDAP配置" name="ladpConfig">
        <ladp-config :id="configData.ladp?.id" :ladp="configData.ladp?.value" />
      </el-tab-pane>
      <el-tab-pane label="RADIUS配置" name="radiusConfig">
        <radius-config :id="configData.radius?.id" :radius="configData.radius?.value" />
      </el-tab-pane>
      <el-tab-pane label="SSO配置" name="ssoConfig">
        <sso-config :id="configData.sso?.id" :sso="configData.sso?.value" />
      </el-tab-pane>
      <el-tab-pane label="安全设置" name="safeConfig">
        <safe-config />
      </el-tab-pane>
      <el-tab-pane label="数据存储时长" name="ecs_data_time">
        <StorageDuration :id="configData.ecsDataTime?.id" :ecs-data-time="configData.ecsDataTime?.value" />
      </el-tab-pane>
      <el-tab-pane label="日历信息维护" name="ecs_date_info">
        <DateInfo :id="configData.ecsDataTime?.id" :ecs-data-time="configData.ecsDataTime?.value" />
      </el-tab-pane>
      <!-- <el-tab-pane label="短信设置" name="phoneConfig">
        <message-config :message="configData.message?.value" />
      </el-tab-pane>
      <el-tab-pane label="外挂设置" name="plugin">
        <external-storage :storage="configData.storage?.value" />
      </el-tab-pane>
      <el-tab-pane label="通知设置" name="notification">
        <notification :notification="configData.notification?.value" />
      </el-tab-pane> -->
      <!-- <el-tab-pane label="认证设置" name="authenticationConfig">
        <authentication-config />
      </el-tab-pane> -->
    </el-tabs>
  </div>
</template>

<style lang="scss" scoped></style>
