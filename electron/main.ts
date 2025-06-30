import { app, BrowserWindow, Menu, ipcMain, dialog, shell, Tray, nativeImage } from "electron"
import { join } from "path"
import { spawn } from "child_process"
import * as fs from "fs"
import * as os from "os"

// 保持对窗口对象的全局引用
let mainWindow: BrowserWindow | null = null
let tray: Tray | null = null
let goBackendProcess: any = null

const isDev = process.env.NODE_ENV === "development"
const isWin = process.platform === "win32"
const isMac = process.platform === "darwin"

// 创建主窗口
function createMainWindow(): void {
  mainWindow = new BrowserWindow({
    width: 1400,
    height: 900,
    minWidth: 1200,
    minHeight: 800,
    show: false,
    icon: join(__dirname, "../assets/icon.png"),
    titleBarStyle: isMac ? "hiddenInset" : "default",
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      enableRemoteModule: false,
      preload: join(__dirname, "preload.js"),
      webSecurity: !isDev,
    },
  })

  // 加载应用
  if (isDev) {
    mainWindow.loadURL("http://localhost:3000")
    mainWindow.webContents.openDevTools()
  } else {
    mainWindow.loadFile(join(__dirname, "../dist/index.html"))
  }

  // 窗口事件处理
  mainWindow.once("ready-to-show", () => {
    mainWindow?.show()

    if (isDev) {
      mainWindow?.webContents.openDevTools()
    }
  })

  mainWindow.on("closed", () => {
    mainWindow = null
  })

  // 处理外部链接
  mainWindow.webContents.setWindowOpenHandler(({ url }) => {
    shell.openExternal(url)
    return { action: "deny" }
  })

  // 阻止导航到外部URL
  mainWindow.webContents.on("will-navigate", (event, navigationUrl) => {
    const parsedUrl = new URL(navigationUrl)

    if (parsedUrl.origin !== "http://localhost:3000" && !isDev) {
      event.preventDefault()
    }
  })
}

// 创建系统托盘
function createTray(): void {
  const iconPath = join(__dirname, "../assets/tray-icon.png")
  const trayIcon = nativeImage.createFromPath(iconPath)

  tray = new Tray(trayIcon.resize({ width: 16, height: 16 }))

  const contextMenu = Menu.buildFromTemplate([
    {
      label: "显示主窗口",
      click: () => {
        if (mainWindow) {
          mainWindow.show()
          mainWindow.focus()
        } else {
          createMainWindow()
        }
      },
    },
    {
      label: "射频扫描",
      click: () => {
        if (mainWindow) {
          mainWindow.webContents.send("navigate-to", "/rf-security")
          mainWindow.show()
          mainWindow.focus()
        }
      },
    },
    {
      label: "设备管理",
      click: () => {
        if (mainWindow) {
          mainWindow.webContents.send("navigate-to", "/assets")
          mainWindow.show()
          mainWindow.focus()
        }
      },
    },
    { type: "separator" },
    {
      label: "关于",
      click: () => {
        dialog.showMessageBox(mainWindow!, {
          type: "info",
          title: "关于",
          message: "电力无人机安全扫描系统",
          detail: "Version 1.0.0\n专业的无人机设备安全漏洞检测与评估平台",
          buttons: ["确定"],
        })
      },
    },
    {
      label: "退出",
      click: () => {
        app.quit()
      },
    },
  ])

  tray.setToolTip("电力无人机安全扫描系统")
  tray.setContextMenu(contextMenu)

  // 双击托盘图标显示窗口
  tray.on("double-click", () => {
    if (mainWindow) {
      mainWindow.show()
      mainWindow.focus()
    } else {
      createMainWindow()
    }
  })
}

