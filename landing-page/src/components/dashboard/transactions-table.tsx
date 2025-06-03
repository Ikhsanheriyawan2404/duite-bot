"use client"

import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { MoreHorizontal, Edit2, Trash2 } from "lucide-react"
import type { Transaction } from "@/lib/types"
import { format } from "date-fns"

interface TransactionsTableProps {
  transactions: Transaction[]
  onEdit: (transaction: Transaction) => void
  onDelete: (transactionId: string) => void
}

export function TransactionsTable({ transactions, onEdit, onDelete }: TransactionsTableProps) {
  if (!transactions || transactions.length === 0) {
    return <p className="text-muted-foreground py-4">Belum ada data transaksi. Silakan tambahkan transaksi pertama Anda untuk memulai.</p>
  }

  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Tanggal</TableHead>
            <TableHead>Keterangan</TableHead>
            <TableHead>Kategori</TableHead>
            <TableHead className="text-right">Nominal</TableHead>
            <TableHead>Tipe</TableHead>
            <TableHead className="text-right">Aksi</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {transactions.map((transaction) => (
            <TableRow key={transaction.id}>
              <TableCell>{format(new Date(transaction.transaction_date), "MMM d, yyyy")}</TableCell>
              <TableCell className="font-medium">{transaction.original_text}</TableCell>
              <TableCell>
                <Badge variant="outline">{transaction.category}</Badge>
              </TableCell>
              <TableCell className={`text-right ${transaction.transaction_type === "INCOME" ? "text-green-600" : "text-red-600"}`}>
                {transaction.transaction_type === "INCOME" ? "+" : "-"}{transaction.amount.toLocaleString('id')}
              </TableCell>
              <TableCell>
                <Badge
                  variant={transaction.transaction_type === "INCOME" ? "default" : "destructive"}
                  className={
                    transaction.transaction_type === "INCOME"
                      ? "bg-green-100 text-green-700 hover:bg-green-200"
                      : "bg-red-100 text-red-700 hover:bg-red-200"
                  }
                >
                  {transaction.transaction_type === "INCOME" ? "Masuk" : "Keluar"}
                </Badge>
              </TableCell>
              <TableCell className="text-right">
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" className="h-8 w-8 p-0">
                      <span className="sr-only">Open menu</span>
                      <MoreHorizontal className="h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem onClick={() => onEdit(transaction)}>
                      <Edit2 className="mr-2 h-4 w-4" />
                      Edit
                    </DropdownMenuItem>
                    <DropdownMenuItem
                      onClick={() => onDelete(transaction.id)}
                      className="text-red-600 focus:text-red-600 focus:bg-red-50"
                    >
                      <Trash2 className="mr-2 h-4 w-4" />
                      Hapus
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
