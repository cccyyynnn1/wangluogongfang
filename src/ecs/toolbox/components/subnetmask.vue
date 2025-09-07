<script lang="ts">
  export default {
    name: 'Subnetmask',
  }
</script>

<script setup lang="ts">
  const numberInputProps = {
    max: 255,
    min: 0,
    valueOnClear: 0,
    controls: false,
  }
  const ipProt = ref([192, 168, 1, 0])
  const maskBit = ref(24)
  const ipTotal = ref()

  const subnetMask = ref('')
  const firstIp = ref('')
  const lastIp = ref('')

  const formatStrToArrBySize = (str: string, size = 8) => {
    const _arr: string[] = []
    let index = 0
    while (index < str.length) {
      _arr.push(str.slice(index, size + index))
      index += size
    }
    return _arr
  }

  const formatBinaryStrToArr = (binaryIp: string) => formatStrToArrBySize(binaryIp).map((i) => parseInt(i, 2))

  const computedHandle = () => {
    if (maskBit.value > 32) maskBit.value = 32
    const length = 32 - maskBit.value

    // ip转二进制
    const binaryIp = ipProt.value.map((i) => i.toString(2).padStart(8, '0')).join('')
    // 子网掩码（二进制）
    const maskBitCode = Array(maskBit.value).fill(1).join('').padEnd(32, '0')
    // 网络号(binaryIp和子网掩码每一项and操作)
    const subnet_mask = [...binaryIp].map((v, i) => +v & +maskBitCode[i]).join('')
    // 网络号10进制
    const _subnetMaskIp = formatBinaryStrToArr(subnet_mask)

    // 广播地址（网络号右侧补0 ）
    const broadcastAddress = subnet_mask.slice(0, maskBit.value).padEnd(32, '1')
    // 广播地址10进制
    const _broadcastAddressIp = formatBinaryStrToArr(broadcastAddress)
    // 可用IP数量
    ipTotal.value = maskBit.value === 32 ? 1 : 2 ** length - 2 // 减去网络号广播地址
    // 子网掩码IP
    subnetMask.value = formatBinaryStrToArr(maskBitCode).join('.')
    if (maskBit.value < 32) {
      // 第一个可用IP
      _subnetMaskIp[3] += 1
      // 最后一个可用IP
      _broadcastAddressIp[3] -= 1
    }
    firstIp.value = _subnetMaskIp.join('.')
    lastIp.value = _broadcastAddressIp.join('.')
  }
  onMounted(() => computedHandle())
</script>

<template>
  <div class="subnetmask">
    <div class="subnetmask-ip">
      <el-input-number v-model="ipProt[0]" v-bind="{ ...numberInputProps }" class="w-66" size="large" />
      <span class="symbol">.</span>
      <el-input-number v-model="ipProt[1]" v-bind="{ ...numberInputProps }" class="w-66" size="large" />
      <span class="symbol">.</span>
      <el-input-number v-model="ipProt[2]" v-bind="{ ...numberInputProps }" class="w-66" size="large" />
      <span class="symbol">.</span>
      <el-input-number v-model="ipProt[3]" v-bind="{ ...numberInputProps }" class="w-66" size="large" />
      <span class="symbol">/</span>
      <el-input-number v-model="maskBit" v-bind="{ ...numberInputProps, min: 0 }" class="w-66" size="large" />
      <el-button class="computed" size="large" type="primary" @click="computedHandle">计算</el-button>
    </div>
    <el-descriptions border :column="1">
      <el-descriptions-item label="子网掩码" label-class-name="subnetmask-label">
        {{ subnetMask }}
      </el-descriptions-item>
      <el-descriptions-item label="第一个可用IP" label-class-name="subnetmask-label">
        {{ firstIp }}
      </el-descriptions-item>
      <el-descriptions-item label="最后一个可用IP" label-class-name="subnetmask-label">
        {{ lastIp }}
      </el-descriptions-item>
      <el-descriptions-item label="可用IP数量" label-class-name="subnetmask-label">{{ ipTotal }}</el-descriptions-item>
    </el-descriptions>
  </div>
</template>

<style scoped lang="scss">
  .subnetmask {
    .subnetmask-ip {
      margin: 70px auto;
      width: 543px;
    }
    :deep() {
      .subnetmask-label {
        width: 166px;
      }
    }
    .symbol {
      margin: 0 10px;
      vertical-align: -webkit-baseline-middle;
    }
    .w-66 {
      width: 66px;
    }
    .computed {
      width: 100px;
      vertical-align: baseline;
      margin-left: 14px;
    }
  }
</style>
