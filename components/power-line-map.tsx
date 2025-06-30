"use client"

import { useEffect, useRef } from "react"
import { MapPin, DrillIcon as Drone, AlertTriangle } from "lucide-react"

export function PowerLineMap() {
  const mapRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    // This would be replaced with actual map initialization code
    // using libraries like Leaflet or Mapbox
    if (mapRef.current) {
      const ctx = document.createElement("canvas").getContext("2d")
      if (!ctx) return

      // Simulate map rendering with a placeholder
      const img = new Image()
      img.crossOrigin = "anonymous"
      img.src = "/placeholder.svg?height=400&width=800"
      img.onload = () => {
        if (!mapRef.current) return
        const placeholder = document.createElement("div")
        placeholder.className = "relative h-[400px] w-full bg-muted"

        // Add map markers
        const markers = [
          { type: "drone", lat: 30, lng: 40, label: "无人机 #1" },
          { type: "drone", lat: 35, lng: 45, label: "无人机 #2" },
          { type: "alert", lat: 32, lng: 42, label: "绝缘子损坏" },
          { type: "powerline", lat: 33, lng: 43, label: "主干线 #A45" },
        ]

        markers.forEach((marker, i) => {
          const el = document.createElement("div")
          el.className = "absolute flex items-center justify-center"
          el.style.top = `${marker.lat}%`
          el.style.left = `${marker.lng}%`

          const icon = document.createElement("div")
          icon.className = "flex h-8 w-8 items-center justify-center rounded-full"

          if (marker.type === "drone") {
            icon.className += " bg-blue-100"
            icon.innerHTML = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" class="text-blue-600"><path d="M12 8V4H8"></path><path d="M12 4h4"></path><path d="M16 8h4"></path><path d="M20 12V8"></path><path d="M20 12h-4"></path><path d="M4 12h4"></path><path d="M16 16h4"></path><path d="M20 20v-4"></path><path d="M12 20v-4"></path><path d="M8 20h4"></path><path d="M4 20v-4"></path><path d="M4 16h4"></path><path d="M4 12V8"></path><path d="M8 8H4"></path></svg>`
          } else if (marker.type === "alert") {
            icon.className += " bg-red-100"
            icon.innerHTML = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" class="text-red-600"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"></path><path d="M12 9v4"></path><path d="M12 17h.01"></path></svg>`
          } else {
            icon.className += " bg-green-100"
            icon.innerHTML = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" class="text-green-600"><path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z"></path><circle cx="12" cy="10" r="3"></circle></svg>`
          }

          const label = document.createElement("span")
          label.className = "absolute top-full mt-1 whitespace-nowrap text-xs font-medium"
          label.textContent = marker.label

          el.appendChild(icon)
          el.appendChild(label)
          placeholder.appendChild(el)
        })

        // Add power lines (simplified representation)
        const lines = [
          { x1: "20%", y1: "30%", x2: "40%", y2: "50%" },
          { x1: "40%", y1: "50%", x2: "60%", y2: "40%" },
          { x1: "60%", y1: "40%", x2: "80%", y2: "60%" },
        ]

        lines.forEach((line) => {
          const el = document.createElement("div")
          el.className = "absolute bg-green-600 h-0.5"
          el.style.top = line.y1
          el.style.left = line.x1
          el.style.width = "calc(" + line.x2 + " - " + line.x1 + ")"
          el.style.transform = `rotate(${Math.atan2(
            Number.parseInt(line.y2) - Number.parseInt(line.y1),
            Number.parseInt(line.x2) - Number.parseInt(line.x1),
          )}rad)`
          el.style.transformOrigin = "0 0"
          placeholder.appendChild(el)
        })

        mapRef.current.innerHTML = ""
        mapRef.current.appendChild(placeholder)
      }
    }
  }, [])

  return (
    <div className="relative h-[400px] overflow-hidden rounded-b-lg">
      <div ref={mapRef} className="h-full w-full bg-muted">
        <div className="flex h-full items-center justify-center">
          <p className="text-sm text-muted-foreground">加载地图中...</p>
        </div>
      </div>
      <div className="absolute bottom-4 right-4 flex gap-2">
        <div className="flex items-center gap-1 rounded-md bg-background/80 px-2 py-1 text-xs backdrop-blur-sm">
          <Drone className="h-3 w-3 text-blue-600" />
          <span>无人机</span>
        </div>
        <div className="flex items-center gap-1 rounded-md bg-background/80 px-2 py-1 text-xs backdrop-blur-sm">
          <MapPin className="h-3 w-3 text-green-600" />
          <span>电力线</span>
        </div>
        <div className="flex items-center gap-1 rounded-md bg-background/80 px-2 py-1 text-xs backdrop-blur-sm">
          <AlertTriangle className="h-3 w-3 text-red-600" />
          <span>告警</span>
        </div>
      </div>
    </div>
  )
}
