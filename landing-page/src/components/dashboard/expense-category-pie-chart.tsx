"use client"

import * as React from "react"
import { Pie, PieChart, ResponsiveContainer, Cell } from "recharts"
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  ChartLegend,
  ChartLegendContent,
} from "@/components/ui/chart"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import type { CategoryExpenseData } from "@/lib/types"

interface ExpenseCategoryPieChartProps {
  data: CategoryExpenseData[]
}

const generateChartConfig = (data: CategoryExpenseData[]) => {
  const config: { [key: string]: { label: string; color: string } } = {}
  data.forEach((item, index) => {
    config[item.category] = {
      label: item.category,
      color: `hsl(var(--chart-${(index % 5) + 1}))`, // Cycle through 5 predefined chart colors
    }
  })
  return config
}

export function ExpenseCategoryPieChart({ data }: ExpenseCategoryPieChartProps) {
  const chartConfig = React.useMemo(() => generateChartConfig(data), [data])

  if (!data || data.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Expense Distribution</CardTitle>
          <CardDescription>Breakdown of expenses by category.</CardDescription>
        </CardHeader>
        <CardContent className="h-[300px] flex items-center justify-center">
          <p className="text-muted-foreground">No expense data available for chart.</p>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Expense Distribution</CardTitle>
        <CardDescription>Breakdown of expenses by category.</CardDescription>
      </CardHeader>
      <CardContent className="flex items-center justify-center">
        <ChartContainer config={chartConfig} className="mx-auto aspect-square h-[250px]">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <ChartTooltip cursor={false} content={<ChartTooltipContent hideLabel />} />
              <Pie data={data} dataKey="totalExpense" nameKey="category" innerRadius={60} strokeWidth={5}>
                {data.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={chartConfig[entry.category]?.color || `hsl(var(--chart-1))`} />
                ))}
              </Pie>
              <ChartLegend content={<ChartLegendContent nameKey="category" />} className="-translate-y-2" />
            </PieChart>
          </ResponsiveContainer>
        </ChartContainer>
      </CardContent>
    </Card>
  )
}
