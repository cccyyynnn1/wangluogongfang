# 专网组件_前端_终端

#### 软件架构
  基于Vue3 + Typescript + Pinia + ElementPlus技术栈开发，需要具备以下相关技术知识：

1. [Vue3](https://vuejs.org/)
2. [Typescript](https://www.typescriptlang.org/)
3. [Pinia](https://pinia.vuejs.org/)
4. [VueUse](https://vueuse.org/)
5. [Element Plus](https://element-plus.org/en-US/)

#### 安装教程

1.  git pull 
2.  cd 项目目录
3.  npm i 或者 yarn （node版本^16）

#### 开发说明
```
├─.vscode                         项目相关格式化配置（已经添加到.gitgnore）
├─library                         Vue Admin Plus公共组件代码
├─mock                            mock数据
├─public                          公共静态资源
├─types                           Vue Admin Plus类型文件
├─vab-icons                       vab Icon
├─src                             项目主文件
| ├─api-ecs                       接口API
| ├─assets                        静态资源
| ├─config                        项目配置
| | ├─cli.config.js               脚手架配置
| | ├─index.js                    配置文件入口
| | ├─net.config.js               网络接口配置
| | ├─setting.config              Vue Admin Plus相关配置
| | └theme.config                 主题相关配置
| ├─ecs                           项目路由页面
| | ├─alert                       告警页面
| | ├─login                       登陆页面
| | ├─register                    注册页面
| | ├─algorithms                  算法页面
| | ├─assets                      资源页面
| | ├─config                      配置页面
| | ├─management                  管理相关页面
| | ├─public                      系统激活页面
| | ├─retrieve                    检索页面
| | ├─site                        站点页面
| | └system                       系统页面
| ├─i18n                          国际化                      
| ├─icon                          雪碧图（library有使用到）
| ├─plugins                       vab二次封装组件
| ├─router                        项目路由（permissions里是权限路由）
| ├─store                         Pinia store
| ├─types                         项目相关类型文件
| └tils                           工具类
├─vue.config.js                   vue-cli配置
└README.md                        优先看
```
#### 打包部署

##### 注意事项

1. 第1步：.env 配置你的github用户名(如：VUE_GITHUB_USER_NAME=nan-ojtime)
2. 第2步：请在项目根目录新建一个.env.local的新文件，切记是新建空的文件不是直接拷贝.env文件的内容
3. 第3步：.env.local的文件只能有一行不可以换行，购买时生成，格式如下：VUE_APP_SECRET_KEY=XXXXXXX
4. 第4步：修改 .env.development里：VUE_APP_BASE_URL
5. 第5步：注释  src > main.ts 里 mockXHR 和 pwa
6. 第6步：项目根目录执行npm run build 或者 yarn build
7. 第七步 把dist里所有文件打包成.zip包。！！！注意：不能包含dist文件夹，只需要dist内部文件。
