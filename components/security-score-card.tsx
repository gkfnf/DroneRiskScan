import { Card, CardContent } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"
import type { ReactNode } from "react"
import { ArrowDown, ArrowUp } from "lucide-react"

interface SecurityScoreCardProps {
  title: string
  score: number
  change: number
  icon: ReactNode
  description: string
}

export function SecurityScoreCard({ title, score, change, icon, description }: SecurityScoreCardProps) {
  const getScoreColor = (score: number) => {
    if (score >= 80) return "text-green-600"
    if (score >= 60) return "text-amber-600"
    return "text-red-600"
  }

  const getProgressColor = (score: number) => {
    if (score >= 80) return "bg-green-600"
    if (score >= 60) return "bg-amber-600"
    return "bg-red-600"
  }

  return (
    <Card>
      <CardContent className="pt-6">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            {icon}
            <span className="font-medium">{title}</span>
          </div>
          <div className="flex items-center gap-1">
            {change !== 0 &&
              (change > 0 ? (
                <ArrowUp className="h-3 w-3 text-green-600" />
              ) : (
                <ArrowDown className="h-3 w-3 text-red-600" />
              ))}
            <span className={`text-xs ${change > 0 ? "text-green-600" : change < 0 ? "text-red-600" : ""}`}>
              {change > 0 ? "+" : ""}
              {change}
            </span>
          </div>
        </div>
        <div className="mt-4 flex items-end justify-between">
          <div className={`text-3xl font-bold ${getScoreColor(score)}`}>{score}</div>
          <div className="text-xs text-muted-foreground">满分100</div>
        </div>
        <Progress value={score} className="mt-2 h-1.5" indicatorClassName={getProgressColor(score)} />
        <p className="mt-2 text-xs text-muted-foreground">{description}</p>
      </CardContent>
    </Card>
  )
}
