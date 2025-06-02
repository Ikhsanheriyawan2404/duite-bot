"use client"

import { Button } from "@/components/ui/button"
import { Calendar } from "@/components/ui/calendar"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { cn } from "@/lib/utils"
import { CalendarIcon, XIcon } from "lucide-react"
import { format } from "date-fns"
import type { DateRange } from "react-day-picker"

interface TransactionFiltersProps {
  dateRange: DateRange | undefined
  onDateRangeChange: (range: DateRange | undefined) => void
  selectedType: string
  onSelectedTypeChange: (type: string) => void
  selectedCategory: string
  onSelectedCategoryChange: (category: string) => void
  categories: string[]
  onResetFilters: () => void
}

export function TransactionFilters({
  dateRange,
  onDateRangeChange,
  selectedType,
  onSelectedTypeChange,
  selectedCategory,
  onSelectedCategoryChange,
  categories,
  onResetFilters,
}: TransactionFiltersProps) {
  return (
    <div className="mb-6 p-4 border rounded-lg bg-card text-card-foreground shadow-sm">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 items-end">
        <div>
          <label htmlFor="date-range" className="block text-sm font-medium mb-1">
            Rentang Tanggal
          </label>
          <Popover>
            <PopoverTrigger asChild>
              <Button
                id="date-range"
                variant={"outline"}
                className={cn("w-full justify-start text-left font-normal", !dateRange && "text-muted-foreground")}
              >
                <CalendarIcon className="mr-2 h-4 w-4" />
                {dateRange?.from ? (
                  dateRange.to ? (
                    <>
                      {format(dateRange.from, "LLL dd, y")} - {format(dateRange.to, "LLL dd, y")}
                    </>
                  ) : (
                    format(dateRange.from, "LLL dd, y")
                  )
                ) : (
                  <span>Pilih periode tanggal</span>
                )}
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-auto p-0" align="start">
              <Calendar
                initialFocus
                mode="range"
                defaultMonth={dateRange?.from}
                selected={dateRange}
                onSelect={onDateRangeChange}
                numberOfMonths={2}
              />
            </PopoverContent>
          </Popover>
        </div>

        <div>
          <label htmlFor="type-filter" className="block text-sm font-medium mb-1">
            Tipe
          </label>
          <Select value={selectedType} onValueChange={onSelectedTypeChange}>
            <SelectTrigger id="type-filter">
              <SelectValue placeholder="Filter by type" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">Semua Tipe</SelectItem>
              <SelectItem value="income">Pemasukan</SelectItem>
              <SelectItem value="expense">Pengeluaran</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div>
          <label htmlFor="category-filter" className="block text-sm font-medium mb-1">
            Kategori
          </label>
          <Select value={selectedCategory} onValueChange={onSelectedCategoryChange}>
            <SelectTrigger id="category-filter">
              <SelectValue placeholder="Filter by category" />
            </SelectTrigger>
            <SelectContent>
              {categories.map((category) => (
                <SelectItem key={category} value={category}>
                  {category === "all" ? "Semua Kategori" : category}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <Button onClick={onResetFilters} variant="outline" className="w-full md:w-auto self-end">
          <XIcon className="mr-2 h-4 w-4" /> Reset Filter
        </Button>
      </div>
    </div>
  )
}
