import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import type { LucideIcon } from "lucide-react"

interface SummaryCardProps {
  title: string
  value: string
  icon: LucideIcon
  description?: string
  descriptionVariant?: 'default' | 'success' | 'danger' | 'info' // Tambahkan variant
}

export function SummaryCard({ title, value, icon: Icon, description, descriptionVariant = "default" }: SummaryCardProps) {
  
  const variantClasses = {
    default: 'text-muted-foreground',
    success: 'text-green-600 dark:text-green-400',
    danger: 'text-red-600 dark:text-red-400',
    info: 'text-blue-600 dark:text-blue-400'
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        <Icon className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{value}</div>
        {description && <p className={`text-xs ${variantClasses[descriptionVariant]}`}>{description}</p>}
      </CardContent>
    </Card>
  )
}
