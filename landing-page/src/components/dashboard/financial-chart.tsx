"use client"

import { Bar, BarChart, CartesianGrid, ResponsiveContainer, XAxis, YAxis } from "recharts"
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  ChartLegend,
  ChartLegendContent,
} from "@/components/ui/chart"
import type { MonthlyData } from "@/lib/types"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"

interface FinancialChartProps {
  data: MonthlyData[]
}

const chartConfig = {
  income: {
    label: "Pemasukan",
    color: "hsl(var(--chart-2))",
  },
  expense: {
    label: "Pengeluaran",
    color: "hsl(var(--chart-1))",
  },
}

export function FinancialChart({ data }: FinancialChartProps) {
  if (!data || data.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Ringkasan Bulanan</CardTitle>
          <CardDescription>Pendapatan dan pengeluaran per bulan.</CardDescription>
        </CardHeader>
        <CardContent className="min-h-[300px] flex items-center justify-center">
          <p className="text-muted-foreground">Tidak ada data yang tersedia untuk grafik.</p>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Ringkasan Bulanan</CardTitle>
        <CardDescription>Pendapatan dan pengeluaran per bulan.</CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig} className="w-full min-h-[300px]">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={data} margin={{ top: 5, right: 20, left: 2, bottom: 5 }}>
              <CartesianGrid strokeDasharray="3 3" vertical={false} />
              <XAxis dataKey="month" tickLine={false} axisLine={false} tickMargin={8} />
              <YAxis tickLine={false} axisLine={false} tickMargin={8} />
              <ChartTooltip cursor={false} content={<ChartTooltipContent indicator="dot" />} />
              <ChartLegend content={<ChartLegendContent />} />
              <Bar dataKey="income" fill={chartConfig.income.color} radius={4} />
              <Bar dataKey="expense" fill={chartConfig.expense.color} radius={4} />
            </BarChart>
          </ResponsiveContainer>
        </ChartContainer>
      </CardContent>
    </Card>
  )
}
