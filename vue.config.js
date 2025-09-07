/* eslint-disable prettier/prettier */
/* eslint-disable @typescript-eslint/no-unused-vars */

/**
 * @description vue.config.js全局配置
 */
const JavaScriptObfuscator = require('webpack-obfuscator');
const TerserPlugin = require('terser-webpack-plugin');
const {
  baseURL,
  title,
  devPort,
  assetsDir,
  outputDir,
  lintOnSave,
  publicPath,
  transpileDependencies,
  proxyNet,
} = require('./src/config')
const dayjs = require('dayjs')
const pkg = require('./package.json')

const { resolve, relative } = require('path')
const { defineConfig } = require('@vue/cli-service')
const { createVuePlugin, createChainWebpack } = require('./library/build/index.ts')

const info = {
  ...pkg,
  lastBuildTime: dayjs().format('YYYY-MM-DD HH:mm:ss'),
}

process.env.VUE_APP_TITLE = title
process.env.VUE_APP_AUTHOR = pkg.author
process.env.VUE_APP_INFO = JSON.stringify(info)
process.env.VUE_APP_UPDATE_TIME = info.lastBuildTime
process.env.VUE_APP_GITHUB_USER_NAME = process.env.VUE_GITHUB_USER_NAME
process.env.VUE_APP_RANDOM = `${info.lastBuildTime}-${process.env.VUE_GITHUB_USER_NAME}`

module.exports = defineConfig({
  publicPath,
  assetsDir,
  outputDir,
  lintOnSave,
  transpileDependencies,
  devServer: {
    compress: false,
    client: {
      progress: false,
      overlay: {
        warnings: false,
        errors: true,
      },
    },
    hot: true,
    open: {
      target: [`http://localhost:${devPort}`],
    },
    port: devPort,
    // setupMiddlewares: require('./mock'),
    // 注释掉的地方是前端配置代理访问后端的示例，如无特别需求，不建议使用！！！
    // baseURL必须为/xxx，而不是后端服务器，请先了解代理逻辑，再设置前端代理
    // ！！！一定要注意！！！
    // 1、这里配置了跨域及代理只针对开发环境生效
    // 2、不建议你在前端配置跨域，建议你后端配置Allow-Origin,Method,Headers，放行token字段，一步到位
    // 3、后端配置了跨域，就不需要前端再配置，会发生Origin冲突
    // 4、webpack5版本前端配置代理无法与mock同时使用，如果一定要用前端代理，需注释setupMiddlewares: require('./mock')
    proxy: {
      [baseURL]: {
        target: `https://${proxyNet}/`,
        ws: false,
        changeOrigin: true,
        pathRewrite: {
          [`^${baseURL}`]: '',
        },
      },
    },
  },
  pwa: {
    workboxOptions: {
      skipWaiting: true,
      clientsClaim: true,
    },
    themeColor: '#ffffff',
    msTileColor: '#ffffff',
    appleMobileWebAppCapable: 'yes',
    appleMobileWebAppStatusBarStyle: 'black',
    manifestOptions: {
      name: 'Vue Admin Better - Admin Plus',
      short_name: 'Admin Plus',
      background_color: '#ffffff',
    },
  },
  configureWebpack(config) {
    const allPlugin = createVuePlugin()
    const _plugins =
      process.env.NODE_ENV === 'production'
        ? [
          ...allPlugin,
          new JavaScriptObfuscator(
            {
              // 压缩代码
              compact: true,
              // 是否启用控制流扁平化(降低1.5倍的运行速度)
              controlFlowFlattening: false,
              // 随机的死代码块(增加了混淆代码的大小)
              deadCodeInjection: false,
              // 此选项几乎不可能使用开发者工具的控制台选项卡
              // debugProtection: false,
              // 如果选中，则会在“控制台”选项卡上使用间隔强制调试模式，从而更难使用“开发人员工具”的其他功能。
              // debugProtectionInterval: false,
              // 通过用空函数替换它们来禁用console.log，console.info，console.error和console.warn。这使得调试器的使用更加困难。
              disableConsoleOutput: true,
              // 标识符的混淆方式 hexadecimal(十六进制) mangled(短标识符)
              identifierNamesGenerator: 'hexadecimal',
              log: false,
              // 是否启用全局变量和函数名称的混淆
              renameGlobals: false,
              // 通过固定和随机（在代码混淆时生成）的位置移动数组。这使得将删除的字符串的顺序与其原始位置相匹配变得更加困难。如果原始源代码不小，建议使用此选项，因为辅助函数可以引起注意。
              rotateStringArray: true,
              // 混淆后的代码,不能使用代码美化,同时需要配置 cpmpat:true;
              selfDefending: true,
              // 删除字符串文字并将它们放在一个特殊的数组中
              stringArray: true,
              // stringArrayEncoding: 'base64',
              stringArrayThreshold: 0.75,
              // 允许启用/禁用字符串转换为unicode转义序列。Unicode转义序列大大增加了代码大小，并且可以轻松地将字符串恢复为原始视图。建议仅对小型源代码启用此选项。
              unicodeEscapeSequence: false,
            },
            []
          ),
        ]
        : allPlugin
    return {
      resolve: {
        alias: {
          '~': resolve(__dirname, '.'),
          '@': resolve(__dirname, 'src'),
          '/#': resolve(__dirname, 'types'),
          '@vab': resolve(__dirname, 'library'),
          '@gp': resolve('library/plugins/vab'),
          'vue-i18n': 'vue-i18n/dist/vue-i18n.cjs.js',
        },
        fallback: {
          fs: false,
          path: require.resolve('path-browserify'),
        },
      },
      plugins: _plugins,
      performance: {
        hints: false,
      },
      module: {
        rules: [
          // 配置读取 *.md 文件的规则
          {
            test: /\.md$/,
            use: [{ loader: 'html-loader' }, { loader: 'markdown-loader', options: {} }],
          },
        ],
      },
      optimization:
        process.env.NODE_ENV === 'production'
          ? {
            minimize: true, // 是否压缩代码
            minimizer: [
              new TerserPlugin({
                terserOptions: {
                  compress: { drop_console: true }, // 移除所有的console.log语句
                  mangle: true, // 开启变量名混淆
                },
              }),
            ],
          }
          : undefined,
    }
  },
  chainWebpack(config) {
    //为了防止忘记配置而造成项目无法打包，请保留以下提示
    if (process.env.NODE_ENV === 'production') {
      if (process['env'].VUE_GITHUB_USER_NAME == 'test' && process['env'].VUE_APP_SECRET_KEY == 'preview')
        console.log('检测到您的用户名和key未配置，key在购买时通过邮件邀请函发放，请仔细阅读文档并进行配置')
    }
    createChainWebpack(process.env.NODE_ENV, config)
  },
  runtimeCompiler: false,
  productionSourceMap: false,
  css: {
    sourceMap: false,
    extract:
      process.env.NODE_ENV === 'production'
        ? {
          ignoreOrder: true,
        }
        : false,
    loaderOptions: {
      sass: {
        sassOptions: { outputStyle: 'expanded' },
        additionalData(content, { rootContext, resourcePath }) {
          const relativePath = relative(rootContext, resourcePath)
          if (relativePath.replace(/\\/g, '/') !== 'library/styles/variables/variables.module.scss')
            return `@use "~@vab/styles/variables/variables.module.scss" as *;${content}`
          return content
        },
      },
    },
  },
})
