import ElectronIntegration from "@/components/electron-integration"

export default function ElectronPage() {
  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-6xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">Electron 桌面应用功能</h1>
          <p className="text-gray-600">桌面应用专有功能和系统集成</p>
        </div>
        <ElectronIntegration />
      </div>
    </div>
  )
}
