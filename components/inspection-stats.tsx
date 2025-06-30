"use client"

import { BarChart, Bar, XAxis, YAxis, CartesianGrid, ResponsiveContainer } from "recharts"
import { ChartContainer, ChartTooltip, ChartTooltipContent } from "@/components/ui/chart"

const data = [
  { date: "6/1", inspections: 12, issues: 2 },
  { date: "6/5", inspections: 18, issues: 3 },
  { date: "6/10", inspections: 15, issues: 1 },
  { date: "6/15", inspections: 22, issues: 4 },
  { date: "6/20", inspections: 20, issues: 2 },
  { date: "6/25", inspections: 25, issues: 5 },
  { date: "6/30", inspections: 17, issues: 3 },
]

export function InspectionStats() {
  return (
    <ChartContainer
      config={{
        inspections: {
          label: "检查次数",
          color: "hsl(var(--chart-1))",
        },
        issues: {
          label: "发现问题",
          color: "hsl(var(--chart-2))",
        },
      }}
      className="h-[300px]"
    >
      <ResponsiveContainer width="100%" height="100%">
        <BarChart data={data}>
          <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
          <XAxis dataKey="date" className="text-xs" />
          <YAxis className="text-xs" />
          <ChartTooltip content={<ChartTooltipContent />} />
          <Bar dataKey="inspections" fill="var(--color-inspections)" radius={[4, 4, 0, 0]} />
          <Bar dataKey="issues" fill="var(--color-issues)" radius={[4, 4, 0, 0]} />
        </BarChart>
      </ResponsiveContainer>
    </ChartContainer>
  )
}