// 创建应用菜单
function createMenu(): void {
  const template: Electron.MenuItemConstructorOptions[] = [
    {
      label: "文件",
      submenu: [
        {
          label: "新建扫描任务",
          accelerator: "CmdOrCtrl+N",
          click: () => {
            mainWindow?.webContents.send("create-scan-task")
          },
        },
        {
          label: "导入资产",
          accelerator: "CmdOrCtrl+I",
          click: async () => {
            const result = await dialog.showOpenDialog(mainWindow!, {
              properties: ["openFile"],
              filters: [
                { name: "CSV Files", extensions: ["csv"] },
                { name: "JSON Files", extensions: ["json"] },
                { name: "All Files", extensions: ["*"] },
              ],
            })

            if (!result.canceled && result.filePaths.length > 0) {
              mainWindow?.webContents.send("import-assets", result.filePaths[0])
            }
          },
        },
        {
          label: "导出报告",
          accelerator: "CmdOrCtrl+E",
          click: async () => {
            const result = await dialog.showSaveDialog(mainWindow!, {
              defaultPath: `security-report-${new Date().toISOString().split("T")[0]}.pdf`,
              filters: [
                { name: "PDF Files", extensions: ["pdf"] },
                { name: "HTML Files", extensions: ["html"] },
                { name: "JSON Files", extensions: ["json"] },
              ],
            })

            if (!result.canceled && result.filePath) {
              mainWindow?.webContents.send("export-report", result.filePath)
            }
          },
        },
        { type: "separator" },
        {
          label: "设置",
          accelerator: "CmdOrCtrl+,",
          click: () => {
            mainWindow?.webContents.send("navigate-to", "/settings")
          },
        },
        { type: "separator" },
        {
          label: "退出",
          accelerator: isMac ? "Cmd+Q" : "Ctrl+Q",
          click: () => {
            app.quit()
          },
        },
      ],
    },
    {
      label: "扫描",
      submenu: [
        {
          label: "开始全面扫描",
          accelerator: "F5",
          click: () => {
            mainWindow?.webContents.send("start-full-scan")
          },
        },
        {
          label: "射频安全扫描",
          accelerator: "F6",
          click: () => {
            mainWindow?.webContents.send("start-rf-scan")
          },
        },
        {
          label: "停止所有扫描",
          accelerator: "Escape",
          click: () => {
            mainWindow?.webContents.send("stop-all-scans")
          },
        },
      ],
    },
    {
      label: "工具",
      submenu: [
        {
          label: "网络诊断",
          click: () => {
            mainWindow?.webContents.send("open-network-diagnostic")
          },
        },
        {
          label: "射频分析器",
          click: () => {
            mainWindow?.webContents.send("open-rf-analyzer")
          },
        },
        {
          label: "日志查看器",
          click: () => {
            mainWindow?.webContents.send("open-log-viewer")
          },
        },
        { type: "separator" },
        {
          label: "数据库管理",
          click: () => {
            mainWindow?.webContents.send("open-database-manager")
          },
        },
      ],
    },
    {
      label: "窗口",
      submenu: [
        {
          label: "最小化",
          accelerator: "CmdOrCtrl+M",
          click: () => {
            mainWindow?.minimize()
          },
        },
        {
          label: "关闭",
          accelerator: "CmdOrCtrl+W",
          click: () => {
            mainWindow?.close()
          },
        },
        { type: "separator" },
        {
          label: "重新加载",
          accelerator: "CmdOrCtrl+R",
          click: () => {
            mainWindow?.reload()
          },
        },
        {
          label: "强制重新加载",
          accelerator: "CmdOrCtrl+Shift+R",
          click: () => {
            mainWindow?.webContents.reloadIgnoringCache()
          },
        },
        {
          label: "开发者工具",
          accelerator: "F12",
          click: () => {
            mainWindow?.webContents.toggleDevTools()
          },
        },
      ],
    },
    {
      label: "帮助",
      submenu: [
        {
          label: "用户手册",
          click: () => {
            shell.openExternal("https://docs.example.com/user-manual")
          },
        },
        {
          label: "技术支持",
          click: () => {
            shell.openExternal("https://support.example.com")
          },
        },
        { type: "separator" },
        {
          label: "检查更新",
          click: () => {
            mainWindow?.webContents.send("check-for-updates")
          },
        },
        {
          label: "关于",
          click: () => {
            dialog.showMessageBox(mainWindow!, {
              type: "info",
              title: "关于电力无人机安全扫描系统",
              message: "电力无人机安全扫描系统",
              detail: `版本: 1.0.0
平台: ${os.platform()} ${os.arch()}
Node.js: ${process.versions.node}
Electron: ${process.versions.electron}
Chrome: ${process.versions.chrome}

专业的无人机设备安全漏洞检测与评估平台
支持网络基础设施、射频安全、设备管理等功能`,
              buttons: ["确定"],
            })
          },
        },
      ],
    },
  ]

  // macOS 特殊处理
  if (isMac) {
    template.unshift({
      label: app.getName(),
      submenu: [
        {
          label: "关于 " + app.getName(),
          role: "about",
        },
        { type: "separator" },
        {
          label: "服务",
          role: "services",
        },
        { type: "separator" },
        {
          label: "隐藏 " + app.getName(),
          accelerator: "Command+H",
          role: "hide",
        },
        {
          label: "隐藏其他",
          accelerator: "Command+Shift+H",
          role: "hideothers",
        },
        {
          label: "显示全部",
          role: "unhide",
        },
        { type: "separator" },
        {
          label: "退出",
          accelerator: "Command+Q",
          click: () => {
            app.quit()
          },
        },
      ],
    })
  }

  const menu = Menu.buildFromTemplate(template)
  Menu.setApplicationMenu(menu)
}

// 启动Go后端服务
function startGoBackend(): void {
  const goExecutable = isDev
    ? join(__dirname, "../backend/main")
    : join(process.resourcesPath, "backend", isWin ? "main.exe" : "main")

  if (fs.existsSync(goExecutable)) {
    goBackendProcess = spawn(goExecutable, [], {
      stdio: "pipe",
      env: {
        ...process.env,
        PORT: "8080",
        DB_PATH: join(app.getPath("userData"), "database.db"),
      },
    })

    goBackendProcess.stdout?.on("data", (data: Buffer) => {
      console.log("Go Backend:", data.toString())
    })

    goBackendProcess.stderr?.on("data", (data: Buffer) => {
      console.error("Go Backend Error:", data.toString())
    })

    goBackendProcess.on("close", (code: number) => {
      console.log("Go Backend process exited with code", code)
    })
  } else {
    console.warn("Go backend executable not found:", goExecutable)
  }
}

// IPC 事件处理
function setupIpcHandlers(): void {
  // 文件操作
  ipcMain.handle("show-save-dialog", async (event, options) => {
    const result = await dialog.showSaveDialog(mainWindow!, options)
    return result
  })

  ipcMain.handle("show-open-dialog", async (event, options) => {
    const result = await dialog.showOpenDialog(mainWindow!, options)
    return result
  })

  ipcMain.handle("write-file", async (event, filePath, data) => {
    try {
      fs.writeFileSync(filePath, data)
      return { success: true }
    } catch (error) {
      return { success: false, error: error.message }
    }
  })

  ipcMain.handle("read-file", async (event, filePath) => {
    try {
      const data = fs.readFileSync(filePath, "utf8")
      return { success: true, data }
    } catch (error) {
      return { success: false, error: error.message }
    }
  })

  // 系统信息
  ipcMain.handle("get-system-info", async () => {
    return {
      platform: os.platform(),
      arch: os.arch(),
      version: os.release(),
      memory: os.totalmem(),
      cpus: os.cpus().length,
      userDataPath: app.getPath("userData"),
      appVersion: app.getVersion(),
    }
  })

  // 通知
  ipcMain.handle("show-notification", async (event, options) => {
    const { Notification } = require("electron")

    if (Notification.isSupported()) {
      const notification = new Notification({
        title: options.title,
        body: options.body,
        icon: options.icon || join(__dirname, "../assets/icon.png"),
      })

      notification.show()
      return { success: true }
    }

    return { success: false, error: "Notifications not supported" }
  })

  // 窗口控制
  ipcMain.handle("minimize-window", () => {
    mainWindow?.minimize()
  })

  ipcMain.handle("maximize-window", () => {
    if (mainWindow?.isMaximized()) {
      mainWindow.unmaximize()
    } else {
      mainWindow?.maximize()
    }
  })

  ipcMain.handle("close-window", () => {
    mainWindow?.close()
  })

  // 硬件设备访问
  ipcMain.handle("get-serial-ports", async () => {
    try {
      const { SerialPort } = require("serialport")
      const ports = await SerialPort.list()
      return { success: true, ports }
    } catch (error) {
      return { success: false, error: error.message }
    }
  })

  // 网络接口信息
  ipcMain.handle("get-network-interfaces", async () => {
    try {
      const interfaces = os.networkInterfaces()
      return { success: true, interfaces }
    } catch (error) {
      return { success: false, error: error.message }
    }
  })
}

// 应用事件处理
app.whenReady().then(() => {
  createMainWindow()
  createTray()
  createMenu()
  setupIpcHandlers()
  startGoBackend()

  app.on("activate", () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createMainWindow()
    }
  })
})

app.on("window-all-closed", () => {
  if (!isMac) {
    app.quit()
  }
})

app.on("before-quit", () => {
  // 关闭Go后端进程
  if (goBackendProcess) {
    goBackendProcess.kill()
  }
})

// 安全设置
app.on("web-contents-created", (event, contents) => {
  contents.on("new-window", (event, navigationUrl) => {
    event.preventDefault()
    shell.openExternal(navigationUrl)
  })

  contents.on("will-attach-webview", (event, webPreferences, params) => {
    // 禁用webview
    event.preventDefault()
  })

  contents.on("will-navigate", (event, navigationUrl) => {
    const parsedUrl = new URL(navigationUrl)

    if (parsedUrl.origin !== "http://localhost:3000" && !isDev) {
      event.preventDefault()
    }
  })
})
